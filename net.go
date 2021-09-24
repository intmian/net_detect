package main

import "fmt"

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
	p := GoProcessBar{}
	p.Init()
	p.Run()
	for i := 0; i < r.num; i++ {
		time := <-r.retChan
		if time != -1 {
			sumTime += time
			nolost += 1
		}
		sumChan <- 1 // 通知上层，进度+1
	}
	go func() {

	}()
	s := ""
	s += r.name
	if r.useProxy {
		s += "(代理)\t"
	} else {
		s += "(直连)\t"
	}
	if nolost == 0 {
		s += "完全丢失"
	} else {
		s += fmt.Sprintf("%3dms\t%d%%\t到达", sumTime/nolost, (nolost * 10000)/(r.num*100))
	}
	s += "\n"
	p.Stop()
	print(s)
}

type AllResData struct {
	retDatas    []RetData
	sumChan     chan int // 用来记录进度
	WebCheckNum int
	ChanNum     int
}

func (d *AllResData) WaitAndPaint() {
	for _,v := range d.retDatas {
		v.WaitAndPrint(d.sumChan)
	}
}

func (d *AllResData) Init(WebCheckNum int, chanNum int) {
	d.sumChan = make(chan int, chanNum*WebCheckNum)
	d.retDatas = make([]RetData,0)
	d.ChanNum = chanNum
	d.WebCheckNum = WebCheckNum
}

func NormalDetect(netHttPing *NetHttping) {
	webCheckNum := gSetting.Data.WebCheckNum
	allResData := AllResData{}
	allResData.Init(webCheckNum, 6+len(gSetting.Data.Websites))
	allResData.retDatas = append(allResData.retDatas, TestOneWeb("中国", netHttPing, gSetting.Data.WebChina, false, webCheckNum))
	allResData.retDatas = append(allResData.retDatas, TestOneWeb("中国", netHttPing, gSetting.Data.WebChina, true, webCheckNum))
	allResData.retDatas = append(allResData.retDatas, TestOneWeb("国外未ban", netHttPing, gSetting.Data.WebForeignUnban, false, webCheckNum))
	allResData.retDatas = append(allResData.retDatas, TestOneWeb("国外未ban", netHttPing, gSetting.Data.WebForeignUnban, true, webCheckNum))
	allResData.retDatas = append(allResData.retDatas, TestOneWeb("国外已ban", netHttPing, gSetting.Data.WebForeignBan, false, webCheckNum))
	allResData.retDatas = append(allResData.retDatas, TestOneWeb("国外已ban", netHttPing, gSetting.Data.WebForeignBan, true, webCheckNum))
	for _, url := range gSetting.Data.Websites {
		allResData.retDatas = append(allResData.retDatas, TestOneWeb(url, netHttPing, url, false, webCheckNum))
		allResData.retDatas = append(allResData.retDatas, TestOneWeb(url, netHttPing, url, true, webCheckNum))
	}
	allResData.WaitAndPaint()
}
