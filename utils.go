package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

func NumberToBytes(number int, size int) []byte {
	// turn number into hex string
	hexString := fmt.Sprintf("%x", number)

	// split hex string into pairs of two
	if len(hexString)%2 == 1 {
		hexString = "0" + hexString
	}

	var hexPairs []string
	for i := 0; i < len(hexString); i += 2 {
		hexPairs = append(hexPairs, hexString[i:i+2])
	}

	// convert each pair into a byte
	var bytes []byte
	for _, hexPair := range hexPairs {
		var b byte
		fmt.Sscanf(hexPair, "%x", &b)
		bytes = append(bytes, b)
	}

	// pad bytes with 0s
	for i := len(bytes); i < size; i++ {
		bytes = append([]byte{0x00}, bytes...)
	}

	return bytes
}

func GetDeltaTimeBytes(input int) []byte {
	// round input to nearest tick
	ticks := math.Round(float64(input))

	wr := new(bytes.Buffer)

	println(fmt.Sprintf("writing metric ticks: %v", ticks))

	binary.Write(wr, binary.BigEndian, uint16(ticks))

	return wr.Bytes()
}
