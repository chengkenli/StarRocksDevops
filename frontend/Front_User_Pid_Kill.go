package frontend

import (
    "fmt"
    "github.com/fatih/color"
    "StarRocksDevops/conn"
    "StarRocksDevops/util"
    sync2 "sync"
)

func SessionUserPidKill() {
    var wait sync2.WaitGroup
    var p util.Process
    for _, h := range util.FrontendsIP {
        wait.Add(1)
        go func(h string) {
            defer wait.Done()
            db, _ := conn.StarRocksSingle(util.P.App, h)
            r := db.Raw("show processlist").Scan(&p)
            if r.Error != nil {
                fmt.Println(r.Error.Error())
                return
            }
            var wg sync2.WaitGroup
            for _, item := range p {
                wg.Add(1)
                item := item
                go func() {
                    defer wg.Done()
                    if len(util.P.Username) == 0 {
                        if item.Id == util.P.Pid {
                            r := db.Exec(fmt.Sprintf("kill %d", util.P.Pid))
                            if r.Error != nil {
                                fmt.Println(r.Error.Error())
                                return
                            }
                            c := color.New()
                            if item.Command == "Sleep" {
                                fmt.Println(c.Add(color.FgHiYellow).Sprint(fmt.Sprintf("%s -> KILL %d(%s) - %d is done!", h, util.P.Pid, item.User, item.Time)))
                            }
                            c = color.New()
                            if item.Command == "Query" {
                                fmt.Println(c.Add(color.FgHiGreen).Sprint(fmt.Sprintf("%s -> KILL %d(%s) - %d is done!", h, util.P.Pid, item.User, item.Time)))
                            }
                        }
                    }
                }()
            }
            wg.Wait()
        }(h)
    }
    wait.Wait()
}
