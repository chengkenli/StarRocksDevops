/*
 *@author  chengkenli
 *@project StarRocksSupport
 *@package fronends
 *@file    frontendQueris
 *@date    2024/9/5 11:44
 */

package resource

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/go-resty/resty/v2"
	"golang.org/x/net/html"
	"starrocks/util"
	"strings"
)

func UriCurrentQueries(client *resty.Client, fe string) util.Queris {
	var qs util.Queris

	uri := fmt.Sprintf(`http://%s:8030/system?path=//current_queries`, fe)
	//创建Resty客户端
	//发送POST请求并处理响应
	respones, err := client.R().Get(uri)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	menu, _ := htmlquery.Parse(strings.NewReader(string(respones.Body())))
	table := htmlquery.Find(menu, `//*[@id="table_id"]/tbody/tr`)
	for _, node := range table {
		tr := htmlquery.Find(node, "td")
		if len(tr) >= 11 {
			var wh string
			if len(tr) == 12 {
				wh = td(tr[11])
			}
			qs = append(qs, util.Querisign{
				StartTime:     td(tr[0]),
				QueryId:       td(tr[1]),
				ConnectionId:  td(tr[2]),
				Database:      td(tr[3]),
				User:          td(tr[4]),
				ScanBytes:     td(tr[5]),
				ScanRows:      td(tr[6]),
				MemoryUsage:   td(tr[7]),
				DiskSpillSize: td(tr[8]),
				CPUTime:       td(tr[9]),
				ExecTime:      td(tr[10]),
				Warehouse:     wh,
			})
		}
	}
	return qs
}

func UriCurrentQueriesStmt(fe, id string) string {
	uri := fmt.Sprintf(`http://%s:8030/system?path=//current_queries/%s`, fe, id)
	//发送POST请求并处理响应
	respones, err := client.R().SetBasicAuth(avg.User, avg.Pass).Get(uri)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	menu, _ := htmlquery.Parse(strings.NewReader(string(respones.Body())))
	var sql *html.Node
	if menu != nil {
		sql = htmlquery.FindOne(menu, `//*[@id="table_id"]/tbody/tr/td/a`)
	}
	var stmt string
	if sql != nil {
		stmt = htmlquery.InnerText(sql)
	}
	return stmt
}

func UriCurrentQueriesHosts(fe, id string) []string {
	var nodes []string

	uri := fmt.Sprintf(`http://%s:8030/system?path=//current_queries/%s/hosts`, fe, id)
	//发送POST请求并处理响应
	respones, err := client.R().SetBasicAuth(avg.User, avg.Pass).Get(uri)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	nodes = append(nodes, fmt.Sprintf("  %-2s %-20s %-15s %-15s %-15s %-15s", "ID", "Host", "ScanBytes", "ScanRows", "CpuCostSeconds", "MemUsageBytes"))

	menu, _ := htmlquery.Parse(strings.NewReader(string(respones.Body())))
	tbody := htmlquery.Find(menu, `//*[@id="table_id"]/tbody/tr`)
	if tbody != nil {
		for i, body := range tbody {
			if body == nil {
				continue
			}
			tr := htmlquery.Find(body, "td")
			if len(tr) == 5 {
				msg := fmt.Sprintf("> %-2d %-20s %-15s %-15s %-15s %-15s ", i, td(tr[0]), td(tr[1]), td(tr[2]), td(tr[3]), td(tr[4]))
				nodes = append(nodes, msg)
			}
		}
	}
	return nodes
}

func td(n *html.Node) string {
	v := htmlquery.InnerText(n)
	return v
}
