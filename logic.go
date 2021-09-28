package main

import (
	"fmt"
	"net_detect/tool"
	"strconv"
)

func changeProxy() {
	println("更改端口为:" + tool.InputStr(4))
	port := 7890
	_, err := fmt.Scanf("%d", &port)
	if err != nil {
		return
	}
	gSetting.Data.Proxy = "127.0.0.1:" + strconv.Itoa(port)
	gSetting.Save()
}

func StartMenu(netHttPing *NetHttping) {
	noSub := make([]*tool.SingleMenu,0)
	changeProxy := tool.SingleMenu{
		Name: "更改端口",
		F:    changeProxy,
		SubMenu: noSub,
	}
	settingSubMenu := []*tool.SingleMenu{&changeProxy}
	changeSetting := tool.SingleMenu{
		Name: "增加配置",
		F:    nil,
		SubMenu: settingSubMenu,
	}
	normalDetect := tool.SingleMenu{
		Name:    "开始检测",
		F: func() {
			NormalDetect(netHttPing)
		},
		SubMenu: noSub,
	}
	rootSubMenu := []*tool.SingleMenu{&changeSetting,&normalDetect}
	root := tool.SingleMenu{
		Name:    "根节点",
		F:       nil,
		SubMenu: rootSubMenu,
	}
	c := tool.CmdMenu{}
	c.Init(&root)
	c.Run()
}
