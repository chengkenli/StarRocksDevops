package backend

import (
	"fmt"
	"github.com/fatih/color"
	"starrocks/util"
	"strings"
	"sync"
)

func BackendsID() {
	var (
		d util.Dbs
	)
	util.Connect.Raw("show databases").Scan(&d)

	fmt.Println(fmt.Sprintf("\nstarrocks(%s) 😀\n\nbrokerload:(-- %d --)", util.P.App, util.P.JobID))

	var wg sync.WaitGroup
	for _, base := range d {
		wg.Add(1)
		base := base

		go func() {
			defer wg.Done()

			if base.Database == "_statistics_" {
				return
			}
			var l util.BrokerLoad
			util.Connect.Raw(fmt.Sprintf("show load from %s", base.Database)).Scan(&l)
			for _, bl := range l {
				c := color.New()

				var status string
				switch strings.ReplaceAll(bl.State, " ", "") {
				case "FINISHED":
					status = c.Add(color.FgHiGreen).Sprint(bl.State)
				case "PENDING", "QUEUEING", "PREPARED":
					status = c.Add(color.FgHiYellow).Sprint(bl.State)
				case "LOADING":
					status = c.Add(color.FgHiMagenta).Sprint(bl.State)
				case "CANCELLED":
					status = c.Add(color.FgHiRed).Sprint(bl.State)
				default:
					status = bl.State
				}

				if bl.JobId != util.P.JobID {
					continue
				}

				broker := fmt.Sprintf("%-20s: %-20s\n%-20s: %-20d\n%-20s: %-20s\n%-20s: %-20s\n%-20s: %-20s\n%-20s: %-20s\n%-20s: %-20s\n%-20s: %-20s\n%-20s: %-20s\n%-20s: %-20s\n%-20s: %-20s\n%-20s: %-20s\n%-20s: %-20s\n%-20s: %-20s\n%-20s: %-20s\n%-20s: %-20s",
					"DataBase", base.Database,
					"JobId", bl.JobId,
					"Label", bl.Label,
					"State", status,
					"Progress", bl.Progress,
					"Type", bl.Type,
					"EtlInfo", bl.EtlInfo,
					"TaskInfo", bl.TaskInfo,
					"ErrorMsg", bl.ErrorMsg,
					"CreateTime", bl.CreateTime,
					"EtlStartTime", bl.EtlStartTime,
					"EtlFinishTime", bl.EtlFinishTime,
					"LoadStartTime", bl.LoadStartTime,
					"LoadFinishTime", bl.LoadFinishTime,
					"URL", bl.URL,
					"JobDetails", bl.JobDetails,
				)
				fmt.Println(broker)
				if util.P.Kill {
					sql := fmt.Sprintf("cancel load from %s where label='%s'", base.Database, bl.Label)
					fmt.Println(c.Add(color.FgHiYellow).Sprint(sql))
					r := util.Connect.Raw(fmt.Sprintf(sql))
					if r.Error != nil {
						fmt.Println(c.Add(color.FgHiRed).Sprint(r.Error.Error()))
						return
					}
					fmt.Println(c.Add(color.FgHiGreen).Sprint("查杀成功!"))
					return
				}
			}
		}()
	}
	wg.Wait()
	fmt.Println()
}
