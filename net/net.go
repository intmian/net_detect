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
	// 由于调度存在时间，可能导致先发起的Test的协程可能没有后面的启动快从而导致发起迟了而卡住，这里的调度是存在一些问题的，之后要修改一下.详见下面日志
	/*
	https://www.baidu.com false 发起
	https://www.baidu.com false 发起
	https://www.baidu.com true 发起
	https://www.intmian.com true 发起
	https://www.baidu.com true 发起
	https://www.baidu.com true 发起
	https://www.baidu.com false 发起
	https://www.youtube.com true 发起
	https://www.baidu.com true 发起
	https://www.baidu.com false 发起
	https://www.baidu.com true 83
	https://www.intmian.com false 发起
	https://www.baidu.com true 86
	https://www.intmian.com false 发起
	https://www.baidu.com false 87
	https://www.intmian.com false 发起
	https://www.baidu.com true 89
	https://www.intmian.com false 发起
	https://www.baidu.com false 90
	https://www.intmian.com false 发起
	https://www.baidu.com false 92
	https://github.com true 发起
	https://www.baidu.com false 95
	https://www.intmian.com true 发起
	https://www.baidu.com true 110
	https://www.intmian.com true 发起
	https://www.youtube.com true 621
	https://www.intmian.com true 发起
	https://github.com true 806
	https://www.google.com false 发起
	https://www.intmian.com true 568
	https://www.google.com false 发起
	https://www.intmian.com true 1343
	https://www.intmian.com true 1454
	https://www.google.com false 发起
	https://www.google.com false 发起
	https://www.intmian.com true 1358
	https://www.google.com false 发起
	https://www.intmian.com false 5002
	https://www.intmian.com false 5005
	https://www.google.com true 发起
	https://www.google.com true 发起
	https://www.intmian.com false 5007
	https://www.intmian.com false 5005
	https://www.google.com true 发起
	https://www.google.com true 发起
	https://www.intmian.com false 5001
	https://www.google.com true 发起
	https://www.google.com false 5008
	https://www.baidu.com false 发起
	https://www.baidu.com false 45
	https://www.baidu.com false 发起
	https://www.baidu.com false 56
	https://www.baidu.com false 发起
	https://www.baidu.com false 65
	https://www.baidu.com false 发起
	https://www.baidu.com false 63
	https://www.baidu.com false 发起
	https://www.google.com false 5005
	https://www.baidu.com true 发起
	https://www.baidu.com false 75
	https://www.baidu.com true 发起
	https://www.baidu.com true 53
	https://www.baidu.com true 发起
	https://www.baidu.com true 70
	https://www.baidu.com true 发起
	https://www.baidu.com true 52
	https://www.baidu.com true 发起
	https://www.baidu.com true 61
	https://www.intmian.com false 发起
	https://www.baidu.com true 52
	https://www.intmian.com false 发起
	https://www.google.com false 5003
	https://www.google.com false 5003
	https://www.intmian.com false 发起
	https://www.intmian.com false 发起
	https://www.google.com false 5002
	https://www.intmian.com false 发起
	https://www.google.com true 1394
	https://www.intmian.com true 发起
	https://www.google.com true 1449
	https://www.intmian.com true 发起
	https://www.intmian.com true 202
	https://www.intmian.com true 发起
	https://www.google.com true 1607
	https://www.google.com true 1609
	https://www.intmian.com true 发起
	https://www.intmian.com true 发起
	https://www.google.com true 1611
	https://github.com false 发起
	https://www.intmian.com true 186
	https://github.com false 发起
	https://www.intmian.com true 179
	https://github.com false 发起
	https://www.intmian.com true 188
	https://github.com false 发起
	https://www.intmian.com true 219
	https://github.com false 发起
	https://www.intmian.com false 2497
	https://www.intmian.com false 2597
	https://github.com true 发起
	https://github.com true 发起
	https://www.intmian.com false 2497
	https://www.intmian.com false 2496
	https://www.baidu.com true 发起
	https://www.intmian.com true 发起
	https://www.intmian.com false 2604
	https://www.baidu.com false 发起
	https://www.baidu.com false 39
	https://store.steampowered.com true 发起
	https://www.baidu.com true 42
	China                              (直连)         80ms  100%到达
	China                              (代理)         82ms  100%到达
	https://github.com true 发起
	USA unban                          (直连)       完全丢失
	https://www.intmian.com true 187
	https://github.com true 发起
	USA unban                          (代理)        982ms  100%到达
	USA baned                          (直连)       完全丢失
	USA baned                          (代理)       1534ms  100%到达
	https://www.baidu.com              (直连)         60ms  100%到达
	https://www.baidu.com              (代理)         57ms  100%到达
	https://www.intmian.com            (直连)       2538ms  100%到达
	https://www.intmian.com            (代理)        194ms  100%到达
	https://github.com true 259
	https://www.youku.com false 发起
	https://github.com true 344
	https://www.youku.com false 发起
	https://github.com true 515
	https://www.youku.com false 发起
	https://github.com true 502
	https://www.youku.com false 发起
	https://www.youku.com false 2071
	https://www.youku.com false 发起
	https://www.youku.com false 2162
	https://www.youku.com true 发起
	https://www.youku.com false 40
	https://www.youku.com true 发起
	https://store.steampowered.com true 2470
	https://www.youku.com true 发起
	https://www.youku.com false 1991
	https://www.youku.com true 发起
	https://www.youku.com true 158
	https://www.youku.com true 发起
	https://www.youku.com true 41
	https://store.steampowered.com false 发起
	https://www.youku.com false 1907
	https://store.steampowered.com false 发起
	https://www.youku.com true 44
	https://store.steampowered.com false 发起
	https://www.youku.com true 226
	https://store.steampowered.com false 发起
	https://github.com false 5013
	https://store.steampowered.com false 发起
	https://github.com false 5005
	https://store.steampowered.com true 发起
	https://www.youku.com true 334
	https://store.steampowered.com true 发起
	https://github.com false 5004
	https://store.steampowered.com true 发起
	https://github.com false 5000
	https://www.google.com true 发起
	https://github.com false 5010
	https://github.com                 (直连)       完全丢失
	https://store.steampowered.com true 发起
	https://github.com                 (代理)        485ms  100%到达
	https://www.youku.com              (直连)       1634ms  100%到达
	https://www.youku.com              (代理)        160ms  100%到达
	https://store.steampowered.com false 453
	https://www.google.com false 发起
	https://store.steampowered.com false 490
	https://www.google.com false 发起
	https://store.steampowered.com false 492
	https://www.google.com false 发起
	https://store.steampowered.com false 421
	https://www.google.com false 发起
	https://store.steampowered.com false 490
	https://www.google.com false 发起
	https://store.steampowered.com     (直连)        469ms  100%到达
	https://store.steampowered.com true 473
	https://www.google.com true 发起
	https://store.steampowered.com true 468
	https://www.google.com true 发起
	https://www.google.com true 651
	https://www.google.com true 发起
	https://store.steampowered.com true 859
	https://www.youtube.com false 发起
	https://store.steampowered.com true 744
	https://www.google.com true 发起
	https://store.steampowered.com     (代理)       1002ms  100%到达
	https://www.google.com true 601
	https://www.youtube.com false 发起
	https://www.google.com true 633
	https://www.youtube.com false 发起
	https://www.google.com true 642
	https://www.youtube.com false 发起
	https://www.google.com true 619
	https://www.youtube.com true 发起
	https://www.youtube.com true 210
	https://www.youtube.com false 发起
	https://www.google.com false 5000
	https://www.google.com false 5005
	https://www.youtube.com true 发起
	https://www.youtube.com true 发起
	https://www.google.com false 5010
	https://www.youtube.com true 发起
	https://www.google.com false 5008
	https://www.google.com false 5011
	https://www.google.com             (直连)       完全丢失
	https://www.google.com             (代理)        629ms  100%到达
	https://www.youtube.com true 213
	https://www.youtube.com true 213
	https://www.youtube.com true 281
	https://www.youtube.com false 5001
	https://www.youtube.com false 5003
	https://www.youtube.com false 5009
	https://www.youtube.com false 5015
	https://www.youtube.com false 5003
	https://www.youtube.com            (直连)       完全丢失
	https://www.youtube.com            (代理)        307ms  100%到达
	*/
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
