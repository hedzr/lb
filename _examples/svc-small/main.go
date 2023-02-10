// Copyright © 2021 Hedzr Yeh.

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/hedzr/lb/pkg/logger"
)

var portArg = flag.Int("port", 8111, "server port")

func main() {
	logger.SetLevel(logger.DebugLevel)
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler1)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", *portArg),
		Handler: mux,
	}

	var err error
	var exitCh = make(chan struct{})
	var wgStartup sync.WaitGroup
	wgStartup.Add(2)
	go func() {
		fmt.Printf("Server started at port %v...\n", *portArg)
		wgStartup.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("listen: %+v\n", err)
		}
	}()
	go func() {
		wgStartup.Done()
		<-exitCh
		if err = srv.Shutdown(context.Background()); err != nil {
			logger.Errorf("server Shutdown Failed: %+v", err)
		}
		if err == http.ErrServerClosed {
			err = nil
		}
	}()

	wgStartup.Wait()
	setupCloseHandler(func() {
		// on finished ...
	})
	enterLoop()
}

func helloHandler1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, i'm :%v\n", *portArg)
}

// setupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func setupCloseHandler(onFinished func()) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		onFinished()
		os.Exit(0)
	}()
}

func enterLoop() {
	for {
		time.Sleep(10 * time.Second)
	}
}
