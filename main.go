/**
 * @title: main
 * @Author ChengKen
 * @Date: 10/11/2022 13:30
 * @Version 1.0
 */

package main

import (
	"fmt"
	"starrocks/conn"
	_ "starrocks/init"
	"starrocks/run"
	"starrocks/util"
	"sync"
)

func main() {
	util.Parms()
	var err error
	util.Connect, err = conn.StarRocks(util.P.App)
	if err != nil {
		fmt.Println("出错了 ", err.Error())
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go initFrontendsIP(&wg)
	go initBackendsIP(&wg)
	wg.Wait()
	fmt.Println("init done.")

	run.Run()
}
