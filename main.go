package main

import (
	"net_detect/logic"
	"net_detect/net"
)

func main() {
	h := net.NetHttping{}
	h.Init()
	logic.StartNetMenu(&h)
	h.Finalize()
}
