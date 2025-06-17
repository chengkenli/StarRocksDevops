/*
 *@author  chengkenli
 *@project starrocks
 *@package frontend
 *@file    Front_Profile_List
 *@date    2025/1/20 10:29
 */

package frontend

import (
    "fmt"
    "github.com/fatih/color"
    "gorm.io/gorm"
    "StarRocksDevops/conn"
    "StarRocksDevops/tools"
    "StarRocksDevops/util"
    "strings"
)

func ProfileList() {
    if tools.Version() < 3.1 {
        return
    }
    var f []map[string]interface{}
    r := util.Connect.Raw("show frontends").Scan(&f)
    if r.Error != nil {
        fmt.Println(r.Error.Error())
        return
    }
    var leader string
    for _, m := range f {
        if m["Role"].(string) == "LEADER" {
            leader = m["IP"].(string)
        }
    }
    if len(leader) == 0 {
        fmt.Println("leader is nil.")
    }
    db, err := conn.StarRocksSingle(util.P.App, leader)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    var m []map[string]interface{}
    r = db.Raw("SHOW PROFILELIST").Scan(&m)
    if r.Error != nil {
        fmt.Println(r.Error.Error())
        return
    }
    if len(util.P.QueryId) == 0 {
        for i, m2 := range m {
            fmt.Println(fmt.Sprintf("#%03d %-15v %-10v %-10v %-40v %-10v", i, m2["StartTime"], m2["Time"], m2["State"], m2["QueryId"], m2["Statement"]))
        }
    } else {
        fmt.Println()
        ExplainStr := ANALYZE_PROFILE(db, util.P.QueryId, util.P.Pid)
        if len(ExplainStr) < 1 {
            ExplainStr = GET_QUERY_PROFILE(db, util.P.QueryId)
        }
        fmt.Println(ExplainStr)
        fmt.Println()
    }
}

// GET_QUERY_PROFILE
// SELECT GET_QUERY_PROFILE
func GET_QUERY_PROFILE(db *gorm.DB, queryid string) string {
    c := color.New()
    var mm map[string]interface{}

    fmt.Println(fmt.Sprintf("SELECT GET_QUERY_PROFILE:%s", c.Add(color.FgHiWhite).Sprint(queryid)))
    r := db.Raw(fmt.Sprintf("SELECT GET_QUERY_PROFILE('%s') AS QUERY_PROFILE", queryid)).Scan(&mm)
    if r.Error != nil {
        fmt.Println(r.Error.Error())
        return ""
    }
    return mm["QUERY_PROFILE"].(string)
}

// ANALYZE_PROFILE
// ANALYZE PROFILE FROM xxxx
func ANALYZE_PROFILE(db *gorm.DB, queryid string, plan_node_id int) string {
    c := color.New()
    var mm []map[string]interface{}

    var sql string
    if plan_node_id > -1 {
        fmt.Println(fmt.Sprintf("ANALYZE PROFILE:%s PLAN_NODE_ID:%s", c.Add(color.FgHiWhite).Sprint(queryid), c.Add(color.FgHiYellow).Sprint(plan_node_id)))
        sql = fmt.Sprintf("ANALYZE PROFILE FROM '%s',%d", queryid, plan_node_id)
    } else {
        fmt.Println(fmt.Sprintf("ANALYZE PROFILE:%s", c.Add(color.FgHiWhite).Sprint(queryid)))
        sql = fmt.Sprintf("ANALYZE PROFILE FROM '%s'", queryid)
    }
    for i := 0; i < 3; i++ {
        r := db.Raw(sql).Scan(&mm)
        if r.Error != nil {
            fmt.Println(i, r.Error.Error())
            continue
        }
        break
    }

    var explainStr []string
    for _, m3 := range mm {
        explainStr = append(explainStr, m3["Explain String"].(string))
    }
    return strings.Join(explainStr, "\n")
}
