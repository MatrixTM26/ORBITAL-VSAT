package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/0xTM7/demon/internal/config"
	"github.com/0xTM7/demon/internal/printer"
)

var L7Methods = map[string]bool{
	"GET": true, "POST": true, "PUT": true, "PATCH": true,
	"DELETE": true, "HEAD": true, "OPTIONS": true, "TRACE": true,
	"CONNECT": true, "RANDOM": true,
}

var L4Methods = map[string]bool{
	"TCP": true, "SYN": true, "ACK": true, "RST": true,
	"FIN": true, "XMAS": true, "UDP": true, "UDP-FRAG": true,
	"DNS-AMP": true, "NTP-AMP": true,
}

var L3Methods = map[string]bool{
	"ICMP": true,
}

func Ask(Reader *bufio.Reader, Label, Default string) string {
	if Default != "" {
		fmt.Printf("    %s %s[ default: %s%s%s ] > %s",
			printer.B(printer.C(Label, printer.Green)),
			printer.Reset,
			printer.Cyan, Default, printer.Reset,
			printer.Cyan,
		)
	} else {
		fmt.Printf("    %s > %s",
			printer.B(printer.C(Label, printer.Green)),
			printer.Cyan,
		)
	}
	Line, _ := Reader.ReadString('\n')
	fmt.Print(printer.Reset)
	Line = strings.TrimSpace(Line)
	if Line == "" {
		return Default
	}
	return Line
}

func RunPrompt() *config.DemonConfig {
	R := bufio.NewReader(os.Stdin)
	Cfg := &config.DemonConfig{
		Protocol: config.ProtoH1,
		JA3:      config.JA3Chrome,
		Headers:  make(map[string]string),
	}

	fmt.Println(printer.B(printer.C("    DEMON CONFIGURATION", printer.Red)))
	fmt.Println()

	Cfg.Target = Ask(R, "TARGET", "")
	if Cfg.Target == "" {
		fmt.Println(printer.C("    ! target is required", printer.Red))
		os.Exit(1)
	}

	MethodStr := strings.ToUpper(Ask(R, "METHOD", "GET"))
	Cfg.Method = config.Method(MethodStr)
	Cfg.Layer = config.DetectLayer(Cfg.Method)

	if _, IsL7 := L7Methods[MethodStr]; IsL7 {
		ProtoStr := strings.ToUpper(Ask(R, "PROTOCOL  [ H1 | H2 ]", "H1"))
		switch ProtoStr {
		case "H2":
			Cfg.Protocol = config.ProtoH2
		default:
			Cfg.Protocol = config.ProtoH1
		}

		JA3Str := strings.ToLower(Ask(R, "JA3 PROFILE  [ chrome | firefox | safari ]", "chrome"))
		switch JA3Str {
		case "firefox":
			Cfg.JA3 = config.JA3Firefox
		case "safari":
			Cfg.JA3 = config.JA3Safari
		default:
			Cfg.JA3 = config.JA3Chrome
		}
	}

	ThreadsStr := Ask(R, "THREADS", "500")
	Cfg.Threads, _ = strconv.Atoi(ThreadsStr)
	if Cfg.Threads < 1 {
		Cfg.Threads = 500
	}

	DurStr := Ask(R, "DURATION (seconds)", "60")
	Cfg.Duration, _ = strconv.Atoi(DurStr)
	if Cfg.Duration < 1 {
		Cfg.Duration = 60
	}

	ClusterStr := strings.ToLower(Ask(R, "CLUSTER MODE  [ y/n ]", "n"))
	Cfg.ClusterMode = ClusterStr == "y"

	fmt.Println()
	return Cfg
}
