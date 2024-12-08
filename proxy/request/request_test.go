package request

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBackupRestore(t *testing.T) {
	// Sample data
	requests := []Request{
		{
			Body:       []byte("Hello, World!"),
			Headers:    map[string][]string{"Content-Type": {"application/json"}},
			Method:     "GET",
			Uri:        "/api/test",
			RespStatus: "200 OK",
			RespCode:   200,
		},
		{
			Body:       []byte("Another request"),
			Headers:    map[string][]string{"Authorization": {"Bearer token"}},
			Method:     "POST",
			Uri:        "/api/create",
			RespStatus: "201 Created",
			RespCode:   201,
		},
	}

	// bytes.Buffer as in-memory file
	var buf bytes.Buffer
	// abckup
	if err := backup_io(&buf, requests); err != nil {
		t.Fatalf("Backup failed %v", err)
	}
	// restore
	restored, err := restore_io(&buf)
	if err != nil {
		t.Fatalf("Restore failed %v", err)
	}
	// check equality
	if !reflect.DeepEqual(requests, restored) {
		t.Errorf("Restored data does not match original.\nOriginal: %+v\nRestored: %+v", requests, restored)
	}
}
