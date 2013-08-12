/**
 * Created with IntelliJ IDEA.
 * User: luosangnanka
 * Date: 13-8-9
 * Time: 下午3:44
 * To change this template use File | Settings | File Templates.
 */
package info

import (
    "io/ioutil"
	"strings"
)

// 读取文件 /proc/net/snmp 获得udp信息
type Udp struct {
	InDatagrams  string
	NoPorts      string
	InErrors     string
	OutDatagrams string
}

func (u *Udp) Udp2String() string {
	return u.InDatagrams + ":" + u.NoPorts + ":" + u.InErrors + ":" + u.OutDatagrams
}

func (a *Agent) Udp() (*Udp, error) {
	var udp = new(Udp)
	var i int = 0
	b, err := ioutil.ReadFile(gTcp)
	if err != nil {
		a.Log.Println("read /proc/net/snmp:", err.Error())
		return nil, err
	}
	s := strings.SplitAfter(string(b), "\n")
	for _, v := range s {
		if len(v) == 0 {
			continue
		}
		if strings.HasPrefix(v, "Udp") {
			if i == 1 {
				udpTmp := strings.TrimSpace(v)
				udpTmp = strings.Trim(udpTmp, "Udp: ")
				udpNewTmp := strings.Fields(udpTmp)
				udp = &Udp{udpNewTmp[0], udpNewTmp[1], udpNewTmp[2], udpNewTmp[3]}
			}
			i++
		} else {
			udp = &Udp{}
		}
	}

	return udp, nil
}
