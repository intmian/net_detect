package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type ProcessBarOneSign struct {
	status     int
	signs      []string
	firstPaint bool
}

func (s *ProcessBarOneSign) Init() {
	s.status = 0
	s.firstPaint = true
	s.signs = []string{"↑", "→", "↓", "←"}
}

func (s *ProcessBarOneSign) Next() {
	if !s.firstPaint {
		print("\b")
	} else {
		s.firstPaint = false
	}
	print(s.signs[s.status])
	s.status++
	s.status = s.status % len(s.signs)
}

func (s *ProcessBarOneSign) End() {
	if s.firstPaint {
		return
	}
	print("\b")
	s.status = 0
	s.firstPaint = true
}

type GoProcessBar struct {
	p    ProcessBarOneSign
	end  chan int
	exit chan int
}

func (b *GoProcessBar) Init() {
	b.p.Init()
	b.end = make(chan int, 0)
	b.exit = make(chan int, 0)
}

func (b *GoProcessBar) Run() {
	go func() {
		for {
			select {
			case <-b.end:
				b.p.End()
				b.exit <- 1
				return
			case <-time.After(300 * time.Millisecond):
				b.p.Next()
			}
		}
	}()
}

func (b *GoProcessBar) Stop() {
	b.end <- 1
	<-b.exit
}

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
	lenSub := len(m.now.SubMenu)
	if lenSub == 0 {
		// 跑到叶节点了，就执行这个的逻辑，停一下，再接着跑
		m.now.F()
		m.stop()
		m.returnToLast()
		return false
	}

	m.clear()
	if m.now == nil {
		return
	}
	for i, s := range m.now.SubMenu {
		println(i, ":", s, "。")
	}
	CanReturn := false
	if m.HisListIndex > 0 {
		CanReturn = true
	}

	homeIndex := -1
	backIndex := -1
	exitIndex := -1

	println("Home:", lenSub)
	homeIndex = lenSub
	if CanReturn {
		println("Back:", lenSub+1)
		backIndex = lenSub + 1
	}
	println("Exit:")
	if CanReturn {
		exitIndex = lenSub + 2
	} else {
		exitIndex = lenSub + 1
	}
	print("请选择下一步:__\b")
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
		return  true
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