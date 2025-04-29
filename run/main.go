/**
 * @title: run
 * @Author ChengKen
 * @Date: 10/11/2022 13:56
 * @Version 1.0
 */

package run

import (
	"fmt"
	"starrocks/backend"
	"starrocks/broker"
	"starrocks/errlog"
	"starrocks/frontend"
	"starrocks/ioprofiler"
	"starrocks/proc"
	"starrocks/recovers"
	"starrocks/resource"
	"starrocks/tablet"
	"starrocks/util"
)

func Run() {
	//frontend.SessionVariables()
	//os.Exit(-1)
	//正文
	if util.P.SubmitTask {
		frontend.Submittask()
		return
	}

	if util.P.Stream {
		broker.StreamloadsStatus()
		return
	}
	if util.P.Profilelist {
		frontend.ProfileList()
		return
	}
	if util.P.Restore {
		recovers.Restore()
		return
	}
	if util.P.Recover {
		recovers.Scan()
		return
	}
	if util.P.Io {
		ioprofiler.IoProfile()
		return
	}
	if util.P.ErrScans {
		errlog.ErrScan()
		return
	}
	if util.P.Alter {
		frontend.SessionAlter()
		return
	}
	if util.P.QuerisOne && util.P.QueryId != "" {
		fmt.Println("查询明细：")
		resource.CurrentQueriesCkId()
		return
	}
	if util.P.QuerisOne {
		resource.CurrentQueries()
		//resource.RunningQueries()
		return
	}
	if util.P.Proc != "" {
		proc.Procfronts()
		return
	}
	if util.P.Front {
		frontend.SessionList()
		return
	}
	if util.P.Primary != "" {
		frontend.SessionList()
		return
	}
	if util.P.File != "" {
		frontend.SessionExtract()
		return
	}
	if util.P.Status != "" {
		broker.BrokerStatusEs()
		return
	}
	if util.P.OlapType != "" {
		frontend.SessionOlap()
		return
	}

	if util.P.License {
		frontend.SessionLicense()
		return
	}

	if util.P.RepairSchemaTable {
		tablet.RepairSchema()
		return
	}

	if util.P.Errors {
		frontend.SessionError()
		//if util.P.TabletIsBad {
		//	return
		//}
		//broker.BrokerError()
		return
	}
	if util.P.List {
		frontend.SessionFronBacks()
		backend.BackMemLimit()
		broker.Brokerends()
		backend.Metrics()
		return
	}
	if util.P.JobID > 0 {
		backend.BackendsID()
		return
	}
	if util.P.Pid > 0 && !util.P.Kill {
		frontend.SessionID()
		return
	}

	if len(util.P.Username) != 0 && !util.P.Kill {
		frontend.User()
		return
	}

	if util.P.RoutineLoad {
		broker.RoutineLoad()
		return
	}

	if util.P.JobID > 0 || util.P.JobName != "" {
		broker.Routineload()
		return
	}

	if util.P.Pid > 0 && util.P.Kill {
		frontend.SessionUserPidKill()
		return
	}
	if len(util.P.Username) != 0 && util.P.Kill {
		fmt.Println("查杀矩阵启动中...")
		frontend.SessionUserKill()
		return
	}
	if util.P.Clear {
		fmt.Println("清场矩阵启动中...")
		frontend.SessionClear()
	}

	frontend.SessionVariables()
	frontend.SessionList()
	//resource.CurrentQueries()
	resource.RunningQueries()
	broker.StreamList()
	broker.BrokerList()
	broker.RoutineLoad()
	resource.GroupResource()
	fmt.Println()
}
