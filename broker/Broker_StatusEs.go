/*
 *@author  chengkenli
 *@project starrocks
 *@package broker
 *@file    BrokerStatusErrors
 *@date    2024/6/20 9:27
 */

package broker

import (
	"fmt"
	"github.com/fatih/color"
	"starrocks/util"
	"strings"
	"sync"
	"time"
)

// BrokerStatusEs 根据状态列出broker所有作业信息
func BrokerStatusEs() {
	var databases []map[string]interface{}
	util.Connect.Raw("show databases").Scan(&databases)

	fmt.Println()
	c := color.New()

	fmt.Println(c.Add(color.FgHiYellow).Sprint("Broker Load:"))
	fmt.Println(fmt.Sprintf("%-9s %-18s %-50s %-15s %-40s %-10s %-20s %-20s %-20s", "JobId", "Database", "Label", "State", "ErrorMsg", "Type", "CreateTime", "LoadFinishTime", "Edtime"))

	var wgs sync.WaitGroup
	for _, database := range databases {
		wgs.Add(1)
		go func(database map[string]interface{}) {
			defer wgs.Done()

			if database["Database"] == "_statistics_" {
				return
			}
			var brokers util.BrokerLoad
			r := util.Connect.Raw(fmt.Sprintf("SHOW LOAD FROM %s WHERE STATE = '%s'", database["Database"], util.P.Status)).Scan(&brokers)
			if r.Error != nil {
				return
			}
			if brokers == nil {
				return
			}
			var wg sync.WaitGroup
			for _, b := range brokers {
				wg.Add(1)
				b := b
				go func() {
					defer wg.Done()

					var errorMsg string
					if util.P.Status == "CANCELLED" {
						errorMsg = strings.NewReplacer("type:ETL_RUN_FAIL; msg:", "", "type:LOAD_RUN_FAIL; msg:", "", "type:ETL_QUALITY_UNSATISFIED; msg:", "", "Cancelled, msg: ", "").Replace(b.ErrorMsg)
						if len(errorMsg) > 32 {
							errorMsg = errorMsg[:30] + "..."
						}
					} else {
						errorMsg = ""
					}

					pool.Put(&job{
						JobId:          b.JobId,
						Label:          b.Label,
						State:          b.State,
						Progress:       b.Progress,
						Type:           b.Type,
						EtlInfo:        b.EtlInfo,
						TaskInfo:       b.TaskInfo,
						ErrorMsg:       errorMsg,
						CreateTime:     b.CreateTime,
						EtlStartTime:   b.EtlStartTime,
						EtlFinishTime:  b.EtlFinishTime,
						LoadStartTime:  b.LoadStartTime,
						LoadFinishTime: b.LoadFinishTime,
						URL:            b.URL,
						JobDetails:     b.JobDetails,
						Database:       database["Database"].(string),
					})

					p := pool.Get().(*job)

					t1, _ := time.Parse("2006-01-02 15:04:05", p.CreateTime)
					t2, _ := time.Parse("2006-01-02 15:04:05", p.LoadFinishTime)
					edtime := t2.Sub(t1).String()
					c := color.New()
					broker := fmt.Sprintf("%-13s %-18s %-50s %-15s %-40s %-10s %-20s %-20s %-20s",
						c.Add(color.FgHiGreen).Sprint(p.JobId),
						p.Database,
						p.Label,
						p.State,
						p.ErrorMsg,
						p.Type,
						p.CreateTime,
						p.LoadFinishTime,
						edtime,
					)
					fmt.Println(broker)
				}()
			}
			wg.Wait()
		}(database)
	}
	wgs.Wait()
}
