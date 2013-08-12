/**
 * Created with IntelliJ IDEA.
 * User: luosangnanka
 * Date: 13-8-8
 * Time: 下午3:23
 * To change this template use File | Settings | File Templates.
 */
package info

import (
    f "fmt"
	"io/ioutil"
	"net"
	// "os/exec"
	"strconv"
	"strings"
	"time"
)

// 上传和下载的流量, 从系统启动之后累加
type Traffic struct {
	Name     string
	Receive  float64
	Transmit float64
	Time     time.Time
}

func (t Traffic) Traffic2String() string {
	recv, trans := ByteSize(t.Receive), ByteSize(t.Transmit)
	ts := t.Time.Format(TimeFarmat)
	return f.Sprintf("%s receive:%s, transmit:%s, %v", t.Name, recv, trans, ts)
}

// 读取 /proc/net/dev
func (a *Agent) Traffic() ([]Traffic, error) {
	bts, err := ioutil.ReadFile(gNetdev)
	if err != nil {
		a.Log.Println("read /proc/net/dev failed", err.Error())
		return nil, err
	}

	date := time.Now()
	lines := strings.Split(string(bts), "\n")
	var traffic []Traffic
	for i := 3; i < len(lines); i++ {
		t := strings.Fields(lines[i])
		if len(t) == 17 {
			name := strings.Trim(t[0], ":")
			recv, _ := strconv.ParseFloat(t[1], 64)
			tran, _ := strconv.ParseFloat(t[10], 64)
			traffic = append(traffic, Traffic{name, recv, tran, date})
		}
	}
	return traffic, nil
}

// 获得主机ip和mac地址 IntranetIp为内网ip ExtranetIp为外网ip
type Network struct {
	IntranetIp string
	ExtranetIp string
	Mac        string
}

func (n *Network) Network2String() string {
	return n.IntranetIp + ":" + n.ExtranetIp + ":" + n.Mac
}

func (a *Agent) Network() (*Network, error) {
	// 获取内网ip
	addr, err := net.InterfaceAddrs()
	if err != nil {
		a.Log.Println("IntranetIp:", err.Error())
		return nil, err
	}
	var network = new(Network)
	network.IntranetIp = strings.Split(addr[1].String(), "/")[0]

	// 获取外网出口ip 通过shell脚本执行 可能需要一点时间
	// var cmd = "./bin/ip"

	// ip, err := exec.Command("/usr/bin/curl", "ifconfig.me").Output()
	// if err != nil {
	//a.Log.Println("ExtranetIp:", err.Error())
	// return network, err
	// }
	// network.ExtranetIp = strings.TrimSpace(string(ip))

	// 获取mac地址
	macInfo, err := net.Interfaces()
	if err != nil {
		a.Log.Println("Mac:", err.Error())
		return network, err
	}
	for _, mac := range macInfo {
		network.Mac = mac.HardwareAddr.String()
	}

	return network, nil
}

// net整体信息
type Net struct {
	Traffic []Traffic
	Network *Network
}

func (a *Agent) Net() (*Net, error) {
	traffic, err1 := a.Traffic()
	network, err2 := a.Network()
	return &Net{traffic, network}, f.Errorf("%s,%s", err1, err2)
}
