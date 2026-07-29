package l7

import (
	"fmt"
	"math/rand"

	"github.com/0xTM7/demon/internal/config"
	"github.com/0xTM7/demon/internal/utils"
)

var AcceptHeaders = []string{
	"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
	"application/json,text/plain,*/*",
	"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
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

func BuildHeaders(Cfg *config.DemonConfig, Path string) map[string]string {
	H := map[string]string{
		"Host":             Cfg.Host,
		"User-Agent":       utils.Pick(Cfg.UserAgents),
		"Accept":           utils.Pick(AcceptHeaders),
		"Accept-Encoding":  utils.Pick(AcceptEncodings),
		"Accept-Language":  utils.Pick(AcceptLanguages),
		"Cache-Control":    utils.Pick(CacheControls),
		"Pragma":           "no-cache",
		"Connection":       "keep-alive",
		"Referer":          utils.Pick(Cfg.Referers),
		"X-Forwarded-For":  utils.RandIP(),
		"X-Real-IP":        utils.RandIP(),
		"X-Request-ID":     utils.RandString(32),
	}

	if rand.Intn(2) == 0 {
		H["Upgrade-Insecure-Requests"] = "1"
	}
	if rand.Intn(3) == 0 {
		H["DNT"] = "1"
	}
	if rand.Intn(4) == 0 {
		H["X-Requested-With"] = "XMLHttpRequest"
	}

	for K, V := range Cfg.Headers {
		H[K] = V
	}

	return H
}

func CacheBypassPath(Path string) string {
	return Path + utils.RandQuery(utils.RandInt(1, 4)) + fmt.Sprintf("&_=%d", rand.Int63())
}
