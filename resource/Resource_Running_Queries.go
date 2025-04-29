/*
 *@author  chengkenli
 *@project starrocks
 *@package resource
 *@file    RunningQueries
 *@date    2024/7/11 10:36
 */

package resource

import (
	"fmt"
	"github.com/fatih/color"
	"starrocks/tools"
	"starrocks/util"
)

func RunningQueries() {
	if tools.Version() < 3.0 {
		return
	}
	fmt.Println()
	var m []map[string]interface{}
	r := util.Connect.Raw("SHOW RUNNING QUERIES").Scan(&m)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
	}

	c := color.New()
	fmt.Println(c.Add(color.FgHiYellow).Sprint("运行队列:"))
	fmt.Println(fmt.Sprintf("%-3s %-40s %-12s %-20s %-20s %-20s %-12s %-8s %-8s %-35s %-15s ",
		"Id", "QueryId", "GroupId", "StartTime", "PendingTimeout", "QueryTimeout",
		"State", "Slots", "DOP", "Frontend", "FeStartTime"))

	var running, penning []string
	for i, m2 := range m {
		c := color.New()
		var state string
		if m2["State"].(string) == "RUNNING" {
			running = append(running, m2["State"].(string))
			state = c.Add(color.FgHiGreen).Sprint(m2["State"].(string))
		}
		if m2["State"].(string) == "PENDING" {
			penning = append(penning, m2["State"].(string))
			state = c.Add(color.FgHiYellow).Sprint(m2["State"].(string))
		}
		msg := fmt.Sprintf("%-3d %-40s %-12s %-20s %-20s %-20s %-21s %-8s %-8s %-35s %-15s ", i,
			m2["QueryId"].(string),
			m2["ResourceGroupId"].(string),
			m2["StartTime"].(string),
			m2["PendingTimeout"].(string),
			m2["QueryTimeout"].(string),
			state,
			m2["Slots"].(string),
			m2["DOP"].(string),
			m2["Frontend"].(string),
			m2["FeStartTime"].(string))
		if m2["State"].(string) == "RUNNING" {
			running = append(running, msg)
			fmt.Println(c.Add(color.FgHiGreen).Sprint(msg))
			continue
		}
		if m2["State"].(string) == "PENDING" {
			penning = append(penning, msg)
			fmt.Println(c.Add(color.FgHiYellow).Sprint(msg))
			continue
		}
		fmt.Println(msg)
	}
	fmt.Println(c.Add(color.FgHiYellow).Sprint(fmt.Sprintf("cnt:(%d)  RUNNING:(%d)、PENDING:(%d)", len(m), len(running), len(penning))))
}
