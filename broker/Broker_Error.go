package broker

import (
	"fmt"
	"github.com/fatih/color"
	"starrocks/util"
	"strings"
	"sync"
)

func BrokerError() {
	c := color.New()
	fmt.Println(c.Add(color.FgHiYellow).Sprint("Broker Session: "))
	var m []map[string]interface{}
	r := util.Connect.Raw("show proc '/dbs'").Scan(&m)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return
	}

	var wg sync.WaitGroup
	var i int
	for _, m2 := range m {
		wg.Add(1)
		go func(m2 map[string]interface{}) {
			defer wg.Done()
			database := strings.ReplaceAll(m2["DbName"].(string), "default_cluster:", "")
			var tbLoads []map[string]interface{}
			r := util.Connect.Raw(fmt.Sprintf("show load from %s where STATE = 'CANCELLED'", database)).Scan(&tbLoads)
			if r.Error != nil {
				fmt.Println(r.Error.Error())
				return
			}
			for _, load := range tbLoads {
				fmt.Println(fmt.Sprintf("%-5d%-22s%-18s%-50s%-40s", i, load["CreateTime"].(string), load["JobId"].(string), fmt.Sprintf("%s.%s", database, load["Label"].(string)), c.Add(color.FgHiRed).Sprint(load["ErrorMsg"].(string))))
				i++
			}
		}(m2)
	}
	wg.Wait()
	fmt.Println()
}

func Backendids(sql string) (int, bool) {
	type PROCS []struct {
		ReplicaId             int    `bson:"ReplicaId"`
		BackendId             int    `bson:"BackendId"`
		Version               int    `bson:"Version"`
		VersionHash           int    `bson:"VersionHash"`
		LstSuccessVersion     int    `bson:"LstSuccessVersion"`
		LstSuccessVersionHash int    `bson:"LstSuccessVersionHash"`
		LstFailedVersion      int    `bson:"LstFailedVersion"`
		LstFailedVersionHash  int    `bson:"LstFailedVersionHash"`
		LstFailedTime         string `bson:"LstFailedTime"`
		SchemaHash            int    `bson:"SchemaHash"`
		DataSize              int    `bson:"DataSize"`
		RowCount              int    `bson:"RowCount"`
		State                 string `bson:"State"`
		IsBad                 bool   `bson:"IsBad"`
		IsSetBadForce         bool   `bson:"IsSetBadForce"`
		VersionCount          string `bson:"VersionCount"`
		PathHash              string `bson:"PathHash"`
		MetaUrl               string `bson:"MetaUrl"`
		CompactionStatus      string `bson:"CompactionStatus"`
	}

	var proc PROCS
	r := util.Connect.Raw(sql).Scan(&proc)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return len(proc), false
	}

	for _, s := range proc {
		if s.BackendId == util.P.BackendId && !s.IsBad {
			return len(proc), true
		}
	}
	return len(proc), false
}
