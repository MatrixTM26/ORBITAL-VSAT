package random

import (
	"encoding/binary"
	"fmt"
	"math/rand"
)

const Charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandomString(Length int) string {
	B := make([]byte, Length)
	for I := range B {
		B[I] = Charset[rand.Intn(len(Charset))]
	}
	return string(B)
}

func RandomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(254)+1,
		rand.Intn(254)+1,
		rand.Intn(254)+1,
		rand.Intn(254)+1,
	)
}

func RandomInt(Min, Max int) int {
	return Min + rand.Intn(Max-Min+1)
}

func RandomBytes(N int) []byte {
	B := make([]byte, N)
	rand.Read(B)
	return B
}

func CalculateChecksum(Data []byte) uint16 {
	Sum := 0
	for I := 0; I+1 < len(Data); I += 2 {
		Sum += int(binary.BigEndian.Uint16(Data[I : I+2]))
	}
	if len(Data)%2 != 0 {
		Sum += int(Data[len(Data)-1]) << 8
	}
	for Sum>>16 != 0 {
		Sum = (Sum & 0xffff) + (Sum >> 16)
	}
	return ^uint16(Sum)
}
