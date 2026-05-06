package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"forceit-v2/go/internal/app"
	"forceit-v2/go/internal/sensor"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	r := app.Runtime{Sensor: sensor.TCPJSONClient{Addr: "127.0.0.1:50051"}}
	if err := r.Run(ctx); err != nil && err != context.Canceled {
		log.Fatal(err)
	}
}
