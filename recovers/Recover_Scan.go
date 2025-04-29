/*
 *@author  chengkenli
 *@project starrocks
 *@package recover
 *@file    recover_scan
 *@date    2025/1/17 10:08
 */

package recovers

import (
	"fmt"
	"github.com/fatih/color"
	"starrocks/util"
	"strconv"
	"strings"
	"time"
)

func Scan() {
	c := color.New()

	var n map[string]interface{}
	r := util.Connect.Raw("ADMIN SHOW FRONTEND CONFIG LIKE 'catalog_trash_expire_second'").Scan(&n)
	if r.Error != nil {
		fmt.Println(c.Add(color.FgHiRed).Sprint(r.Error.Error()))
		return
	}
	expire, _ := strconv.Atoi(n["Value"].(string))
	fmt.Println(expire)

	fmt.Println(c.Add(color.FgHiBlue).Sprint("正在解析"))
	var m []map[string]interface{}
	r = util.Connect.Raw(fmt.Sprintf(`select * from audit.starrocks_audit_log where to_date(timestamp)='%s' 
and lower(stmt) like '%%drop%%' 
and lower(stmt) not like '%%load label%%' 
and lower(stmt) not like '%%explain costs%%' 
and lower(stmt) not like '%%drop temporary partition%%' 
and lower(stmt) not like '%%#tableau_%%' 
and lower(stmt) not like '%%-%%' 
and lower(stmt) not like '%%insert%%' 
and lower(stmt) not like '%%show%%' 
and lower(stmt) not like '%%select%%' order by timestamp`,
		util.P.Date)).Scan(&m)
	if r.Error != nil {
		fmt.Println(c.Add(color.FgHiRed).Sprint(r.Error.Error()))
		return
	}

	fmt.Println(fmt.Sprintf("%-2s %-20s %-26s %-63s %-14s %s", "序号", "时间", "用户", "表名", "分区", "过期时间"))

	for i, m2 := range m {
		c := color.New()

		ts := c.Add(color.FgHiYellow).Sprint(m2["timestamp"].(time.Time).Format("2006-01-02 15:04:05"))
		user := c.Add(color.FgHiWhite).Sprint(m2["user"].(string))

		//expire time
		lastTime := m2["timestamp"].(time.Time).Add(time.Duration(expire) * time.Second)
		expireTime := fmt.Sprintf("%s(%s)", lastTime.Format("2006-01-02 15:04:05"), lastTime.Sub(time.Now()).String())

		var msg string
		if ok, v := ddldrop(m2); ok {
			if strings.Contains(v, "^") {
				msg = fmt.Sprintf("%-4d %-31s %-40s %-80s %-34s %s", i, ts, user, c.Add(color.FgHiGreen).Sprint(strings.Split(v, "^")[0]), c.Add(color.FgHiCyan).Sprint(strings.Split(v, "^")[1]), expireTime)
			} else {
				msg = fmt.Sprintf("%-4d %-31s %-40s %-80s %-16s %s", i, ts, user, c.Add(color.FgHiGreen).Sprint(v), " ", expireTime)
			}
		} else {
			msg = fmt.Sprintf("%-4d %-31s %-40s %-80s %-16s %s", i, ts, user, c.Add(color.FgHiRed).Sprint(m2["stmt"].(string)), " ", expireTime)
		}

		fmt.Println(msg)
	}
}

func Restore() {
	c := color.New()

	if len(util.P.Table) != 0 && len(util.P.Partition) != 0 {
		sql := fmt.Sprintf("RECOVER PARTITION %s FROM %s", util.P.Partition, util.P.Table)
		r := util.Connect.Exec(sql)
		if r.Error != nil {
			fmt.Println(c.Add(color.FgHiRed).Sprint("Fail>", sql))
			fmt.Println(c.Add(color.FgHiRed).Sprint(r.Error.Error()))
			return
		}
		fmt.Println(c.Add(color.FgHiGreen).Sprint("Ok  >", sql))
		fmt.Println()

		return
	}

	if len(util.P.Table) != 0 && len(util.P.Partition) == 0 {
		sql := fmt.Sprintf("RECOVER TABLE %s", util.P.Table)
		r := util.Connect.Exec(sql)
		if r.Error != nil {
			fmt.Println(c.Add(color.FgHiRed).Sprint("Fail>", sql))
			fmt.Println(c.Add(color.FgHiRed).Sprint(r.Error.Error()))
			return
		}
		fmt.Println(c.Add(color.FgHiGreen).Sprint("Ok  >", sql))
		fmt.Println()
		return
	}
}
