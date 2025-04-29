/**
 * @title: def
 * @Author ChengKen
 * @Date: 10/11/2022 13:35
 * @Version 1.0
 */

package util

import "gorm.io/gorm"

type CustomLogger struct{}

func (l *CustomLogger) Errorf(format string, v ...interface{}) {}
func (l *CustomLogger) Warnf(format string, v ...interface{})  {} // 忽略 WARN
func (l *CustomLogger) Debugf(format string, v ...interface{}) {}

var (
	P           ArvgParms
	FrontendsIP []string
	BackendsIP  []string
	Connect     *gorm.DB
	MetaLink    []map[string]interface{}
)

type ConnectParms struct {
	Host       string
	Port       int
	User       string
	Pass       string
	Base       string
	Uri        string
	AaccessKey string
	SecretKey  string
}

type ArvgParms struct {
	Kill              bool
	Pid               int
	BrokerID          int
	Username          string
	Primary           string
	App               string
	Help              bool
	List              bool
	Clear             bool
	ALterName         string
	RoutineLoad       bool
	JobID             int
	JobName           string
	Errors            bool
	ErrorsTime        int
	License           bool
	TabletIsBad       bool
	RepairSchemaTable bool
	Thread            int
	Dbname            string
	BackendId         int
	OlapType          string
	Status            string
	File              string
	QueryId           string
	Front             bool
	Proc              string
	QuerisOne         bool
	Alter             bool
	// 慢查询搜索
	ErrScans    bool
	StartTime   string
	EndTime     string
	CommitType  string
	CommitLimit int
	ErrKey      string
	//io profile
	IP      string
	Seconds int
	Io      bool
	//回收站
	Date      string
	Recover   bool
	Restore   bool
	Table     string
	Partition string
	//profile list
	Profilelist bool
	//stream load
	Stream     bool
	Label      string
	SubmitTask bool
}

type Fronends []struct {
	Name              string `bson:"Name"`
	IP                string `bson:"IP"`
	EditLogPort       int    `bson:"EditLogPort"`
	HttpPort          int    `bson:"HttpPort"`
	QueryPort         int    `bson:"QueryPort"`
	RpcPort           int    `bson:"RpcPort"`
	Role              string `bson:"Role"`
	IsMaster          string `bson:"IsMaster"`
	ClusterId         int    `bson:"ClusterId"`
	Join              string `bson:"Join"`
	Alive             string `bson:"Alive"`
	ReplayedJournalId int    `bson:"ReplayedJournalId"`
	LastHeartbeat     string `bson:"LastHeartbeat"`
	IsHelper          string `bson:"IsHelper"`
	ErrMsg            string `bson:"ErrMsg"`
	StartTime         string `bson:"StartTime"`
	Version           string `bson:"Version"`
}

type Backends []struct {
	BackendId             string `bson:"BackendId"`
	Cluster               string `bson:"Cluster"`
	IP                    string `bson:"IP"`
	HeartbeatPort         int    `bson:"HeartbeatPort"`
	BePort                int    `bson:"BePort"`
	HttpPort              int    `bson:"HttpPort"`
	BrpcPort              int    `bson:"BrpcPort"`
	LastStartTime         string `bson:"LastStartTime"`
	LastHeartbeat         string `bson:"LastHeartbeat"`
	Alive                 string `bson:"Alive"`
	SystemDecommissioned  string `bson:"SystemDecommissioned"`
	ClusterDecommissioned string `bson:"ClusterDecommissioned"`
	TabletNum             int    `bson:"TabletNum"`
	DataUsedCapacity      string `bson:"DataUsedCapacity"`
	AvailCapacity         string `bson:"AvailCapacity"`
	TotalCapacity         string `bson:"TotalCapacity"`
	UsedPct               string `bson:"UsedPct"`
	MaxDiskUsedPct        string `bson:"MaxDiskUsedPct"`
	ErrMsg                string `bson:"ErrMsg"`
	Version               string `bson:"Version"`
	Status                string `bson:"Status"`
	DataTotalCapacity     string `bson:"DataTotalCapacity"`
	DataUsedPct           string `bson:"DataUsedPct"`
}
type Process []struct {
	Id      int    `bson:"Id"`
	User    string `bson:"User"`
	Host    string `bson:"Host"`
	Cluster string `bson:"Cluster"`
	Db      string `bson:"Db"`
	Command string `bson:"Command"`
	Time    int    `bson:"Time"`
	State   string `bson:"State"`
	Info    string `bson:"Info"`
}
type Process2 []struct {
	Id        int    `bson:"Id"`
	User      string `bson:"User"`
	Host      string `bson:"Host"`
	Cluster   string `bson:"Cluster"`
	Db        string `bson:"Db"`
	Command   string `bson:"Command"`
	Time      int    `bson:"Time"`
	State     string `bson:"State"`
	Info      string `bson:"Info"`
	IsPending string `bson:"IsPending"`
	Warehouse string `bson:"Warehouse"`
}

