package backend

import (
	"StarRocksDevops/tools"
	"StarRocksDevops/util"
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func BackMemLimit() {
	var b util.Backends
	r := util.Connect.Raw("show backends").Scan(&b)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return
	}

	beMem := []string{
		"starrocks_be_chunk_allocator_mem_bytes",
		"starrocks_be_clone_mem_bytes",
		"starrocks_be_column_pool_mem_bytes",
		"starrocks_be_compaction_mem_bytes",
		"starrocks_be_consistency_mem_bytes",
		"starrocks_be_load_mem_bytes",
		"starrocks_be_process_mem_bytes",
		"starrocks_be_query_mem_bytes",
		"starrocks_be_schema_change_mem_bytes",
		"starrocks_be_storage_page_cache_mem_bytes",
		"starrocks_be_tablet_meta_mem_bytes",
		"starrocks_be_tcmalloc_bytes_in_use",
		"starrocks_be_update_mem_bytes",
	}

	c := color.New()
	fmt.Println(c.Add(color.FgHiGreen).Sprint("mem_limit:"))
	fmt.Println(fmt.Sprintf("%-15s%-10s%-10s%-10s%-10s%-10s%-10s%-10s%-10s%-10s%-10s%-10s%-18s%-13s%-20s",
		"BE",
		"CPU缓存",
		"克隆副本",
		"字段线程",
		"合并内存",
		"CheckSum",
		"导入内存",
		"进程内存",
		"查询内存",
		"结构变更",
		"页面缓存",
		"元数据",
		"TcMalloc统计BE消耗",
		"主键模型",
		"总消耗",
	))

	var wg sync.WaitGroup
	for _, s := range b {
		wg.Add(1)
		s := s
		go func() {
			defer wg.Done()

			r := tools.Post("GET", fmt.Sprintf("http://%s:8040/metrics", s.IP), nil)
			filename := fmt.Sprintf("%s.%d.log", s.IP, time.Now().Unix())
			tools.Writefile(filename, string(r))

			file, err := os.Open(filename)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer file.Close()

			mn := make(map[string]string)
			scanner := bufio.NewScanner(file)

			var model int
			for scanner.Scan() {
				for _, s2 := range beMem {
					if strings.Contains(scanner.Text(), s2) && !strings.Contains(scanner.Text(), "#") {
						k := strings.Split(scanner.Text(), " ")[0]
						//v := strings.Split(scanner.Text(), " ")[1]
						v, _ := strconv.ParseFloat(strings.Split(scanner.Text(), " ")[1], 64)
						if v <= 0 {
							mn[k] = fmt.Sprintf("%0.0f", v)
						} else {
							if v > 64424509440 && v < 107374182400 {
								model = 7
							}
							if v > 107374182400 && v < 214748364800 {
								model = 9
							}
							if v > 214748364800 {
								model = 8
							}
							//mn[k] = fmt.Sprintf("%0.0f (%0.3fGB)", v, v/1024/1024/1024)
							mn[k] = fmt.Sprintf("%0.3f GB", v/1024/1024/1024)
						}
						if s2 == "starrocks_be_tcmalloc_bytes_in_use" {
							continue
						}
					}
				}
			}
			if model == 7 {
				c := color.New()
				msg := c.Add(color.FgHiBlue).Sprint(fmt.Sprintf("%-15s%-12s%-14s%-14s%-14s%-10s%-14s%-14s%-14s%-14s%-14s%-13s%-22s%-17s%s",
					s.IP,
					mn["starrocks_be_chunk_allocator_mem_bytes"],
					mn["starrocks_be_clone_mem_bytes"],
					mn["starrocks_be_column_pool_mem_bytes"],
					mn["starrocks_be_compaction_mem_bytes"],
					mn["starrocks_be_consistency_mem_bytes"],
					mn["starrocks_be_load_mem_bytes"],
					mn["starrocks_be_process_mem_bytes"],
					mn["starrocks_be_query_mem_bytes"],
					mn["starrocks_be_schema_change_mem_bytes"],
					mn["starrocks_be_storage_page_cache_mem_bytes"],
					mn["starrocks_be_tablet_meta_mem_bytes"],
					mn["starrocks_be_tcmalloc_bytes_in_use"],
					mn["starrocks_be_update_mem_bytes"],
					mn["starrocks_be_process_mem_bytes"],
				))
				fmt.Println(msg)
				return
			}
			if model == 9 {
				c := color.New()
				msg := c.Add(color.FgHiYellow).Sprint(fmt.Sprintf("%-15s%-12s%-14s%-14s%-14s%-10s%-14s%-14s%-14s%-14s%-14s%-13s%-22s%-17s%s",
					s.IP,
					mn["starrocks_be_chunk_allocator_mem_bytes"],
					mn["starrocks_be_clone_mem_bytes"],
					mn["starrocks_be_column_pool_mem_bytes"],
					mn["starrocks_be_compaction_mem_bytes"],
					mn["starrocks_be_consistency_mem_bytes"],
					mn["starrocks_be_load_mem_bytes"],
					mn["starrocks_be_process_mem_bytes"],
					mn["starrocks_be_query_mem_bytes"],
					mn["starrocks_be_schema_change_mem_bytes"],
					mn["starrocks_be_storage_page_cache_mem_bytes"],
					mn["starrocks_be_tablet_meta_mem_bytes"],
					mn["starrocks_be_tcmalloc_bytes_in_use"],
					mn["starrocks_be_update_mem_bytes"],
					mn["starrocks_be_process_mem_bytes"],
				))
				fmt.Println(msg)
				return
			}
			if model == 8 {
				c := color.New()
				msg := c.Add(color.FgHiRed).Sprint(fmt.Sprintf("%-15s%-12s%-14s%-14s%-14s%-10s%-14s%-14s%-14s%-14s%-14s%-13s%-22s%-17s%s",
					s.IP,
					mn["starrocks_be_chunk_allocator_mem_bytes"],
					mn["starrocks_be_clone_mem_bytes"],
					mn["starrocks_be_column_pool_mem_bytes"],
					mn["starrocks_be_compaction_mem_bytes"],
					mn["starrocks_be_consistency_mem_bytes"],
					mn["starrocks_be_load_mem_bytes"],
					mn["starrocks_be_process_mem_bytes"],
					mn["starrocks_be_query_mem_bytes"],
					mn["starrocks_be_schema_change_mem_bytes"],
					mn["starrocks_be_storage_page_cache_mem_bytes"],
					mn["starrocks_be_tablet_meta_mem_bytes"],
					mn["starrocks_be_tcmalloc_bytes_in_use"],
					mn["starrocks_be_update_mem_bytes"],
					mn["starrocks_be_process_mem_bytes"],
				))
				fmt.Println(msg)
				return
			}
			fmt.Println(fmt.Sprintf("%-15s%-12s%-14s%-14s%-14s%-10s%-14s%-14s%-14s%-14s%-14s%-13s%-22s%-17s%s",
				s.IP,
				mn["starrocks_be_chunk_allocator_mem_bytes"],
				mn["starrocks_be_clone_mem_bytes"],
				mn["starrocks_be_column_pool_mem_bytes"],
				mn["starrocks_be_compaction_mem_bytes"],
				mn["starrocks_be_consistency_mem_bytes"],
				mn["starrocks_be_load_mem_bytes"],
				mn["starrocks_be_process_mem_bytes"],
				mn["starrocks_be_query_mem_bytes"],
				mn["starrocks_be_schema_change_mem_bytes"],
				mn["starrocks_be_storage_page_cache_mem_bytes"],
				mn["starrocks_be_tablet_meta_mem_bytes"],
				mn["starrocks_be_tcmalloc_bytes_in_use"],
				mn["starrocks_be_update_mem_bytes"],
				mn["starrocks_be_process_mem_bytes"],
			))
		}()
	}
	wg.Wait()
	fmt.Println()
}
