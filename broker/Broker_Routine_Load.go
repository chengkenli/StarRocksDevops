package broker

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"starrocks/util"
	"strconv"
	"strings"
	"sync"
)

func Routineload() {
	fmt.Println()
	fmt.Println("routine load...")

	type Routine struct {
		Id                   int    `bson:"Id"`
		Name                 string `bson:"Name"`
		CreateTime           string `bson:"CreateTime"`
		PauseTime            string `bson:"PauseTime"`
		EndTime              string `bson:"EndTime"`
		DbName               string `bson:"DbName"`
		TableName            string `bson:"TableName"`
		State                string `bson:"State"`
		DataSourceType       string `bson:"DataSourceType"`
		CurrentTaskNum       string `bson:"CurrentTaskNum"`
		JobProperties        string `bson:"JobProperties"`
		DataSourceProperties string `bson:"DataSourceProperties"`
		CustomProperties     string `bson:"CustomProperties"`
		Statistic            string `bson:"Statistic"`
		Progress             string `bson:"Progress"`
		ReasonOfStateChanged string `bson:"ReasonOfStateChanged"`
		ErrorLogUrls         string `bson:"ErrorLogUrls"`
		OtherMsg             string `bson:"OtherMsg"`
	}

	if util.P.JobID <= 0 && util.P.JobName == "" {
		return
	}

	fmt.Println()

	var m []map[string]interface{}
	var DbName string
	r := util.Connect.Raw("show proc '/routine_loads'").Scan(&m)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return
	}
	for _, m2 := range m {
		id, _ := strconv.Atoi(m2["Id"].(string))
		if id == util.P.JobID || m2["Name"].(string) == util.P.JobName {
			DbName = m2["DbName"].(string)
			break
		}
	}

	if DbName == "" {
		fmt.Println("DbName is nil.")
		return
	} else {
		DbName = strings.ReplaceAll(DbName, "default_cluster:", "")
	}

	var routine Routine
	if util.P.JobID > 0 {
		r := util.Connect.Raw(fmt.Sprintf("show routine load from %s where id=%d", DbName, util.P.JobID)).Scan(&routine)
		if r.Error != nil {
			fmt.Println(r.Error.Error())
			return
		}
	}

	if util.P.JobName != "" {
		r := util.Connect.Raw(fmt.Sprintf("show routine load from %s where name='%s'", DbName, util.P.JobName)).Scan(&routine)
		if r.Error != nil {
			fmt.Println(r.Error.Error())
			return
		}
	}

	c := color.New()
	fmt.Println(c.Add(color.FgHiYellow).Sprint("Id                       : ") + c.Add(color.FgHiGreen).Sprint(routine.Id))
	fmt.Println(c.Add(color.FgHiYellow).Sprint("Name                     : ") + c.Add(color.FgHiGreen).Sprint(routine.Name))
	fmt.Println(c.Add(color.FgHiYellow).Sprint("CreateTime               : ") + c.Add(color.FgHiGreen).Sprint(routine.CreateTime))
	fmt.Println(c.Add(color.FgHiYellow).Sprint("PauseTime                : ") + c.Add(color.FgHiGreen).Sprint(routine.PauseTime))
	fmt.Println(c.Add(color.FgHiYellow).Sprint("EndTime                  : ") + c.Add(color.FgHiGreen).Sprint(routine.EndTime))
	fmt.Println(c.Add(color.FgHiYellow).Sprint("DbName                   : ") + c.Add(color.FgHiGreen).Sprint(routine.DbName))
	fmt.Println(c.Add(color.FgHiYellow).Sprint("TableName                : ") + c.Add(color.FgHiGreen).Sprint(routine.TableName))
	fmt.Println(c.Add(color.FgHiYellow).Sprint("State                    : ") + c.Add(color.FgHiGreen).Sprint(routine.State))
	fmt.Println(c.Add(color.FgHiYellow).Sprint("DataSourceType           : ") + c.Add(color.FgHiGreen).Sprint(routine.DataSourceType))
	fmt.Println(c.Add(color.FgHiYellow).Sprint("CurrentTaskNum           : ") + c.Add(color.FgHiGreen).Sprint(routine.CurrentTaskNum))
	fmt.Println(c.Add(color.FgHiYellow).Sprint("JobProperties            : ") + c.Add(color.FgHiGreen).Sprint(routine.JobProperties))
	fmt.Println(c.Add(color.FgHiYellow).Sprint("DataSourceProperties     : ") + c.Add(color.FgHiGreen).Sprint(routine.DataSourceProperties))
	fmt.Println(c.Add(color.FgHiYellow).Sprint("CustomProperties         : ") + c.Add(color.FgHiGreen).Sprint(routine.CustomProperties))
	fmt.Println(c.Add(color.FgHiYellow).Sprint("Statistic                : ") + c.Add(color.FgHiGreen).Sprint(routine.Statistic))
	fmt.Println(c.Add(color.FgHiYellow).Sprint("Progress                 : ") + c.Add(color.FgHiGreen).Sprint(routine.Progress))
	fmt.Println(c.Add(color.FgHiYellow).Sprint("ReasonOfStateChanged     : ") + c.Add(color.FgHiGreen).Sprint(routine.ReasonOfStateChanged))
	fmt.Println(c.Add(color.FgHiYellow).Sprint("ErrorLogUrls             : ") + c.Add(color.FgHiGreen).Sprint(routine.ErrorLogUrls))
	fmt.Println(c.Add(color.FgHiYellow).Sprint("OtherMsg                 : ") + c.Add(color.FgHiGreen).Sprint(routine.OtherMsg))

	fmt.Println()
}

