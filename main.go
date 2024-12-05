package main

import (
	"caching-proxy/proxy"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	port := flag.Int("port", 3000, "the port on which the caching proxy server will run")
	origin := flag.String("origin", "", "the URL of the server to which the requests will be forwarded")
	flag.Parse()

	// start proxy in another goroutine (thread)
	go proxy.Start(*port, *origin)

	// add signal handler
	quit := make(chan os.Signal, 1)                    // create a channel for signals
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // relay SIGINT, SIGTERM signals to quit channel
	// wait for signal
	<-quit
	proxy.ShutDown()
}
