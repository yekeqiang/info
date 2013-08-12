/**
 * Created with IntelliJ IDEA.
 * User: zhangye
 * Date: 13-8-8
 * Time: 下午2:40
 * To change this template use File | Settings | File Templates.
 * 获取硬盘和CPU温度
 */
package info

import (
    "bytes"
	"errors"
	f "fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"time"
)

// 硬盘温度
type HddTemp struct {
	Dev  string
	Desc string
	Temp string
}

func (t HddTemp) HddTemp2String() string {
	return t.Dev + ":" + t.Desc + ":" + t.Temp
}

func newHddTemp(dev, desc, temp string) HddTemp {
	return HddTemp{dev, desc, temp}
}

// 需要HddTmep以守护进程运行, 如: "sudo hddtemp -d /dev/sd[a-d]"
// 需要安装工具hddtemp Ubuntu 下 sudo apt-get install hddtemp
func (a *Agent) HddTemp() (temps []HddTemp, err error) {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:7634", 2*time.Second)
	if err != nil {
		a.Log.Println("tcp://127.0.0.1:7634", err.Error())
		return
	}
	defer conn.Close()

	bts, err := ioutil.ReadAll(conn)
	if err != nil {
		a.Log.Println("reading from tcp://127.0.0.1:7643", err.Error())
		return
	}

	line := bytes.Split(bts, []byte("\n"))
	for i := 0; i < len(line); i++ {
		s := bytes.Split(line[i], []byte("|"))
		temps = append(temps, newHddTemp(string(s[1]), string(s[2]), string(s[3])))
	}
	return temps, nil
}

// CPU温度
// 使用命令sensors Ubuntu需要安装 sudo apt-get install lm-sensors
type Sensor string

func (a *Agent) Sensors() (Sensor, error) {
	bts, err := exec.Command("sensors").Output()
	if err != nil {
		a.Log.Println(`sensors error, try to run commod "sensors" or maybe you should install lm-sensors first`)
		return "", errors.New("Sensors error")
	}
	return Sensor(bts), nil
}

type Temp struct {
	Disks   []HddTemp
	Sensors Sensor
}

func (temp *Temp) Temp2String() string {
	head := "System temperature"
	var hdd string
	for _, disk := range temp.Disks {
		hdd += f.Sprintf("%s\n", disk)
	}
	return f.Sprintf("%s\nHddTemp:\n%s\nSensors:\n%s", head, hdd, temp.Sensors)
}

func (a *Agent) Temp() (*Temp, error) {
	var templ = new(Temp)
	sensors, err := a.Sensors()
	if err != nil {
		sensors = Sensor("No sensors found or maybe you should install lm-sensors first")
	}
	templ.Sensors = sensors
	disks, err := a.HddTemp()
	if err != nil {
		return templ, errors.New("func HddTemp failed")
	}
	templ.Disks = disks
	return templ, nil
}
