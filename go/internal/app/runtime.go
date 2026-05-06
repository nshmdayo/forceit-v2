package app

import (
	"context"
	"log"
	"time"

	"forceit-v2/go/internal/domain"
	"forceit-v2/go/internal/sensor"
)

type Runtime struct {
	Sensor sensor.TCPJSONClient
}

func (r Runtime) Run(ctx context.Context) error {
	frames := make(chan domain.SkeletonFrame, 4)
	go func() {
		if err := r.Sensor.Stream(ctx, frames); err != nil {
			log.Printf("sensor stream stopped: %v", err)
		}
		close(frames)
	}()

	tick := time.NewTicker(time.Second / 60)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case f, ok := <-frames:
			if !ok {
				return nil
			}
			_ = f // TODO: map joints to domain forces.
		case <-tick.C:
			// TODO: fixed-step physics update and render sync.
		}
	}
}
