package broker

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"os"
	"starrocks/util"
	"strconv"
	"strings"
	"sync"
	"time"
)

func StreamList() {
	if getv() < 3.3 {
		return
	}
	fmt.Println()
	var (
		COMMITED    []string
		BEFORE_LOAD []string
		FINISHED    []string
		CANCELLED   []string
		LOADING     []string
	)
	var streamData []map[string]interface{}
	r := util.Connect.Raw("show proc '/stream_loads'").Scan(&streamData)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return
	}
	c := color.New()
	fmt.Println(c.Add(color.FgHiYellow).Sprint("STREAM LOAD:"))
	fmt.Println(fmt.Sprintf("%-3s %-90s %-10s %-50s %-10s %-10v %-21s %-21s %-10s", "N", "Label", "Id", "TableName", "State", "Edtime", "CreateTimeMs", "BeforeLoadTimeMs", "EndTimeMs"))

	ch := make(chan struct{}, 20)
	var wg sync.WaitGroup
	for i, i2 := range streamData {
		wg.Add(1)
		i2 := i2
		i := i
		go func() {
			defer func() {
				<-ch
				wg.Done()
			}()

			ch <- struct{}{}

			var sdata map[string]interface{}
			sql := fmt.Sprintf("show proc '/stream_loads/%s'", i2["Label"].(string))
			r := util.Connect.Raw(sql).Scan(&sdata)
			if r.Error != nil {
				fmt.Println(r.Error.Error())
				return
			}

			var edtime string
			if v, ok := sdata["CreateTimeMs"]; ok {
				if v2, ok2 := sdata["EndTimeMs"]; ok2 {
					if v != nil && v2 != nil {
						t1, _ := time.Parse("2006-01-02 15:04:05", v.(string))
						t2, _ := time.Parse("2006-01-02 15:04:05", v2.(string))
						edtime = t2.Sub(t1).String()
					}
				}
			}
			msg := fmt.Sprintf("%-3d %-90s %-10s %-50s %-10v %-10v %-21v %-21v %-10v",
				i, i2["Label"].(string), i2["Id"].(string), fmt.Sprintf("%s.%s", i2["DbName"].(string), i2["TableName"].(string)),
				i2["State"].(string), edtime, sdata["CreateTimeMs"], sdata["BeforeLoadTimeMs"], sdata["EndTimeMs"])

			//msg := fmt.Sprintf("%-3d %-90s %-10s %-50s %s", i, i2["Label"].(string), i2["Id"].(string), fmt.Sprintf("%s.%s", i2["DbName"].(string), i2["TableName"].(string)), i2["State"].(string))
			switch i2["State"].(string) {
			case "FINISHED":
				FINISHED = append(FINISHED, i2["Label"].(string))
			case "CANCELLED":
				CANCELLED = append(CANCELLED, i2["Label"].(string))
			case "BEFORE_LOAD":
				BEFORE_LOAD = append(BEFORE_LOAD, i2["Label"].(string))
			case "COMMITED":
				COMMITED = append(COMMITED, i2["Label"].(string))
				fmt.Println(msg)
			case "LOADING":
				LOADING = append(LOADING, i2["Label"].(string))
				fmt.Println(msg)
			default:
				fmt.Println(msg)
			}
		}()
	}
	wg.Wait()

	fmt.Println(c.Add(color.FgHiYellow).Sprint(fmt.Sprintf("LOADING:(%d)、FINISHED:(%d)、CANCELLED:(%d)、BEFORE_LOAD:(%d)、COMMITED:(%d)", len(LOADING), len(FINISHED), len(CANCELLED), len(BEFORE_LOAD), len(COMMITED))))
}

