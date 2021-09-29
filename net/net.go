package net

import (
	"fmt"
	"net_detect/setting"
)

func TestOneWeb(printStr string, netHttPing *NetHttping, url string, useProxy bool, num int) RetData {
	c := make(chan int, num)
	for i := 0; i < num; i++ {
		netHttPing.Httping(url, useProxy, c)
	}
	return RetData{
		name:     printStr,
		retChan:  c,
		num:      num,
		useProxy: useProxy,
	}
}

type RetData struct {
	name     string // 返回时打印的开头
	retChan  <-chan int
	num      int
	useProxy bool
}

func (r *RetData) WaitAndPrint(sumChan chan<- int) {
	// 这个最好还是设计成多线程非阻塞的，因为存在一个并发的问题，可能会提前跑完下一个请求的数据，让用户看起来比较奇怪。。。（这组数据等半天，下组数据瞬间好，）
	sumTime := 0
	nolost := 0
	// bar := progressbar.Default(int64(r.num))
	for i := 0; i < r.num; i++ {
		t := <-r.retChan
		if t != -1 {
			sumTime += t
			nolost += 1
		}
		// bar.Add(1)
		sumChan <- 1 // 通知上层，进度+1
	}
	s := ""
	s += r.name
	for len(s) < 35 {
		s += " "
	}
	if r.useProxy {
		s += "(代理)\t"
	} else {
		s += "(直连)\t"
	}
	if nolost == 0 {
		s += "完全丢失"
	} else {
		s += fmt.Sprintf("%4dms\t%d%%到达", sumTime/nolost, (nolost*10000)/(r.num*100))
	}
	s += "\n"
	print(s)
}

type AllResData struct {
	retDatas    []RetData
	sumChan     chan int // 用来记录进度
	WebCheckNum int
	ChanNum     int
}

func (d *AllResData) WaitAndPaint() {
	for _, v := range d.retDatas {
		v.WaitAndPrint(d.sumChan)
	}
}

func (d *AllResData) Init(WebCheckNum int, chanNum int) {
	d.sumChan = make(chan int, chanNum*WebCheckNum)
	d.retDatas = make([]RetData, 0)
	d.ChanNum = chanNum
	d.WebCheckNum = WebCheckNum
}

func NormalDetect(netHttPing *NetHttping) {
	webCheckNum := setting.GSetting.Data.WebCheckNum
	allResData := AllResData{}
	allResData.Init(webCheckNum, 6+len(setting.GSetting.Data.Websites)*2)
	allResData.retDatas = append(allResData.retDatas, TestOneWeb("China", netHttPing, setting.GSetting.Data.WebChina, false, webCheckNum))
	allResData.retDatas = append(allResData.retDatas, TestOneWeb("China", netHttPing, setting.GSetting.Data.WebChina, true, webCheckNum))
	allResData.retDatas = append(allResData.retDatas, TestOneWeb("USA unban", netHttPing, setting.GSetting.Data.WebForeignUnban, false, webCheckNum))
	allResData.retDatas = append(allResData.retDatas, TestOneWeb("USA unban", netHttPing, setting.GSetting.Data.WebForeignUnban, true, webCheckNum))
	allResData.retDatas = append(allResData.retDatas, TestOneWeb("USA baned", netHttPing, setting.GSetting.Data.WebForeignBan, false, webCheckNum))
	allResData.retDatas = append(allResData.retDatas, TestOneWeb("USA baned", netHttPing, setting.GSetting.Data.WebForeignBan, true, webCheckNum))
	for _, url := range setting.GSetting.Data.Websites {
		allResData.retDatas = append(allResData.retDatas, TestOneWeb(url, netHttPing, url, false, webCheckNum))
		allResData.retDatas = append(allResData.retDatas, TestOneWeb(url, netHttPing, url, true, webCheckNum))
	}
	allResData.WaitAndPaint()
}