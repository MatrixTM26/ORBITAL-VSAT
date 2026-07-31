package main

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

type Resolved struct {
	Target string
	Scheme string
	Host   string
	IP     string
	Port   int
	Path   string
}

func Resolve(Target string) (Resolved, error) {
	Raw := Target
	if !strings.HasPrefix(Raw, "http://") && !strings.HasPrefix(Raw, "https://") {
		Raw = "https://" + Raw
	}

	U, Err := url.Parse(Raw)
	if Err != nil {
		return Resolved{}, fmt.Errorf("invalid URL: %w", Err)
	}

	R := Resolved{
		Target: Raw,
		Scheme: U.Scheme,
		Host:   U.Hostname(),
		Path:   U.Path,
	}
	if R.Path == "" {
		R.Path = "/"
	}

	if U.Port() != "" {
		fmt.Sscanf(U.Port(), "%d", &R.Port)
	} else {
		if R.Scheme == "https" {
			R.Port = 443
		} else {
			R.Port = 80
		}
	}

	Addrs, Err := net.LookupHost(R.Host)
	if Err != nil {
		return Resolved{}, fmt.Errorf("cannot resolve %s: %w", R.Host, Err)
	}
	R.IP = Addrs[0]
	return R, nil
}

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
