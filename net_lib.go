package main

import (
	"net/http"
	"net/url"
	"time"
)

type netTest struct {
	addr   string
	ifPort bool
	result chan uint32
}

type NetHttping struct {
	maxChan chan int
	client *http.Client
	clientProxy *http.Client
}

func (h *NetHttping) Init() {
	h.maxChan = make(chan int, gSetting.Data.MaxParallel)
	h.clientProxy = buildHTTPClient(true)
	h.client = buildHTTPClient(false)
}

func (h *NetHttping) Httping(url string, useProxy bool, result chan<- int) {
	go func() {  // 为了避免阻塞Httping，所以在此处起一个
		h.maxChan <- 1 // 用channel做并发控制
		result <- h.httping(url, useProxy)
		<-h.maxChan
	}()
}

func (h *NetHttping) httping(url string, useProxy bool) int {
	req, _ := http.NewRequest("GET", url, nil)
	var r *http.Response
	start := time.Now() // 获取当前时间
	if useProxy {
		r, _ = h.clientProxy.Do(req)
	} else {
		r, _ = h.client.Do(req)
	}
	elapsed := time.Since(start)
	if r == nil {
		return -1
	}
	return int(elapsed.Milliseconds())
}


func buildHTTPClient(isProxy bool) *http.Client {
	var proxy func(*http.Request) (*url.URL, error) = nil
	if isProxy {
		proxy = func(_ *http.Request) (*url.URL, error) {
			return url.Parse("sock5://" + gSetting.Data.Proxy)
		}
	}
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{Transport: transport,Timeout: time.Second * 5}  // 暂定三秒，避免有些注定收不到的请求完全占用了多个线程而卡住了
	return client
}