package app

import (
	"context"
	"log"
	"time"

	"forceit-v2/go/internal/domain"
	"forceit-v2/go/internal/sensor"
)

const fixedFPS = 60

type Runtime struct {
	Sensor sensor.TCPJSONClient
}

func (r Runtime) Run(ctx context.Context) error {
	frames := r.startSensorStream(ctx)
	tick := time.NewTicker(time.Second / fixedFPS)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case frame, ok := <-frames:
			if !ok {
				return nil
			}
			r.handleFrame(frame)
		case <-tick.C:
			r.step(1.0 / fixedFPS)
		}
	}
}

func (r Runtime) startSensorStream(ctx context.Context) <-chan domain.SkeletonFrame {
	frames := make(chan domain.SkeletonFrame, 4)
	go func() {
		if err := r.Sensor.Stream(ctx, frames); err != nil {
			log.Printf("sensor stream stopped: %v", err)
		}
		log.Print("sensor stream finished")
		close(frames)
	}()
	return frames
}

func (r Runtime) handleFrame(frame domain.SkeletonFrame) {
	_ = frame // TODO: map joints to domain forces.
}

func (r Runtime) step(dt float64) {
	_ = dt // TODO: fixed-step physics update and render sync.
}
