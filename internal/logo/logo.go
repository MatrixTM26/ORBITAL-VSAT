package logo

import (
	"fmt"
	"strings"
	"time"

	"github.com/MatrixTM26/vsat/internal/ansi"
	"github.com/MatrixTM26/vsat/internal/color"
	"github.com/MatrixTM26/vsat/internal/monitor"
)

const Banner = `
     __   __   ______   ________   _________   
    /_/\ /_/\ /_____/\ /_______/\ /________/\  
    \:\ \\ \ \\::::_\/_\::: _  \ \\__.::.__\/  
     \:\ \\ \ \\:\/___/\\::(_)  \ \  \::\ \    
      \:\_/.:\ \\_::._\:\\:: __  \ \  \::\ \   
       \ ..::/ /  /____\:\\:.\ \  \ \  \::\ \  
        \___/_/   \_____\/ \__\/\__\/   \__\/  
`

func PrintBanner() {
	ansi.Clear()
	fmt.Println(color.C(Banner, color.Red))
	time.Sleep(100 * time.Millisecond)
	ansi.TypewriterFast(color.C("	    Volumetric Socket Artillery", color.Red))
	ansi.TypewriterFast(color.C("  Author  :  MatrixTM26", color.Dim))
	ansi.TypewriterFast(color.C("  Version :  1.2.0", color.Dim))
	fmt.Println()
}

func PrintSeparator() {
	fmt.Println(color.C("  "+strings.Repeat(" ", 48), color.Dim))
}

func PrintConfig(
	Target, ResolvedIP string,
	Port int,
	Layer, Method, Protocol, JAProfile string,
	Threads, Duration int,
	ClusterMode bool,
	Workers int,
	UACount, RefCount int,
) {
	fmt.Println(color.B(color.C("  TARGET", color.Cyan)))
	fmt.Printf("  %-16s  %s\n", color.C("url", color.Dim), Target)
	fmt.Printf("  %-16s  %s:%d\n", color.C("resolved", color.Dim), ResolvedIP, Port)
	fmt.Printf("  %-16s  %s\n", color.C("layer", color.Dim), Layer)
	fmt.Printf("  %-16s  %s\n", color.C("method", color.Dim), Method)
	if Layer == "L7" {
		fmt.Printf("  %-16s  %s\n", color.C("protocol", color.Dim), Protocol)
		fmt.Printf("  %-16s  %s\n", color.C("ja3", color.Dim), JAProfile)
		fmt.Printf("  %-16s  %d\n", color.C("user-agents", color.Dim), UACount)
		fmt.Printf("  %-16s  %d\n", color.C("referers", color.Dim), RefCount)
	}
	fmt.Printf("  %-16s  %d\n", color.C("threads", color.Dim), Threads)
	fmt.Printf("  %-16s  %ds\n", color.C("duration", color.Dim), Duration)
	if ClusterMode {
		fmt.Printf("  %-16s  %s\n",
			color.C("cluster", color.Dim),
			color.C(fmt.Sprintf("%d total goroutines", Workers), color.Yellow),
		)
	}
	fmt.Println()
}

func PrintLive(S monitor.Snapshot, Duration int) {
	Elapsed := int(S.Elapsed.Seconds())
	Pct := 0
	if Duration > 0 {
		Pct = int(float64(Elapsed) / float64(Duration) * 100)
		if Pct > 100 {
			Pct = 100
		}
	}

	Width := 40
	Filled := int(float64(Pct) / 100.0 * float64(Width))
	Bar := strings.Repeat("█", Filled) + strings.Repeat("░", Width-Filled)

	RPS := fmtRate(S.RPS)
	BPS := fmtBytes(S.BPS)

	fmt.Printf(
		"\r  %s  %s%%  %s  %s  %s",
		color.C(Bar, color.Red),
		color.C(fmt.Sprintf("%3d", Pct), color.White),
		color.C(RPS+" req/s", color.Yellow),
		color.C(BPS+"/s", color.Cyan),
		color.C(fmt.Sprintf("%d sent", S.Total), color.Green),
	)
}

func PrintSummary(
	Target, Method, Layer, Protocol, JAProfile string,
	Threads, Duration int,
	S monitor.Snapshot,
) {
	fmt.Println()
	fmt.Println()
	fmt.Println(color.B(color.C("  SUMMARY", color.Magenta)))

	Elapsed := S.Elapsed.Seconds()
	RPS := 0.0
	if Elapsed > 0 {
		RPS = float64(S.Total) / Elapsed
	}

	fmt.Printf("  %-18s  %d\n", "total requests", S.Total)
	fmt.Printf("  %-18s  %s\n", "throughput", color.C(fmtRate(RPS)+" req/s", color.Yellow))
	fmt.Printf("  %-18s  %s\n", "data sent", fmtBytes(float64(S.BytesSent)))
	fmt.Printf("  %-18s  %.2fs\n", "elapsed", Elapsed)
	fmt.Println()
	fmt.Println(color.B(color.C("  CONFIG", color.Magenta)))
	fmt.Printf("  %-18s  %s\n", "target", Target)
	fmt.Printf("  %-18s  %s\n", "method", Method)
	fmt.Printf("  %-18s  %s\n", "layer", Layer)
	if Layer == "L7" {
		fmt.Printf("  %-18s  %s\n", "protocol", Protocol)
		fmt.Printf("  %-18s  %s\n", "ja3", JAProfile)
	}
	fmt.Printf("  %-18s  %d\n", "threads", Threads)
	fmt.Printf("  %-18s  %ds\n", "duration", Duration)
	fmt.Println()
}

func fmtBytes(B float64) string {
	switch {
	case B < 1024:
		return fmt.Sprintf("%.0fB", B)
	case B < 1048576:
		return fmt.Sprintf("%.1fKB", B/1024)
	case B < 1073741824:
		return fmt.Sprintf("%.1fMB", B/1048576)
	default:
		return fmt.Sprintf("%.2fGB", B/1073741824)
	}
}

func fmtRate(N float64) string {
	if N < 1000 {
		return fmt.Sprintf("%.1f", N)
	}
	return fmt.Sprintf("%.2fk", N/1000)
}
