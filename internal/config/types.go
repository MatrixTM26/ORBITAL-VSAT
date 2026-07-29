package config

import "sync/atomic"

type Layer string
type Method string
type Protocol string
type JA3Profile string

const (
	LayerL7 Layer = "L7"
	LayerL4 Layer = "L4"
	LayerL3 Layer = "L3"
)

const (
	MethodGET     Method = "GET"
	MethodPOST    Method = "POST"
	MethodPUT     Method = "PUT"
	MethodPATCH   Method = "PATCH"
	MethodDELETE  Method = "DELETE"
	MethodHEAD    Method = "HEAD"
	MethodOPTIONS Method = "OPTIONS"
	MethodTRACE   Method = "TRACE"
	MethodCONNECT Method = "CONNECT"
	MethodRANDOM  Method = "RANDOM"

	MethodTCP     Method = "TCP"
	MethodSYN     Method = "SYN"
	MethodACK     Method = "ACK"
	MethodRST     Method = "RST"
	MethodFIN     Method = "FIN"
	MethodXMAS    Method = "XMAS"
	MethodUDP     Method = "UDP"
	MethodUDPFRAG Method = "UDP-FRAG"
	MethodDNSAMP  Method = "DNS-AMP"
	MethodNTPAMP  Method = "NTP-AMP"

	MethodICMP Method = "ICMP"
)

const (
	ProtoH1 Protocol = "H1"
	ProtoH2 Protocol = "H2"
)

const (
	JA3Chrome  JA3Profile = "chrome"
	JA3Firefox JA3Profile = "firefox"
	JA3Safari  JA3Profile = "safari"
)

type Stats struct {
	Requests  atomic.Int64
	BytesSent atomic.Int64
	Success   atomic.Int64
	Failed    atomic.Int64
}

type DemonConfig struct {
	Target      string
	Host        string
	IP          string
	Port        int
	Scheme      string
	Path        string
	Layer       Layer
	Method      Method
	Protocol    Protocol
	JA3         JA3Profile
	Threads     int
	Duration    int
	ClusterMode bool
	Workers     int
	UserAgents  []string
	Referers    []string
	Body        string
	Headers     map[string]string
	Stats       *Stats
}

var DefaultUserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:125.0) Gecko/20100101 Firefox/125.0",
	"Mozilla/5.0 (X11; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_4_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_4_1 like Mac OS X) AppleWebKit/605.1.15 Version/17.4.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 Chrome/124.0.0.0 Mobile Safari/537.36",
}

var DefaultReferers = []string{
	"https://www.google.com",
	"https://www.bing.com",
	"https://duckduckgo.com",
	"https://search.yahoo.com",
	"https://www.baidu.com",
	"https://yandex.com",
	"https://www.reddit.com",
	"https://twitter.com",
	"https://github.com",
	"https://www.facebook.com",
	"https://www.youtube.com",
}
