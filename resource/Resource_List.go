package resource

import (
    "fmt"
    "github.com/fatih/color"
    "gorm.io/gorm"
    "os"
    "StarRocksDevops/util"
    "strconv"
    "strings"
)

func show(db *gorm.DB) {
    c := color.New(color.Bold)
    v, b := version(db)
    var o string
    if !b {
        o = c.Add(color.FgHiRed).Sprint("OFF")
    } else {
        o = c.Add(color.FgHiGreen).Sprint("ON")
    }
    fmt.Println(fmt.Sprintf("%s %s", c.Add(color.FgHiWhite).Sprintf("StarRocks 资源隔离管理"), o))
    fmt.Println()

    if !b {
        fmt.Println(fmt.Sprintf("%s集群 %0.1f暂未开启资源隔离功能，以下展示的配置将不会生效! 设置：/*set global enable_resource_group=true*/\n", strings.ToUpper(util.P.App), v))
    }

    switch v {
    case 3.0, 3.1, 3.2:
        list3(db)
    case 3.3:
        list33(db)
    default:
        fmt.Println("nil")
    }
    fmt.Println()
}

func list33(db *gorm.DB) {
    var m []map[string]interface{}
    row := db.Raw("SHOW RESOURCE GROUPS ALL").Scan(&m)
    if row.Error != nil {
        fmt.Println(row.Error.Error())
        return
    }
    //  "spill_mem_limit_threshold",
    //	"big_query_mem_limit",
    //	"big_query_cpu_second_limit",
    //	"big_query_scan_rows_limit",
    fmt.Println("spill_mem_limit_threshold,big_query_mem_limit,big_query_cpu_second_limit,big_query_scan_rows_limit")
    fmt.Println()
    fmt.Println(fmt.Sprintf("%-3s%-7s%-16s%-13s%-13s%-20s%-25s%-20s%-15s%-20s%-20s%-20s",
        "n",
        "id",
        "name",
        "mem_limit",
        "cpu_weight",
        "concurrency_limit",
        "exclusive_cpu_cores",
        "spill",
        "<b>mem",
        "<b>cpu_second",
        "<b>scan_rows",
        "classifiers",
    ),
    )
    for i, s := range m {
        var data []string
        //序号
        data = append(data, fmt.Sprintf("%-2d", i))
        //id
        v, ok := s["id"]
        if ok {
            data = append(data, fmt.Sprintf("%-6s", v.(string)))
        }
        //name
        v, ok = s["name"]
        if ok {
            data = append(data, fmt.Sprintf("%-15s", v.(string)))
        }
        //mem_limit
        v, ok = s["mem_limit"]
        if ok {
            data = append(data, fmt.Sprintf("%-12s", v.(string)))
        }
        //cpu_weight
        v, ok = s["cpu_weight"]
        if ok {
            data = append(data, fmt.Sprintf("%-12s", v.(string)))
        }
        //concurrency_limit
        v, ok = s["concurrency_limit"]
        if ok {
            data = append(data, fmt.Sprintf("%-19s", v.(string)))
        }
        //exclusive_cpu_cores
        v, ok = s["exclusive_cpu_cores"]
        if ok {
            data = append(data, fmt.Sprintf("%-24s", v.(string)))
        }
        //spill_mem_limit_threshold
        v, ok = s["spill_mem_limit_threshold"]
        if ok {
            data = append(data, fmt.Sprintf("%-19s", v.(string)))
        }
        //big_query_mem_limit
        v, ok = s["big_query_mem_limit"]
        if ok {
            data = append(data, fmt.Sprintf("%-14s", v.(string)))
        }
        //big_query_cpu_second_limit
        v, ok = s["big_query_cpu_second_limit"]
        if ok {
            data = append(data, fmt.Sprintf("%-19s", v.(string)))
        }
        //big_query_scan_rows_limit
        v, ok = s["big_query_scan_rows_limit"]
        if ok {
            data = append(data, fmt.Sprintf("%-19s", v.(string)))
        }
        //classifiers
        v, ok = s["classifiers"]
        if ok {
            data = append(data, fmt.Sprintf("%-19s", v.(string)))
        }

        c := color.New()
        fmt.Println(c.Add(color.FgHiWhite).Sprint(strings.Join(data, " ")))
    }
}

