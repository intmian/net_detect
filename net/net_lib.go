package net

import (
	"math/rand"
	"net/http"
	"net/url"
	"net_detect/setting"
	"time"
)

type NetHttping struct {
	maxChan     chan int
	client      *http.Client
	clientProxy *http.Client
}

func (h *NetHttping) Init() {
	h.maxChan = make(chan int, setting.GSetting.Data.MaxParallel)
	h.clientProxy = buildHTTPClient(true)
	h.client = buildHTTPClient(false)
}

func (h *NetHttping) Httping(url string, useProxy bool, result chan<- int) {
	go func() { // 为了避免阻塞Httping，所以在此处起一个
		h.maxChan <- 1 // 用channel做并发控制
		randT := rand.Uint64() % uint64(setting.GSetting.Data.HttpRequestRandTimeOutMillisecond) // 避免阻塞做一个随机延迟
		time.Sleep(time.Millisecond * time.Duration(randT))
		result <- h.httping(url, useProxy)
		<-h.maxChan
	}()
}

func (h *NetHttping) httping(url string, useProxy bool) int {
	req, _ := http.NewRequest("GET", url, nil)
	var r *http.Response
	//println(url, useProxy, "发起")
	start := time.Now() // 获取当前时间
	if useProxy {
		r, _ = h.clientProxy.Do(req)
	} else {
		r, _ = h.client.Do(req)
	}
	elapsed := time.Since(start)
	//println(url, useProxy, elapsed.Milliseconds())
	if r == nil {
		return -1
	}
	return int(elapsed.Milliseconds())
}

func buildHTTPClient(isProxy bool) *http.Client {
	var proxy func(*http.Request) (*url.URL, error) = nil
	if isProxy {
		proxy = func(_ *http.Request) (*url.URL, error) {
			return url.Parse("sock5://" + setting.GSetting.Data.Proxy)
		}
	}
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{Transport: transport, Timeout: time.Second * time.Duration(setting.GSetting.Data.HttpTimeOutSecond)} // 暂定三秒，避免有些注定收不到的请求完全占用了多个线程而卡住了
	return client
}
