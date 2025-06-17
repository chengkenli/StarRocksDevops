package frontend

import (
    "fmt"
    "github.com/fatih/color"
    "sort"
    "StarRocksDevops/conn"
    "StarRocksDevops/util"
    "sync"
)

func User() {
    var (
        p util.Process
        u util.Property
    )

    fmt.Println(fmt.Sprintf("\nstarrocks(%s) 😀\n\nuser:\n%-15s %-25s %-20s %-20s %-20s %-8s %-8s %-8s", util.P.App, "Id", "User", "Host", "Cluster", "Db", "Command", "Time", "State"))

    /*当前用户连接数*/
    c := 0
    for _, h := range util.FrontendsIP {
        db, err := conn.StarRocksSingle(util.P.App, h)
        if err != nil {
            fmt.Println(err.Error())
            continue
        }
        db.Raw("show processlist").Scan(&p)
        fmt.Println(h + " 👇")
        for _, item := range p {
            if item.User == util.P.Username {
                front := fmt.Sprintf("%-15d %-25s %-20s %-20s %-20s %-8s %-8d %-8s",
                    item.Id,
                    item.User,
                    item.Host,
                    item.Cluster,
                    item.Db,
                    item.Command,
                    item.Time,
                    item.State,
                )
                c++
                fmt.Println(front)
            }
        }
    }

    cr := color.New(color.Bold)
    /*获取默认连接数*/
    util.Connect.Raw(fmt.Sprintf("SHOW PROPERTY FOR '%s'", util.P.Username)).Scan(&u)
    fmt.Println()
    for _, s := range u {
        if s.Key == "max_user_connections" {
            fmt.Println(fmt.Sprintf("user:[%s] current:[%s] default:[%s]", util.P.Username, cr.Add(color.FgHiYellow).Sprint(c), cr.Add(color.FgHiGreen).Sprint(s.Value)))
            break
        }
    }
    fmt.Println()
}

func fulls() {
    var (
        p     util.Process
        users []string
    )

    /*当前用户连接数*/
    var wg sync.WaitGroup
    for _, h := range util.FrontendsIP {
        wg.Add(1)
        go func(h string) {
            wg.Done()
            db, err := conn.StarRocksSingle(util.P.App, h)
            if err != nil {
                fmt.Println(err.Error())
                return
            }
            r := db.Raw("show processlist").Scan(&p)
            if r.Error != nil {
                fmt.Println(r.Error.Error())
                return
            }
            for _, item := range p {
                users = append(users, item.User)
            }
        }(h)
    }
    wg.Wait()
    in, str := arrsUser(users)
    fmt.Println(in)
    fmt.Println(str)
    fmt.Println()
}

func arrsUser(arr []string) ([]int, []string) {
    var vvs []string
    var vv []int
    type kv struct {
        key   string
        value int
    }
    tmplist := make([]kv, 0)
    ch := make(chan map[string]int, 0)
    go func() {
        m := make(map[string]int)
        for _, v := range arr {
            if m[v] == 0 {
                m[v] = 1
            } else {
                m[v]++
            }
        }
        ch <- m
    }()

    m := <-ch
    for k, v := range m {
        tmplist = append(tmplist, kv{key: k, value: v})
    }

    for _, k := range tmplist {
        vvs = append(vvs, fmt.Sprintf("(%s:%d)", k.key, k.value))
        vv = append(vv, k.value)
    }
    sort.Sort(sort.Reverse(sort.IntSlice(vv)))
    sort.Sort(sort.StringSlice(vvs))

    return vv, vvs
}
