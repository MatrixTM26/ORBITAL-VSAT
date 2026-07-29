package utils

import (
	"math/rand"
	"net"
	"strings"
)

const Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandString(N int) string {
	B := make([]byte, N)
	for I := range B {
		B[I] = Charset[rand.Intn(len(Charset))]
	}
	return string(B)
}

func RandInt(Min, Max int) int {
	return Min + rand.Intn(Max-Min+1)
}

func Pick(S []string) string {
	if len(S) == 0 {
		return ""
	}
	return S[rand.Intn(len(S))]
}

func RandIP() string {
	return net.IPv4(
		byte(rand.Intn(223)+1),
		byte(rand.Intn(255)),
		byte(rand.Intn(255)),
		byte(rand.Intn(254)+1),
	).String()
}

func RandQuery(N int) string {
	Parts := make([]string, N)
	for I := range Parts {
		Parts[I] = RandString(RandInt(3, 8)) + "=" + RandString(RandInt(4, 12))
	}
	return "?" + strings.Join(Parts, "&")
}

func RandPayload(Size int) []byte {
	B := make([]byte, Size)
	rand.Read(B)
	return B
}

func Checksum(Data []byte) uint16 {
	Sum := 0
	for I := 0; I+1 < len(Data); I += 2 {
		Sum += int(Data[I])<<8 | int(Data[I+1])
	}
	if len(Data)%2 != 0 {
		Sum += int(Data[len(Data)-1]) << 8
	}
	for Sum>>16 != 0 {
		Sum = (Sum & 0xffff) + (Sum >> 16)
	}
	return ^uint16(Sum)
}
