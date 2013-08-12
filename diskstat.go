/**
 * Created with IntelliJ IDEA.
 * User: luosangnanka
 * Date: 13-8-8
 * Time: 下午9:47
 * To change this template use File | Settings | File Templates.
 */
package info

import (
  	"io/ioutil"
	"strings"
)

// 读取磁盘空间使用率
type DiskStat struct {
	ID     string // 编号
	SdName string // 设备名称

	RCompleNum      string // 读完成次数
	RCompleMergeNum string // 合并完成就
	RSectorsNum     string // 读扇区次数
	RSpentMill      string // 读操作花费毫秒数

	WConpleNum      string // 写完成次数
	WCompleMergeNum string // 合并写完成次数
	WSectirsNum     string // 写扇区次数
	WSpentMill      string // 写操作花费毫秒数

	RWResquestNum string // 正在处理的输入/输出请求书
	RWSpentMill   string // 输入/输出操作花费的毫秒数
	RWSpentMillW  string // 输入/输出操作花费的加权毫秒数
}

func (d *DiskStat) Diskstat2String() string {
	return d.ID + ":" + d.SdName + ":" + d.RCompleNum + ":" + d.RCompleMergeNum + ":" + d.RSectorsNum + ":" + d.RSpentMill + ":" + d.WConpleNum + ":" + d.WCompleMergeNum + ":" + d.WSectirsNum + ":" + d.WSpentMill + ":" + d.RWResquestNum + ":" + d.RWSpentMill + ":" + d.RWSpentMillW
}

// 读取 /proc/diskstat 对应结构体 第一个为整个硬盘统计信息 sd1为第一个分区信息
//  8    1 sda1 131 955 2188 752 2 0 4 2 0 609 754
func (a *Agent) DiskStat() ([]DiskStat, error) {
	var diskstats []DiskStat
	b, err := ioutil.ReadFile(gDiskstat)
	if err != nil {
		a.Log.Println("read /proc/diskstat:", err.Error())
		return nil, err
	}
	s := strings.SplitAfter(string(b), "\n")
	// f.Println(s)
	for _, disk := range s {
		// f.Println(disk)
		diskTmp := strings.Fields(disk)
		if len(diskTmp) == 0 {
			continue
		}
		if strings.HasPrefix(diskTmp[2], "sd") {
			// RCompleNum, _ := strconv.ParseInt(diskTmp[3], 10, 64)
			// RCompleMergeNum, _ := strconv.ParseInt(diskTmp[4], 10, 64)
			// RSectorsNum, _ := strconv.ParseInt(diskTmp[5], 10, 64)
			// RSpentMill, _ := strconv.ParseInt(diskTmp[6], 10, 64)

			// WConpleNum, _ := strconv.ParseInt(diskTmp[7], 10, 64)
			// WCompleMergeNum, _ := strconv.ParseInt(diskTmp[8], 10, 64)
			// WSectirsNum, _ := strconv.ParseInt(diskTmp[9], 10, 64)

			// WSpentMill, _ := strconv.ParseInt(diskTmp[10], 10, 64)
			// RWResquestNum, _ := strconv.ParseInt(diskTmp[11], 10, 64)
			// RWSpentMill, _ := strconv.ParseInt(diskTmp[12], 10, 64)
			// RWSpentMillW, _ := strconv.ParseInt(diskTmp[13], 10, 64)

			diskstats = append(diskstats, DiskStat{diskTmp[1], diskTmp[2], diskTmp[3], diskTmp[4], diskTmp[5], diskTmp[6], diskTmp[7], diskTmp[8], diskTmp[9], diskTmp[10], diskTmp[11], diskTmp[12], diskTmp[13]})
		}
	}
	return diskstats, nil
}

