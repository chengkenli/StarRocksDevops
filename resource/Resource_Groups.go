/*
 *@author  chengkenli
 *@project starrocks
 *@package resource
 *@file    GroupResource
 *@date    2024/6/14 11:37
 */

package resource

import (
	"fmt"
	"github.com/fatih/color"
	"starrocks/tools"
	"starrocks/util"
	"strconv"
)

// 定义单位常量
const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
	PB = 1024 * TB
)

// ByteSizeToString 将字节数转换为人类可读的字符串表示形式（如 KB、MB、GB）
func byteSizeToString(s string) string {
	size, _ := strconv.ParseFloat(s, 64)
	if size < KB {
		return fmt.Sprintf("%.2fB", size)
	} else if size < MB {
		return fmt.Sprintf("%.2fKB", size/KB)
	} else if size < GB {
		return fmt.Sprintf("%.2fMB", size/MB)
	} else if size < TB {
		return fmt.Sprintf("%.2fGB", size/GB)
	} else if size < PB {
		return fmt.Sprintf("%.2fTB", size/TB)
	} else {
		return fmt.Sprintf("%.2fPB", size/PB)
	}
}

func GroupResource() {
	type group2 []struct {
		Name             string `bson:"Name"`
		Id               int    `bson:"Id"`
		CPUCoreLimit     int    `bson:"CPUCoreLimit"`
		MemLimit         string `bson:"MemLimit"`
		ConcurrencyLimit string `bson:"ConcurrencyLimit"`
		Type             string `bson:"Type"`
		Classifiers      string `bson:"Classifiers"`
	}

	c := color.New()
	fmt.Println()
	fmt.Println(c.Add(color.FgHiYellow).Sprint("资源隔离:"))
	switch tools.Version() {
	case 2.2, 2.3, 2.5:
		var m map[string]interface{}
		r := util.Connect.Raw("SHOW VARIABLES LIKE 'enable_resource_group'").Scan(&m)
		if r.Error != nil {
			return
		}
		if m == nil {
			return
		}

		fmt.Println(m["Value"].(string))
		var g group2
		r = util.Connect.Raw("SHOW RESOURCE GROUPS ALL").Scan(&g)
		if r.Error != nil {
			return
		}
		fmt.Println(fmt.Sprintf("%-5s%-20s%-15s%-18s%-12s%-22s%-10s%-20s", "On", "Name", "Id", "CpuCoreLimit", "MemLimit", "ConcurrencyLimit", "Type", "Classifiers"))

		for i, s := range g {
			msg := fmt.Sprintf("%-5d%-20s%-15d%-18d%-12s%-22s%-10s%-20s", i,
				s.Name,
				s.Id,
				s.CPUCoreLimit,
				s.MemLimit,
				s.ConcurrencyLimit,
				s.Type,
				s.Classifiers,
			)
			fmt.Println(msg)
		}
	case 3.0, 3.1, 3.2:
		var m2 map[string]interface{}
		r := util.Connect.Raw("SHOW VARIABLES LIKE 'enable_pipeline_engine'").Scan(&m2)
		if r.Error != nil {
			return
		}
		fmt.Println(c.Add(color.FgHiGreen).Sprint(m2["Value"].(string)))
		var m []map[string]interface{}
		row := util.Connect.Raw("SHOW RESOURCE GROUPS ALL").Scan(&m)
		if row.Error != nil {
			fmt.Println(row.Error.Error())
			return
		}
		fmt.Println(fmt.Sprintf("%-6s%-25s%-13s%-10s%-10s%-15s%-15s%-13s%-15s%-15s%-10s%-10s%-20s",
			"num", "name", "id", "core", "mem", "max_cores", "<b>cpu_second", "<b>scan_rows", "<b>mem", "concurrency", "spill", "type", "classifiers"),
		)

		for i, s := range m {
			msg := fmt.Sprintf("%-6d%-25s%-13s%-10s%-10s%-15s%-15s%-13s%-15s%-15s%-10s%-10s%-20s",
				i,
				s["name"].(string),
				s["id"].(string),
				s["cpu_core_limit"].(string),
				s["mem_limit"].(string),
				s["max_cpu_cores"].(string),
				s["big_query_cpu_second_limit"].(string),
				s["big_query_scan_rows_limit"].(string),
				byteSizeToString(s["big_query_mem_limit"].(string)),
				s["concurrency_limit"].(string),
				s["spill_mem_limit_threshold"].(string),
				s["type"].(string),
				s["classifiers"].(string),
			)
			if s["name"].(string) == "ap" {
				fmt.Println(c.Add(color.FgHiBlue).Sprint(msg))
				continue
			}
			if s["name"].(string) == "public" {
				c := color.New(color.BgGreen)
				fmt.Println(c.Add(color.FgHiRed).Sprint(msg))
				continue
			}
			fmt.Println(msg)
			if i >= 5 {
				fmt.Println(fmt.Sprintf("ReSourceGroup:[%d]...", len(m)))
				break
			}
		}
	case 3.3:
		list33(util.Connect)
	}
	fmt.Println()
}
