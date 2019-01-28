package client

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func helperLoadBytes(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}
