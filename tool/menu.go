package tool

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type SingleMenu struct {
	Name    string // 项目名
	F       func() // 可以被call的节点
	SubMenu []*SingleMenu
}

type CmdMenu struct {
	root         *SingleMenu
	now          *SingleMenu
	HisList      []*SingleMenu // 方便完成返回到上一步
	HisListIndex int
}

func (m *CmdMenu) Init(root *SingleMenu) {
	m.root = root
	m.now = root
	m.HisList = make([]*SingleMenu, 100)
	m.HisListIndex = 0
}

func (m *CmdMenu) returnToLast() {
	index := m.HisListIndex
	if index == 0 {
		return
	}
	if m.HisList[index-1] == nil {
		return
	}
	m.HisListIndex -= 1
	m.now = m.HisList[m.HisListIndex]
}

func (m *CmdMenu) gotoSub(index int) bool {
	if m.now == nil {
		return false
	}
	if len(m.now.SubMenu)-1 < index {
		return false
	}
	m.HisList[m.HisListIndex] = m.now
	m.HisListIndex++
	m.now = m.now.SubMenu[index]
	return true
}

func (m *CmdMenu) gotoRoot() {
	m.now = m.root
	m.HisListIndex = 0
}

func (m *CmdMenu) clear() {
	cmd := exec.Command("cmd.exe", "/c", "cls")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		print("clear fail")
	}
}

func (m *CmdMenu) stop() {
	cmd := exec.Command("cmd.exe", "/c", "pause")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		print("clear fail")
	}
}

func (m *CmdMenu) do() (exit bool) {
	m.clear()
	lenSub := len(m.now.SubMenu)
	if lenSub == 0 {
		// 跑到叶节点了，就执行这个的逻辑，停一下，再接着跑
		println(m.now.Name)
		m.now.F()
		m.stop()
		m.returnToLast()
		return false
	}
	if m.now == nil {
		return
	}
	for i, s := range m.now.SubMenu {
		println(strconv.Itoa(i) + ":" + s.Name)
	}
	CanReturn := false
	if m.HisListIndex > 0 {
		CanReturn = true
	}

	homeIndex := -1
	backIndex := -1
	exitIndex := -1

	println(strconv.Itoa(lenSub) + ":Home")
	homeIndex = lenSub
	if CanReturn {
		println(strconv.Itoa(lenSub+1) + ":Back")
		backIndex = lenSub + 1
	}

	if CanReturn {
		exitIndex = lenSub + 2
	} else {
		exitIndex = lenSub + 1
	}
	println(strconv.Itoa(exitIndex) + ":Exit")
	print("请选择下一步:" + InputStr(1))
	input := 0
	_, err := fmt.Scanln(&input)
	if err != nil {
	}
	switch {
	case input < 0:
		return false
	case input < lenSub:
		m.gotoSub(input)
		return false
	case input == homeIndex:
		m.gotoRoot()
		return false
	case input == backIndex:
		m.returnToLast()
		return false
	case input == exitIndex:
		return true
	default:
		return false
	}
}

func (m *CmdMenu) Run() {
	for {
		if m.do() {
			return
		}
	}
}
