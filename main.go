package main

import (
	"github.com/kirinlabs/HttpRequest"
	"net/http"
	"net/url"
	"time"
)

var proxyConf = "112.195.81.161:8118"

func main() {
	req := HttpRequest.Request{}
	req.SetTimeout(5 * time.Second)
	req.Proxy(func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxyConf)
	})
}
