/**
 * Created with IntelliJ IDEA.
 * User: luosangnanka
 * Date: 13-8-8
 * Time: 上午9:47
 * To change this template use File | Settings | File Templates.
 */
package info

import (
    "io/ioutil"
	"os"
	"time"
)

// 服务器主机名、启动时间以及运行时间
type HostName struct {
	Name   string
	Boot   time.Time
	Uptime string
}

func (h *HostName) Host2String() string {
	return h.Name + ":" + h.Boot.Format(TimeFarmat) + ":" + h.Uptime
}

// 读取/proc/uptime, 第一数值即为系统运行时间, 单位为(s)
func (a *Agent) HostName() (*HostName, error) {
	b, err := ioutil.ReadFile(gUptime)
	if err != nil {
		a.Log.Println("ReadFile /proc/uptime:", err.Error())
	}
	for i := 0; i < len(b); i++ {
		if b[i] == ' ' {
			b = b[0:i]
			break
		}
	}

	t := string(b) + "s"
	// 获得已经运行时间
	d, err := time.ParseDuration(t)
	if err != nil {
		a.Log.Println("time.ParseDuration:", err.Error())
		return nil, err
	}
	// 获得登陆时间
	boot := time.Now().Add(-d)

	// 获得主机名
	hostname, err := os.Hostname()
	if err != nil {
		a.Log.Println("hostname:", err.Error())
		return nil, err
	}

	return &HostName{hostname, boot, d.String()}, nil
}
