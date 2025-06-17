/*
 *@author  chengkenli
 *@project starrocks
 *@package ioprofiler
 *@file    scan_ioprofile
 *@date    2024/12/30 14:22
 */

package ioprofiler

import (
    "fmt"
    "github.com/fatih/color"
    "github.com/go-resty/resty/v2"
    "gorm.io/gorm"
    "regexp"
    "StarRocksDevops/conn"
    "StarRocksDevops/tools"
    "StarRocksDevops/util"
    "strconv"
    "strings"
    "sync"
    "time"
)

func IoProfile() {

    c := color.New()

    fmt.Println()
    go func() {
        fmt.Println(fmt.Sprintf("请稍等 %s Seconds，正在采集中", c.Add(color.FgHiYellow).Sprint(util.P.Seconds)))
        for i := 1; i <= util.P.Seconds; i++ {
            tools.PrintProgress(i, util.P.Seconds)
            time.Sleep(time.Second * 1)
        }

        fmt.Println(fmt.Sprintf("\n\n%-12s %-12s %-25s %-10s %-15s %-10s %-12s %-10s %-60v %-10v",
            "Tablet",
            "TAG",
            "read_bytes",
            "ops",
            "write_bytes",
            "ops",
            "total_bytes",
            "ops",
            "Schema",
            "Partition",
        ))
    }()

    //创建Resty客户端
    Client := resty.New().SetDisableWarn(true)

    // 适配数组
    var IpList []string
    if len(util.P.IP) == 0 {
        fmt.Println("数组为空")
        IpList = util.BackendsIP
    } else {
        IpList = strings.Split(replace(util.P.IP), ",")
    }

    // 多例
    var wg sync.WaitGroup
    for i, ip := range IpList {
        wg.Add(1)
        go func(ip string, i int) {
            defer wg.Done()
            IPS(Client, util.Connect, ip)
            //pipeline_blocking_drivers(Client, ip)
        }(ip, i)
    }
    wg.Wait()
    fmt.Println()
}

func replace(s string) string {
    return strings.NewReplacer(" ", "", "\r", "", "\n", "").Replace(s)
}

func IPS(Client *resty.Client, db *gorm.DB, ip string) {
    if len(ip) == 0 {
        fmt.Println("ip is nil.")
        return
    }

    c := color.New()
    uri := fmt.Sprintf("http://%s:8040/ioprofile?seconds=%d", ip, util.P.Seconds)
    //发送POST请求并处理响应
    respones, err := Client.R().Get(uri)
    if err != nil {
        fmt.Println(c.Add(color.FgHiRed).Sprint(err.Error()))
        return
    }
    var data []string
    for _, body := range strings.Split(string(respones.Body()), "\n") {
        if strings.Contains(body, "Tablet") {
            continue
        }
        // 创建正则表达式，匹配一个或多个空格
        result := regexp.MustCompile(`\s+`).ReplaceAllString(body, " ")
        d := strings.Split(result, " ")
        if len(d) < 7 {
            continue
        }
        // tablet
        var m map[string]interface{}
        r := db.Raw("show tablet " + replace(d[1])).Scan(&m)
        if r.Error != nil {
            fmt.Println(c.Add(color.FgHiRed).Sprint(r.Error.Error()))
            continue
        }

        tbname := fmt.Sprintf("%s.%s", m["DbName"], m["TableName"])
        if len(tbname) >= 50 {
            tbname = tbname[0:49] + "..."
        }
        var ptname string
        if len(m["PartitionName"].(string)) > 10 {
            ptname = m["PartitionName"].(string)[0:9] + "..."
        } else {
            ptname = m["PartitionName"].(string)
        }

        atoi, _ := strconv.Atoi(d[3])
        msg := fmt.Sprintf("%-12s %-12s %-25s %-10s %-15s %-10s %-12s %-10s %-60v %-20v %-10v",
            replace(d[1]),
            replace(d[2]),
            fmt.Sprintf("%s(%s)", replace(d[3]), tools.ByteSizeToString(int64(atoi))),
            replace(d[4]),
            replace(d[5]),
            replace(d[6]),
            replace(d[7]),
            replace(d[8]),
            tbname,
            ptname,
            strings.Join(tools.RmDupSlice(cdsn(m["TableName"].(string))), ","),
        )

        size, _ := strconv.ParseFloat(strconv.FormatInt(int64(atoi), 10), 64)
        if size < KB {
            data = append(data, msg)
        } else if size < MB {
            data = append(data, c.Add(color.FgHiWhite).Sprint(msg))
        } else if size < GB {
            data = append(data, c.Add(color.FgHiBlue).Sprint(msg))
        } else if size < TB {
            if size >= 200*GB {
                data = append(data, c.Add(color.FgHiRed).Sprint(msg))
            } else {
                data = append(data, c.Add(color.FgHiYellow).Sprint(msg))
            }
        } else if size < PB {
            data = append(data, c.Add(color.FgHiRed).Sprint(msg))
        } else {
            data = append(data, c.Add(color.FgHiCyan).Sprint(msg))
        }
    }

    fmt.Println(c.Add(color.FgHiWhite).Sprint(ip))
    fmt.Println(strings.Join(data, "\n"))
    fmt.Println()
}

func cdsn(tablename string) []string {
    // 查询分发
    var frmsg []string
    for _, h := range util.FrontendsIP {
        db, err := conn.StarRocksSingle(util.P.App, h)
        if err != nil {
            fmt.Println(err.Error())
            return nil
        }

        var p util.Process
        db.Raw("show full processlist").Scan(&p)
        for _, item := range p {
            if item.Command == "Sleep" {
                continue
            }
            if !strings.Contains(item.Info, tablename) {
                continue
            }
            frmsg = append(frmsg, item.User)
        }
    }
    return frmsg
}
