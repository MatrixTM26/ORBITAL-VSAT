package main

var DefaultUserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 Chrome/124.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:125.0) Gecko/20100101 Firefox/125.0",
	"Mozilla/5.0 (X11; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_4_1) AppleWebKit/605.1.15 Version/17.4.1 Safari/605.1.15",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_4_1 like Mac OS X) AppleWebKit/605.1.15 Version/17.4.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 Chrome/124.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/124.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Linux; Android 13; SM-S918B) AppleWebKit/537.36 Chrome/124.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 17_4_1 like Mac OS X) AppleWebKit/605.1.15 Version/17.4.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0",
}

var DefaultReferers = []string{
	"https://www.google.com/search?q=" ,
	"https://www.bing.com/search?q=",
	"https://duckduckgo.com/?q=",
	"https://search.yahoo.com/search?p=",
	"https://www.baidu.com/s?wd=",
	"https://yandex.com/search/?text=",
	"https://www.reddit.com/",
	"https://twitter.com/",
	"https://github.com/",
	"https://www.facebook.com/",
	"https://www.youtube.com/",
	"https://news.ycombinator.com/",
}

var L7Methods = map[string]bool{
	"GET": true, "POST": true, "PUT": true, "HEAD": true,
	"DELETE": true, "PATCH": true, "OPTIONS": true,
	"CONNECT": true, "TRACE": true, "RANDOM": true,
	"SLOWLORIS": true, "SLOWPOST": true, "SLOWREAD": true,
	"RUDY": true,
}

var L4Methods = map[string]bool{
	"TCP": true, "SYN": true, "ACK": true, "RST": true,
	"FIN": true, "XMAS": true, "PSH": true, "URG": true,
	"NULL": true, "SYNACK": true, "UDP": true, "UDP-FRAG": true,
	"DNS-AMP": true, "NTP-AMP": true, "SSDP-AMP": true,
	"MEMCACHED": true, "CHARGEN": true, "FRAGGLE": true,
}

var L3Methods = map[string]bool{
	"ICMP": true, "SMURF": true, "PING": true,
}
