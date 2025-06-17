/*
 *@author  chengkenli
 *@project starrocks
 *@package resource
 *@file    CurrentQueries
 *@date    2024/7/11 10:08
 */

package resource

import (
    "fmt"
    "github.com/fatih/color"
    "math/rand"
    "StarRocksDevops/conn"
    "StarRocksDevops/tools"
    "StarRocksDevops/util"
    "strings"
)

func CurrentQueriesCkId() {
    if tools.Version() < 3.0 {
        return
    }

    fmt.Println()
    for _, h := range util.FrontendsIP {
        db, err := conn.StarRocksSingle(util.P.App, h)
        if err != nil {
            fmt.Println(err.Error())
            return
        }
        var m []map[string]interface{}
        r := db.Raw("SHOW PROC '/current_queries'").Scan(&m)
        if r.Error != nil {
            fmt.Println(r.Error.Error())
        }

        for i, m2 := range m {
            if m2["QueryId"].(string) != util.P.QueryId {
                continue
            }
            fmt.Println(fmt.Sprintf("Id              :   %d", i))
            fmt.Println(fmt.Sprintf("StartTime       :   %s", m2["StartTime"].(string)))
            fmt.Println(fmt.Sprintf("QueryId         :   %s", m2["QueryId"].(string)))
            fmt.Println(fmt.Sprintf("ConnectionId    :   %s", m2["ConnectionId"].(string)))
            fmt.Println(fmt.Sprintf("Database        :   %s", m2["Database"].(string)))
            fmt.Println(fmt.Sprintf("User            :   %s", m2["User"].(string)))
            fmt.Println(fmt.Sprintf("ScanBytes       :   %s", m2["ScanBytes"].(string)))
            fmt.Println(fmt.Sprintf("ScanRows        :   %s", strings.NewReplacer("rows", "").Replace(m2["ScanRows"].(string))))
            fmt.Println(fmt.Sprintf("MemoryUsage     :   %s", m2["MemoryUsage"].(string)))
            fmt.Println(fmt.Sprintf("DiskSpillSize   :   %s", m2["DiskSpillSize"].(string)))
            fmt.Println(fmt.Sprintf("CPUTime         :   %s", m2["CPUTime"].(string)))
            fmt.Println(fmt.Sprintf("ExecTime        :   %s", m2["ExecTime"].(string)))

            var cm []map[string]interface{}
            r := db.Raw(fmt.Sprintf("SHOW PROC '/current_queries/%s/hosts'", m2["QueryId"].(string))).Scan(&cm)
            if r.Error != nil {
                fmt.Println(r.Error.Error())
                continue
            }

            fmt.Println()
            fmt.Println(fmt.Sprintf("  %-2s %-20s %-15s %-15s %-15s %-15s", "ID", "Host", "ScanBytes", "ScanRows", "MemUsageBytes", "CpuCostSeconds"))
            for i, m3 := range cm {
                msg := fmt.Sprintf("> %-2d %-20s %-15s %-15s %-15s %-15s ", i,
                    m3["Host"].(string),
                    m3["ScanBytes"].(string),
                    strings.NewReplacer(" rows", "").Replace(m3["ScanRows"].(string)),
                    m3["MemUsageBytes"].(string),
                    m3["CpuCostSeconds"].(string),
                )
                fmt.Println(msg)

            }

            var dm map[string]interface{}
            r = db.Raw(fmt.Sprintf("SHOW PROC '/current_queries/%s'", m2["QueryId"].(string))).Scan(&dm)
            if r.Error != nil {
                fmt.Println(r.Error.Error())
                continue
            }
            fmt.Println()
            filename := fmt.Sprintf("%s", m2["QueryId"].(string))
            if len(dm["Sql"].(string)) >= 300 {
                tools.Writefile(filename, dm["Sql"].(string))
                fmt.Println(dm["Sql"].(string)[0:299] + " ...")
                fmt.Println(filename)
            } else {
                fmt.Println(dm["Sql"].(string))
            }
        }
    }
    fmt.Println()
}

func Color(s interface{}) string {
    c := color.New()
    str := []string{
        c.Add(color.FgHiYellow).Sprint(s),
        c.Add(color.FgHiWhite).Sprint(s),
        c.Add(color.FgHiBlack).Sprint(s),
        c.Add(color.FgHiCyan).Sprint(s),
        c.Add(color.FgHiMagenta).Sprint(s),
        c.Add(color.FgHiBlue).Sprint(s),
        c.Add(color.FgHiRed).Sprint(s),
        c.Add(color.FgHiGreen).Sprint(s),
    }
    return str[rand.Intn(len(str))]
}
