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

func ask(R *bufio.Reader, Label, Default string) string {
	if Default != "" {
		fmt.Printf("    %s %s[%s%s%s] %s»%s ",
			color.B(color.C(Label, color.BrightGreen)),
			color.Dim, color.Cyan, Default, color.Reset+color.Dim,
			color.Yellow, color.Cyan,
		)
	} else {
		fmt.Printf("    %s %s»%s ",
			color.B(color.C(Label, color.BrightGreen)),
			color.Yellow, color.Cyan,
		)
	}
	Line, _ := R.ReadString('\n')
	fmt.Print(color.Reset)
	Line = strings.TrimSpace(Line)
	if Line == "" {
		return Default
	}
	return Line
}

func RunPrompt() DemonArgs {
	R := bufio.NewReader(os.Stdin)
	Args := DemonArgs{
		Protocol:  "H1",
		JAProfile: "Chrome",
	}

	fmt.Println(color.B(color.C("    CONFIGURATION", color.Red)))
	fmt.Println()

	Args.Target = ask(R, "TARGET", "")
	if Args.Target == "" {
		fmt.Println(color.C("    ! target is required", color.Red))
		os.Exit(1)
	}

	MethodRaw := strings.ToUpper(ask(R, "METHOD", "GET"))
	Args.Method = MethodRaw

	if L7Methods[MethodRaw] {
		Args.Layer = "L7"

		ProtoRaw := strings.ToUpper(ask(R, "PROTOCOL [H1/H2]", "H1"))
		if ProtoRaw == "H2" {
			Args.Protocol = "H2"
		}

		JARaw := ask(R, "JA3 [Chrome/Firefox/Safari]", "Chrome")
		switch strings.ToLower(JARaw) {
		case "firefox":
			Args.JAProfile = "Firefox"
		case "safari":
			Args.JAProfile = "Safari"
		default:
			Args.JAProfile = "Chrome"
		}

	} else if L4Methods[MethodRaw] {
		Args.Layer = "L4"
	} else if L3Methods[MethodRaw] {
		Args.Layer = "L3"
	} else {
		fmt.Printf("    %s unknown method: %s\n", color.C("!", color.Red), MethodRaw)
		os.Exit(1)
	}

	ThreadsRaw := ask(R, "THREADS", "500")
	N, Err := strconv.Atoi(ThreadsRaw)
	if Err != nil || N < 1 {
		N = 500
	}
	Args.Threads = N

	DurRaw := ask(R, "DURATION (seconds)", "60")
	D, Err := strconv.Atoi(DurRaw)
	if Err != nil || D < 1 {
		D = 60
	}
	Args.Duration = D

	ClusterRaw := strings.ToLower(ask(R, "CLUSTER MODE [y/n]", "n"))
	Args.ClusterMode = ClusterRaw == "y"

	fmt.Println()
	return Args
}
