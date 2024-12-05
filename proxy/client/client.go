package client

import (
	"bytes"
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
func New(remote string) *Client {
	host, scheme := validateOrigin(remote)
	return &Client{
		handle: &http.Client{},
		origin: host,
		scheme: scheme,
	}
}

// send client request
func (client *Client) SendRequest(request *ClientReqRes) ClientReqRes {
	request.Uri = client.scheme + client.origin + request.Uri

	req, err := http.NewRequest(request.Method, request.Uri, bytes.NewBuffer(request.Body))
	if err != nil {
		panic(fmt.Sprintf("Error creating %s request: %v\n", request.Method, err))
	}
	// set headers
	req.Header = request.Headers
	req.Header.Set("Host", client.origin)

	resp, err := client.handle.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Error sending %s request to %s: %v\n", request.Method, client.origin, err))
	}

	return ClientReqRes{
		Body:       helpers.ReadBody(resp.Body),
		Headers:    resp.Header,
		RespStatus: resp.Status,
		RespCode:   resp.StatusCode,
	}
}
