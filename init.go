package main

import (
	"fmt"
	"starrocks/util"
	"sync"
)

func initFrontendsIP(wg *sync.WaitGroup) {
	defer wg.Done()
	type frontends []struct {
		Name              string `json:"Name"`
		IP                string `json:"IP"`
		EditLogPort       int    `json:"EditLogPort"`
		HttpPort          int    `json:"HttpPort"`
		QueryPort         int    `json:"QueryPort"`
		RpcPort           int    `json:"RpcPort"`
		Role              string `json:"Role"`
		IsMaster          bool   `json:"IsMaster"`
		ClusterId         int    `json:"ClusterId"`
		Join              bool   `json:"Join"`
		Alive             bool   `json:"Alive"`
		ReplayedJournalId int    `json:"ReplayedJournalId"`
		LastHeartbeat     string `json:"LastHeartbeat"`
		IsHelper          bool   `json:"IsHelper"`
		ErrMsg            string `json:"ErrMsg"`
		StartTime         string `json:"StartTime"`
		Version           string `json:"Version"`
	}
	var f frontends
	r := util.Connect.Raw("show frontends").Scan(&f)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return
	}
	for _, s := range f {
		util.FrontendsIP = append(util.FrontendsIP, s.IP)
	}
}

func initBackendsIP(wg *sync.WaitGroup) {
	defer wg.Done()
	var f []map[string]interface{}
	r := util.Connect.Raw("show backends").Scan(&f)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return
	}
	for _, s := range f {
		util.BackendsIP = append(util.BackendsIP, s["IP"].(string))
	}
}
