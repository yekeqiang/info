info
====
// info 使用说明
// author zhangye mailTo{jiangyeziwh@gmail.com}

一、
1、info.go 用于获取
    HostName 系统运行情况, hostname， 登陆开始时间, 登陆时间
    Load 系统负载情况 包含了cpu、内存、系统负载、以及io状态
    Traffic 系统上传下载流量监控
    Temp系统温度监控 硬盘和CPU温度

    而具体函数对应与对应的文件中

2、host.go获取hostname、登陆时间
3、load.go获取cpu负载等信息
4、net.go 获取上传下载流量以及mac地址ip信息等
5、temperature.go 用于获取温度信息
6、diskstats.go 用于读取磁盘空间使用率
7、tcp.go 用于读取tcp信息
8、time.go 用于获取当前机器时间
9、udp.go 用于获取udp信息

二、使用示例
package main

import (
  f "fmt"
	"info"
)

var agent = info.DefaultAgent

func main() {
    // 获取系统信息
	system := agent.System()
	f.Println(system.System2String())

    // 获取网络信息
	ip, err := agent.Network()
	if err != nil {
		f.Println(err.Error())
	}
	f.Println(ip)

    // 获取时间
	timenow, err := agent.GetTime()
	if err != nil {
		f.Println(err.Error())
	}
	f.Println(timenow)

	// 获取磁盘空间使用率
	diskstat, err := agent.DiskStat()
    if err != nil {
    	f.Println(err.Error())
    }
    f.Println(diskstat)
}

