package monitor

import (
	"time"

	"github.com/MatrixTM26/vsat/internal/methods"
)

type Snapshot struct {
	Total     int64
	BytesSent int64
	RPS       float64
	BPS       float64
	Elapsed   time.Duration
}

func Run(Cfg *methods.Config, OnTick func(Snapshot), Stop <-chan struct{}) {
	Ticker  := time.NewTicker(time.Second)
	Start   := time.Now()
	var PrevTotal, PrevBytes int64

	defer Ticker.Stop()

	for {
		select {
		case <-Stop:
			return
		case <-Ticker.C:
			CurTotal := Cfg.Stats.RequestsCount.Load()
			CurBytes := Cfg.Stats.BytesSent.Load()
			Elapsed  := time.Since(Start)

			OnTick(Snapshot{
				Total:     CurTotal,
				BytesSent: CurBytes,
				RPS:       float64(CurTotal - PrevTotal),
				BPS:       float64(CurBytes - PrevBytes),
				Elapsed:   Elapsed,
			})

			PrevTotal = CurTotal
			PrevBytes = CurBytes
		}
	}
}
