package proxy

import (
	"caching-proxy/logger"
	"caching-proxy/proxy/client"
	"caching-proxy/proxy/helpers"
	"context"
	"fmt"
	"net/http"
	"time"
)

var (
	server  *http.Server
	oClient *client.Client
	plog    *logger.Logger
)

// recover function
func recover_hdl(w http.ResponseWriter) {
	if r := recover(); r != nil {
		plog.Errorf("%v\n", r)
		http.Error(w, fmt.Sprintf("Internal error: %v", r), http.StatusInternalServerError) //500
	}
}

// send response to the caller
func send_response(w http.ResponseWriter, resp *client.ClientReqRes) {

	for key, values := range resp.Headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.RespCode)
	w.Write(resp.Body)
}

// main handler
func handler(w http.ResponseWriter, r *http.Request) {
	defer recover_hdl(w)

	url := r.URL.Path
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}

	resp := oClient.SendRequest(&client.ClientReqRes{
		Method:  r.Method,
		Uri:     url,
		Headers: r.Header,
		Body:    helpers.ReadBody(r.Body),
	})

	send_response(w, &resp)
}

// start proxy
func Start(port int, origin string, log *logger.Logger) {

	if port <= 0 {
		panic("Invalid port")
	}
	if origin == "" {
		panic("Empty origin")
	}
	if log == nil {
		panic("Logger is nil")
	}

	oClient = client.New(origin)
	plog = log

	plog.Infof("port: %d", port)
	plog.Infof("origin: %s", origin)

	mux := http.NewServeMux()
	// Register handler func
	mux.HandleFunc("/", handler)

	server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		plog.Errorf("%v", err)
	}
}

// shutdown proxy
func ShutDown() {
	plog.Debugf("Shutting down gracefully\n")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		plog.Errorf("%v\n", err)
	}
	plog.Debugf("Server shut down\n")
}
