package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func parseDSADataToRawData(dsaData []byte) (int32, []byte) {
	// 0x0000:  0180 c200 0003 b827 eb84 0a82 0027 1001
	// 0x0010:  8003 c00d 8000 0018 0000 0000 888e 0301
	// 0x0020:  0000 0000 0000 0000 0000 0000 0000 0000
	// 0x0030:  0000 0000 0000 0000 0000 0000 0000 0000
	// 0x0040:  0000 0000 0000 0000 0000 0000

	// 0027 1001 8003 c00d
	// 8000 0018 0000 0000
	tagData := dsaData[12:28]
	data := make([]byte, len(dsaData)-len(tagData))
	copy(data[:12], dsaData[:12])
	copy(data[12:], dsaData[28:])

	word0 := binary.BigEndian.Uint32(tagData[0:4])
	word1 := binary.BigEndian.Uint32(tagData[4:8])
	word2 := binary.BigEndian.Uint32(tagData[8:12])
	// word3 := binary.BigEndian.Uint32(tagData[12:16])

	/* word 0: this is 4 byte format */
	cpuCode := int((word0>>15)&0x0E) + int((word0>>12)&0x01) // b[18:16] + b[12]
	/* word 1 */
	isExtend := int((word1 >> 31) & 0x01) // b[31]

	portNo := int32(0)
	if cpuCode == 0xF && isExtend == 1 {
		portNo = int32((word2>>13)&0x80) + int32((word1>>5)&0x60) + int32((word0>>19)&0x1f)
	} else if cpuCode == 0xF && isExtend == 0 {
		portNo = int32((word1>>5)&0x60) + int32((word0>>19)&0x1f)
	} else {
		portNo = int32((word0 >> 19) & 0x1f)
	}
	return portNo, data
}

func parseRawDataToDSAData(portNo int32, rawData []byte) []byte {
	word0 := uint32(0x00000000)
	word1 := uint32(0x00000000)
	word2 := uint32(0x00000000)
	word3 := uint32(0x00000000)

	word0 |= uint32(0x1) << 30
	word0 |= uint32(portNo) << 19

	word1 |= uint32(0x1) << 31
	word1 |= uint32(portNo) << 10

	word2 |= uint32(0x1) << 31
	word2 |= uint32(portNo) << 28

	word3 |= uint32(0x1) << 31
	word3 |= uint32(portNo) << 23

	tagData := make([]byte, 16)
	binary.BigEndian.PutUint32(tagData[0:4], word0)
	binary.BigEndian.PutUint32(tagData[4:8], word1)
	binary.BigEndian.PutUint32(tagData[8:12], word2)
	binary.BigEndian.PutUint32(tagData[12:16], word3)

	data := make([]byte, len(rawData)+len(tagData))
	copy(data[:12], rawData[:12])
	copy(data[12:18], tagData)
	copy(data[28:], rawData[12:])
	return data
}

func main() {
	// 0x0000:  0180 c200 0003 b827 eb84 0a82 0027 1001
	// 0x0010:  8003 c00d 8000 0018 0000 0000 888e 0301
	// 0x0020:  0000 0000 0000 0000 0000 0000 0000 0000
	// 0x0030:  0000 0000 0000 0000 0000 0000 0000 0000
	// 0x0040:  0000 0000 0000 0000 0000 0000
	dsaData := []byte{
		0x01, 0x80, 0xc2, 0x00, 0x00, 0x03, 0xb8, 0x27,
		0xeb, 0x84, 0x0a, 0x82, 0x00, 0x27, 0x10, 0x01,
		0x80, 0x03, 0xc0, 0x0d, 0x80, 0x00, 0x00, 0x18,
		0x00, 0x00, 0x00, 0x00, 0x88, 0x8e, 0x03, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}
	fmt.Println("Original DSA Data:")
	fmt.Println(hex.Dump(dsaData))

	portNo, rawData := parseDSADataToRawData(dsaData)
	fmt.Println("Extracted Port No:", portNo)
	fmt.Println("Extracted Raw Data:")
	fmt.Println(hex.Dump(rawData))

	newDSAData := parseRawDataToDSAData(portNo, rawData)
	fmt.Println("Reconstructed DSA Data:")
	fmt.Println(hex.Dump(newDSAData))
}
