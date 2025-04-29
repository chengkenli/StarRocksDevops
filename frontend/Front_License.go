package frontend

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

func SessionLicense() {
	avg := init_meta()

	type license struct {
		Code int `json:"code"`
		List []struct {
			Cores    int   `json:"cores"`
			ExpireAt int64 `json:"expire_at"`
			Hosts    int   `json:"hosts"`
		} `json:"list"`
		Total int `json:"total"`
	}
	type key struct {
		Code int `json:"code"`
		Data struct {
			Cores int    `json:"cores"`
			D     string `json:"d"`
		} `json:"data"`
	}
	fmt.Println()

	//创建Resty客户端
	Client := resty.New()
	//发送POST请求并处理响应
	respones1, err := Client.R().SetBody(map[string]string{
		"name":     avg.User,
		"password": avg.Pass,
	}).Post(avg.Uri + "/api/user/login")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(respones1.Body()))
	/*------------------------------------------------*/
	respones2, err := Client.R().Get(avg.Uri + "/api/license/list")
	if err != nil {
		fmt.Println("请求失败:", err.Error())
		return
	}
	fmt.Println(string(respones2.Body()))
	/*------------------------------------------------*/
	respones3, err := Client.R().Get(avg.Uri + "/api/license/collect-hosts-info")
	if err != nil {
		fmt.Println("请求失败:", err.Error())
		return
	}
	fmt.Println(string(respones3.Body()))
	/*------------------------------------------------*/
	var k key
	err = json.Unmarshal(respones3.Body(), &k)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var l license
	err = json.Unmarshal(respones2.Body(), &l)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println()
	c := color.New()
	for _, s := range l.List {
		fmt.Println(fmt.Sprintf("%-10s: %s", "主机数量", c.Add(color.FgHiRed).Sprint(s.Hosts)))
		fmt.Println(fmt.Sprintf("%-10s: %s", "核心数量", c.Add(color.FgHiRed).Sprint(s.Cores)))
		fmt.Println(fmt.Sprintf("%-10s: %s", "过期时间", c.Add(color.FgHiYellow).Sprint(UnixToTime(strconv.FormatInt(s.ExpireAt, 10)).Format("2006-01-02 15:04:05"))))
		fmt.Println(fmt.Sprintf("%-14s: %s", "license", c.Add(color.FgHiBlue).Sprint(k.Data.D)))
		fmt.Println(fmt.Sprintf("%-10s: %s", "手机号码", c.Add(color.FgHiWhite).Sprint("15217032776")))
		fmt.Println(fmt.Sprintf("%-10s: %s", "邮箱号码", c.Add(color.FgHiWhite).Sprint("Chengken.Li@walmart.com")))
	}
	fmt.Println()
}

func UnixToTime(e string) (datatime time.Time) {
	data, _ := strconv.ParseInt(e, 10, 64)
	datatime = time.Unix(data/1000, 0)
	return
}
