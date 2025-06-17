/*
 *@author  chengkenli
 *@project starrocks
 *@package errlog
 *@file    ErrScan
 *@date    2024/10/16 9:38
 */

package errlog

import (
    "fmt"
    "github.com/fatih/color"
    "StarRocksDevops/tools"
    "StarRocksDevops/util"
    "strings"
    "time"
)

func ErrScan() {
    fmt.Println()
    c := color.New()
    fmt.Println("高消耗" + c.Add(color.FgHiBlue).Sprint("审计日志") + "ScanRows倒序TOP")

    var m []map[string]interface{}
    var where string
    if util.P.CommitType != "" {
        where = fmt.Sprintf("and lower(stmt) like '%%%s%%'", util.P.CommitType)
    }
    sql := fmt.Sprintf("select * from audit.starrocks_audit_log where timestamp>='%s' and timestamp<'%s' %s order by %s desc limit %d",
        util.P.StartTime,
        util.P.EndTime,
        where,
        util.P.ErrKey,
        util.P.CommitLimit,
    )
    fmt.Println(sql)

    r := util.Connect.Raw(sql).Scan(&m)
    if r.Error != nil {
        fmt.Println(r.Error.Error())
        return
    }

    fmt.Println(fmt.Sprintf("%s %-7s %-19s %-37s %-20s %-14s %-26s %-13s %-13s %-23s %-12s %-7s %-10v %-10s",
        "Id", "Type", "TimeStamp", "QueryId", "User", "ScanRows", "ScanBytes", "ReturnRows", "CpuCostNs", "MemCostBytes", "QueryTime", "State", "PTimeMs", "resourceGroup"))
    for i, m2 := range m {
        c := color.New()

        if strings.Contains(strings.ToLower(m2["stmt"].(string)), "insert") {
            util.P.CommitType = "insert"
        } else {
            util.P.CommitType = "select"
        }

        var commitType string
        switch util.P.CommitType {
        case "select", "query":
            commitType = c.Add(color.FgHiGreen).Sprint(util.P.CommitType)
        case "insert":
            commitType = c.Add(color.FgHiYellow).Sprint(util.P.CommitType)
        default:
            commitType = util.P.CommitType
        }
        //ScanBytes
        var ScanBytes string
        scanBytes := tools.ByteSizeToString(m2["scanBytes"].(int64))
        ScanBytes = scanBytes
        if strings.Contains(scanBytes, "TB") {
            ScanBytes = fmt.Sprintf("%d(%s)", m2["scanBytes"].(int64), c.Add(color.FgHiRed).Sprint(scanBytes))
        }
        if strings.Contains(scanBytes, "GB") {
            ScanBytes = fmt.Sprintf("%d(%s)", m2["scanBytes"].(int64), c.Add(color.FgHiYellow).Sprint(scanBytes))
        }
        if strings.Contains(scanBytes, "MB") {
            ScanBytes = fmt.Sprintf("%d(%s)", m2["scanBytes"].(int64), c.Add(color.FgHiGreen).Sprint(scanBytes))
        }
        if strings.Contains(scanBytes, "KB") {
            ScanBytes = fmt.Sprintf("%d(%s)", m2["scanBytes"].(int64), c.Add(color.FgHiBlue).Sprint(scanBytes))
        }
        //memCostBytes
        var MemCostBytes string
        memCostBytes := tools.ByteSizeToString(m2["memCostBytes"].(int64))
        MemCostBytes = memCostBytes
        if strings.Contains(memCostBytes, "TB") {
            MemCostBytes = fmt.Sprintf("%d(%s)", m2["memCostBytes"].(int64), c.Add(color.FgHiRed).Sprint(memCostBytes))
        }
        if strings.Contains(memCostBytes, "GB") {
            MemCostBytes = fmt.Sprintf("%d(%s)", m2["memCostBytes"].(int64), c.Add(color.FgHiYellow).Sprint(memCostBytes))
        }
        if strings.Contains(memCostBytes, "MB") {
            MemCostBytes = fmt.Sprintf("%d(%s)", m2["memCostBytes"].(int64), c.Add(color.FgHiGreen).Sprint(memCostBytes))
        }
        if strings.Contains(memCostBytes, "KB") {
            MemCostBytes = fmt.Sprintf("%d(%s)", m2["memCostBytes"].(int64), c.Add(color.FgHiBlue).Sprint(memCostBytes))
        }
        //state
        var State string
        state := m2["state"].(string)
        switch state {
        case "ERR":
            State = c.Add(color.FgHiRed).Sprint(state)
        case "EOF":
            State = c.Add(color.FgGreen).Sprint(state)
        case "OK":
            State = c.Add(color.FgHiGreen).Sprint(state)
        default:
            State = state
        }
        // default_wg
        var DefaultWg string
        rg, ok := m2["resourceGroup"].(string)
        if ok {
            if rg == "default_wg" {
                DefaultWg = c.Add(color.FgHiWhite).Sprint(rg)
            } else {
                DefaultWg = c.Add(color.FgHiMagenta).Sprint(rg)
            }
        } else {
            DefaultWg = rg
        }
        //index
        var pendingTimeMs int64
        v, ok := m2["pendingTimeMs"]
        if ok {
            pendingTimeMs = v.(int64)
        }
        msg := fmt.Sprintf("%02d %-15s %-20s %-37s %-20s %-14d %-38s %-13d %-13d %-38s %-12s %-25s %-10s %-10s",
            i,
            commitType,
            m2["timestamp"].(time.Time).Format("2006-01-02 15:04:05"),
            m2["queryId"].(string),
            m2["user"].(string),
            m2["scanRows"].(int64),
            ScanBytes,
            m2["returnRows"].(int64),
            m2["cpuCostNs"].(int64),
            MemCostBytes,
            tools.GetHour(int(m2["queryTime"].(int64))/1000),
            State,
            tools.GetHour(int(pendingTimeMs/1000)),
            DefaultWg,
        )
        fmt.Println(msg)
    }
    fmt.Println()
}
