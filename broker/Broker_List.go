package broker

import (
	"fmt"
	"github.com/fatih/color"
	"starrocks/util"
	"sync"
	"time"
)

var pool *sync.Pool

func init() {
	pool = &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating a new Pool ...")
			return new(job)
		},
	}
}

type job struct {
	JobId          int
	Label          string
	State          string
	Progress       string
	Type           string
	EtlInfo        string
	TaskInfo       string
	ErrorMsg       string
	CreateTime     string
	EtlStartTime   string
	EtlFinishTime  string
	LoadStartTime  string
	LoadFinishTime string
	URL            string
	JobDetails     string
	Database       string
}

func BrokerList() {
	var (
		PENDING   []string
		ETL       []string
		LOADING   []string
		FINISHED  []string
		CANCELLED []string
		QUEUEING  []string
	)

	var databases []map[string]interface{}
	util.Connect.Raw("show databases").Scan(&databases)

	fmt.Println()
	c := color.New()

	fmt.Println(c.Add(color.FgHiYellow).Sprint("BROKER LOAD:"))
	fmt.Println(fmt.Sprintf("%-9s %-30s %-50s %-10s %-20s %-10s %-20s %-20s", "JobID", "DB", "Label", "State", "Progress", "Type", "CreateTime", "Edtime"))

	var wgs sync.WaitGroup
	for _, database := range databases {
		wgs.Add(1)
		go func(database map[string]interface{}) {
			defer wgs.Done()

			if database["Database"] == "_statistics_" {
				return
			}
			var brokers util.BrokerLoad
			r := util.Connect.Raw(fmt.Sprintf("SHOW LOAD FROM %s", database["Database"])).Scan(&brokers)
			if r.Error != nil {
				return
			}
			if brokers == nil {
				return
			}
			var wg sync.WaitGroup
			for _, b := range brokers {
				wg.Add(1)
				b := b
				go func() {
					defer wg.Done()

					pool.Put(&job{
						JobId:          b.JobId,
						Label:          b.Label,
						State:          b.State,
						Progress:       b.Progress,
						Type:           b.Type,
						EtlInfo:        b.EtlInfo,
						TaskInfo:       b.TaskInfo,
						ErrorMsg:       b.ErrorMsg,
						CreateTime:     b.CreateTime,
						EtlStartTime:   b.EtlStartTime,
						EtlFinishTime:  b.EtlFinishTime,
						LoadStartTime:  b.LoadStartTime,
						LoadFinishTime: b.LoadFinishTime,
						URL:            b.URL,
						JobDetails:     b.JobDetails,
						Database:       database["Database"].(string),
					})

					p := pool.Get().(*job)

					c := color.New()
					var state string
					if p.State == "PENDING" {
						PENDING = append(PENDING, p.Label)
						state = c.Add(color.FgHiBlue).Sprint(p.State)
					}
					if p.State == "ETL" {
						ETL = append(ETL, p.Label)
						state = c.Add(color.FgHiYellow).Sprint(p.State)
					}
					if p.State == "FINISHED" {
						FINISHED = append(FINISHED, p.Label)
						return
					}
					if p.State == "CANCELLED" {
						CANCELLED = append(CANCELLED, p.Label)
						return
					}
					if p.State == "QUEUEING" {
						QUEUEING = append(QUEUEING, p.Label)
						state = c.Add(color.FgHiMagenta).Sprint(p.State)
					}
					if p.State == "LOADING" {
						LOADING = append(LOADING, p.Label)
						state = c.Add(color.FgHiGreen).Sprint(p.State)
					}

					// 解析日期时间字符串
					time1, err := time.Parse("2006-01-02 15:04:05", p.CreateTime)
					if err != nil {
						fmt.Println("Error parsing time1:", err)
						fmt.Println(p)
						return
					}
					switch util.P.App {
					case "cdp", "api", "ma":
						time1 = time1.Add(8 * time.Hour)
					default:

					}
					time2, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
					if err != nil {
						fmt.Println("Error parsing time2:", err)
						fmt.Println(p)
						return
					}
					// 计算两个时间之间的差异
					edtime := time2.Sub(time1)

					broker := fmt.Sprintf("%-9s %-30s %-50s %-19s %-20s %-10s %-20s %-20s",
						c.Add(color.FgHiGreen).Sprint(p.JobId),
						p.Database,
						p.Label,
						state,
						p.Progress,
						p.Type,
						p.CreateTime,
						edtime.String(),
					)
					fmt.Println(broker)
				}()
			}
			wg.Wait()
		}(database)
	}
	wgs.Wait()
	fmt.Println(c.Add(color.FgHiYellow).Sprint(fmt.Sprintf("LOADING:(%d)、PENDING:(%d)、ETL:(%d)、FINISHED:(%d)、CANCELLED:(%d)、QUEUEING:(%d)", len(LOADING), len(PENDING), len(ETL), len(FINISHED), len(CANCELLED), len(QUEUEING))))
}
