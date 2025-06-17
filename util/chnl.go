/*
 *@author  chengkenli
 *@project srAnr
 *@package util
 *@file    channel
 *@date    2025/6/11 21:14
 */

package util

var Instglobal = Inst() // 初始化全局实例

// Instantiation 用于传递信号的 channel
type Instantiation struct {
	Begin chan struct{}
	vaild chan struct{}
}

// Inst 初始化通道实例
func Inst() *Instantiation {
	return &Instantiation{
		Begin: make(chan struct{}, 1),
		vaild: make(chan struct{}),
	}
}
