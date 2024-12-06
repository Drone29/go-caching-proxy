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
	Method string `json:"-"` // don't need to save this to json
	Uri    string `json:"-"` // don't need to save this to json
	// response only
	RespStatus string `json:"status"`
	RespCode   int    `json:"code"`
}

// create an alias to ensure MarshalJSON() is not called recursively
type aliasRequest Request

type auxRequest struct {
	*aliasRequest
	Body string `json:"body"` //serialize body as string (json will utilize base64)
}

func (r *Request) MarshalJSON() ([]byte, error) {
	aux := &auxRequest{
		aliasRequest: (*aliasRequest)(r),
		Body:         string(r.Body),
	}
	return json.Marshal(aux)
}

func (r *Request) UnmarshalJSON(data []byte) error {
	aux := &auxRequest{
		aliasRequest: (*aliasRequest)(r),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	// convert body to []byte
	r.Body = []byte(aux.Body)
	return nil
}

// backup requests to json file
func Backup(filename string, requests []Request) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	return encoder.Encode(requests)
}

func Restore(filename string) (requests []Request, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&requests)
	return requests, err
}
