package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/MatrixTM26/vsat/internal/color"
)

type DemonArgs struct {
	Target      string
	Method      string
	Layer       string
	Protocol    string
	JAProfile   string
	Threads     int
	Duration    int
	ClusterMode bool
}

func PrintHelp() {
	fmt.Println(color.B(color.C("  AVAILABLE METHODS", color.Cyan)))
	fmt.Println()
	fmt.Println(color.B(color.C("  Layer 7", color.Yellow)))
	fmt.Printf("  %-14s  HTTP flood (GET POST PUT PATCH DELETE HEAD OPTIONS TRACE CONNECT)\n", color.C("HTTP", color.Green))
	fmt.Printf("  %-14s  Random method per request\n", color.C("RANDOM", color.Green))
	fmt.Printf("  %-14s  Slow headers flood\n", color.C("SLOWLORIS", color.Green))
	fmt.Printf("  %-14s  Slow POST body flood\n", color.C("SLOWPOST", color.Green))
	fmt.Printf("  %-14s  Slow response read\n", color.C("SLOWREAD", color.Green))
	fmt.Printf("  %-14s  R-U-Dead-Yet byte-per-byte POST\n", color.C("RUDY", color.Green))
	fmt.Println()
	fmt.Println(color.B(color.C("  Layer 4", color.Yellow)))
	fmt.Printf("  %-14s  TCP connection flood\n", color.C("TCP", color.Green))
	fmt.Printf("  %-14s  Raw TCP flag floods\n", color.C("SYN ACK RST FIN", color.Green))
	fmt.Printf("  %-14s  Raw TCP flag floods\n", color.C("XMAS PSH URG NULL", color.Green))
	fmt.Printf("  %-14s  SYN+ACK flood\n", color.C("SYNACK", color.Green))
	fmt.Printf("  %-14s  UDP datagram flood\n", color.C("UDP", color.Green))
	fmt.Printf("  %-14s  UDP fragmented flood\n", color.C("UDP-FRAG", color.Green))
	fmt.Printf("  %-14s  DNS amplification\n", color.C("DNS-AMP", color.Green))
	fmt.Printf("  %-14s  NTP monlist amplification\n", color.C("NTP-AMP", color.Green))
	fmt.Printf("  %-14s  SSDP amplification\n", color.C("SSDP-AMP", color.Green))
	fmt.Printf("  %-14s  Memcached amplification\n", color.C("MEMCACHED", color.Green))
	fmt.Printf("  %-14s  CHARGEN amplification\n", color.C("CHARGEN", color.Green))
	fmt.Printf("  %-14s  UDP broadcast flood\n", color.C("FRAGGLE", color.Green))
	fmt.Println()
	fmt.Println(color.B(color.C("  Layer 3", color.Yellow)))
	fmt.Printf("  %-14s  Raw ICMP echo flood\n", color.C("ICMP", color.Green))
	fmt.Printf("  %-14s  ICMP broadcast flood\n", color.C("SMURF", color.Green))
	fmt.Printf("  %-14s  Controlled ping flood\n", color.C("PING", color.Green))
	fmt.Println()
}

func ask(R *bufio.Reader, Label, Default string) string {
	if Default != "" {
		fmt.Printf("  %s %s[%s%s%s] %s»%s ",
			color.B(color.C(Label, color.BrightGreen)),
			color.Dim, color.Cyan, Default, color.Reset+color.Dim,
			color.Yellow, color.Cyan,
		)
	} else {
		fmt.Printf("  %s %s»%s ",
			color.B(color.C(Label, color.BrightGreen)),
			color.Yellow, color.Cyan,
		)
	}
	Line, _ := R.ReadString('\n')
	fmt.Print(color.Reset)
	return strings.TrimSpace(Line)
}

func RunPrompt() DemonArgs {
	R := bufio.NewReader(os.Stdin)
	Args := DemonArgs{
		Protocol:  "H1",
		JAProfile: "Chrome",
	}

	fmt.Printf("  %s  %s  %s\n",
		color.B(color.C("Y", color.Green))+" start",
		color.B(color.C("N", color.Red))+" exit",
		color.B(color.C("H", color.Cyan))+" help",
	)
	fmt.Println()

	for {
		Raw := strings.ToLower(ask(R, "Continue? [Y/N/H]", ""))
		switch Raw {
		case "y":
		case "n":
			fmt.Println(color.C("  * bye", color.Dim))
			os.Exit(0)
		case "h":
			fmt.Println()
			PrintHelp()
		default:
			fmt.Println(color.C("  ! enter Y, N, or H", color.Red))
			continue
		}
		if Raw == "y" {
			break
		}
	}

	fmt.Println()
	fmt.Println(color.B(color.C("  CONFIGURATION", color.Red)))
	fmt.Println()

	for {
		Raw := ask(R, "TARGET", "")
		if Raw == "" {
			fmt.Println(color.C("  ! target is required", color.Red))
			continue
		}
		Args.Target = Raw
		break
	}

	var MethodRaw string
	for {
		Raw := ask(R, "METHOD", "GET")
		MethodRaw = strings.ToUpper(Raw)
		if L7Methods[MethodRaw] || L4Methods[MethodRaw] || L3Methods[MethodRaw] {
			break
		}
		fmt.Println(color.C("  ! unknown method: "+MethodRaw, color.Red))
	}
	Args.Method = MethodRaw

	if L7Methods[MethodRaw] {
		Args.Layer = "L7"

		for {
			Raw := strings.ToUpper(ask(R, "PROTOCOL [H1/H2]", "H1"))
			if Raw == "H1" || Raw == "H2" {
				if Raw == "H2" {
					Args.Protocol = "H2"
				}
				break
			}
			fmt.Println(color.C("  ! choose H1 or H2", color.Red))
		}

		for {
			Raw := ask(R, "JA3 [Chrome/Firefox/Safari]", "Chrome")
			switch strings.ToLower(Raw) {
			case "chrome", "":
				Args.JAProfile = "Chrome"
				goto DoneJA3
			case "firefox":
				Args.JAProfile = "Firefox"
				goto DoneJA3
			case "safari":
				Args.JAProfile = "Safari"
				goto DoneJA3
			}
			fmt.Println(color.C("  ! choose Chrome, Firefox, or Safari", color.Red))
		}
	DoneJA3:

	} else if L4Methods[MethodRaw] {
		Args.Layer = "L4"
	} else {
		Args.Layer = "L3"
	}

	for {
		Raw := ask(R, "THREADS", "500")
		N, Err := strconv.Atoi(Raw)
		if Err != nil || N < 1 {
			fmt.Println(color.C("  ! must be a number >= 1", color.Red))
			continue
		}
		Args.Threads = N
		break
	}

	for {
		Raw := ask(R, "DURATION (seconds)", "60")
		D, Err := strconv.Atoi(Raw)
		if Err != nil || D < 1 {
			fmt.Println(color.C("  ! must be a number >= 1", color.Red))
			continue
		}
		Args.Duration = D
		break
	}

	for {
		Raw := strings.ToLower(ask(R, "CLUSTER MODE [y/n]", "n"))
		if Raw == "y" || Raw == "n" {
			Args.ClusterMode = Raw == "y"
			break
		}
		fmt.Println(color.C("  ! enter y or n", color.Red))
	}

	fmt.Println()
	return Args
}
