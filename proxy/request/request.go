package request

import (
	"encoding/json"
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

// store in json file using this format
type auxRequestWithHost struct {
	Host     string    `json:"host"`
	Requests []Request `json:"requests"`
}

type reqMap = map[string][]Request

// backup requests for single host to json file
func BackupOne(filename, host string, requests []Request) error {
	// read all currently backed up records (ignore errors)
	records, _ := RestoreAll(filename)
	// update map
	records[host] = requests
	return BackupMany(filename, records)
}

// backup multiple requests for multiple hosts
// requests - a map of host-[]requests pairs
func BackupMany(filename string, requests reqMap) error {
	// convert to intermediary types
	backup := []auxRequestWithHost{}
	for k, v := range requests {
		backup = append(backup, auxRequestWithHost{
			Host:     k,
			Requests: v,
		})
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	return encoder.Encode(backup)
}

// restore all requests for all hosts
// requests - a map of host-[]requests pairs
func RestoreAll(filename string) (requests reqMap, err error) {
	requests = make(reqMap)
	records := []auxRequestWithHost{}
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&records)
	for _, v := range records {
		requests[v.Host] = v.Requests
	}
	return requests, err
}
