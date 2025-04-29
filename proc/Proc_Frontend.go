/*
 *@author  chengkenli
 *@project starrocks
 *@package proc
 *@file    ProcFrontend
 *@date    2024/9/2 20:48
 */

package proc

import (
	"fmt"
	"github.com/chengkenli/proc"
	"starrocks/util"
)

// Procfronts 进程查看逻辑
func Procfronts() {
	var command string
	switch util.P.Proc {
	case "fe":
		command = fmt.Sprintf("ps aux | grep 'fe-' | grep -v grep")
	case "be":
		command = fmt.Sprintf("ps aux | grep 'starrocks_be' | grep -v grep")
	case "broker":
		command = fmt.Sprintf("ps aux | grep 'broker' | grep -v grep")
	}

	switch util.P.Proc {
	case "fe":
		for _, host := range util.FrontendsIP {
			result := proc.Connect(
				&proc.ConSSH{
					User:           "starrocks",
					Host:           host,
					Port:           22,
					PrivateKeyFile: "/u/users/svccndlopsns/.ssh/id_rsa",
					Command:        command,
				})
			fmt.Println(host, result.(string))
		}
	case "be", "broker":
		for _, host := range util.BackendsIP {
			result := proc.Connect(
				&proc.ConSSH{
					User:           "starrocks",
					Host:           host,
					Port:           22,
					PrivateKeyFile: "/u/users/svccndlopsns/.ssh/id_rsa",
					Command:        command,
				})
			fmt.Println(host, result.(string))
		}
	}
}
