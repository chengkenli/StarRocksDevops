/*
 *@author  chengkenli
 *@project srAnr
 *@package util
 *@file    auth
 *@date    2025/6/11 23:08
 */

package util

import (
    "fmt"
    "os"
)

func init() {
    go func() {
        for {
            select {
            case <-Instglobal.vaild:
                // 校验配置文件中的元数据配置信息
                if Read.Metadb.Host == "" || Read.Metadb.Port == 0 || Read.Metadb.User == "" || Read.Metadb.Password == "" || Read.Metadb.Base == "" {
                    fmt.Println("metadata configuration error")
                    os.Exit(-1)
                }
            }
        }
    }()
}
