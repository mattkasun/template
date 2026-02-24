package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// setup logging
	log.SetFlags(log.Lshortfile) // date/time added by journald
	slog.SetDefault(slog.Default())
	// signals, waitgroups and contexts
	wg := &sync.WaitGroup{} //nolint:varnamelen
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	// strart goroutines
	wg.Go(func() {
		web(ctx)
	})
	// wait for signals
	<-quit
	log.Println("quitting ...")
	cancel()
	wg.Wait()
}
