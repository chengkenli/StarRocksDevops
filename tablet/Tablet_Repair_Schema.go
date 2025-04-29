package tablet

import (
	"fmt"
	"github.com/fatih/color"
	"starrocks/broker"
	"starrocks/util"
	"strconv"
	"strings"
	"sync"
	"time"
)

func RepairSchema() {

	if util.P.BackendId <= 0 {
		fmt.Println("BackendId is nil, 请指定 -bi <BackendId>")
		return
	}

	var dbs []map[string]interface{}
	r := util.Connect.Raw("show proc '/dbs'").Scan(&dbs)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return
	}
	/*获取所有表数量*/
	var TableNum, TableC int
	for _, m := range dbs {
		dbname := strings.ReplaceAll(m["DbName"].(string), "default_cluster:", "")
		if util.P.Dbname != "" {
			if util.P.Dbname != dbname {
				continue
			}
		}
		atoi, _ := strconv.Atoi(m["TableNum"].(string))
		TableNum += atoi
	}
	fmt.Println(fmt.Sprintf("%s scan:%d，thread:%d，time:%s", time.Now().Format("2006-01-02 15:04:05"), TableNum, util.P.Thread, time.Now().Format("2006-01-02 15:04:05")))

	fmt.Println("正在扫描...")
	/*start*/
	var doneC = make(chan int)
	// 进度条
	tic := time.Tick(3 * time.Second)
	go func(c chan int) {
		for {
			select {
			case <-doneC:
				return
			case <-tic:
				processRate := float64(TableC) / float64(TableNum) * 100
				rate := strconv.FormatFloat(processRate, 'f', 2, 64)
				c := color.New()
				fmt.Println(c.Add(color.FgHiGreen).Sprint(fmt.Sprintf("%s %s", fmt.Sprintf("scaning.."), strings.Repeat("#", int(processRate/2))+rate+"%")))
				if TableC >= TableNum {
					return
				}
			}
		}
	}(doneC)

	ch := make(chan struct{}, util.P.Thread)

	for _, m := range dbs {
		var tableName []map[string]interface{}
		r := util.Connect.Raw(fmt.Sprintf("show proc '/dbs/%s'", m["DbId"].(string))).Scan(&tableName)
		if r.Error != nil {
			fmt.Println(r.Error.Error())
			return
		}
		dbname := strings.ReplaceAll(m["DbName"].(string), "default_cluster:", "")

		if util.P.Dbname != "" {
			if util.P.Dbname != dbname {
				continue
			}
		}

		var wait sync.WaitGroup
		for _, tb := range tableName {
			wait.Add(1)
			tb := tb

			go func() {
				defer func() {
					<-ch
					wait.Done()
				}()

				ch <- struct{}{}

				if tb["Type"].(string) != "OLAP" {
					TableC++
					return
				}
				schemaTablename := fmt.Sprintf("%s.%s", dbname, tb["TableName"].(string))
				var m map[string]interface{}

				for i := 0; i < 5; i++ {
					r := util.Connect.Raw(fmt.Sprintf("select count(*) from %s", schemaTablename)).Scan(&m)
					if r.Error != nil {
						fmt.Println(r.Error.Error())
						if strings.Contains(r.Error.Error(), "tablet_id") {
							tabletID := strings.Split(strings.Split(r.Error.Error(), "tablet_id: ")[1], ",")[0]
							checkTablet(tabletID)
							break
						}
						if strings.Contains(r.Error.Error(), "no queryable replica found in tablet:") {
							tabletID := strings.Split(r.Error.Error(), "tablet: ")[1]
							checkTablet(tabletID)
							break
						}
					} else {
						if i < 5 {
							fmt.Println(fmt.Sprintf("%s 检索矩阵探测: %d ", time.Now().Format("2006-01-02 15:04:05"), i) + schemaTablename)
							i++
							continue
						}
						break
					}
				}

				/*----------------------------------------------------------------------------------*/
				TableC++
			}()
		}
		wait.Wait()
	}

	fmt.Println(fmt.Sprintf("%s %s %s", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("done....."), strings.Repeat("#", 50)+"100%"))
	fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " done.")
}

func checkTablet(tabletID string) {
	c := color.New()

	type TABLET struct {
		DbName        string `bson:"DbName"`
		TableName     string `bson:"TableName"`
		PartitionName string `bson:"PartitionName"`
		IndexName     string `bson:"IndexName"`
		DbId          int    `bson:"DbId"`
		TableId       int    `bson:"TableId"`
		PartitionId   int    `bson:"PartitionId"`
		IndexId       int    `bson:"IndexId"`
		IsSync        bool   `bson:"IsSync"`
		DetailCmd     string `bson:"DetailCmd"`
	}
	var tablet TABLET
	r := util.Connect.Raw(fmt.Sprintf("show tablet %s", tabletID)).Scan(&tablet)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return
	}
	length, on := broker.Backendids(tablet.DetailCmd)
	if length < 2 {
		sql := fmt.Sprintf("truncate table %s.%s partition(%s)", strings.ReplaceAll(tablet.DbName, "default_cluster:", ""), tablet.TableName, tablet.PartitionName)
		r := util.Connect.Exec(sql)
		if r.Error != nil {
			fmt.Println(r.Error.Error())
			return
		}
		fmt.Println(c.Add(color.FgHiMagenta).Sprint(fmt.Sprintf("%s 单副本重建列: %s.%s (%s) [%s]", time.Now().Format("2006-01-02 15:04:05"), strings.ReplaceAll(tablet.DbName, "default_cluster:", ""), tablet.TableName, tablet.PartitionName, tabletID)))
		return
	}

	if on {
		sql := fmt.Sprintf(`ADMIN SET REPLICA STATUS PROPERTIES("tablet_id" = "%s", "backend_id" = "%d", "status" = "bad")`, tabletID, util.P.BackendId)
		r := util.Connect.Exec(sql)
		if r.Error != nil {
			fmt.Println(r.Error.Error())
			return
		}

		fmt.Println(c.Add(color.FgHiBlue).Sprint(fmt.Sprintf("%s 修复矩阵提交: %s.%s (%s) [%s]", time.Now().Format("2006-01-02 15:04:05"), strings.ReplaceAll(tablet.DbName, "default_cluster:", ""), tablet.TableName, tablet.PartitionName, tabletID)))
	}
}
