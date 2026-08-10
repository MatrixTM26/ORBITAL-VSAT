package fingerprint

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

type JA3Profile struct {
	Ciphers []uint16
	Curves  []tls.CurveID
	ALPN    []string
}

var Profiles = map[string]JA3Profile{
	"Chrome": {
		Ciphers: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		},
		Curves: []tls.CurveID{tls.X25519, tls.CurveP256, tls.CurveP384},
	},
	"Firefox": {
		Ciphers: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
		Curves: []tls.CurveID{tls.X25519, tls.CurveP256, tls.CurveP384, tls.CurveP521},
	},
	"Safari": {
		Ciphers: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		Curves: []tls.CurveID{tls.CurveP256, tls.X25519, tls.CurveP384},
	},
}

func CreateJA3Socket(IP string, Port int, Scheme, Host, Protocol, JAProfile string) (net.Conn, error) {
	Addr := fmt.Sprintf("%s:%d", IP, Port)
	RawConn, Err := net.DialTimeout("tcp", Addr, 3*time.Second)
	if Err != nil {
		return nil, Err
	}

	TCP, _ := RawConn.(*net.TCPConn)
	if TCP != nil {
		TCP.SetNoDelay(true)
		TCP.SetKeepAlive(true)
		TCP.SetKeepAlivePeriod(30 * time.Second)
	}

	if Scheme != "https" {
		return RawConn, nil
	}

	Profile, Exists := Profiles[JAProfile]
	if !Exists {
		Profile = Profiles["Chrome"]
	}

	ALPN := []string{"http/1.1"}
	if Protocol == "H2" {
		ALPN = []string{"h2", "http/1.1"}
	}

	TLSCfg := &tls.Config{
		ServerName:         Host,
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS13,
		CipherSuites:       Profile.Ciphers,
		CurvePreferences:   Profile.Curves,
		NextProtos:         ALPN,
	}

	TLSConn := tls.Client(RawConn, TLSCfg)
	TLSConn.SetDeadline(time.Now().Add(5 * time.Second))
	if Err := TLSConn.Handshake(); Err != nil {
		RawConn.Close()
		return nil, Err
	}
	TLSConn.SetDeadline(time.Time{})
	return TLSConn, nil
}
