package frontend

import (
    "fmt"
    "github.com/fatih/color"
    "StarRocksDevops/tools"
    "StarRocksDevops/util"
    "strconv"
    "strings"
    "sync"
)

func SessionFronBacks() {
    var (
        f util.Fronends
        b util.Backends
    )
    util.Connect.Raw("show frontends").Scan(&f)
    util.Connect.Raw("show backends").Scan(&b)

    c := color.New()
    var scheduledTabletNum, tabletCount, tabletBytes int
    /*打印fe信息*/
    fmt.Println(fmt.Sprintf("\nstarrocks(%s) 😀\n\n%s\n%-33s %-13s %-5s %-5s %-5s %-5s %-9s %-6s %-11s %-5s %-5s %-10s %-30s %-20s %-20s %-20s %-11s", util.P.App, c.Add(color.FgHiGreen).Sprint("frontends:"),
        "name", "ip", "eport", "hport", "qport", "rport", "role", "master", "clusterid", "join", "alive", "rjournalid", "ishelper", "lastheartbeat", "starttime", "version", "errmsg"))
    var wg sync.WaitGroup
    for _, item := range f {
        wg.Add(1)
        item := item
        go func() {
            defer wg.Done()

            if item.Role == "LEADER" || item.IsMaster == "true" {
                r := tools.Post("GET", fmt.Sprintf("http://%s:8030/metrics", item.IP), nil)
                for _, s := range strings.Split(string(r), "\n") {
                    if s == "" {
                        continue
                    }
                    if strings.Contains(s, "starrocks_fe_scheduled_tablet_num") {
                        scheduledTabletNum, _ = strconv.Atoi(strings.Split(s, " ")[1])
                    }
                    if strings.Contains(s, `starrocks_fe_memory{type="tablet_count"}`) {
                        tabletCount, _ = strconv.Atoi(strings.Split(s, " ")[1])
                    }
                    if strings.Contains(s, `starrocks_fe_memory{type="tablet_bytes"}`) {
                        tabletBytes, _ = strconv.Atoi(strings.Split(s, " ")[1])
                    }
                }
            }
            front := fmt.Sprintf("%-33s %-13s %-5d %-5d %-5d %-5d %-9s %-6s %-11d %-5s %-5s %-10d %-30s %-20s %-20s %-11s %-11s",
                item.Name,
                item.IP,
                item.EditLogPort,
                item.HttpPort,
                item.QueryPort,
                item.RpcPort,
                item.Role,
                item.IsMaster,
                item.ClusterId,
                item.Join,
                item.Alive,
                item.ReplayedJournalId,
                item.IsHelper,
                item.LastHeartbeat,
                item.StartTime,
                item.Version,
                item.ErrMsg,
            )
            fmt.Println(front)
        }()
    }
    wg.Wait()

    /*打印be信息*/
    fmt.Println(fmt.Sprintf("\n%s\n%-10s %-10s %-10s %-13s %-5s %-5s %-5s %-5s %-5s %-5s %-7s %-10s %-10s %-10s %-10s %-8s %-8s %-20s %-20s %-20s %-11s", c.Add(color.FgHiGreen).Sprint("backends:"),
        "id", "datatotal", "datausedpct", "ip", "heart", "be", "http", "brpc", "alive", "system", "cluster", "tabletnum", "dataused", "avail", "total", "usedpct", "maxdisk", "lastheartbeat", "starttime", "version", "errmsg"),
    )

    var tabletNum int
    var wg1 sync.WaitGroup
    for _, item := range b {
        wg1.Add(1)
        item := item
        go func() {
            defer wg1.Done()
            tabletNum = tabletNum + item.TabletNum
            back := fmt.Sprintf("%-10s %-10s %-11s %-13s %-5d %-5d %-5d %-5d %-5s %-6s %-7s %-10d %-10s %-10s %-10s %-8s %-8s %-20s %-20s %-20s %s",
                item.BackendId,
                item.DataTotalCapacity,
                item.DataUsedPct,
                item.IP,
                item.HeartbeatPort,
                item.BePort,
                item.HttpPort,
                item.BrpcPort,
                item.Alive,
                item.SystemDecommissioned,
                item.ClusterDecommissioned,
                item.TabletNum,
                item.DataUsedCapacity,
                item.AvailCapacity,
                item.TotalCapacity,
                item.UsedPct,
                item.MaxDiskUsedPct,
                item.LastHeartbeat,
                item.LastStartTime,
                item.Version,
                item.ErrMsg,
            )
            fmt.Println(back)
        }()
    }
    wg1.Wait()
    fmt.Println(fmt.Sprintf("tabletNum:(%s)、tabletCount:(%s)、tabletBytes:(%s)、unhealthyTabletNum:(%s)",
        c.Add(color.FgHiGreen).Sprint(tabletNum),
        c.Add(color.FgHiGreen).Sprint(tabletCount),
        c.Add(color.FgHiGreen).Sprint(tabletBytes),
        c.Add(color.FgHiRed).Sprint(scheduledTabletNum)))
    fmt.Println()
}
