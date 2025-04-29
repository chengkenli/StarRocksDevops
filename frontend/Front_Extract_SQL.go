/*
 *@author  chengkenli
 *@project starrocks
 *@package frontend
 *@file    SessionSQL
 *@date    2024/7/22 16:24
 */

package frontend

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"regexp"
	"starrocks/util"
	"strconv"
	"strings"
	"sync"
	"time"
)

func SessionExtract() {
	data, err := SessionExtractSQL(util.P.File)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c := color.New()
	fmt.Println()
	fmt.Println(c.Add(color.FgHiYellow).Sprint("分析SQL，提取表名:"))
	fmt.Println()

	for _, d := range data {
		fmt.Println(fmt.Sprintf("%-5s %-4s", "", d))
	}
	fmt.Println()

	fmt.Println(fmt.Sprintf("%-5s %-4s %-19s %-3s %-10s %-15s %-10s %-15s %-70s",
		"",
		"id",
		"create time",
		"",
		"bsize",
		"size",
		"replica",
		"rowcount",
		"schema",
	))

	var wg sync.WaitGroup
	for i, i2 := range data {
		wg.Add(1)
		go func(i int, i2 string) {
			defer wg.Done()
			schema := strings.Split(i2, ".")
			//提取创建日期
			var m map[string]interface{}
			sql := fmt.Sprintf("select `TABLE_CATALOG`,`TABLE_SCHEMA`,`TABLE_NAME`,`TABLE_TYPE`,`ENGINE`,`CREATE_TIME`,`TABLE_COMMENT`,`TABLE_ROWS` from information_schema.tables where TABLE_SCHEMA='%s' and TABLE_NAME='%s'", schema[0], schema[1])
			r := util.Connect.Raw(sql).Scan(&m)
			if r.Error != nil {
				return
			}
			if m == nil {
				return
			}
			var maxDataSize float64
			var mbucket, msize, mreplicaCount, mrowCount string
			if m["TABLE_TYPE"].(string) == "BASE TABLE" {
				//提取表总容量和行数
				var cm []map[string]interface{}
				r = util.Connect.Raw(fmt.Sprintf("show data from %s.%s", schema[0], schema[1])).Scan(&cm)
				if r.Error != nil {
					return
				}
				msize = cm[0]["Size"].(string)
				mreplicaCount = cm[0]["ReplicaCount"].(string)
				mrowCount = cm[0]["RowCount"].(string)
				//提取分桶
				var mb map[string]interface{}
				r = util.Connect.Raw(fmt.Sprintf("show partitions from %s.%s order by LastConsistencyCheckTime,DataSize desc limit 1", schema[0], schema[1])).Scan(&mb)
				if r.Error != nil {
					return
				}
				mbucket = mb["Buckets"].(string)
				//提取最大的分区容量
				var n []map[string]interface{}
				r = util.Connect.Raw(fmt.Sprintf("show partitions from %s.%s", schema[0], schema[1])).Scan(&n)
				if r.Error != nil {
					return
				}
				var Max []float64
				for _, m := range n {
					Max = append(Max, size(m["DataSize"].(string)))
				}
				maxDataSize = Max[0]
				for i := 0; i < len(Max); i++ {
					if Max[i] >= maxDataSize {
						maxDataSize = Max[i]
					}
				}

			}

			msg := fmt.Sprintf("%-5s %04d %-10s %-3s %-10s %-15s %-10s %-15s %-70s",
				"",
				i,
				m["CREATE_TIME"].(time.Time).Format("2006-01-02 15:04:05"),
				mbucket,
				fmt.Sprintf("%0.2fGB", maxDataSize/1024/1024/1024),
				msize,
				mreplicaCount,
				mrowCount,
				i2,
			)
			if m["CREATE_TIME"].(time.Time).Format("2006-01-02") == time.Now().Format("2006-01-02") {
				c := color.New()
				fmt.Println(c.Add(color.FgHiGreen).Sprint(msg))
			} else {
				fmt.Println(msg)
			}

		}(i, i2)
	}
	wg.Wait()
	fmt.Println()
}

