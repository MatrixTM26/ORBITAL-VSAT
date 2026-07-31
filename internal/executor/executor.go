package executor

import (
	"runtime"
	"sync"
	"time"

	"github.com/MatrixTM26/vsat/internal/methods"
	"github.com/MatrixTM26/vsat/internal/methods/l3"
	"github.com/MatrixTM26/vsat/internal/methods/l4"
	"github.com/MatrixTM26/vsat/internal/methods/l7"
)

type WorkerFunc func(Cfg *methods.Config, Stop <-chan struct{})

type Executor struct {
	Cfg         *methods.Config
	Threads     int
	Duration    int
	ClusterMode bool
}

func NewExecutor(Cfg *methods.Config, Threads, Duration int, ClusterMode bool) *Executor {
	return &Executor{
		Cfg:         Cfg,
		Threads:     Threads,
		Duration:    Duration,
		ClusterMode: ClusterMode,
	}
}

func (E *Executor) resolveWorker() WorkerFunc {
	switch E.Cfg.Method {
	case "GET", "POST", "PUT", "HEAD", "DELETE", "PATCH", "OPTIONS", "CONNECT", "TRACE", "RANDOM":
		return l7.HTTPExecutor
	case "SLOWLORIS":
		return l7.SlowlorisExecutor
	case "SLOWPOST":
		return l7.SlowPostExecutor
	case "SLOWREAD":
		return l7.SlowReadExecutor
	case "RUDY":
		return l7.RUDYExecutor
	case "TCP":
		return l4.TCPConnectExecutor
	case "SYN":
		return l4.SYNExecutor
	case "ACK":
		return l4.ACKExecutor
	case "RST":
		return l4.RSTExecutor
	case "FIN":
		return l4.FINExecutor
	case "XMAS":
		return l4.XMASExecutor
	case "PSH":
		return l4.PSHExecutor
	case "URG":
		return l4.URGExecutor
	case "NULL":
		return l4.NULLExecutor
	case "SYNACK":
		return l4.SYNACKExecutor
	case "UDP":
		return l4.UDPExecutor
	case "UDP-FRAG":
		return l4.UDPFragExecutor
	case "DNS-AMP":
		return l4.DNSAMPExecutor
	case "NTP-AMP":
		return l4.NTPAMPExecutor
	case "SSDP-AMP":
		return l4.SSDAMPExecutor
	case "MEMCACHED":
		return l4.MEMCACHEDExecutor
	case "CHARGEN":
		return l4.CHARGENExecutor
	case "FRAGGLE":
		return l4.FRAGGLEExecutor
	case "ICMP":
		return l3.ICMPExecutor
	case "SMURF":
		return l3.SMURFExecutor
	case "PING":
		return l3.PINGExecutor
	default:
		return l7.HTTPExecutor
	}
}

func (E *Executor) Run() {
	Stop := make(chan struct{})
	var WG sync.WaitGroup

	Total := E.Threads
	if E.ClusterMode {
		Total = E.Threads * runtime.NumCPU()
	}

	Worker := E.resolveWorker()

	for I := 0; I < Total; I++ {
		WG.Add(1)
		go func() {
			defer WG.Done()
			Worker(E.Cfg, Stop)
		}()
	}

	time.Sleep(time.Duration(E.Duration) * time.Second)
	close(Stop)
	WG.Wait()
}
