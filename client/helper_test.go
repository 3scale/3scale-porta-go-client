package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
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

func fetchBasicAuthHeader(req *http.Request) (string, error) {
	auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		return "", fmt.Errorf("Basic auth header missing or not valid")
	}

	return auth[1], nil
}