func SessionExtractSQL(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("打开文件失败：", err)
		return nil, err
	}
	// 使用defer语句确保文件在函数返回前被关闭
	defer file.Close()
	// 读取文件内容
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	input := string(data)

	database := []string{
		"adhoc", "ads", "ads_dev", "ads_dev_secure", "ads_rt", "ads_rt_dev", "ads_rt_dev_secure", "ads_rt_secure", "ads_secure",
		"algo", "algo_dev", "ap_secure", "audit", "bi_item", "bi_realty", "bi_realty_secure", "bi_sams_secure", "bi_scm", "bi_sc_secure",
		"cdp", "cdp_api", "cloud_fcst_dm", "cn_backup_secure", "cn_chilled_data", "cn_core_dim_vm", "cn_di_data", "cn_ec_bi_secure",
		"cn_ec_wmdj_user_action", "cn_mdse_dm_dl_tables", "cn_po_home_system", "cn_po_home_system_dev", "cn_pricing_dl_tables",
		"cn_sams_dl_secure", "cn_wc_highsecure", "cn_wc_mb_secure", "cn_wc_mb_vm", "cn_wc_repl_vm", "cn_wc_vm", "cn_wid_dl_secure",
		"cn_wm_mb_secure", "cn_wm_mb_vm", "cn_wm_repl_vm", "cn_wm_vm", "data_test", "demo", "dim", "dim_dev", "dim_dev_secure", "dim_rt",
		"dim_rt_dev", "dim_rt_dev_secure", "dim_rt_secure", "dim_secure", "dm", "dm_dev", "dm_dev_secure", "dm_secure", "dw", "dwd", "dwd_dev",
		"dwd_dev_secure", "dw_dev", "dw_dev_secure", "dwd_rt", "dwd_rt_dev", "dwd_rt_dev_secure", "dwd_rt_secure", "dwd_secure", "dw_rt",
		"dws", "dws_dev", "dws_dev_secure", "dw_secure", "dws_rt", "dws_rt_dev", "dws_rt_dev_secure", "dws_rt_secure", "dws_secure",
		"euclid_scn_forecast_prod", "finance_kettle", "fin_sox", "fin_sox_dev", "flash_report", "flash_report_dev", "flash_report_sit",
		"hyper_bi_secure", "hyper_ec_secure", "hyper_mdse_dm_secure", "information_schema", "ma_test", "mbrship_secure", "mcfc_report",
		"mcfc_report_dev", "o2o_datacubes_secure", "ods", "ods_app_dev", "ods_app_dev_secure", "ods_app_test", "ods_app_test_secure",
		"ods_archive", "ods_dev", "ods_dev_secure", "ods_gray", "ods_migration_td_gray", "ods_rt", "ods_rt_dev", "ods_rt_dev_secure",
		"ods_rt_secure", "ods_secure", "ods_secure_rt", "ods_sox", "ods_sox_app_dev", "ods_sox_app_test", "ods_sox_dev", "ods_sox_test",
		"ods_test", "ods_test_secure", "ops", "pro_dgtmkt_data", "pro_scct_dev", "sams_finance", "scct_inv_monitor", "scct_logis", "scct_logis_dev",
		"scm_dcqe_secure", "scm_network_secure", "scm_secure", "scm_uihealth", "starrocks_monitor", "_statistics_", "supply_kettle", "svccn_logis",
		"svccn_logis_query", "svcdordgtmkt", "sys", "wm_ad_hoc", "wm_cn_util", "wm_common_vm", "ww_core_dim_vm"}

	if strings.Contains(input, "hadoop") && !strings.Contains(strings.ToLower(input), "outfile") {
		return nil, nil
	}
	var schema []string
	/*schema.table*/
	re := regexp.MustCompile(`([a-zA-Z][^\s=,'.]+)\.([^\s=,'.]+)`)
	var result []string
	for _, s := range re.FindAllString(input, -1) {
		b := regexp.MustCompile(`[\\/\(\),:|+><~!@#%^&*='";?-]`).FindString(s) != ""
		if !b {
			data := strings.Split(s, ".")
			for _, s2 := range database {
				if data[0] == s2 {
					result = append(result, s)
				}
			}
		}
	}
	/*catalog.schema.table*/
	re2 := regexp.MustCompile(`([a-zA-Z][^\s=,'.]+)\.([^\s=,'.]+)\.([^\s=,'.]+)`)
	var result2 []string
	for _, s := range re2.FindAllString(input, -1) {
		b := regexp.MustCompile(`[\\/\(\),:|+><~!@#%^&*='";?-]`).FindString(s) != ""
		if !b {
			data := strings.Split(s, ".")
			for _, s2 := range database {
				if data[1] == s2 {
					result2 = append(result2, s)
				}
			}
		}
	}

	schema = append(schema, result...)
	schema = append(schema, result2...)

	return removeDuplicateStrings(schema), nil
}

// RemoveDuplicateStrings /*数组去重*/
func removeDuplicateStrings(strs []string) []string {
	result := []string{}
	tempMap := map[string]byte{} // 存放不重复字符串
	for _, e := range strs {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

func size(s string) float64 {
	s = strings.ToLower(s)
	//正则
	re := regexp.MustCompile(`\d+\.?\d*`)

	float, _ := strconv.ParseFloat(re.FindString(s), 64)

	if strings.Contains(s, "kb") {
		return float * 1024
	}
	if strings.Contains(s, "mb") {
		return float * 1024 * 1024
	}
	if strings.Contains(s, "gb") {
		return float * 1024 * 1024 * 1024
	}
	if strings.Contains(s, "tb") {
		return float * 1024 * 1024 * 1024 * 1024
	}
	return 0
}
