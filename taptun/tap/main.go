package main

import (
	"encoding/hex"
	"fmt"
	"net"

	"github.com/songgao/water"
)

func main() {
	const (
		ifaceName = "tap99"
	)

	// Create TAP interface
	cfg := water.Config{}
	cfg.Name = ifaceName
	cfg.DeviceType = water.TAP
	tapIface, err := water.New(cfg)
	if err != nil {
		panic(err)
	}
	defer tapIface.Close()

	// Verify interface creation
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		panic(err)
	}
	_ = iface

	// Read packets from the TAP interface
	buf := make([]byte, 10240)
	for {
		bufn, err := tapIface.Read(buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}
		if bufn <= 0 || bufn > 10240 {
			fmt.Println("invalid length:", bufn)
			continue
		}
		fmt.Printf("Read %d bytes from %s\n", bufn, tapIface.Name())
	}

	// ICMP echo request packet in hex format
	hexstrICMP := "080027e29fa6080027fc6ac90800450001c4b5d0007a4001bcea0201010202010101c8c9cacbcccdcecfd0d1d2d3d4d5d6d7d8d9dadbdcdddedfe0e1e2e3e4e5e6e7e8e9eaebecedeeeff0f1f2f3f4f5f6f7f8f9fafbfcfdfeff000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4c5c6c7c8c9cacbcccdcecfd0d1d2d3d4d5d6d7d8d9dadbdcdddedfe0e1e2e3e4e5e6e7e8e9eaebecedeeeff0f1f2f3f4f5f6f7f8f9fafbfcfdfeff000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f606162636465666768696a6b6c6d6e6f7071727374757677"

	// Decode hex string to bytes
	bytes, err := hex.DecodeString(hexstrICMP)
	if err != nil {
		panic(err)
	}

	// Write the ICMP packet to the TAP interface
	tapIface.Write(bytes)
}
