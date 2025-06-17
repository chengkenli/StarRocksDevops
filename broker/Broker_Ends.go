/*
 *@author  chengkenli
 *@project starrocks
 *@package broker
 *@file    BrokerEnds
 *@date    2024/7/1 17:40
 */

package broker

import (
    "fmt"
    "github.com/fatih/color"
    "StarRocksDevops/util"
)

func Brokerends() {
    c := color.New()

    var m []map[string]interface{}
    r := util.Connect.Raw("show broker").Scan(&m)
    if r.Error != nil {
        fmt.Println(r.Error.Error())
        return
    }

    fmt.Println(c.Add(color.FgHiGreen).Sprint("brokers:"))
    fmt.Println(fmt.Sprintf("%-12v %-15v %-7v %-7v %-20v %-20v %-10v", "Name", "IP", "Port", "Alive", "LastStartTime", "LastUpdateTime", "ErrMsg"))
    for _, m2 := range m {
        fmt.Println(fmt.Sprintf("%-12v %-15v %-7v %-7v %-20v %-20v %-10v", m2["Name"], m2["IP"], m2["Port"], m2["Alive"], m2["LastStartTime"], m2["LastUpdateTime"], m2["ErrMsg"]))
    }
    fmt.Println()
}