func RoutineLoad() {
	type Routine []struct {
		Name          string `bson:"Name"`
		Id            int    `bson:"Id"`
		DbName        string `bson:"DbName"`
		Statistic     string `bson:"Statistic"`
		TaskStatistic string `bson:"TaskStatistic"`
	}
	type Statistic struct {
		ReceivedBytes     int64 `json:"receivedBytes"`
		ErrorRows         int   `json:"errorRows"`
		CommittedTaskNum  int   `json:"committedTaskNum"`
		LoadedRows        int   `json:"loadedRows"`
		LoadRowsRate      int   `json:"loadRowsRate"`
		AbortedTaskNum    int   `json:"abortedTaskNum"`
		TotalRows         int   `json:"totalRows"`
		UnselectedRows    int   `json:"unselectedRows"`
		ReceivedBytesRate int   `json:"receivedBytesRate"`
		TaskExecuteTimeMs int64 `json:"taskExecuteTimeMs"`
	}

	var routine Routine
	r := util.Connect.Raw("show proc '/routine_loads'").Scan(&routine)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return
	}

	c := color.New()
	fmt.Println()
	fmt.Println(c.Add(color.FgHiYellow).Sprint(fmt.Sprintf("ROUTINE LOAD(%d):", len(routine))))
	fmt.Println(fmt.Sprintf("%-5s %-10s %-80s %-15s %-7s %-15s %-10s %-10s %-10s %-10s %-15s %-15s %-10s",
		"Num", "JobID", "JobName", "ReceivedBytes", "ErrRows", "CommittedTaskNum", "LoadedRows", "LoadRowsRate", "AbortedTaskNum", "TotalRows", "UnselectedRows", "ReceivedBytesRate", "TaskExecuteTimeMs"))

	var wg sync.WaitGroup
	for i, s := range routine {
		s := s
		if !util.P.RoutineLoad {
			if i >= 5 {
				break
			}
		}
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			var statistic Statistic
			err := json.Unmarshal([]byte(s.Statistic), &statistic)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			fmt.Println(fmt.Sprintf("%-5d %-10d %-80s %-15d %-7d %-16d %-10d %-12d %-14d %-10d %-15d %-17d %-10d",
				i, s.Id, fmt.Sprintf("%s.%s", strings.ReplaceAll(s.DbName, "default_cluster:", ""), s.Name),
				statistic.ReceivedBytes, statistic.ErrorRows, statistic.CommittedTaskNum, statistic.LoadedRows, statistic.LoadRowsRate, statistic.AbortedTaskNum, statistic.TotalRows, statistic.UnselectedRows, statistic.ReceivedBytesRate, statistic.TaskExecuteTimeMs))

		}(i)
	}
	wg.Wait()
	fmt.Println()
}
