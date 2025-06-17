/*
 *@author  chengkenli
 *@project starrocks
 *@package util
 *@file    hex
 *@date    2025/4/8 10:51
 */

package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func HexEncrypt(cryted, key string) string {
	// 创建一个新的HMAC对象，使用SHA256算法和密钥
	hmacHash := hmac.New(sha256.New, []byte(key))
	// 写入数据
	hmacHash.Write([]byte(cryted))
	// 计算HMAC值
	sha := hmacHash.Sum(nil)
	// 将HMAC值转换为十六进制字符串
	hexSha := hex.EncodeToString(sha)
	// 打印结果
	return hexSha
}

// HexEncrypt2 calculates HMAC-SHA256 for a given data and hex-encoded key
func HexEncrypt2(data, key string) string {
	// 将十六进制编码的密钥解码为字节切片
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		fmt.Println(err.Error())
	}
	// 创建一个新的HMAC对象，使用SHA256算法和密钥
	hmacHash := hmac.New(sha256.New, keyBytes)
	// 写入数据
	hmacHash.Write([]byte(data))
	// 计算HMAC值
	sha := hmacHash.Sum(nil)
	// 将HMAC值转换为十六进制字符串
	hexSha := hex.EncodeToString(sha)
	// 返回十六进制字符串
	return hexSha
}
