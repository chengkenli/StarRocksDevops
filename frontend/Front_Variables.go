/*
 *@author  chengkenli
 *@project starrocks
 *@package frontend
 *@file    SessionVariables
 *@date    2024/6/11 9:18
 */

package frontend

import (
	"fmt"
	"github.com/fatih/color"
	"starrocks/util"
)

func SessionVariables() {
	fmt.Println()
	c := color.New()
	fmt.Println(c.Add(color.FgHiYellow).Sprint("系统参数:"))
	//第一批次
	fmt.Println(fmt.Sprintf("%-35s %-35s %-35s %-35s %-35s",
		"exec_mem_limit",
		"load_mem_limit",
		"query_mem_limit",
		"query_timeout",
		"new_planner_optimize_timeout",
	))
	fmt.Println(fmt.Sprintf("%-44s %-44s %-44s %-44s %-44s",
		feVars("exec_mem_limit"),
		feVars("load_mem_limit"),
		feVars("query_mem_limit"),
		feVars("query_timeout"),
		feVars("new_planner_optimize_timeout"),
	))
	//第二批次
	fmt.Println(fmt.Sprintf("%-35s %-35s %-35s %-35s %-35s",
		"disable_balance",
		"max_scheduling_tablets",
		"max_balancing_tablets",
		"schedule_slot_num_per_path",
		"enable_insert_strict",
	))
	fmt.Println(fmt.Sprintf("%-44s %-44s %-44s %-44s %-44s",
		feVars("disable_balance"),
		feVars("max_scheduling_tablets"),
		feVars("max_balancing_tablets"),
		feVars("schedule_slot_num_per_path"),
		feVars("enable_insert_strict"),
	))
	//第三批次
	fmt.Println(fmt.Sprintf("%-35s %-35s %-35s %-35s %-35s %-35s",
		"enable_query_queue_load",
		"enable_query_queue_select",
		"enable_query_queue_statistic",
		"enable_group_level_query_queue",
		"query_queue_max_queued_queries",
		"query_queue_pending_timeout_second",
	))
	fmt.Println(fmt.Sprintf("%-44s %-44s %-44s %-44s %-44s %-44s",
		feVars("enable_query_queue_load"),
		feVars("enable_query_queue_select"),
		feVars("enable_query_queue_statistic"),
		feVars("enable_group_level_query_queue"),
		feVars("query_queue_max_queued_queries"),
		feVars("query_queue_pending_timeout_second"),
	))
	//第四批次
	fmt.Println(fmt.Sprintf("%-35s %-35s %-35s",
		"query_queue_concurrency_limit",
		"query_queue_mem_used_pct_limit",
		"query_queue_cpu_used_permille_limit",
	))
	fmt.Println(fmt.Sprintf("%-44s %-44s %-44s",
		feVars("query_queue_concurrency_limit"),
		feVars("query_queue_mem_used_pct_limit"),
		feVars("query_queue_cpu_used_permille_limit"),
	))

	//第五批次
	fmt.Println(fmt.Sprintf("%-35s",
		"catalog_trash_expire_second",
	))
	fmt.Println(fmt.Sprintf("%-44s",
		feVars("catalog_trash_expire_second"),
	))

}

func feVars(Key string) string {
	var m map[string]interface{}
	r := util.Connect.Raw(fmt.Sprintf("show variables like '%s'", Key)).Scan(&m)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return r.Error.Error()
	}
	c := color.New()
	if m == nil {
		var m map[string]interface{}
		r := util.Connect.Raw(fmt.Sprintf("ADMIN SHOW FRONTEND CONFIG LIKE '%s'", Key)).Scan(&m)
		if r.Error != nil {
			fmt.Println(r.Error.Error())
			return r.Error.Error()
		}
		if m == nil {
			return c.Add(color.FgHiRed).Sprint("nil")
		}
		return c.Add(color.FgHiBlue).Sprint(m["Value"].(string))
	}
	return c.Add(color.FgHiGreen).Sprint(m["Value"].(string))
}
