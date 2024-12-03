package proxy

import (
	"caching-proxy/logger"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	server *http.Server
	errLog = logger.Err
	dbgLog = logger.Dbg
)

// recover function
func recover_hdl(w http.ResponseWriter) {
	if r := recover(); r != nil {
		errLog.Printf("%v", r)
		http.Error(w, fmt.Sprintf("Internal error: %v", r), http.StatusInternalServerError) //500
	}
}

// main handler
func handler(w http.ResponseWriter, r *http.Request) {
	defer recover_hdl(w)

}

// validate origin url
func validateOrigin(origin string) error {
	parsedURL, err := url.Parse(origin)
	if err != nil {
		return err
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("invalid scheme")
	}
	if parsedURL.Path != "" {
		return errors.New("origin shouldn't contain any path")
	}
	if parsedURL.RawQuery != "" {
		return errors.New("origin shouldn't contain any queries")
	}
	return nil
}

// start proxy
func Start(port int, origin string) {

	if port <= 0 {
		panic("Invalid port")
	}
	if origin == "" {
		panic("Empty origin")
	}
	if !strings.HasPrefix("http://", origin) &&
		!strings.HasPrefix("https://", origin) {
		origin = "https://" + origin
	}

	if err := validateOrigin(origin); err != nil {
		panic(err)
	}

	dbgLog.Printf("port: %d", port)
	dbgLog.Printf("origin: %s", origin)

	mux := http.NewServeMux()
	// Register handler func
	mux.HandleFunc("/", handler)

	server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		errLog.Printf("%v", err)
	}
}

// shutdown proxy
func ShutDown() {
	dbgLog.Println("Shutting down gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		errLog.Printf("%v", err)
	}
	dbgLog.Println("Server shut down")
}
