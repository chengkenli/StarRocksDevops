package frontend

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "github.com/fatih/color"
    "github.com/go-resty/resty/v2"
    "os"
    "StarRocksDevops/tools"
    "StarRocksDevops/util"
    "strings"
    "sync"
    "time"
)

func init_meta() util.ConnectParms {
    var avg util.ConnectParms
    // connect
    if len(util.MetaLink) == 0 {
        return util.ConnectParms{}
    }
    for _, m := range util.MetaLink {
        if m["app"].(string) == util.P.App {
            var manager_access_key, manager_secret_key string
            if m["manager_access_key"] != nil {
                manager_access_key = m["manager_access_key"].(string)
            }
            if m["manager_secret_key"] != nil {
                manager_secret_key = m["manager_secret_key"].(string)
            }
            avg = util.ConnectParms{
                Host:       m["feip"].(string),
                Port:       int(m["feport"].(int32)),
                User:       m["user"].(string),
                Pass:       m["password"].(string),
                Uri:        m["address"].(string),
                AaccessKey: manager_access_key,
                SecretKey:  manager_secret_key,
            }
        }
    }
    return avg
}

func SessionError() {
    avg := init_meta()

    var mv map[string]interface{}
    util.Connect.Raw("select current_version() as version").Scan(&mv)
    version := tools.VToFloat(mv["version"].(string))

    c := color.New()
    fmt.Println(c.Add(color.FgHiYellow).Sprint("Query Session: "))

    var queryid, starttime string
    if version >= 2.5 {
        queryid = "queryId"
        starttime = "timestamp"
    } else {
        queryid = "query_id"
        starttime = "time"
    }
    os.Setenv("TZ", "Asia/Shanghai")
    var m []map[string]interface{}
    sql := fmt.Sprintf("select user,%s as queryid,%s as starttime from audit.starrocks_audit_log where %s >= date_sub(now(), interval %d hour) and user != 'root' and state='ERR' order by starttime desc", queryid, starttime, starttime, util.P.ErrorsTime)
    r := util.Connect.Raw(sql).Scan(&m)
    if r.Error != nil {
        fmt.Println(util.P.App + " -> " + r.Error.Error())
        return
    }
    if len(m) == 0 {
        return
    }

    fmt.Println(version)
    if version >= 3.3 {
        switch version {
        case 3.311:
            scan1(avg, m)
        case 3.39:
            scan2(avg, m)
        default:
            scan0(avg, m)
        }
    } else {
        if util.P.App == "sr-app" && version == 3.21 {
            scan1(avg, m)
        } else if util.P.App == "sccts" && version == 3.21 {
            scan2(avg, m)
        } else {
            scan0(avg, m)
        }
    }
}

// 这是最原始的方式，支持2.x-3.x
func scan0(avg util.ConnectParms, m []map[string]interface{}) {
    /*-------------------------登录manager web-------------------------------*/
    //创建Resty客户端
    Client := resty.New()
    //发送POST请求并处理响应
    result, err := Client.R().SetBody(map[string]string{
        "name":     avg.User,
        "password": avg.Pass,
    }).Post(avg.Uri + "/api/user/login")
    if err != nil {
        fmt.Println("err:", err.Error())
        return
    }
    if strings.Contains(string(result.Body()), "Failed") {
        fmt.Println("login is failed.")
        return
    }
    /*-------------------------登录manager web end----------------------------*/
    fmt.Println()
    fmt.Println(fmt.Sprintf("%-18s%-14s%-48s%-40s", "提交时间", "提交用户", "查询ID", "报错信息"))
    ch := make(chan struct{}, 3)
    var wg sync.WaitGroup
    for _, m2 := range m {
        wg.Add(1)
        go func(m2 map[string]interface{}) {
            defer func() {
                <-ch
                wg.Done()
            }()

            ch <- struct{}{}
            /*新的合并模块---------重新利用Cookies请求license*/
            respones3, err := Client.R().Get(avg.Uri + "/api/query/detail/" + m2["queryid"].(string))
            if err != nil {
                fmt.Println("请求失败:", err.Error())
                return
            }
            type mm struct {
                Code int `json:"code"`
                Data struct {
                    QueryDetail struct {
                        QueryID      string   `json:"queryId"`
                        User         string   `json:"user"`
                        Status       string   `json:"status"`
                        ErrorMessage string   `json:"errorMessage"`
                        StartTime    int      `json:"startTime"`
                        EndTime      int      `json:"endTime"`
                        TimeUsed     int      `json:"timeUsed"`
                        SQL          string   `json:"sql"`
                        ExplainText  []string `json:"explainText"`
                        ProfileText  []string `json:"profileText"`
                    } `json:"queryDetail"`
                } `json:"data"`
            }

            var m mm
            err = json.Unmarshal(respones3.Body(), &m)
            if err != nil {
                fmt.Println(err.Error())
                return
            }
            if m.Data.QueryDetail.ErrorMessage == "" {
                return
            }
            c := color.New()
            fmt.Println(fmt.Sprintf("%-22s%-18s%-50s%-40s", m2["starttime"].(time.Time).Format("2006-01-02 15:04:05"), strings.ReplaceAll(m2["user"].(string), "default_cluster:", ""), m2["queryid"].(string), c.Add(color.FgHiRed).Sprint(m.Data.QueryDetail.ErrorMessage)))

        }(m2)
    }
    wg.Wait()
    fmt.Println()
}

