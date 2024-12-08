package proxy

import (
	"caching-proxy/logger"
	"caching-proxy/proxy/cache"
	"caching-proxy/proxy/client"
	"caching-proxy/proxy/helpers"
	"caching-proxy/proxy/request"
	"context"
	"fmt"
	"net/http"
	"time"
)

var (
	server  *http.Server
	oClient *client.Client
	plog    *logger.Logger
	pcache  *cache.Cache
	stopMon chan struct{} // stop channel for backup monitor
)

// recover function
func recover_hdl(w http.ResponseWriter) {
	if r := recover(); r != nil {
		plog.Errorf("%v\n", r)
		http.Error(w, fmt.Sprintf("Internal error: %v", r), http.StatusInternalServerError) //500
	}
}

// send response to the caller
func send_response(w http.ResponseWriter, resp *request.Request) {
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

	// TODO: check also method
	cache_key := r.Method + "::" + url
	// check if exists in cache
	resp, ok := pcache.Get(cache_key)
	if !ok {
		// doesn't exist
		plog.Debugf("sending %s request to %s\n", r.Method, url)
		resp = oClient.SendRequest(&request.Request{
			Method:  r.Method,
			Uri:     url,
			Headers: r.Header,
			Body:    helpers.ReadBody(r.Body),
		})
		plog.Debugf("saving %s to cache\n", url)
		pcache.Put(cache_key, resp)
		resp.Headers["X-Cache"] = []string{"MISS"}
	} else {
		// exists
		plog.Debugf("found request existing in cache\n")
		resp.Headers["X-Cache"] = []string{"HIT"}
	}

	send_response(w, &resp)
}

// do some setup work at start
func setup(port int, origin string, backup string, log *logger.Logger) {
	if port <= 0 {
		panic("Invalid port")
	}
	if origin == "" {
		panic("Empty origin")
	}
	if log == nil {
		panic("Logger is nil")
	}

	plog = log
	oClient = client.New(origin, log)
	pcache = cache.New(origin, backup)
	stopMon = make(chan struct{})
}

// cache backup monitor
func backup_monitor(seconds time.Duration) {
	ticker := time.NewTicker(seconds * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// periodically backup cache to file
			if pcache.HasChanged() {
				plog.Debugf("Backing up cache\n")
				if err := pcache.Backup(); err != nil {
					plog.Errorf("backup error %v\n", err)
				}
			}
		case <-stopMon:
			// return if explicitly stopped
			plog.Debugf("Backup monitor stopped")
			return
		}
	}
}

// start proxy
func Start(port int, origin string, backup string, log *logger.Logger) {

	setup(port, origin, backup, log)

	plog.Infof("port: %d", port)
	plog.Infof("origin: %s", origin)

	// run backup monitor in a separate thread
	go backup_monitor(1)

	mux := http.NewServeMux()
	// Register handler func
	mux.HandleFunc("/", handler)

	server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		plog.Errorf("%v\n", err)
	}

	// backup cache to file
	if err := pcache.Backup(); err != nil {
		plog.Errorf("%v\n", err)
	}
}

// shutdown proxy
func ShutDown() {
	plog.Debugf("Shutting down gracefully\n")
	// close stop channel to stop monitor thread
	close(stopMon)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		plog.Errorf("%v\n", err)
	}
	plog.Debugf("Server shut down\n")
}
