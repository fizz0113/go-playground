package main

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"syscall"
	"unsafe"

	"github.com/google/gopacket/afpacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"golang.org/x/sys/unix"
)

// setsockopt 封裝 unix.Syscall6 用於設置套接字選項
func setsockopt(fd, level, name int, val unsafe.Pointer, vallen uintptr) error {
	_, _, errno := unix.Syscall6(
		unix.SYS_SETSOCKOPT,
		uintptr(fd),
		uintptr(level),
		uintptr(name),
		uintptr(val),
		vallen,
		0,
	)
	if errno != 0 {
		return error(errno)
	}
	return nil
}

// getTPacketFD 使用反射獲取 afpacket.TPacket 結構體中的 fd 欄位值
func getTPacketFD(sock *afpacket.TPacket) int {
	// 使用 reflect 獲取私有欄位
	val := reflect.ValueOf(sock).Elem()
	field := val.FieldByName("fd") // 假設私有欄位名稱是 "fd"

	// 檢查欄位是否存在
	if !field.IsValid() {
		panic("Field 'fd' not found in afpacket.TPacket")
	}

	// 確保可以訪問私有欄位
	field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()

	// 返回欄位值作為 int
	return int(field.Int())
}

const (
	PACKET_MASK_ANY  = 0xffffffff
	PACKET_RECV_TYPE = 18
)

func main() {
	const (
		ifaceName = "tap0"
		expr      = "ether[28:2] = 0x888e"
	)

	sock, err := afpacket.NewTPacket(
		afpacket.OptTPacketVersion(1),
		afpacket.OptInterface(ifaceName),
		afpacket.OptFrameSize(2048),
		afpacket.OptNumBlocks(1),
		afpacket.OptBlockSize(2048*1024))
	if err != nil {
		panic(err)
	}
	defer sock.Close()

	// get sockFD
	sockFD := getTPacketFD(sock)

	// setBPF
	bpfList, err := pcap.CompileBPFFilter(layers.LinkTypeEthernet, 65535, expr)
	if err != nil {
		panic(err)
	}
	fprog := unix.SockFprog{
		Len:    uint16(len(bpfList)),
		Filter: (*unix.SockFilter)(unsafe.Pointer(&bpfList[0])),
	}
	err = setsockopt(
		sockFD,
		unix.SOL_SOCKET,
		unix.SO_ATTACH_FILTER,
		unsafe.Pointer(&fprog),
		unix.SizeofSockFprog)
	if err != nil {
		panic(err)
	}

	// setApply
	v := PACKET_MASK_ANY & ^(1 << syscall.PACKET_OUTGOING) & ^(1 << syscall.PACKET_LOOPBACK)
	err = setsockopt(
		sockFD,
		syscall.SOL_PACKET,
		PACKET_RECV_TYPE,
		unsafe.Pointer(&v),
		unix.SizeofInt)
	if err != nil {
		// fmt.Println("set PACKET_RECV_TYPE error:", err)
	}

	// setQdiscBypass
	v = 1
	err = setsockopt(
		sockFD,
		syscall.SOL_PACKET,
		unix.PACKET_QDISC_BYPASS,
		unsafe.Pointer(&v),
		unix.SizeofInt)
	if err != nil {
		// fmt.Println("set PACKET_QDISC_BYPASS error:", err)
	}

	// Read packet from raw socket
	buf := make([]byte, 10240)
	for {
		ci, err := sock.ReadPacketDataTo(buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}
		if ci.Length <= 0 || ci.Length > 10240 {
			fmt.Println("invalid length:", ci.Length)
			continue
		}
		fmt.Printf("Read %d bytes from %s\n", ci.Length, ifaceName)
	}

	// ICMP echo request packet in hex format
	hexstrICMP := "080027e29fa6080027fc6ac90800450001c4b5d0007a4001bcea0201010202010101c8c9cacbcccdcecfd0d1d2d3d4d5d6d7d8d9dadbdcdddedfe0e1e2e3e4e5e6e7e8e9eaebecedeeeff0f1f2f3f4f5f6f7f8f9fafbfcfdfeff000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4c5c6c7c8c9cacbcccdcecfd0d1d2d3d4d5d6d7d8d9dadbdcdddedfe0e1e2e3e4e5e6e7e8e9eaebecedeeeff0f1f2f3f4f5f6f7f8f9fafbfcfdfeff000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f606162636465666768696a6b6c6d6e6f7071727374757677"

	// Decode hex string to bytes
	bytes, err := hex.DecodeString(hexstrICMP)
	if err != nil {
		panic(err)
	}

	// Write packet to raw socket
	sock.WritePacketData(bytes)
}
