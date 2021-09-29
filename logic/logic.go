package logic

import (
	"net_detect/net"
	"net_detect/setting"
	"net_detect/tool"
	"strconv"
)

func changeProxy() {
	port := 7890
	err := tool.Input("更改端口为", 4, &port)
	if err != nil {
		return
	}
	setting.GSetting.Data.Proxy = "127.0.0.1:" + strconv.Itoa(port)
	setting.GSetting.Save()
}

func showSetting() {
	//bytes, err := ioutil.ReadFile("config\\setting.json")
	//if err != nil {
	//	println("读取配置失败")
	//	return
	//}
	//print(string(bytes))
	println("代理地址:", setting.GSetting.Data.Proxy)
	println("单网站检测次数:", setting.GSetting.Data.WebCheckNum)
	println("同时并行请求:", setting.GSetting.Data.MaxParallel)
	println("国内网址:", setting.GSetting.Data.WebChina)
	println("国外未ban网址:", setting.GSetting.Data.WebForeignUnban)
	println("国外已ban网址:", setting.GSetting.Data.WebForeignBan)
	println("自定义网址:")
	for _, s := range setting.GSetting.Data.Websites {
		println("  " + s)
	}
}

func StartNetMenu(netHttPing *net.NetHttping) {
	noSub := make([]*tool.SingleMenu, 0)
	changeProxy := tool.SingleMenu{
		Name:    "更改端口",
		F:       changeProxy,
		SubMenu: noSub,
	}
	showSetting := tool.SingleMenu{
		Name:    "查看配置",
		F:       showSetting,
		SubMenu: noSub,
	}
	settingSubMenu := []*tool.SingleMenu{&showSetting, &changeProxy}
	changeSetting := tool.SingleMenu{
		Name:    "配置相关",
		F:       nil,
		SubMenu: settingSubMenu,
	}
	normalDetect := tool.SingleMenu{
		Name: "开始检测",
		F: func() {
			net.NormalDetect(netHttPing)
		},
		SubMenu: noSub,
	}
	rootSubMenu := []*tool.SingleMenu{&normalDetect, &changeSetting}
	root := tool.SingleMenu{
		Name:    "根节点",
		F:       nil,
		SubMenu: rootSubMenu,
	}
	c := tool.CmdMenu{}
	c.Init(&root)
	c.Run()
}
