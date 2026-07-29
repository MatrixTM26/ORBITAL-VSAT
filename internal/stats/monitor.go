package stats

import (
	"time"

	"github.com/0xTM7/demon/internal/config"
)

type Snapshot struct {
	Requests  int64
	BytesSent int64
	Success   int64
	Failed    int64
	RPS       float64
	BPS       float64
	Elapsed   time.Duration
}

func Monitor(Cfg *config.DemonConfig, OnTick func(Snapshot), Stop <-chan struct{}) {
	Ticker := time.NewTicker(time.Second)
	defer Ticker.Stop()

	Start := time.Now()
	var PrevReq, PrevBytes int64

	for {
		select {
		case <-Stop:
			return
		case <-Ticker.C:
			CurReq   := Cfg.Stats.Requests.Load()
			CurBytes := Cfg.Stats.BytesSent.Load()
			Elapsed  := time.Since(Start)

			DeltaReq   := CurReq - PrevReq
			DeltaBytes := CurBytes - PrevBytes
			PrevReq    = CurReq
			PrevBytes  = CurBytes

			OnTick(Snapshot{
				Requests:  CurReq,
				BytesSent: CurBytes,
				Success:   Cfg.Stats.Success.Load(),
				Failed:    Cfg.Stats.Failed.Load(),
				RPS:       float64(DeltaReq),
				BPS:       float64(DeltaBytes),
				Elapsed:   Elapsed,
			})
		}
	}
}
