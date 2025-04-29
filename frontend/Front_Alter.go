package frontend

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"starrocks/util"
	"sync"
)

func SessionAlter() {
	type alter []struct {
		JobId         string `json:"JobId"`
		TableName     string `json:"TableName"`
		CreateTime    string `json:"CreateTime"`
		FinishTime    string `json:"FinishTime"`
		IndexName     string `json:"IndexName"`
		IndexId       string `json:"IndexId"`
		OriginIndexId string `json:"OriginIndexId"`
		SchemaVersion string `json:"SchemaVersion"`
		TransactionId string `json:"TransactionId"`
		State         string `json:"State"`
		Msg           string `json:"Msg"`
		Progress      string `json:"Progress"`
		Timeout       string `json:"Timeout"`
	}

	if len(util.P.ALterName) != 0 && util.P.Kill {
		r := util.Connect.Exec(fmt.Sprintf("CANCEL ALTER TABLE COLUMN FROM %s", util.P.ALterName))
		if r.Error != nil {
			fmt.Println(r.Error.Error())
			os.Exit(-1)
		}
		c := color.New()
		fmt.Println(fmt.Sprintf("cancel alter table column from %s is %s!", util.P.ALterName, c.Add(color.FgHiGreen).Sprint("done")))
		os.Exit(-1)
	}

	var m []map[string]interface{}
	r := util.Connect.Raw("show databases").Scan(&m)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return
	}

	fmt.Println()
	c := color.New()
	fmt.Println(c.Add(color.FgHiYellow).Sprint("Alter Table:"))

	fmt.Println(fmt.Sprintf("%-70s%-25s%-25s%-15s%-15s%-20s%-15s%-15s", "TableName", "CreateTime", "FinishTime", "OriginIndexId", "SchemaVersion", "TransactionId", "State", "Progress"))
	var wg sync.WaitGroup
	for _, m2 := range m {
		wg.Add(1)
		m2 := m2
		go func() {
			defer wg.Done()
			var a alter
			states := []string{"RUNNING", "PENDING"}
			for _, state := range states {
				r = util.Connect.Raw(fmt.Sprintf("USE %s;SHOW ALTER TABLE COLUMN where State='%s'", m2["Database"], state)).Scan(&a)
				if r.Error != nil {
					fmt.Println(r.Error.Error())
					continue
				}
				if a == nil {
					continue
				}
				for _, s := range a {
					c := color.New()
					if s.State == "RUNNING" {
						fmt.Println(c.Add(color.FgHiGreen).Sprint(fmt.Sprintf("%-70s%-25s%-25s%-15s%-15s%-20s%-15s%-15s", fmt.Sprintf("%s.%s", m2["Database"].(string), s.TableName), s.CreateTime, s.FinishTime, s.OriginIndexId, s.SchemaVersion, s.TransactionId, s.State, s.Progress)))
					} else {
						fmt.Println(fmt.Sprintf("%-70s%-25s%-25s%-15s%-15s%-20s%-15s%-15s", fmt.Sprintf("%s.%s", m2["Database"].(string), s.TableName), s.CreateTime, s.FinishTime, s.OriginIndexId, s.SchemaVersion, s.TransactionId, s.State, s.Progress))
					}
				}
			}
		}()
	}
	wg.Wait()
	fmt.Println()
}
