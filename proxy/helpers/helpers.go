package helpers

import (
	"fmt"
	"io"
)

func ReadBody(body io.ReadCloser) []byte {
	content, err := io.ReadAll(body)
	defer body.Close()
	if err != nil {
		panic(fmt.Sprintf("Error reading body:\n%v", err))
	}
	return content
}
