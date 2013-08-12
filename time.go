/**
 * Created with IntelliJ IDEA.
 * User: luosangnanka
 * Date: 13-8-8
 * Time: 下午8:08
 * To change this template use File | Settings | File Templates.
 * 用于获取系统时间
 */
package info

import (
    "time"
)

// 获取系统时间
type Time struct {
	Time string
}

func (t *Time) Time2String() string {
	return t.Time
}

func (a *Agent) GetTime() (*Time, error) {
	time := time.Now().Format(TimeFarmat)
	var myTime = new(Time)
	myTime.Time = time

	return myTime, nil
}
