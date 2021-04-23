package main

import (
	"fmt"
	"log"
	"net"
	"syscall"
	"unsafe"
)

func newNetLinkFD() int {
	fd, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, syscall.NETLINK_ROUTE)
	if err != nil {
		log.Fatal(err)
	}

	sa := &syscall.SockaddrNetlink{
		Family: syscall.AF_NETLINK,
		// syscall.RTM_GETLINK
		Groups: syscall.RTNLGRP_LINK, /* | syscall.RTNLGRP_IPV4_IFADDR | syscall.RTNLGRP_IPV4_ROUTE */
		Pid:    0,
	}

	err = syscall.Bind(fd, sa)
	if err != nil {
		log.Fatal(err)
	}
	return fd
}

func makeFlags(rawFlags uint32) net.Flags {
	var f net.Flags
	if rawFlags&syscall.IFF_UP != 0 {
		f |= net.FlagUp
	}
	if rawFlags&syscall.IFF_BROADCAST != 0 {
		f |= net.FlagBroadcast
	}
	if rawFlags&syscall.IFF_LOOPBACK != 0 {
		f |= net.FlagLoopback
	}
	if rawFlags&syscall.IFF_POINTOPOINT != 0 {
		f |= net.FlagPointToPoint
	}
	if rawFlags&syscall.IFF_MULTICAST != 0 {
		f |= net.FlagMulticast
	}
	return f
}

func updateInterfaceStatus(ifimsg *syscall.IfInfomsg, nlattrs []syscall.NetlinkRouteAttr) {
	netInterface := &net.Interface{Index: int(ifimsg.Index), Flags: makeFlags(ifimsg.Flags)}
	for _, attr := range nlattrs {
		switch attr.Attr.Type {
		case syscall.IFLA_ADDRESS:
			netInterface.HardwareAddr = attr.Value[:]
		case syscall.IFLA_IFNAME:
			netInterface.Name = string(attr.Value[:len(attr.Value)-1])
		case syscall.IFLA_MTU:
			netInterface.MTU = int(*(*uint32)(unsafe.Pointer(&attr.Value[:4][0])))
		}
	}

	if netInterface.MTU == 0 || len(netInterface.HardwareAddr) == 0 {
		return
	}
	// bs, _ := json.Marshal(netInterface)
	log.Println(netInterface)
}

func main() {
	fmt.Println("Network Interface Link Change Listening")
	log.SetFlags(log.Llongfile)

	fd := newNetLinkFD()
	buf := make([]byte, 2048)
	for {

		bufn, err := syscall.Read(fd, buf)
		if err != nil {
			log.Println(err)
			continue
		}
		if bufn <= 0 {
			continue
		}
		nlmsgs, err := syscall.ParseNetlinkMessage(buf[:bufn])
		if err != nil {
			log.Println(err)
			continue
		}
		for _, nlmsg := range nlmsgs {
			switch nlmsg.Header.Type {
			case syscall.NLMSG_DONE:
			case syscall.NLMSG_ERROR:
			case syscall.RTM_DELLINK:
			case syscall.RTM_DELADDR, syscall.RTM_NEWADDR:
			case syscall.RTM_DELROUTE, syscall.RTM_NEWROUTE:
			case syscall.RTM_NEWLINK:
				ifimsg := (*syscall.IfInfomsg)(unsafe.Pointer(&nlmsg.Data[0]))
				nlattrs, err := syscall.ParseNetlinkRouteAttr(&nlmsg)
				if err != nil {
					log.Println(err)
					continue
				}
				updateInterfaceStatus(ifimsg, nlattrs)
			}
		}

	}
}