func StreamloadsStatus() {
	if getv() < 3.3 {
		return
	}
	var (
		COMMITED    []string
		BEFORE_LOAD []string
		FINISHED    []string
		CANCELLED   []string
	)

	var streamData []map[string]interface{}
	r := util.Connect.Raw("show proc '/stream_loads'").Scan(&streamData)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return
	}
	fmt.Println(fmt.Sprintf("%-3s %-90s %-10s %-50s %-10s %-10v %-21s %-21s %-10s", "N", "Label", "Id", "TableName", "State", "Edtime", "CreateTimeMs", "BeforeLoadTimeMs", "EndTimeMs"))

	ch := make(chan struct{}, 5)
	var wg sync.WaitGroup
	for i, i2 := range streamData {
		wg.Add(1)
		i2 := i2
		i := i
		go func() {
			defer func() {
				<-ch
				wg.Done()
			}()

			ch <- struct{}{}
			var sdata map[string]interface{}
			sql := fmt.Sprintf("show proc '/stream_loads/%s'", i2["Label"].(string))
			r := util.Connect.Raw(sql).Scan(&sdata)
			if r.Error != nil {
				fmt.Println(r.Error.Error())
				return
			}

			var edtime string
			if v, ok := sdata["CreateTimeMs"]; ok {
				if v2, ok2 := sdata["EndTimeMs"]; ok2 {
					if v != nil && v2 != nil {
						t1, _ := time.Parse("2006-01-02 15:04:05", v.(string))
						t2, _ := time.Parse("2006-01-02 15:04:05", v2.(string))
						edtime = t2.Sub(t1).String()
					}
				}
			}
			c := color.New()
			var state string
			msg := fmt.Sprintf("%-3d %-90s %-10s %-50s %-10v %-10v %-21v %-21v %-10v",
				i, i2["Label"].(string), i2["Id"].(string), fmt.Sprintf("%s.%s", i2["DbName"].(string), i2["TableName"].(string)),
				i2["State"].(string), edtime, sdata["CreateTimeMs"], sdata["BeforeLoadTimeMs"], sdata["EndTimeMs"])

			switch i2["State"].(string) {
			case "FINISHED":
				state = c.Add(color.FgHiGreen).Sprint(msg)
				FINISHED = append(FINISHED, state)
			case "CANCELLED":
				state = c.Add(color.FgHiRed).Sprint(msg)
				CANCELLED = append(CANCELLED, state)
			case "BEFORE_LOAD":
				state = c.Add(color.FgHiYellow).Sprint(msg)
				BEFORE_LOAD = append(BEFORE_LOAD, state)
			case "COMMITED":
				state = c.Add(color.FgHiBlue).Sprint(msg)
				COMMITED = append(COMMITED, state)
			default:
				state = msg
			}

			if len(util.P.Label) != 0 {
				if util.P.Label == i2["Label"].(string) {
					fmt.Println(state)
					marshal, _ := json.MarshalIndent(sdata, "", "  ")
					fmt.Println(string(marshal))
					os.Exit(-1)
				}
				return
			}

			if len(util.P.Status) == 0 {
				atoi, _ := strconv.Atoi(i2["Id"].(string))
				if util.P.JobID == atoi || util.P.JobName == i2["Label"].(string) {
					fmt.Println(state)
					marshal, _ := json.MarshalIndent(sdata, "", "  ")
					fmt.Println(string(marshal))
					return
				}

				if util.P.JobID < 1 {
					fmt.Println(state)
				}
			}
		}()
	}
	wg.Wait()

	switch util.P.Status {
	case "FINISHED":
		fmt.Println(strings.Join(FINISHED, "\n"))
		os.Exit(-1)
	case "CANCELLED":
		fmt.Println(strings.Join(CANCELLED, "\n"))
		os.Exit(-1)
	case "BEFORE_LOAD":
		fmt.Println(strings.Join(BEFORE_LOAD, "\n"))
		os.Exit(-1)
	case "COMMITED":
		fmt.Println(strings.Join(COMMITED, "\n"))
		os.Exit(-1)
	}

}

func getv() float64 {
	/*匹配starrocks版本*/
	var m map[string]interface{}
	util.Connect.Raw("select current_version() as version").Scan(&m)
	v := strings.Join(strings.Split(m["version"].(string), ".")[0:2], ".")
	version, _ := strconv.ParseFloat(v, 64)
	return version
}
