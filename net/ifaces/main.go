package main

import (
	"fmt"
	"net"
)

func getAllInterfaces() ([]net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	return ifaces, nil
}

func findInterfaceByName(name string) (*net.Interface, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}
	return iface, nil
}

func main() {

	ifaces, err := getAllInterfaces()
	if err != nil {
		panic(err)
	}
	for _, iface := range ifaces {
		fmt.Printf("%7s: %d\n", "Index", iface.Index)
		fmt.Printf("%7s: %s\n", "  Name", iface.Name)
		fmt.Printf("%7s: %s\n", "  MAC", iface.HardwareAddr)
		fmt.Printf("%7s: %d\n", "  MTU", iface.MTU)
		fmt.Printf("%7s: %s\n", "  Flag", iface.Flags)
	}

	ifaceName := "eth0"
	iface, err := findInterfaceByName(ifaceName)
	if err != nil {
		fmt.Printf("Interface %s not found: %v\n", ifaceName, err)
	} else {
		fmt.Printf("Found interface %s: %+v\n", ifaceName, iface)
	}

}
