package frontend

import (
    "fmt"
    "github.com/fatih/color"
    "StarRocksDevops/conn"
    "StarRocksDevops/tools"
    "StarRocksDevops/util"
    "strings"
    "sync"
)

var pool *sync.Pool

type Item struct {
    Id      int
    User    string
    Host    string
    Cluster string
    Db      string
    Command string
    Time    int
    State   string
    Info    string
}
type Item2 struct {
    Id        int
    User      string
    Host      string
    Cluster   string
    Db        string
    Command   string
    Time      int
    State     string
    Info      string
    IsPending string
    Warehouse string
}

func init() {
    pool = &sync.Pool{
        New: func() interface{} {
            fmt.Println("Creating a new Pool ...")
            return new(Item2)
        },
    }
}

func SessionList() {
    var p util.Process2

    fmt.Println()
    c := color.New()
    fmt.Println(c.Add(color.FgHiYellow).Sprint("运行语句:"))
    fmt.Println(util.P.Primary)

    fmt.Println(fmt.Sprintf("%-10s %-15s %-25s %-20s %-20s %-20s %-8s %-18s %-12s %-11s %-8s", "Num", "Id", "User", "Host", "Cluster", "Db", "Command", "Time", "State", "IsPending", "Warehouse"))
    var nodecnts []string
    var sum int
    for _, h := range util.FrontendsIP {
        db, err := conn.StarRocksSingle(util.P.App, h)
        if err != nil {
            fmt.Println(err.Error())
            continue
        }
        db.Raw("show full processlist").Scan(&p)
        c := color.New(color.Bold)
        fmt.Println(c.Add(color.FgHiGreen).Sprint(h) + " 👇")

        var querycnt []string
        var wg sync.WaitGroup
        for i, item := range p {
            wg.Add(1)
            item := item
            go func(i int) {
                defer wg.Done()

                pool.Put(&Item2{
                    Id:        item.Id,
                    User:      item.User,
                    Host:      item.Host,
                    Cluster:   item.Cluster,
                    Db:        item.Db,
                    Command:   item.Command,
                    Time:      item.Time,
                    State:     item.State,
                    Info:      item.Info,
                    IsPending: item.IsPending,
                    Warehouse: item.Warehouse,
                })
                p := pool.Get().(*Item2)

                if p.Command == "Query" {
                    querycnt = append(querycnt, p.User)
                    /*超级管理员系列*/
                    var once sync.Once
                    var mm []map[string]interface{}
                    once.Do(
                        func() {
                            util.Connect.Raw("show grants for " + p.User).Scan(&mm)
                        })
                    if p.Time >= 1 && !strings.Contains(strings.ToLower(p.User), "starrocks") {
                        if strings.Contains(strings.ToLower(fmt.Sprintf("%v", mm)), "admin") || strings.Contains(strings.ToLower(fmt.Sprintf("%v", mm)), "root") {
                            c := color.New()
                            front := fmt.Sprintf("%-10d %-15d %-25s %-20s %-20s %-20s %-8s %-4d %-13s %-8s %-10s %-13s %-10s", i,
                                p.Id,
                                p.User,
                                p.Host,
                                p.Cluster,
                                p.Db,
                                p.Command,
                                p.Time,
                                tools.GetHour(p.Time),
                                p.State,
                                p.IsPending,
                                p.Warehouse,
                                "超级管理员",
                            )
                            if len(util.P.Primary) != 0 {
                                if strings.Contains(p.Info, util.P.Primary) {
                                    fmt.Println(c.Add(color.FgHiBlue).Sprint(front))
                                }
                                return
                            }
                            fmt.Println(c.Add(color.FgHiBlue).Sprint(front))
                            return
                        }
                    }
                    c := color.New()
                    if p.Time >= 1000 {
                        front := fmt.Sprintf("%-10d %-15d %-25s %-20s %-20s %-20s %-8s %-4d %-13s %-8s %-10s %-13s %-10s", i,
                            p.Id,
                            p.User,
                            p.Host,
                            p.Cluster,
                            p.Db,
                            p.Command,
                            p.Time,
                            tools.GetHour(p.Time),
                            p.State,
                            p.IsPending,
                            p.Warehouse,
                            "高耗时慢查询",
                        )
                        if len(util.P.Primary) != 0 {
                            if strings.Contains(p.Info, util.P.Primary) {
                                fmt.Println(c.Add(color.FgHiRed).Sprint(front))
                            }
                            return
                        }
                        fmt.Println(c.Add(color.FgHiRed).Sprint(front))
                        return
                    }

                    if p.Time >= 300 {
                        front := fmt.Sprintf("%-10d %-15d %-25s %-20s %-20s %-20s %-8s %-4d %-13s %-8s %-10s %-13s %-10s", i,
                            p.Id,
                            p.User,
                            p.Host,
                            p.Cluster,
                            p.Db,
                            p.Command,
                            p.Time,
                            tools.GetHour(p.Time),
                            p.State,
                            p.IsPending,
                            p.Warehouse,
                            "慢查询",
                        )
                        if len(util.P.Primary) != 0 {
                            if strings.Contains(p.Info, util.P.Primary) {
                                fmt.Println(c.Add(color.FgHiYellow).Sprint(front))
                            }
                            return
                        }
                        fmt.Println(c.Add(color.FgHiYellow).Sprint(front))
                    } else {
                        front := fmt.Sprintf("%-10d %-15d %-25s %-20s %-20s %-20s %-8s %-4d %-13s %-8s %-10s %-13s", i,
                            p.Id,
                            p.User,
                            p.Host,
                            p.Cluster,
                            p.Db,
                            p.Command,
                            p.Time,
                            tools.GetHour(p.Time),
                            p.State,
                            p.IsPending,
                            p.Warehouse,
                        )
                        if len(util.P.Primary) != 0 {
                            if strings.Contains(p.Info, util.P.Primary) {
                                fmt.Println(front)
                            }
                            return
                        }
                        fmt.Println(front)
                    }
                }
            }(i)
        }
        wg.Wait()
        sum = sum + len(querycnt)
        nodecnts = append(nodecnts, fmt.Sprintf("%s:Query(%d)", h, len(querycnt)))
    }
    fmt.Println(c.Add(color.FgHiYellow).Sprint(fmt.Sprintf("cnt:(%d)  %v", sum, strings.Join(nodecnts, "、"))))
}
