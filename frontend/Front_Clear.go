package frontend

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/patrickmn/go-cache"
	"os"
	"starrocks/conn"
	"starrocks/util"
	"strings"
	sync2 "sync"
	"time"
)

func SessionClear() {
	i := 1
	ticker := time.NewTicker(time.Second * 1)
	cas := cache.New(2*time.Hour, 4*time.Hour)
	var cacheLise []string
	var wait sync2.WaitGroup
	wait.Add(1)
	go func() {
		for {
			select {
			case <-ticker.C:
				for _, h := range util.FrontendsIP {
					_, b := cas.Get(h)
					if b {
						checkgos(cacheLise)
						continue
					}
					db, _ := conn.StarRocksSingle(util.P.App, h)
					c := color.New()
					fmt.Println(c.Add(color.FgHiCyan).Sprint(fmt.Sprintf("第%d次清扫~", i)))

					go func(h string) {
						i++
						var p util.Process
						r := db.Raw("show processlist").Scan(&p)
						if r.Error != nil {
							fmt.Println(r.Error.Error())
							return
						}
						if len(p) == 0 {
							fmt.Println(c.Add(color.FgHiGreen).Sprint(fmt.Sprintf("%s已经没有发现存活踪迹~", h)))
							cas.Set(h, h, cache.DefaultExpiration)
							cacheLise = append(cacheLise, h)
							checkgos(cacheLise)
							return
						}

						var x int
						done := make(chan struct{})
						for _, item := range p {
							item := item
							go func() {
								defer func() {
									x++
									if x == len(p) {
										done <- struct{}{}
									}
								}()

								c := color.New()
								r := db.Exec(fmt.Sprintf("kill %d", item.Id))
								if r.Error != nil {
									fmt.Println(r.Error.Error())
									return
								}
								c = color.New()
								if item.Command == "Sleep" {
									fmt.Println(c.Add(color.FgHiBlue).Sprint(fmt.Sprintf("清场矩阵%s -> KILL %d(%s)[%s] - %d is done!", h, item.Id, item.User, item.Command, item.Time)))
								}
								c = color.New()
								if item.Command == "Query" {
									fmt.Println(c.Add(color.FgHiRed).Sprint(fmt.Sprintf("清场矩阵%s -> KILL %d(%s)[%s] - %d is done!", h, item.Id, item.User, item.Command, item.Time)))
								}
							}()
						}
						<-done
					}(h)
				}
			}
		}
	}()
	wait.Wait()
}

func checkgos(cacheLise []string) {
	if len(cacheLise) == len(util.FrontendsIP) {
		c := color.New()
		fmt.Println(c.Add(color.FgHiGreen).Sprint(fmt.Sprintf("%s 都已经没有发现存活踪迹~", strings.Join(cacheLise, ","))))
		os.Exit(-1)
	}
}
