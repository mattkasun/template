package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// setup logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logFile, err := os.OpenFile(os.Args[0]+".log", os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err == nil {
		w := io.MultiWriter(os.Stderr, logFile)
		log.SetOutput(w)
	}
	// signals, waitgroups and contexts
	wg := &sync.WaitGroup{} //nolint:varnamelen
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	// strart goroutines
	wg.Add(1)
	go web(ctx, wg)
	// wait for signals
	for { //nolint:staticcheck
		select {
		case <-quit:
			log.Println("quitting ...")
			cancel()
			wg.Wait()
			return
		}
	}
}
