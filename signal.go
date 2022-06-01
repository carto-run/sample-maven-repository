package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func signalHandlingContext(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		var sig os.Signal

		sig = <-ch
		log.Println("terminal signal received", sig)
		cancel()

		sig = <-ch
		log.Println("last terminal signal received", sig)
		os.Exit(1)
	}()

	return ctx
}
