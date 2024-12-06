package request

// acts as both request and response
type Request struct {
	Body    []byte
	Headers map[string][]string
	// request only
	Method string
	Uri    string
	// response only
	RespStatus string
	RespCode   int
}
