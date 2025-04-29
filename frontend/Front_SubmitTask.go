/*
 *@author  chengkenli
 *@project starrocks
 *@package frontend
 *@file    Front_SubmitTask
 *@date    2025/4/28 15:29
 */

package frontend

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"starrocks/util"
)

func Submittask() {
	var m []map[string]interface{}
	r := util.Connect.Raw("select * from information_schema.task_runs").Scan(&m)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return
	}
	var Pending, Running, Failed, Success []string
	for _, m2 := range m {
		marshal, _ := json.Marshal(m2)
		switch m2["STATE"].(string) {
		case "PENDING":
			Pending = append(Pending, string(marshal))
		case "RUNNING":
			Running = append(Running, string(marshal))
		case "FAILED":
			Failed = append(Failed, string(marshal))
		case "SUCCESS":
			Success = append(Success, string(marshal))
		}
	}
	if util.P.Status != "" {
		switch util.P.Status {
		case "PENDING":
			for i, item := range Pending {
				fmt.Println(i, item)
			}
		case "RUNNING":
			for i, item := range Running {
				fmt.Println(i, item)
			}
		case "FAILED":
			for i, item := range Failed {
				fmt.Println(i, item)
			}
		case "SUCCESS":
			for i, item := range Success {
				fmt.Println(i, item)
			}
		}
	} else {
		for i, item := range Running {
			fmt.Println(i, item)
		}
	}

	c := color.New()
	fmt.Println(fmt.Sprintf("Submit Task:(%d) PENDING:(%s),RUNNING:(%s),FAILED:(%s),SUCCESS:(%s)",
		len(m),
		c.Add(color.FgHiYellow).Sprint(len(Pending)),
		c.Add(color.FgHiCyan).Sprint(len(Running)),
		c.Add(color.FgHiRed).Sprint(len(Failed)),
		c.Add(color.FgHiGreen).Sprint(len(Success)),
	))
}
