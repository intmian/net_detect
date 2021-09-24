package main

import "time"

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
	print('\b')
	s.status = 0
	s.firstPaint = true
}

type GoProcessBar struct {
	p   ProcessBarOneSign
	end chan int
	exit chan int
}

func (b *GoProcessBar) Init() {
	b.p.Init()
	b.end = make(chan int, 0)
	b.exit = make(chan int,0)
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
	<- b.exit
}
