/**
 * Created with IntelliJ IDEA.
 * User: luosangnanka
 * Date: 13-8-8
 * Time: 上午9:47
 * To change this template use File | Settings | File Templates.
 */
package info

import (
  	"log"
	"os"
)


var TimeFarmat = "2006-01-02 15:04:05 -0700 MST"

// 要读取的 proc 路径
var (
	gUptime   = "/proc/uptime"
	gLoadavg  = "/proc/loadavg"
	gNetdev   = "/proc/net/dev"
	gDiskstat = "/proc/diskstats"
	gCpuinfo  = "/proc/stat"
	gMeminfo  = "/proc/meminfo"
	gTcp      = "/proc/net/snmp"
)

type Agent struct {
	Log *log.Logger
}

func NewAgent(ll *log.Logger) *Agent {
	return &Agent{ll}

}

var DefaultAgent = NewAgent(log.New(os.Stderr, "", log.LstdFlags))

func NewLogger(logname string) *log.Logger {
	file, err := os.OpenFile(logname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("opening log: ", err.Error())
	}

	return log.New(file, "", log.LstdFlags)
}

// 系统信息
// HostName 系统运行情况, hostname， 登陆开始时间, 登陆时间
// Load 系统负载情况 包含了cpu、内存、系统负载、以及io状态
// Traffic 系统上传下载流量监控
// Network 获得主机ip和mac地址 IntranetIp为内网ip ExtranetIp为外网ip
// Temp系统温度监控 硬盘和CPU温度
// DiskStat为读取磁盘空间使用率
// Tcp为Tcp信息
// Udp为Udp信息
// Time为系统时间
type System struct {
	HostName *HostName
	Load     *Load
	Traffic  []Traffic
	NetWork  *Network
	Temp     *Temp
	DiskStat []DiskStat
	Tcp      *Tcp
	Udp      *Udp
	Time     *Time
}

func (s *System) System2String() string {
	s.Load.Free = s.Load.Free.Format()
	Traffic := "System traffic\n"
	for k, v := range s.Traffic {
		Traffic += v.Traffic2String()
		if k < len(s.Traffic)-1 {
			Traffic += "\n"
		}
	}

	DiskStat := "Diskstat\n"
	for k, v := range s.DiskStat {
		DiskStat += v.Diskstat2String()
		if k < len(s.DiskStat)-1 {
			Traffic += "\n"
		}
	}
	res := s.HostName.Host2String() + "\n" + s.Load.Load2String() + "\n" + Traffic + "\n" + s.NetWork.Network2String() + "\n" + s.Temp.Temp2String() + "\n" + DiskStat + "\n" + s.Tcp.Tcp2String() + s.Udp.Udp2String() + s.Time.Time2String()
	return res
}

// when use info pkg use these:
// var agent = info.DefaultAgent
// system := agent.System()
// f.Println(system.System2String())
// then echo the info of system
func (a *Agent) System() *System {
	host, _ := a.HostName()
	load, _ := a.Load()
	traffic, _ := a.Traffic()
	network, _ := a.Network()
	temp, _ := a.Temp()
	diskstat, _ := a.DiskStat()
	tcp, _ := a.Tcp()
	udp, _ := a.Udp()
	time, _ := a.GetTime()
	return &System{host, load, traffic, network, temp, diskstat, tcp, udp, time}
}