// 这里也是最原始的方式，但这里的密码是加了base64的
func scan2(avg util.ConnectParms, m []map[string]interface{}) {
    /*-------------------------登录manager web-------------------------------*/
    //创建Resty客户端
    Client := resty.New()
    //发送POST请求并处理响应
    result, err := Client.R().SetBody(map[string]string{
        "name":     avg.User,
        "password": base64.StdEncoding.EncodeToString([]byte(avg.Pass)),
    }).Post(avg.Uri + "/api/user/login")
    if err != nil {
        fmt.Println("err:", err.Error())
        return
    }
    if strings.Contains(string(result.Body()), "Failed") {
        fmt.Println("login is failed.")
        return
    }
    /*-------------------------登录manager web end----------------------------*/
    fmt.Println()
    fmt.Println(fmt.Sprintf("%-18s%-14s%-48s%-40s", "提交时间", "提交用户", "查询ID", "报错信息"))
    ch := make(chan struct{}, 3)
    var wg sync.WaitGroup
    for _, m2 := range m {
        wg.Add(1)
        go func(m2 map[string]interface{}) {
            defer func() {
                <-ch
                wg.Done()
            }()

            ch <- struct{}{}
            /*新的合并模块---------重新利用Cookies请求license*/
            respones3, err := Client.R().Get(avg.Uri + "/api/query/detail/" + m2["queryid"].(string))
            if err != nil {
                fmt.Println("请求失败:", err.Error())
                return
            }
            type mm struct {
                Code int `json:"code"`
                Data struct {
                    QueryDetail struct {
                        QueryID      string   `json:"queryId"`
                        User         string   `json:"user"`
                        Status       string   `json:"status"`
                        ErrorMessage string   `json:"errorMessage"`
                        StartTime    int      `json:"startTime"`
                        EndTime      int      `json:"endTime"`
                        TimeUsed     int      `json:"timeUsed"`
                        SQL          string   `json:"sql"`
                        ExplainText  []string `json:"explainText"`
                        ProfileText  []string `json:"profileText"`
                    } `json:"queryDetail"`
                } `json:"data"`
            }

            var m mm
            err = json.Unmarshal(respones3.Body(), &m)
            if err != nil {
                fmt.Println(err.Error())
                return
            }
            if m.Data.QueryDetail.ErrorMessage == "" {
                return
            }
            c := color.New()
            fmt.Println(fmt.Sprintf("%-22s%-18s%-50s%-40s", m2["starttime"].(time.Time).Format("2006-01-02 15:04:05"), strings.ReplaceAll(m2["user"].(string), "default_cluster:", ""), m2["queryid"].(string), c.Add(color.FgHiRed).Sprint(m.Data.QueryDetail.ErrorMessage)))

        }(m2)
    }
    wg.Wait()
    fmt.Println()
}

func scan1(avg util.ConnectParms, m []map[string]interface{}) {
    if len(m) == 0 {
        return
    }
    client := resty.New()
    fmt.Println(fmt.Sprintf("%-18s%-14s%-48s%-40s", "提交时间", "提交用户", "查询ID", "报错信息"))
    ch := make(chan struct{}, 3)
    var wg sync.WaitGroup
    for _, m2 := range m {
        wg.Add(1)
        go func(m2 map[string]interface{}) {
            defer func() {
                <-ch
                wg.Done()
            }()

            ch <- struct{}{}

            queryid := m2["queryid"].(string)
            getOpenApi(client, avg, queryid)
        }(m2)
    }
    wg.Wait()
    fmt.Println()
}

func getOpenApi(client *resty.Client, avg util.ConnectParms, queryid string) {
    nonce := tools.RandomPassWord(32)
    unix := time.Now().Unix()

    AuthCredential := fmt.Sprintf("%s/%d/%s", avg.AaccessKey, unix, nonce)
    AuthCredentialEncrypted := util.HexEncrypt(AuthCredential, avg.SecretKey)
    AuthContent := fmt.Sprintf(`HTTPMethod:GET
CanonicalURI:/openapi/v1/query/record/%s
CanonicalQueryString:
CanonicalForm:`, queryid)
    AuthSignature := util.HexEncrypt(AuthContent, AuthCredentialEncrypted)
    AuthHeader := fmt.Sprintf("MANAGER-HMAC-SHA256 Credential=%s,Signature=%s", AuthCredential, AuthSignature)

    uri := fmt.Sprintf("%s/openapi/v1/query/record/%s", avg.Uri, queryid)
    //发送POST请求并处理响应
    respones, err := client.
        R().
        SetHeaders(map[string]string{
            "Content-Type":  "application/x-www-form-urlencoded",
            "Authorization": AuthHeader,
        }).
        Get(strings.NewReplacer(" ", "").Replace(uri))
    if err != nil {
        fmt.Println("报错：", err.Error())
        return
    }
    if strings.Contains(string(respones.Body()), "404 page not found") {
        return
    }
    var item util.HexData
    err = json.Unmarshal(respones.Body(), &item)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    c := color.New()
    if item.Data.StartTime > 0 {
        fmt.Println(fmt.Sprintf("%-22s%-18s%-50s%-40s", time.Unix(int64(item.Data.StartTime), 0).Format("2006-01-02 15:04:05"), item.Data.User, item.Data.QueryID, c.Add(color.FgHiRed).Sprint(item.Data.FailedReason)))
    }
}
