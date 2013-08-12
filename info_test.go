/**
 * Created with IntelliJ IDEA.
 * User: zhangye
 * Date: 13-8-8
 * Time: 下午3:53
 * To change this template use File | Settings | File Templates.
 */
package info

import (
    f "fmt"
	"testing"
)

var agent = DefaultAgent

func TestHostName(t *testing.T) {
	h, err := agent.HostName()
	if err != nil {
		t.Error(err)
	}
	f.Println(h)
}

func TestPcpu(t *testing.T) {
	pcpu, err := agent.Pcpu()
	if err != nil {
		t.Error(err)
	}
	for k, v := range pcpu {
		if k == 0 {
			f.Printf("cpu ALL:\n%s\n", v)
		} else {
			f.Printf("cpu %d:\n%s\n", k-1, v)
		}
	}
}

func TestLoadavg(t *testing.T) {
	loadavg, err := agent.Loadavg()
	if err != nil {
		t.Error(err)
	}
	f.Println(loadavg)
}

func TestFree(t *testing.T) {
	free, err := agent.Free()
	if err != nil {
		t.Error(err)
	}
	f.Println(free)
}

func TestHddtemp(t *testing.T) {
	hddtemp, err := agent.HddTemp()
	if err != nil {
		t.Error(err)
	}
	f.Println(hddtemp)
}

func TestTraffic(t *testing.T) {
	traffic, err := agent.Traffic()
	if err != nil {
		t.Error(err)
	}
	f.Println(traffic)
}

func TestAgentSystem(t *testing.T) {
	f.Println("%s\n", agent.System())
}
