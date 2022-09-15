package main

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
)

type PingSupervisorStatus struct {
	IsConnected  bool
	LastSendTime time.Time
	LastRecvTime time.Time
}

type PingSupervisor struct {
	pinger             *ping.Pinger
	status             PingSupervisorStatus
	connectChangeFuncs []func(bool)
}

func NewPingSupervisor(addr string, interval time.Duration) (*PingSupervisor, error) {
	p := ping.New(addr)
	p.Interval = interval
	x := &PingSupervisor{
		pinger: p,
		status: PingSupervisorStatus{},
	}
	p.OnSend = x.onSend
	p.OnRecv = x.onRecv
	return x, nil
}

func (x *PingSupervisor) onSend(pkt *ping.Packet) {
	// fmt.Printf("[send] icmp_seq=%d\n", pkt.Seq)
	diff := x.status.LastSendTime.UnixNano() - x.status.LastRecvTime.UnixNano()
	if diff > int64(x.pinger.Interval) {
		x.onConnectChange(false)
	}
	x.status.LastSendTime = time.Now()
}

func (x *PingSupervisor) onRecv(pkt *ping.Packet) {
	// fmt.Printf("[recv] %d bytes from %s: icmp_seq=%d time=%v\n", pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	x.onConnectChange(true)
	x.status.LastRecvTime = time.Now()
}

func (x *PingSupervisor) onConnectChange(state bool) {
	if x.status.IsConnected != state {
		// fmt.Println("[conn] ping state changed:", state)
		x.status.IsConnected = state
		for _, f := range x.connectChangeFuncs {
			f(state)
		}
	}
}

func (x *PingSupervisor) Start() {
	go func() {
		for {
			err := x.pinger.Resolve()
			if err != nil {
				fmt.Println("ping resolve error:", err)
				time.Sleep(x.pinger.Interval)
				continue
			}
			break
		}
		x.pinger.Run()
	}()
}

func (x *PingSupervisor) Close() {
	x.pinger.Stop()
}

func (x *PingSupervisor) RegisterConnectChangeFunc(f func(bool)) {
	x.connectChangeFuncs = append(x.connectChangeFuncs, f)
}

func (x *PingSupervisor) GetStatus() *PingSupervisorStatus {
	return &x.status
}

func (x *PingSupervisor) GetStatistics() *ping.Statistics {
	return x.pinger.Statistics()
}

func main() {
	const (
		addr  = "8.8.8.8"
		delay = 1000 // milliseconds
	)
	pingSupervisor, err := NewPingSupervisor(addr, time.Duration(delay)*time.Millisecond)
	if err != nil {
		// panic(err)
	}

	triggerXXXFunc := func(state bool) {
		if state {
			println("Connected: update or trigger something")
		} else {
			println("Disconnected: update or trigger something")
		}
	}
	pingSupervisor.RegisterConnectChangeFunc(triggerXXXFunc)

	pingSupervisor.Start()
	defer pingSupervisor.Close()

	for {
		status := pingSupervisor.GetStatus()
		fmt.Printf("ping Status: %+v\n", status)
		time.Sleep(time.Second)
	}
}
