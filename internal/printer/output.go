package printer

import (
	"fmt"
	"runtime"
	"time"

	"github.com/0xTM7/demon/internal/config"
	"github.com/0xTM7/demon/internal/stats"
)

const Banner = `
    ██████╗ ███████╗███╗   ███╗ ██████╗ ███╗   ██╗
    ██╔══██╗██╔════╝████╗ ████║██╔═══██╗████╗  ██║
    ██║  ██║█████╗  ██╔████╔██║██║   ██║██╔██╗ ██║
    ██║  ██║██╔══╝  ██║╚██╔╝██║██║   ██║██║╚██╗██║
    ██████╔╝███████╗██║ ╚═╝ ██║╚██████╔╝██║ ╚████║
    ╚═════╝ ╚══════╝╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝
`

func PrintBanner() {
	fmt.Println(C(Banner, Red))
	fmt.Println(C("    Multi-Layer Network Stress Testing Framework   v2.0.0   by 0xTM7", Dim))
	fmt.Println()
}

func PrintConfig(Cfg *config.DemonConfig) {
	fmt.Println(B(C("    TARGET", Cyan)))
	fmt.Printf("    %s  %s\n", C(PadRight("url", 14), Dim), Cfg.Target)
	fmt.Printf("    %s  %s:%d\n", C(PadRight("resolved", 14), Dim), Cfg.IP, Cfg.Port)
	fmt.Printf("    %s  %s\n", C(PadRight("layer", 14), Dim), string(Cfg.Layer))
	fmt.Printf("    %s  %s\n", C(PadRight("method", 14), Dim), string(Cfg.Method))
	if Cfg.Layer == config.LayerL7 {
		fmt.Printf("    %s  %s\n", C(PadRight("protocol", 14), Dim), string(Cfg.Protocol))
		fmt.Printf("    %s  %s\n", C(PadRight("ja3", 14), Dim), string(Cfg.JA3))
		fmt.Printf("    %s  %d\n", C(PadRight("user-agents", 14), Dim), len(Cfg.UserAgents))
		fmt.Printf("    %s  %d\n", C(PadRight("referers", 14), Dim), len(Cfg.Referers))
	}
	fmt.Printf("    %s  %d\n", C(PadRight("threads", 14), Dim), Cfg.Threads)
	fmt.Printf("    %s  %ds\n", C(PadRight("duration", 14), Dim), Cfg.Duration)
	if Cfg.ClusterMode {
		Total := Cfg.Threads * runtime.NumCPU()
		fmt.Printf("    %s  %d cores x %d = %s total\n",
			C(PadRight("cluster", 14), Dim),
			runtime.NumCPU(), Cfg.Threads,
			C(fmt.Sprintf("%d", Total), Yellow),
		)
	}
	fmt.Println()
}

func PrintLiveStats(S stats.Snapshot, Duration int) {
	Elapsed := int(S.Elapsed.Seconds())
	Bar := ProgressBar(Elapsed, Duration, 26)
	Pct := 0
	if Duration > 0 {
		Pct = int(float64(Elapsed) / float64(Duration) * 100)
	}
	fmt.Printf(
		"\r    %s  %d%%  %s  %s  %s ok  %s err",
		C(Bar, Red),
		Pct,
		C(FmtRate(S.RPS)+" req", Yellow),
		C(FmtBytes(S.BPS)+"/s", Cyan),
		C(fmt.Sprintf("%d", S.Success), Green),
		C(fmt.Sprintf("%d", S.Failed), Red),
	)
}

func PrintSummary(Cfg *config.DemonConfig, S stats.Snapshot) {
	fmt.Println()
	fmt.Println()
	fmt.Println(B(C("    SUMMARY", Magenta)))

	TotalSec := S.Elapsed.Seconds()
	RPS := 0.0
	if TotalSec > 0 {
		RPS = float64(S.Requests) / TotalSec
	}

	SuccPct := 0.0
	if S.Requests > 0 {
		SuccPct = float64(S.Success) / float64(S.Requests) * 100
	}

	fmt.Printf("    %s  %d\n", PadRight("total", 16), S.Requests)
	fmt.Printf("    %s  %s  %.1f%%\n", PadRight("success", 16), C(fmt.Sprintf("%d", S.Success), Green), SuccPct)
	fmt.Printf("    %s  %s\n", PadRight("failed", 16), C(fmt.Sprintf("%d", S.Failed), Red))
	fmt.Printf("    %s  %s\n", PadRight("throughput", 16), C(FmtRate(RPS)+" req", Yellow))
	fmt.Printf("    %s  %s\n", PadRight("data sent", 16), FmtBytes(float64(S.BytesSent)))
	fmt.Printf("    %s  %s\n", PadRight("elapsed", 16), fmt.Sprintf("%.2fs", TotalSec))
	fmt.Println()

	fmt.Println(B(C("    RUNTIME", Magenta)))
	fmt.Printf("    %s  %s\n", PadRight("target", 16), Cfg.Target)
	fmt.Printf("    %s  %s\n", PadRight("method", 16), string(Cfg.Method))
	fmt.Printf("    %s  %s\n", PadRight("layer", 16), string(Cfg.Layer))
	if Cfg.Layer == config.LayerL7 {
		fmt.Printf("    %s  %s\n", PadRight("protocol", 16), string(Cfg.Protocol))
		fmt.Printf("    %s  %s\n", PadRight("ja3", 16), string(Cfg.JA3))
	}
	fmt.Printf("    %s  %d\n", PadRight("threads", 16), Cfg.Threads)
	fmt.Printf("    %s  %s\n", PadRight("started", 16), time.Now().Format("15:04:05"))
	fmt.Println()
}
