package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/0xTM7/demon/internal/config"
	"github.com/0xTM7/demon/internal/engine"
	"github.com/0xTM7/demon/internal/printer"
	"github.com/0xTM7/demon/internal/stats"
)

func main() {
	printer.PrintBanner()

	Cfg := RunPrompt()

	if Err := config.Resolve(Cfg); Err != nil {
		fmt.Println(printer.C("    ! "+Err.Error(), printer.Red))
		os.Exit(1)
	}

	printer.PrintConfig(Cfg)
	fmt.Println(printer.C("    * launching demon...", printer.Cyan))
	fmt.Println()

	StopMon := make(chan struct{})
	var LastSnap stats.Snapshot
	var SnapMu sync.Mutex

	go stats.Monitor(Cfg, func(S stats.Snapshot) {
		SnapMu.Lock()
		LastSnap = S
		SnapMu.Unlock()
		printer.PrintLiveStats(S, Cfg.Duration)
	}, StopMon)

	Sig := make(chan os.Signal, 1)
	signal.Notify(Sig, syscall.SIGINT, syscall.SIGTERM)

	Done := make(chan struct{})
	go func() {
		Exec := engine.NewExecutor(Cfg)
		Exec.Run()
		close(Done)
	}()

	select {
	case <-Done:
	case <-Sig:
		fmt.Println()
		fmt.Println(printer.C("    ! interrupted", printer.Yellow))
	}

	close(StopMon)

	SnapMu.Lock()
	Final := LastSnap
	SnapMu.Unlock()

	printer.PrintSummary(Cfg, Final)
}
