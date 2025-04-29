/*
 *@author  chengkenli
 *@project starrocks
 *@package backend
 *@file    BackendMetrics
 *@date    2024/6/10 9:48
 */

package backend

import (
	"fmt"
	"github.com/fatih/color"
	"regexp"
	"sort"
	"starrocks/tools"
	"starrocks/util"
	"strconv"
	"strings"
	"sync"
)

func Metrics() {
	var back util.Backends
	r := util.Connect.Raw("show backends").Scan(&back)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return
	}

	beMem := []string{
		"starrocks_be_chunk_allocator_mem_bytes",
		"starrocks_be_column_pool_mem_bytes",
		"starrocks_be_compaction_mem_bytes",
		"starrocks_be_load_mem_bytes",
		"starrocks_be_process_mem_bytes",
		"starrocks_be_query_mem_bytes",
		"starrocks_be_storage_page_cache_mem_bytes",
		"starrocks_be_tablet_meta_mem_bytes",
		"starrocks_be_tcmalloc_bytes_in_use",
		"starrocks_be_update_mem_bytes",
	}

	c := color.New()
	fmt.Println(c.Add(color.FgHiGreen).Sprint("mem_limit:"))
	fmt.Println(fmt.Sprintf("%-14s%-7s%-7s%-7s%-7s%-7s%-7s%-7s%-7s%-11s%-11s",
		"BE",
		"CPU缓存",
		"字段线程",
		"合并内存",
		"导入内存",
		"进程内存",
		"查询内存",
		"页面缓存",
		"元数据",
		"TcMalloc",
		"主键模型，内存消耗",
	))

	var schema []string
	var wg sync.WaitGroup
	for _, b := range back {
		wg.Add(1)
		b := b
		go func() {
			defer wg.Done()

			mm := make(map[string]interface{}, 300)

			//metrics接口
			var model int
			r := tools.Post("GET", fmt.Sprintf("http://%s:8040/metrics", b.IP), nil)
			for _, s := range strings.Split(string(r), "\n") {
				if strings.Contains(s, "#") {
					continue
				}
				v := strings.Split(s, " ")
				for _, s2 := range beMem {
					if v[0] != s2 {
						continue
					}
					f, _ := strconv.ParseFloat(v[1], 64)
					if f <= 0 {
						mm[v[0]] = v[1]
					} else {
						if f > 64424509440 && f < 107374182400 {
							model = 7
						}
						if f > 107374182400 && f < 214748364800 {
							model = 9
						}
						if f > 214748364800 {
							model = 8
						}
						mm[v[0]] = fmt.Sprintf("%0.2fGB", f/1024/1024/1024)
					}
				}
			}
			//memz接口
			var matches []string
			i := false
			o := tools.Post("GET", fmt.Sprintf("http://%s:8040/memz", b.IP), nil)
			for _, s := range strings.Split(string(o), "\n") {
				if strings.Contains(s, "primary index stats") {
					re := regexp.MustCompile(`total:\d+\s*memory:\d+`)
					matches = re.FindStringSubmatch(s)
				}
				if strings.Contains(s, "tabletid") {
					i = true
				}

				if i {
					regex, err := regexp.Compile(`[1-9]\d{5,}`)
					if err != nil {
						continue
					}
					matches := regex.FindAllString(s, -1)
					if matches == nil {
						continue
					}
					var m map[string]interface{}
					r := util.Connect.Raw("show tablet " + matches[0]).Scan(&m)
					if r.Error != nil {
						fmt.Println(r.Error.Error())
						continue
					}
					if m == nil || m["TableName"] == nil {
						continue
					}
					table := fmt.Sprintf("%s.%s", strings.NewReplacer("default_cluster:", "", " ", "").Replace(m["DbName"].(string)), m["TableName"].(string))
					schema = append(schema, table)
				}
			}

			switch model {
			case 7:
				c := color.New()
				msg := c.Add(color.FgHiBlue).Sprint(fmt.Sprintf("%-14v%-9v%-11v%-11v%-11v%-11v%-11v%-11v%-10v%-11v%-8v%v",
					b.IP,
					mm["starrocks_be_chunk_allocator_mem_bytes"],
					mm["starrocks_be_column_pool_mem_bytes"],
					mm["starrocks_be_compaction_mem_bytes"],
					mm["starrocks_be_load_mem_bytes"],
					mm["starrocks_be_process_mem_bytes"],
					mm["starrocks_be_query_mem_bytes"],
					mm["starrocks_be_storage_page_cache_mem_bytes"],
					mm["starrocks_be_tablet_meta_mem_bytes"],
					mm["starrocks_be_tcmalloc_bytes_in_use"],
					mm["starrocks_be_update_mem_bytes"],
					matches,
				))
				fmt.Println(msg)
			case 8:
				c := color.New()
				msg := c.Add(color.FgHiRed).Sprint(fmt.Sprintf("%-14v%-9v%-11v%-11v%-11v%-11v%-11v%-11v%-10v%-11v%-8v%v",
					b.IP,
					mm["starrocks_be_chunk_allocator_mem_bytes"],
					mm["starrocks_be_column_pool_mem_bytes"],
					mm["starrocks_be_compaction_mem_bytes"],
					mm["starrocks_be_load_mem_bytes"],
					mm["starrocks_be_process_mem_bytes"],
					mm["starrocks_be_query_mem_bytes"],
					mm["starrocks_be_storage_page_cache_mem_bytes"],
					mm["starrocks_be_tablet_meta_mem_bytes"],
					mm["starrocks_be_tcmalloc_bytes_in_use"],
					mm["starrocks_be_update_mem_bytes"],
					matches,
				))
				fmt.Println(msg)
			case 9:
				c := color.New()
				msg := c.Add(color.FgHiYellow).Sprint(fmt.Sprintf("%-14v%-9v%-11v%-11v%-11v%-11v%-11v%-11v%-10v%-11v%-8v%v",
					b.IP,
					mm["starrocks_be_chunk_allocator_mem_bytes"],
					mm["starrocks_be_column_pool_mem_bytes"],
					mm["starrocks_be_compaction_mem_bytes"],
					mm["starrocks_be_load_mem_bytes"],
					mm["starrocks_be_process_mem_bytes"],
					mm["starrocks_be_query_mem_bytes"],
					mm["starrocks_be_storage_page_cache_mem_bytes"],
					mm["starrocks_be_tablet_meta_mem_bytes"],
					mm["starrocks_be_tcmalloc_bytes_in_use"],
					mm["starrocks_be_update_mem_bytes"],
					matches,
				))
				fmt.Println(msg)
			default:
				c := color.New()
				msg := c.Add(color.FgHiWhite).Sprint(fmt.Sprintf("%-14v%-9v%-11v%-11v%-11v%-11v%-11v%-11v%-10v%-11v%-8v%v",
					b.IP,
					mm["starrocks_be_chunk_allocator_mem_bytes"],
					mm["starrocks_be_column_pool_mem_bytes"],
					mm["starrocks_be_compaction_mem_bytes"],
					mm["starrocks_be_load_mem_bytes"],
					mm["starrocks_be_process_mem_bytes"],
					mm["starrocks_be_query_mem_bytes"],
					mm["starrocks_be_storage_page_cache_mem_bytes"],
					mm["starrocks_be_tablet_meta_mem_bytes"],
					mm["starrocks_be_tcmalloc_bytes_in_use"],
					mm["starrocks_be_update_mem_bytes"],
					matches,
				))
				fmt.Println(msg)
			}
		}()
	}
	wg.Wait()

	fmt.Println()
	fmt.Println(c.Add(color.FgHiGreen).Sprint("primary key:"))
	fmt.Println(fmt.Sprintf("%-75s %s", "表名", "tablet"))

	_, count := tools.DuplicatesAndCount(schema)
	// 创建一个用于存储键值对的切片
	var pairs []struct {
		Key   string
		Value int
	}
	// 将map中的所有键值对添加到切片中
	for k, v := range count {
		pairs = append(pairs, struct {
			Key   string
			Value int
		}{k, v})
	}
	// 对pairs切片进行降序排序，通过为sort.Slice指定比较函数
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value // 按值降序排序
	})
	// 打印排序后的键值对
	for _, pair := range pairs {
		fmt.Println(fmt.Sprintf("%-77s %d", pair.Key, pair.Value))
	}

	fmt.Println()
}