type Dbs []struct {
	Database string `bson:"Database"`
}
type BrokerLoad []struct {
	JobId          int    `bson:"JobId"`
	Label          string `bson:"Label"`
	State          string `bson:"State"`
	Progress       string `bson:"Progress"`
	Type           string `bson:"Type"`
	EtlInfo        string `bson:"EtlInfo"`
	TaskInfo       string `bson:"TaskInfo"`
	ErrorMsg       string `bson:"ErrorMsg"`
	CreateTime     string `bson:"CreateTime"`
	EtlStartTime   string `bson:"EtlStartTime"`
	EtlFinishTime  string `bson:"EtlFinishTime"`
	LoadStartTime  string `bson:"LoadStartTime"`
	LoadFinishTime string `bson:"LoadFinishTime"`
	URL            string `bson:"URL"`
	JobDetails     string `bson:"JobDetails"`
}

type CreateSQL struct {
	Table       string `bson:"Table"`
	CreateTable string `bson:"Create Table"`
}

type Property []struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}
type Grants struct {
	UserIdentity      string `bson:"UserIdentity"`
	Password          string `bson:"Password"`
	AuthPlugin        string `bson:"AuthPlugin"`
	UserForAuthPlugin string `bson:"UserForAuthPlugin"`
	GlobalPrivs       string `bson:"GlobalPrivs"`
	DatabasePrivs     string `bson:"DatabasePrivs"`
	TablePrivs        string `bson:"TablePrivs"`
	ResourcePrivs     string `bson:"ResourcePrivs"`
}

type Queris []struct {
	StartTime     string `bson:"StartTime"`
	QueryId       string `bson:"QueryId"`
	ConnectionId  string `bson:"ConnectionId"`
	Database      string `bson:"Database"`
	User          string `bson:"User"`
	ScanBytes     string `bson:"ScanBytes"`
	ScanRows      string `bson:"ScanRows"`
	MemoryUsage   string `bson:"MemoryUsage"`
	DiskSpillSize string `bson:"DiskSpillSize"`
	CPUTime       string `bson:"CPUTime"`
	ExecTime      string `bson:"ExecTime"`
	Warehouse     string `bson:"Warehouse"`
}

type Querisign struct {
	StartTime     string `bson:"StartTime"`
	QueryId       string `bson:"QueryId"`
	ConnectionId  string `bson:"ConnectionId"`
	Database      string `bson:"Database"`
	User          string `bson:"User"`
	ScanBytes     string `bson:"ScanBytes"`
	ScanRows      string `bson:"ScanRows"`
	MemoryUsage   string `bson:"MemoryUsage"`
	DiskSpillSize string `bson:"DiskSpillSize"`
	CPUTime       string `bson:"CPUTime"`
	ExecTime      string `bson:"ExecTime"`
	Warehouse     string `bson:"Warehouse"`
}

type StreamData struct {
	BeforeLoadTimeMs      string      `json:"BeforeLoadTimeMs"`
	ChannelNum            string      `json:"ChannelNum"`
	ChannelState          string      `json:"ChannelState"`
	CreateTimeMs          string      `json:"CreateTimeMs"`
	DbName                string      `json:"DbName"`
	EndTimeMs             string      `json:"EndTimeMs"`
	ErrorMsg              string      `json:"ErrorMsg"`
	FinishPreparingTimeMs interface{} `json:"FinishPreparingTimeMs"`
	ID                    string      `json:"Id"`
	Label                 string      `json:"Label"`
	LoadID                string      `json:"LoadId"`
	NumLoadBytes          string      `json:"NumLoadBytes"`
	NumRowsAbNormal       string      `json:"NumRowsAbNormal"`
	NumRowsNormal         string      `json:"NumRowsNormal"`
	NumRowsUnselected     string      `json:"NumRowsUnselected"`
	PreparedChannelNum    string      `json:"PreparedChannelNum"`
	StartLoadingTimeMs    interface{} `json:"StartLoadingTimeMs"`
	StartPreparingTimeMs  interface{} `json:"StartPreparingTimeMs"`
	State                 string      `json:"State"`
	TableName             string      `json:"TableName"`
	TimeoutSecond         string      `json:"TimeoutSecond"`
	TrackingSQL           string      `json:"TrackingSQL"`
	TrackingURL           interface{} `json:"TrackingURL"`
	TxnID                 string      `json:"TxnId"`
	Type                  string      `json:"Type"`
}

type HexData struct {
	Status string `json:"status"`
	Data   struct {
		StartTime    int    `json:"startTime"`
		EndTime      int    `json:"endTime"`
		TimeUsedMs   int    `json:"timeUsedMs"`
		QueryID      string `json:"queryId"`
		ClientIP     string `json:"clientIp"`
		State        string `json:"state"`
		SQL          string `json:"sql"`
		User         string `json:"user"`
		CPUCostNs    int    `json:"cpuCostNs"`
		MemCostBytes int    `json:"memCostBytes"`
		ScanRows     int    `json:"scanRows"`
		ScanBytes    int    `json:"scanBytes"`
		Digest       string `json:"digest"`
		Warehouse    string `json:"warehouse"`
		FailedReason string `json:"failedReason"`
		Database     string `json:"database"`
	} `json:"data"`
}
