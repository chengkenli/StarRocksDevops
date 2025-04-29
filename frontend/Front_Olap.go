package frontend

import (
	"fmt"
	"starrocks/util"
	"strings"
)

func SessionOlap() {

	var dbs []map[string]interface{}
	r := util.Connect.Raw("show proc '/dbs'").Scan(&dbs)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return
	}

	var i int
	fmt.Println()
	fmt.Println(fmt.Sprintf("%-4s %-12s %-18s %-20s %s", "ID", "TableId", "Type", "Schema", "Tablename"))
	for _, m := range dbs {
		var tableName []map[string]interface{}
		r := util.Connect.Raw(fmt.Sprintf("show proc '/dbs/%s'", m["DbId"].(string))).Scan(&tableName)
		if r.Error != nil {
			fmt.Println(r.Error.Error())
			continue
		}
		dbname := strings.ReplaceAll(m["DbName"].(string), "default_cluster:", "")
		if dbname == "information_schema" || dbname == "starrocks_monitor" || dbname == "_statistics_" {
			continue
		}

		for _, tb := range tableName {
			for _, t := range strings.Split(util.P.OlapType, ",") {
				if tb["Type"].(string) == t {
					fmt.Println(fmt.Sprintf("%-4d %-12s %-18s %-20s %s.%s", i, tb["TableId"].(string), t, dbname, dbname, tb["TableName"].(string)))
					i++
				}
			}
		}
	}
	fmt.Println()
}
