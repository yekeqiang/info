/**
 * Created with IntelliJ IDEA.
 * User: luosangnanka
 * Date: 13-8-8
 * Time: 上午10:29
 * To change this template use File | Settings | File Templates.
 * 获取系统负载情况 包含了cpu、内存、系统负载、以及io状态
 */
package info

import (
    "errors"
	f "fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var NumCpu = runtime.NumCPU()

// 使用mpstat命令获得CPU使用率
// ID CPU id 第一个为All
// Us 在用户级别（应用程序）执行时发生的物理处理器使用率的百分比
// Sy 在系统级别（内核）执行时发生的物理处理器使用率的百分比
// Wa 逻辑处理器空闲并且没有未完成的磁盘 I/O 请求时的时间百分比
// Idle 逻辑处理器空闲但其间有未完成的磁盘 I/O 请求时的时间百分比
type Pcpu struct {
	ID   string
	Us   float64
	Sy   float64
	Wa   float64
	Idle float64
}

func (p Pcpu) Cpu2String() string {
	return f.Sprintf("ID\tus\tsy\twa\tidle\n%s\t%.2f\t%.2f\t%.2f\t%.2f", p.ID, p.Us, p.Sy, p.Wa, p.Idle)
}

// CPU使用率, 多核CPU每一核使用情况, 第一个All是平均值
// 使用 /usr/bin/mpstat Ubuntu下面需要 sudo apt-get install sysstat
func (a *Agent) Pcpu() ([]Pcpu, error) {
	var pcpus []Pcpu
	all, err := exec.Command("/usr/bin/mpstat", "-P", "ALL").Output()
	if err != nil {
		a.Log.Println("/usr/bin/mpsata -P ALL", err.Error())
		return nil, err
	}

	s := strings.SplitAfter(string(all), "\n")
	for i := 3; i <= NumCpu+3; i++ {
		var cpu = s[i]
		cur := strings.Fields(cpu)

		/* 输出格式
		11时00分14秒  CPU    %usr   %nice    %sys %iowait    %irq   %soft  %steal  %guest   %idle
		11时00分14秒  all    5.30    0.45    2.05    0.39    0.00    0.01    0.00    0.00   91.80
		11时00分14秒    0    6.04    0.83    3.09    2.39    0.00    0.01    0.00    0.00   87.65
		11时00分14秒    1   13.65    0.43    3.45    0.34    0.00    0.00    0.00    0.00   82.13
		11时00分14秒    2    6.99    1.59    3.11    0.19    0.00    0.00    0.00    0.00   88.11
		11时00分14秒    3    8.78    0.72    3.07    0.12    0.00    0.00    0.00    0.00   87.31
		11时00分14秒    4    1.85    0.01    0.79    0.05    0.00    0.04    0.00    0.00   97.26
		11时00分14秒    5    1.65    0.01    1.05    0.02    0.00    0.00    0.00    0.00   97.27
		11时00分14秒    6    1.72    0.01    1.09    0.02    0.00    0.00    0.00    0.00   97.16
		11时00分14秒    7    1.78    0.01    0.82    0.01    0.00    0.00    0.00    0.00   97.38
		*/

		us, _ := strconv.ParseFloat(cur[2], 64)
		sy, _ := strconv.ParseFloat(cur[4], 64)
		wa, _ := strconv.ParseFloat(cur[5], 64)
		idle, _ := strconv.ParseFloat(cur[10], 64)
		pcpus = append(pcpus, Pcpu{cur[1], us, sy, wa, idle})
	}
	return pcpus, nil
}

// 使用 /proc/stat读取Cpu信息 第一个为其他之和
// cpu  15200090 1231646 5689604 256976673 973430 31 14544
type Procpu struct {
	Total  int64
	User   int64
	Nice   int64
	System int64
	Idle   int64
	Iowait int64
	Irq    int64
	Sftirq int64
}

func (p *Procpu) Procpu2String() string {
	return f.Sprintf("total\tuser\tnice\tsystem\tidle\tiowait\tirq\tsftirq\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t", p.Total, p.User, p.Nice, p.System, p.Idle, p.Iowait, p.Irq, p.Sftirq)
}

func (a *Agent) Procpu() ([]Procpu, error) {
	var procpus []Procpu
	b, err := ioutil.ReadFile(gCpuinfo)
	if err != nil {
		a.Log.Println("read /proc/stat:", err.Error())
		return nil, err
	}
	s := strings.SplitAfter(string(b), "\n")
	for _, pcpu := range s {
		pcpuTmp := strings.Fields(pcpu)
		if len(pcpuTmp) == 0 {
			continue
		}
		if strings.HasPrefix(pcpuTmp[0], "cpu") {
			User, _ := strconv.ParseInt(pcpuTmp[1], 10, 64)
			Nice, _ := strconv.ParseInt(pcpuTmp[2], 10, 64)
			System, _ := strconv.ParseInt(pcpuTmp[3], 10, 64)
			Idle, _ := strconv.ParseInt(pcpuTmp[4], 10, 64)
			Iowait, _ := strconv.ParseInt(pcpuTmp[5], 10, 64)
			Irq, _ := strconv.ParseInt(pcpuTmp[6], 10, 64)
			Sftirq, _ := strconv.ParseInt(pcpuTmp[7], 10, 64)
			Total := (User + Nice + System + Idle + Iowait + Irq + Sftirq)

			procpus = append(procpus, Procpu{Total, User, Nice, System, Idle, Iowait, Irq, Sftirq})
		}
	}

	return procpus, nil
}

// 获取系统IO状态
type IOstat string

// 使用命令 iostat -kdx获得系统IO状态
func (a *Agent) IOstat() (IOstat, error) {
	out, err := exec.Command("iostat", "-kdx").Output()
	if err != nil {
		a.Log.Println("iostat -kdx:", err.Error())
		return "", errors.New("iostat commod error")
	}

	str := strings.Replace(string(out), "\n\n", "\n", -1)
	return IOstat(str), nil
}

// 系统负载情况
type Loadavg struct {
	La1, La5, La15 string
	Processes      string
}

func (l *Loadavg) Load2String() string {
	return f.Sprintf("%s %s %s\t%s", l.La1, l.La5, l.La15, l.Processes)
}

// 读取/proc/loadavg
func (a *Agent) Loadavg() (*Loadavg, error) {
	b, err := ioutil.ReadFile(gLoadavg)
	if err != nil {
		a.Log.Println("reading /proc/loadavg:", err.Error())
		return nil, err
	}
	s := strings.Fields(string(b))

	return &Loadavg{s[0], s[1], s[2], s[3]}, nil
}

func (la *Loadavg) Loadavg5() float64 {
	la5, _ := strconv.ParseFloat(la.La5, 64)
	return la5
}

// ByteSize格式化内存或者流量数据为易读的格式
type ByteSize float64

const (
	_           = iota
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
)

func (b ByteSize) Mem2String() string {
	switch {
	case b >= TB:
		return f.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return f.Sprintf("%2.fGB", b/GB)
	case b >= MB:
		return f.Sprintf("%2.fMB", b/MB)
	case b >= KB:
		return f.Sprintf("%2.fKB", b/KB)
	}
	return f.Sprintf("%.2fB", b)
}

// 物理内存
type Mem struct {
	Total   string
	Used    string
	Free    string
	Buffers string
	Cached  string
}

// swap 交换分区
type Swap struct {
	Total string
	Used  string
	Free  string
}

// free -o
type Free struct {
	Mem  Mem
	Swap Swap
}

func (fr *Free) Free2String() string {
	mem, swap := fr.Mem, fr.Swap
	s := f.Sprintf("Mem:\t%s\t%s\t%s\t%s\t%s\nSwap:\t%s\t%s\t%s",
		mem.Total, mem.Used, mem.Free, mem.Buffers, mem.Cached,
		swap.Total, swap.Used, swap.Free)
	return s
}

// 获得内存信息方法一
// 使用”free -o -b"命令 获取内存 需要安装
func (a *Agent) Free() (*Free, error) {
	var free = new(Free)
	bts, err := exec.Command("free", "-o", "-b").Output()
	if err != nil {
		a.Log.Println("free -o -b", err.Error())
		return nil, err
	}
	lines := strings.Split(string(bts), "\n")
	m, s := strings.Fields(lines[1]), strings.Fields(lines[2])

	free.Mem = Mem{Total: m[1], Used: m[2], Free: m[3], Buffers: m[5], Cached: m[6]}
	free.Swap = Swap{Total: s[1], Used: s[2], Free: s[3]}
	return free, nil
}

// 获得内存信息方法二
// 通过读取/proc/meminfo 获得内存信息
func (a *Agent) Meminfo() (*Free, error) {
	var meminfo = new(Free)
	b, err := ioutil.ReadFile(gMeminfo)
	if err != nil {
		a.Log.Println("read /proc/meminfo:", err.Error())
		return nil, err
	}
	s := strings.SplitAfter(string(b), "\n")
	var m []string
	var mm []string
	for _, mem := range s {
		if mem == "" {
			continue
		}
		memTmp := strings.Split(mem, ":")
		ss := strings.TrimSpace(memTmp[1])
		ss = strings.Trim(ss, "kB ")
		m = append(m, strings.TrimSpace(ss))
	}
	// f.Println(m)
	// 计算mem使用
	total, err := strconv.Atoi(m[0])
	if err != nil {
		a.Log.Println("total strconv:", err.Error())
		return nil, err
	}
	free, err := strconv.Atoi(m[1])
	if err != nil {
		a.Log.Println("free strconv:", err.Error())
		return nil, err
	}
	used := (total - free) * 1024
	newUsed := strconv.Itoa(used)

	// 计算swap使用
	swaptotal, err := strconv.Atoi(m[13])
	if err != nil {
		a.Log.Println("swaptotal strconv:", err.Error())
		return nil, err
	}
	swapfree, err := strconv.Atoi(m[14])
	if err != nil {
		a.Log.Println("swapfree strconv:", err.Error())
		return nil, err
	}
	swapused := (swaptotal - swapfree) * 1024
	newSwapused := strconv.Itoa(swapused)

	for _, v := range m {
		mmTmp, _ := strconv.Atoi(v)
		mmTmp *= 1024
		mm = append(mm, strconv.Itoa(mmTmp))
	}

	meminfo.Mem = Mem{Total: mm[0], Used: newUsed, Free: mm[1], Buffers: mm[2], Cached: mm[3]}
	meminfo.Swap = Swap{Total: mm[13], Used: newSwapused, Free: mm[14]}
	return meminfo, nil
}

// 未使用的内存加上缓存
func (m *Free) Real() float64 {
	r := m.Mem
	free, _ := strconv.ParseFloat(r.Free, 32)
	buffers, _ := strconv.ParseFloat(r.Buffers, 32)
	cached, _ := strconv.ParseFloat(r.Cached, 32)

	rf := free + buffers + cached
	return rf
}

func (m *Free) Total() float64 {
	total, _ := strconv.ParseFloat(m.Mem.Total, 32)
	return total
}

// 内存信息转换成TB GB MB KB 格式
func format(s string) ByteSize {
	bys, _ := strconv.ParseFloat(s, 64)
	return ByteSize(bys)
}

func (m *Free) Format() *Free {
	var r, s = m.Mem, m.Swap

	fr := Mem{
		format(r.Total).Mem2String(),
		format(r.Used).Mem2String(),
		format(r.Free).Mem2String(),
		format(r.Buffers).Mem2String(),
		format(r.Cached).Mem2String(),
	}
	fs := Swap{
		format(s.Total).Mem2String(),
		format(s.Used).Mem2String(),
		format(s.Free).Mem2String(),
	}
	return &Free{fr, fs}
}

// Load struct 包含了cpu、内存、系统负载、以及io状态
type Load struct {
	Cpu     []Procpu
	Free    *Free
	Loadavg *Loadavg
	IO      IOstat
}

func (l *Load) Load2String() string {
	head := "System Load status"
	var cpu string
	for k, v := range l.Cpu {
		if k == 0 {
			cpu += f.Sprintf("CPU ALL:\n%s\n", v)
			continue
		}
		cpu += f.Sprintf("CPU %d:\n%s\n", k-1, v)
	}
	return f.Sprintf("%s\nCPU status: %s\nMemory status:\n%s\n\nLoadavg: %s\n\nIostat:\n%s",
		head, cpu, l.Free, l.Loadavg, l.IO)
}

func (a *Agent) Load() (*Load, error) {
	pcpu, err1 := a.Procpu()
	free, err2 := a.Meminfo()
	loadavg, err3 := a.Loadavg()
	iostat, err4 := a.IOstat()
	return &Load{pcpu, free, loadavg, iostat}, f.Errorf("%s,%s,%s,%s", err1, err2, err3, err4)
}
