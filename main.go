package main

import (
	"caching-proxy/logger"
	"caching-proxy/proxy"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func recover_handler() {
	if r := recover(); r != nil {
		log.Printf("error %v\n", r)
	}
}

func main() {

	defer recover_handler()

	port := flag.Int("port", 0, "the port on which the caching proxy server will run")
	origin := flag.String("origin", "", "the URL of the server to which the requests will be forwarded")
	dbg := flag.Bool("debug", false, "turn on debug logs")
	flag.Parse()
	// launch in another thread
	go func() {
		// catch proxy errors
		defer recover_handler()
		// start proxy
		proxy.Start(*port, *origin, logger.New(*dbg))
	}()

	// add signal handler
	quit := make(chan os.Signal, 1)                    // create a channel for signals
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // relay SIGINT, SIGTERM signals to quit channel
	// wait for signal
	<-quit
	proxy.ShutDown()
}
