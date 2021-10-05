package setting

import "net_detect/misc"

type SettingJson struct {
	Proxy                             string   `json:"proxy"`         // 代理
	WebCheckNum                       int      `json:"web_check_num"` // 单个网页的测试次数
	MaxParallel                       int      `json:"max_parallel"`  // 最大测试并发，过大会因为阻塞测不准
	Websites                          []string `json:"websites"`      // 自定义的一些网页
	WebChina                          string   `json:"web_china"`
	WebForeignUnban                   string   `json:"web_foreign_unban"`
	WebForeignBan                     string   `json:"web_foreign_ban"`
	HttpTimeOutSecond                 int      `json:"http_time_out_second"`                   // http client 设置的timeout
	HttpRequestRandTimeOutMillisecond int      `json:"http_request_rand_time_out_millisecond"` // http client 发送请求的随机延迟
}

type Setting struct {
	Data *SettingJson
	misc.TJsonTool
}

func NewSetting() *Setting {
	// 当配置文件没有的时候会panic，这个需要提示
	j := SettingJson{}
	return &Setting{
		Data:      &j,
		TJsonTool: *misc.NewTJsonTool("config\\setting.json", &j),
	}
}

var GSetting = *NewSetting()
