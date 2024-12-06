package client

import (
	"bytes"
	"caching-proxy/logger"
	"caching-proxy/proxy/helpers"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	handle *http.Client
	origin string
	scheme string
	clog   *logger.Logger
}

type ClientReqRes struct {
	Body    []byte
	Headers map[string][]string
	// request only
	Method string
	Uri    string
	// response only
	RespStatus string
	RespCode   int
}

const default_scheme = "http://"

// validate origin url
func validateOrigin(origin string) (string, string) {
	if !strings.HasPrefix(origin, "http") {
		origin = default_scheme + origin
	}
	parsedURL, err := url.Parse(origin)
	if err != nil {
		panic(err)
	}
	if parsedURL.Path != "" {
		panic("origin shouldn't contain any path")
	}
	if parsedURL.RawQuery != "" {
		panic("origin shouldn't contain any queries")
	}

	return parsedURL.Host, parsedURL.Scheme + "://"
}

// create new client
func New(remote string, log *logger.Logger) *Client {
	if log == nil {
		panic("Logger is nil")
	}
	host, scheme := validateOrigin(remote)
	return &Client{
		handle: &http.Client{},
		origin: host,
		scheme: scheme,
		clog:   log,
	}
}

// send client request
func (client *Client) SendRequest(request *ClientReqRes) ClientReqRes {

	defer func() {
		if r := recover(); r != nil {
			client.clog.Errorf("client error: %v\n", r)
		}
	}()

	request.Uri = client.scheme + client.origin + request.Uri

	client.clog.Debugf("request uri: %s\n", request.Uri)

	req, err := http.NewRequest(request.Method, request.Uri, bytes.NewBuffer(request.Body))
	if err != nil {
		panic(fmt.Sprintf("Error creating %s request: %v\n", request.Method, err))
	}
	// set headers
	req.Header = request.Headers
	req.Header.Set("Host", client.origin)

	// client.clog.Debugf("request headers: %v\n", req.Header)

	resp, err := client.handle.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Error sending %s request to %s: %v\n", request.Method, client.origin, err))
	}

	client.clog.Debugf("response status %s\n", resp.Status)
	// client.clog.Debugf("response headers %v\n", resp.Header)

	return ClientReqRes{
		Body:       helpers.ReadBody(resp.Body),
		Headers:    resp.Header,
		RespStatus: resp.Status,
		RespCode:   resp.StatusCode,
	}
}
