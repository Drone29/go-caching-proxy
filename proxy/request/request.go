package request

import (
	"encoding/json"
	"io"
	"os"
)

// acts as both request and response
type Request struct {
	Body    []byte              `json:"body"`
	Headers map[string][]string `json:"headers"`
	// request only
	Method string `json:"method"`
	Uri    string `json:"uri"`
	// response only
	RespStatus string `json:"status"`
	RespCode   int    `json:"code"`
}

// backup multiple requests for multiple hosts
// requests - a map of host-[]requests pairs
func Backup(filename string, requests []Request) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return backup_io(f, requests)
}

// takes generic io.Writer, useful for tests
func backup_io(w io.Writer, requests []Request) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(requests)
}

// restore all requests for all hosts
// requests - a map of host-[]requests pairs
func Restore(filename string) (requests []Request, err error) {
	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// create new file
			// populate file with empty data
			err = Backup(filename, []Request{})
		}
		return
	}
	defer f.Close()
	return restore_io(f)
}

// takes generic io.Reader, useful for tests
func restore_io(r io.Reader) (requests []Request, err error) {
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&requests)
	return requests, err
}
