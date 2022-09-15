package main

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
)

func main() {
	const (
		addr  = "8.8.8.8"
		delay = 1000 // milliseconds
	)
	p, err := ping.NewPinger(addr)
	if err != nil {
		panic(err)
	}
	p.Interval = time.Duration(delay) * time.Millisecond
	p.OnSend = func(pkt *ping.Packet) {
		fmt.Printf("[send] icmp_seq=%d\n", pkt.Seq)
	}
	p.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("[recv] %d bytes from %s: icmp_seq=%d time=%v\n", pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}
	p.Run()
}
