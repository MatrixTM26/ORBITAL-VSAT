package l3

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/0xTM7/demon/internal/config"
	"github.com/0xTM7/demon/internal/utils"
)

func BuildICMP(PayloadSize int) []byte {
	Payload := utils.RandPayload(PayloadSize)
	Pkt := make([]byte, 8+len(Payload))
	Pkt[0] = 8
	Pkt[1] = 0
	binary.BigEndian.PutUint16(Pkt[4:6], uint16(utils.RandInt(1, 65535)))
	binary.BigEndian.PutUint16(Pkt[6:8], 1)
	copy(Pkt[8:], Payload)
	Cksum := utils.Checksum(Pkt)
	binary.BigEndian.PutUint16(Pkt[2:4], Cksum)
	return Pkt
}

func RunICMPWorker(Cfg *config.DemonConfig, Stop <-chan struct{}) {
	Conn, Err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if Err != nil {
		FallbackPing(Cfg, Stop)
		return
	}
	defer Conn.Close()

	DstAddr, _ := net.ResolveIPAddr("ip4", Cfg.IP)

	for {
		select {
		case <-Stop:
			return
		default:
			Pkt := BuildICMP(utils.RandInt(64, 1400))
			N, Err := Conn.WriteTo(Pkt, DstAddr)
			if Err != nil {
				Cfg.Stats.Failed.Add(1)
			} else {
				Cfg.Stats.Requests.Add(1)
				Cfg.Stats.BytesSent.Add(int64(N))
				Cfg.Stats.Success.Add(1)
			}
		}
	}
}

func FallbackPing(Cfg *config.DemonConfig, Stop <-chan struct{}) {
	Addr := Cfg.IP + ":0"
	for {
		select {
		case <-Stop:
			return
		default:
			Conn, Err := net.DialTimeout("udp", Addr, 2*time.Second)
			if Err != nil {
				Cfg.Stats.Failed.Add(1)
				time.Sleep(10 * time.Millisecond)
				continue
			}
			Conn.Close()
			Cfg.Stats.Requests.Add(1)
			Cfg.Stats.Success.Add(1)
		}
	}
}