func list3(db *gorm.DB) {
    var m []map[string]interface{}
    row := db.Raw("SHOW RESOURCE GROUPS ALL").Scan(&m)
    if row.Error != nil {
        fmt.Println(row.Error.Error())
        return
    }
    c := color.New()
    msg := `
	num           :  no                            /*序号*/
	name          :  name                          /*资源组的名称*/
	id            :  id                            /*资源组的ID*/
	core          :  cpu_core_limit                /*该资源组在当前 BE 节点可使用的 CPU 核数软上限*/
	mem           :  mem_limit                     /*(大于0生效)该资源组在当前 BE 节点可使用于查询的内存（query_pool）占总内存的百分比（%）*/
	max_cores     :  max_cpu_cores                 /*(大于0生效)资源组在单个 BE 节点中使用的 CPU 核数上限*/
	<b>cpu_second :  big_query_cpu_second_limit    /*(大于0生效)大查询任务在每个 BE 上可以使用 CPU 的时间上限*/
	<b>scan_rows  :  big_query_scan_rows_limit     /*(大于0生效)大查询任务在每个 BE 上可以扫描的行数上限*/
	<b>mem        :  big_query_mem_limit           /*(大于0生效)大查询任务在每个 BE 上可以使用的内存上限*/
	concurrency   :  concurrency_limit             /*(大于0生效)资源组中并发查询数的上限*/
	spill         :  spill_mem_limit_threshold     /*(大于0生效)当前资源组触发落盘的内存占用阈值*/
	type          :  type                          /*资源组的类型 type 支持 short_query 与 normal*/
	classifiers   :  classifiers                   /*类*/
`
    fmt.Println(c.Add(color.FgHiCyan).Sprint(msg))
    fmt.Println()
    fmt.Println(fmt.Sprintf("%-6s%-18s%-13s%-10s%-10s%-15s%-15s%-13s%-15s%-15s%-10s%-10s%-20s",
        "num", "name", "id", "core", "mem", "max_cores", "<b>cpu_second", "<b>scan_rows", "<b>mem", "concurrency", "spill", "type", "classifiers"),
    )

    for i, s := range m {
        msg := fmt.Sprintf("%-6d%-18s%-13s%-10s%-10s%-15s%-15s%-13s%-15s%-15s%-10s%-10s%-20s",
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
        c := color.New()
        fmt.Println(c.Add(color.FgHiBlue).Sprint(msg))
    }
}

/*获取集群版本*/
func version(db *gorm.DB) (float64, bool) {
    type V struct {
        Version string `bson:"version"`
    }
    /*匹配starrocks版本*/
    var v V
    sql := fmt.Sprintf("select current_version() as version")
    db.Raw(sql).Scan(&v)
    version, _ := strconv.ParseFloat(fmt.Sprintf("%s.%s", strings.Split(strings.Split(v.Version, " ")[0], ".")[0], strings.Split(strings.Split(v.Version, " ")[0], ".")[1]), 64)
    fmt.Println(fmt.Sprintf("\n%s version:%s -> %f", strings.ToUpper(util.P.App), strings.Split(v.Version, " ")[0], version))
    /*end*/
    if version <= 2.1 {
        fmt.Println(fmt.Sprintf("%0.1f版本暂不支持资源隔离功能!\n", version))
        os.Exit(-1)
    }

    var m map[string]interface{}
    row := db.Raw("SHOW VARIABLES LIKE 'enable_resource_group'").Scan(&m)
    if row.Error != nil {
        fmt.Println(row.Error.Error())
        return version, false
    }

    if version >= 3.1 {
        return version, true
    }

    b, _ := strconv.ParseBool(m["Value"].(string))
    return version, b
}
