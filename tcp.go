/**
 * Created with IntelliJ IDEA.
 * User: luosangnanka
 * Date: 13-8-9
 * Time: 下午2:59
 * To change this template use File | Settings | File Templates.
 */
package info

import (
    f "fmt"
	"io/ioutil"
	"strings"
)

// 通过读取/proc/net/snmp 获取tcp信息
type Tcp struct {
	ActiveOpens  string
	PassiveOpens string
	InSegs       string
	OutSegs      string
	RetransSegs  string
}

func (t *Tcp) Tcp2String() string {
	return f.Sprintf("ActiveOpens\tPassiveOpens\tInSegs\tOutSegs\tRetransSegs\n%d\t%d\t%d\t%d\t%d\n", t.ActiveOpens, t.PassiveOpens, t.InSegs, t.OutSegs, t.RetransSegs)
}

func (a *Agent) Tcp() (*Tcp, error) {
	var i int = 0
	var tcpAll = new(Tcp)
	b, err := ioutil.ReadFile(gTcp)
	if err != nil {
		a.Log.Println("read /proc/net/snmp:", err.Error())
		return nil, err
	}
	s := strings.SplitAfter(string(b), "\n")
	for _, tcp := range s {
		if len(tcp) == 0 {
			continue
		}
		if strings.HasPrefix(tcp, "Tcp") {
			if i == 1 {
				tcp = strings.TrimSpace(tcp)
				tcp = strings.Trim(tcp, "Tcp: ")
				t := strings.Fields(tcp)
				tcpAll = &Tcp{ActiveOpens: t[4], PassiveOpens: t[5], InSegs: t[9], OutSegs: t[10], RetransSegs: t[11]}
			}
			i++
		} else {
			tcpAll = &Tcp{}
		}
	}

	return tcpAll, nil
}
