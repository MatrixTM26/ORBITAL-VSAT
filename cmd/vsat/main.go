package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/MatrixTM26/vsat/internal/color"
	"github.com/MatrixTM26/vsat/internal/executor"
	"github.com/MatrixTM26/vsat/internal/logo"
	"github.com/MatrixTM26/vsat/internal/methods"
	"github.com/MatrixTM26/vsat/internal/monitor"
)

func main() {
	logo.PrintBanner()

	Args := RunPrompt()

	R, Err := Resolve(Args.Target)
	if Err != nil {
		fmt.Println(color.C("  ! "+Err.Error(), color.Red))
		os.Exit(1)
	}

	UAs  := LoadFile("config/UA.txt", DefaultUserAgents)
	Refs := LoadFile("config/Referers.txt", DefaultReferers)

	Cfg := &methods.Config{
		IP:        R.IP,
		Port:      R.Port,
		Host:      R.Host,
		Path:      R.Path,
		Scheme:    R.Scheme,
		Protocol:  Args.Protocol,
		JAProfile: Args.JAProfile,
		Method:    Args.Method,
		UserAgents: UAs,
		Referers:   Refs,
		Stats:     &methods.Stats{},
	}

	Workers := Args.Threads
	if Args.ClusterMode {
		Workers = Args.Threads * runtime.NumCPU()
	}

	logo.PrintConfig(
		R.Target, R.IP, R.Port,
		Args.Layer, Args.Method, Args.Protocol, Args.JAProfile,
		Args.Threads, Args.Duration,
		Args.ClusterMode, Workers,
		len(UAs), len(Refs),
	)

	fmt.Println(color.C("  launching vsat...", color.Green))
	fmt.Println()

	StopMon := make(chan struct{})
	var LastSnap monitor.Snapshot
	var SnapMu sync.Mutex

	go monitor.Run(Cfg, func(S monitor.Snapshot) {
		SnapMu.Lock()
		LastSnap = S
		SnapMu.Unlock()
		logo.PrintLive(S, Args.Duration)
	}, StopMon)

	Sig := make(chan os.Signal, 1)
	signal.Notify(Sig, syscall.SIGINT, syscall.SIGTERM)

	Done := make(chan struct{})
	go func() {
		Exec := executor.NewExecutor(Cfg, Args.Threads, Args.Duration, Args.ClusterMode)
		Exec.Run()
		close(Done)
	}()

	select {
	case <-Done:
	case <-Sig:
		fmt.Println()
		fmt.Println(color.C("  ! interrupted", color.Yellow))
	}

	close(StopMon)

	SnapMu.Lock()
	Final := LastSnap
	SnapMu.Unlock()

	logo.PrintSummary(
		R.Target, Args.Method, Args.Layer, Args.Protocol, Args.JAProfile,
		Args.Threads, Args.Duration, Final,
	)
}
