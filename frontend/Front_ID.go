package frontend

import (
	"fmt"
	"github.com/fatih/color"
	"starrocks/conn"
	"starrocks/tools"
	"starrocks/util"
	"time"
)

func SessionID() {
	var p util.Process
	fmt.Println(fmt.Sprintf("\nstarrocks(%s) 😀\n\nprocesspid:\n%-15s %-25s %-20s %-20s %-20s %-8s %-8s %-8s\n", util.P.App, "Id", "User", "Host", "Cluster", "Db", "Command", "Time", "State"))
	filename := fmt.Sprintf("%d.%d.log", util.P.Pid, time.Now().Unix())

	for _, h := range util.FrontendsIP {
		db, err := conn.StarRocksSingle(util.P.App, h)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		db.Raw("show full processlist").Scan(&p)
		for _, item := range p {
			if item.Id == util.P.Pid {
				c := color.New()
				front := fmt.Sprintf("%-15d %-25s %-20s %-20s %-20s %-8s %-8d %-8s \n\n%-8s",
					item.Id,
					item.User,
					item.Host,
					item.Cluster,
					item.Db,
					item.Command,
					item.Time,
					item.State,
					item.Info,
				)
				tools.Writefile(filename, front)
				fmt.Println(c.Add(color.FgHiYellow).Sprint(front))
			}
		}
	}
	fmt.Println()
}
