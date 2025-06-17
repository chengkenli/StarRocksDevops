package util

import (
    "encoding/json"
    "flag"
    "fmt"
    "github.com/fatih/color"
    "github.com/spf13/viper"
    "os"
    "path/filepath"
    "time"
)

func usage() {
    fmt.Printf("\nUsage: %s [-s starrocks] [-h]\n%s tool version: [1.2]\n\n", filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
    flag.PrintDefaults()
    fmt.Println()
}

func init() {
    var (
        conf        string
        defaultConf string
    )

    execDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
    defaultConf = fmt.Sprintf("%s/.%s.yaml", execDir, filepath.Base(os.Args[0]))

    flag.StringVar(&conf, "c", defaultConf, "conf file")

    c := color.New()
    flag.BoolVar(&P.Clear, "c", false, fmt.Sprintf("%s: %s", c.Add(color.FgHiRed).Sprint("慎用"), c.Add(color.FgHiYellow).Sprint("清场")))
    flag.BoolVar(&P.List, "l", false, "展示集群fe，be节点")
    flag.StringVar(&P.Username, "u", "", "查看指定用户连接数")
    flag.StringVar(&P.Primary, "key", "", "筛选出集群包含'关键字'的语句")
    flag.StringVar(&P.ALterName, "n", "", "填写查看add column的表名 该参数与-b相配合kill")
    flag.BoolVar(&P.Kill, "k", false, "是否执行杀死PID")
    flag.IntVar(&P.Pid, "p", -1, "PID 语句")
    flag.StringVar(&P.App, "s", "", "集群名称: adhoc、app、sr-adhoc、sr-app、cdp...")
    flag.BoolVar(&P.Help, "h", false, "帮助")
    flag.BoolVar(&P.RoutineLoad, "routine", false, "展示集群所有routine load任务信息")
    flag.BoolVar(&P.Errors, "e", false, c.Add(color.FgHiBlue).Sprint("展示集群报错的所有Session"))
    flag.BoolVar(&P.License, "license", false, c.Add(color.FgHiBlue).Sprint("查看集群license"))
    flag.IntVar(&P.ErrorsTime, "etime", 1, c.Add(color.FgHiBlue).Sprint("展示集群多少小时内报错的所有Session"))
    flag.IntVar(&P.JobID, "jobid", -1, "routine load/broker load --JobId")
    flag.BoolVar(&P.TabletIsBad, "isbad", false, "临时使用")
    flag.BoolVar(&P.RepairSchemaTable, "rs", false, c.Add(color.FgHiYellow).Sprint("启动副本探测模式(修复异常)"))
    flag.IntVar(&P.Thread, "tt", 10, c.Add(color.FgHiYellow).Sprint("副本探测模式: [Thread]"))
    flag.IntVar(&P.BackendId, "b", 0, c.Add(color.FgHiYellow).Sprint("副本探测模式: <BackendId>"))
    flag.StringVar(&P.JobName, "jobname", "", "routine load --JobName")
    flag.StringVar(&P.OlapType, "o", "", "Type: OLAP,MYSQL,OLAP_EXTERNAL,VIEW,MATERIALIZED_VIEW")
    flag.StringVar(&P.Dbname, "db", "", "副本探测模式: 库名")
    flag.StringVar(&P.Status, "status", "", "打印相关状态的作业信息: [PENDING|QUEUEING|LOADING|PREPARED|FINISHED|CANCELLED]")
    flag.StringVar(&P.File, "f", "", "输入sql文件，分析SQL，提取表名")
    flag.StringVar(&P.QueryId, "id", "", "根据QueryID查询当天队列的查询语句在BE节点上的资源消耗")
    flag.BoolVar(&P.Front, "front", false, c.Add(color.FgHiYellow).Sprint("仅查看session query."))
    flag.BoolVar(&P.QuerisOne, "q", false, c.Add(color.FgHiYellow).Sprint("仅查看队列"))
    flag.StringVar(&P.Proc, "proc", "", "fe/be进程")
    flag.BoolVar(&P.Alter, "alter", false, c.Add(color.FgHiYellow).Sprint("仅查看修改字段操作"))

    flag.BoolVar(&P.ErrScans, "es", false, c.Add(color.FgHiYellow).Sprint("搜索审计日志"))
    flag.StringVar(&P.StartTime, "st", time.Now().Add(-10*time.Minute).Format("2006-01-02 15:04:05"), "搜索审计日志：开始-日期时间")
    flag.StringVar(&P.EndTime, "et", time.Now().Format("2006-01-02 15:04:05"), "搜索审计日志：结束-日期时间")
    flag.StringVar(&P.CommitType, "ct", "", "搜索审计日志：提交类型[insert/select]")
    flag.IntVar(&P.CommitLimit, "cl", 20, "搜索审计日志：显示行数")
    flag.StringVar(&P.ErrKey, "ek", "scanRows", "搜索审计日志：排序关键字")

    flag.BoolVar(&P.Io, "io", false, c.Add(color.FgHiYellow).Sprint("查看be io/cpu/mem"))
    flag.StringVar(&P.IP, "ip", "", "Backends IP")
    flag.IntVar(&P.Seconds, "seconds", 10, "采集时间")

    flag.StringVar(&P.Date, "d", time.Now().Format("2006-01-02"), "解析【某天】审计日志drop行为")
    flag.BoolVar(&P.Recover, "v", false, c.Add(color.FgHiYellow).Sprint("解析审计日志drop行为"))
    flag.BoolVar(&P.Restore, "r", false, c.Add(color.FgHiYellow).Sprint("是否从回收站中恢复某个表/分区"))
    flag.StringVar(&P.Table, "t", "", "从回收站中恢复某个表,表名")
    flag.StringVar(&P.Partition, "pi", "", "从回收站中恢复某个分区,分区名")

    flag.BoolVar(&P.Profilelist, "profile", false, "查看profile列表")

    flag.BoolVar(&P.Stream, "stream", false, "查看streamload列表")
    flag.StringVar(&P.Label, "label", "", "label name")

    flag.BoolVar(&P.SubmitTask, "submit", false, "查看submit task作业列表")

    flag.Parse()
    flag.Usage = usage

    if P.Help || len(P.App) == 0 {
        flag.Usage()
        os.Exit(-1)
    }

    paths, name := filepath.Split(conf)
    Config := viper.New()
    Config.SetConfigFile(fmt.Sprintf("%s%s", paths, name))
    if err := Config.ReadInConfig(); err != nil {
        fmt.Println(err.Error())
    }
    // 读取配置文件所有内容放入内存
    marshal, err := json.Marshal(Config.AllSettings())
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(-1)
    }
    err = json.Unmarshal(marshal, &Read)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(-1)
    }
    // auth metadb
    Instglobal.vaild <- struct{}{}
    // send begin channel
    Instglobal.Begin <- struct{}{}

}

func Parms() {
}
