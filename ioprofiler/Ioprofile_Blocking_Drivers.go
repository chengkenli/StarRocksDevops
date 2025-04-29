/*
 *@author  chengkenli
 *@project starrocks
 *@package ioprofiler
 *@file    scan_pipeline_blocking_drivers
 *@date    2024/12/31 15:10
 */

package ioprofiler

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"regexp"
	"starrocks/util"
	"strings"
)

type blocking struct {
	QueriesInWorkgroup []struct {
		QueryID   string `json:"query_id"`
		Fragments []struct {
			FragmentID     string `json:"fragment_id"`
			FragmentStatus string `json:"fragment_status"`
			Drivers        []struct {
				DriverID   int    `json:"driver_id"`
				State      string `json:"state"`
				DriverDesc string `json:"driver_desc"`
			} `json:"drivers"`
		} `json:"fragments"`
	} `json:"queries_in_workgroup"`
}

func pipeline_blocking_drivers(Client *resty.Client, ip string) {
	c := color.New()

	uri := fmt.Sprintf("http://%s:8040/api/pipeline_blocking_drivers/stat", ip)
	//发送POST请求并处理响应
	respones, err := Client.R().Get(uri)
	if err != nil {
		fmt.Println(c.Add(color.FgHiRed).Sprint(err.Error()))
		return
	}
	var b blocking
	err = json.Unmarshal(respones.Body(), &b)
	if err != nil {
		fmt.Println(c.Add(color.FgHiRed).Sprint(err.Error()))
		return
	}

	for o, item := range b.QueriesInWorkgroup {

		if len(util.P.QueryId) != 0 {
			if item.QueryID != util.P.QueryId {
				continue
			}
			util.P.List = true
		}

		for _, fragment := range item.Fragments {

			c := color.New()
			fmt.Println(fmt.Sprintf("%-2d %-3s %-49s %-40s:", o, fragment.FragmentStatus, c.Add(color.FgHiCyan).Sprint(item.QueryID), c.Add(color.FgHiBlue).Sprint(fragment.FragmentID)))

			for _, driver := range fragment.Drivers {
				// 正则表达式
				name := regexp.MustCompile(`driver=([a-zA-Z0-9_]+)`).FindStringSubmatch(driver.DriverDesc)
				//status := regexp.MustCompile(`status=([a-zA-Z_]+)`).FindStringSubmatch(driver.DriverDesc)
				desc := regexp.MustCompile(`operator-chain: \[([^\]]+)\]`).FindStringSubmatch(driver.DriverDesc)
				detail := strings.Split(desc[1], "->")

				var msg string
				if util.P.List {
					var data []string
					for i, s := range detail {
						if i == 0 {
							data = append(data, fmt.Sprintf("%.2d > %s", i, s))
							continue
						}
						data = append(data, fmt.Sprintf("%s %.2d > %s", strings.Repeat(" ", 47), i, s))
					}

					msg = fmt.Sprintf("%-5d %-20s %-20s %-10s",
						driver.DriverID,
						driver.State,
						name[1],
						strings.Join(data, c.Add(color.FgHiYellow).Sprint("->")+"\n"),
					)
				} else {
					msg = fmt.Sprintf("%-5d %-20s %-20s",
						driver.DriverID,
						driver.State,
						name[1],
					)
				}
				fmt.Println(msg)

			}
		}
	}

	//fmt.Println(string(respones.Body()))
}
