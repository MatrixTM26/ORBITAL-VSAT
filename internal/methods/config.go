package methods

import "sync/atomic"

type Stats struct {
	RequestsCount atomic.Int64
	BytesSent     atomic.Int64
}

type Config struct {
	IP         string
	Port       int
	Host       string
	Path       string
	Scheme     string
	Protocol   string
	JAProfile  string
	Method     string
	UserAgents []string
	Referers   []string
	Stats      *Stats
}
