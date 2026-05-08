package app

import (
	"context"
	"errors"
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
	frames, streamErrs := r.startSensorStream(ctx)
	tick := time.NewTicker(time.Second / fixedFPS)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.Canceled) {
				return nil
			}
			return ctx.Err()
		case err, ok := <-streamErrs:
			if !ok {
				return nil
			}
			if errors.Is(ctx.Err(), context.Canceled) {
				return nil
			}
			return err
		case frame, ok := <-frames:
			if !ok {
				frames = nil
				continue
			}
			r.handleFrame(frame)
		case <-tick.C:
			r.step(1.0 / fixedFPS)
		}
	}
}

func (r Runtime) startSensorStream(ctx context.Context) (<-chan domain.SkeletonFrame, <-chan error) {
	frames := make(chan domain.SkeletonFrame, 4)
	errCh := make(chan error, 1)
	go func() {
		defer close(frames)
		defer close(errCh)

		err := r.Sensor.Stream(ctx, frames)
		switch {
		case err == nil:
			log.Print("sensor stream finished")
		case errors.Is(err, context.Canceled) || errors.Is(ctx.Err(), context.Canceled):
			log.Print("sensor stream canceled")
		default:
			log.Printf("sensor stream stopped: %v", err)
			errCh <- err
		}
	}()
	return frames, errCh
}

func (r Runtime) handleFrame(frame domain.SkeletonFrame) {
	_ = frame // TODO: map joints to domain forces.
}

func (r Runtime) step(dt float64) {
	_ = dt // TODO: fixed-step physics update and render sync.
}
