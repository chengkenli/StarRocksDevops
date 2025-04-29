/*
 *@author  chengkenli
 *@project starrocks
 *@package resource
 *@file    init
 *@date    2025/4/28 14:22
 */

package resource

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"starrocks/conn"
	"starrocks/util"
)

var avg util.ConnectParms
var leader string
var client *resty.Client

func init() {
	//匹配元数据
	for _, m := range util.MetaLink {
		if m["app"].(string) == util.P.App {
			avg = util.ConnectParms{
				Host: m["feip"].(string),
				Port: int(m["feport"].(int32)),
				User: m["user"].(string),
				Pass: m["password"].(string),
			}
		}
	}
	_init()
	//获取leader
	var f []map[string]interface{}
	db, err := conn.StarRocks(util.P.App)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	r := db.Raw("show frontends").Scan(&f)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return
	}
	for _, m := range f {
		if m["Role"].(string) == "LEADER" {
			leader = m["IP"].(string)
		}
	}
	if len(leader) == 0 {
		fmt.Println("leader is nil.")
	}
}

func _init() {
	//创建Resty客户端
	client = resty.New().SetLogger(&util.CustomLogger{}).SetBasicAuth(avg.User, avg.Pass)
}
