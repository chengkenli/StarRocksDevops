/*
 *@author  chengkenli
 *@project starrocks
 *@package resource
 *@file    CurrentQueries
 *@date    2024/7/11 10:08
 */

package resource

import (
	"fmt"
	"github.com/fatih/color"
	"starrocks/tools"
	"strconv"
	"strings"
)

func CurrentQueries() {
	if tools.Version() < 3.0 {
		return
	}

	RunningQueries()

	var sum int
	var nodecnts []string
	c := color.New()
	fmt.Println()
	fmt.Println(c.Add(color.FgHiYellow).Sprint("当前队列:"))
	fmt.Println(fmt.Sprintf("%-3s %-22s %-40s %-15s %-20s %-25s %-15s %-15s %-15s %-15s %-15s %-15s",
		"Id", "StartTime", "QueryId", "ConnectionId", "Database", "User", "ScanBytes", "ScanRows", "MemoryUsage", "DiskSpillSize", "CPUTime", "ExecTime"))

	queries := UriCurrentQueries(client, leader)
	sum = sum + len(queries)
	nodecnts = append(nodecnts, fmt.Sprintf("Leader %s:Query(%d)", leader, len(queries)))

	for i, m2 := range queries {
		c := color.New()
		msg := fmt.Sprintf("%-3d %-22s %-40s %-15s %-20s %-25s %-15s %-15s %-15s %-15s %-15s %-15s ", i,
			m2.StartTime,
			m2.QueryId,
			m2.ConnectionId,
			m2.Database,
			m2.User,
			m2.ScanBytes,
			strings.NewReplacer(" rows", "").Replace(m2.ScanRows),
			m2.MemoryUsage,
			m2.DiskSpillSize,
			m2.CPUTime,
			m2.ExecTime)

		sr, _ := strconv.ParseInt(strings.Split(m2.ScanRows, " ")[0], 10, 64)
		if sr >= 200000000 && sr < 10000000000 {
			fmt.Println(c.Add(color.FgHiMagenta).Sprint(msg))
			continue
		}
		if sr >= 10000000000 {
			fmt.Println(c.Add(color.FgHiRed).Sprint(msg))
			continue
		}
		c = color.New()
		if strings.Contains(m2.ScanBytes, "GB") || strings.Contains(m2.MemoryUsage, "GB") {
			scanByte, _ := strconv.ParseInt(strings.Split(m2.ScanBytes, " ")[0], 10, 64)
			memoryUsage, _ := strconv.ParseInt(strings.Split(m2.MemoryUsage, " ")[0], 10, 64)
			if scanByte >= 256 || memoryUsage >= 256 {
				fmt.Println(c.Add(color.FgHiRed).Sprint(msg))
				continue
			}
		}
		if strings.Contains(m2.ScanBytes, "TB") || strings.Contains(m2.MemoryUsage, "TB") {
			fmt.Println(c.Add(color.FgHiRed).Sprint(msg))
			continue
		}
		fmt.Println(c.Add(color.FgHiWhite).Sprint(msg))
	}
	c = color.New()
	fmt.Println(fmt.Sprintf("cnt:(%s)  %v", c.Add(color.FgHiCyan).Sprint(sum), strings.Join(nodecnts, "、")))
}
