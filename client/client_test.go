package client

import (
	"bytes"
	"github.com/3scale/3scale-porta-go-client/fake"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestNewAdminPortal(t *testing.T) {
	_, err := NewAdminPortal("https", "www.test.com", 443)
	if err != nil {
		t.Fatalf("unexpected error when creating client")
	}
}

func TestHandleJsonResp(t *testing.T) {
	var pce ProxyConfigElement
	resp := fake.GetProxyConfigLatestSuccess()
	err := handleJsonResp(resp, http.StatusOK, &pce)
	if err != nil {
		t.Fatal("unexpected error")
	}
	if pce.ProxyConfig.ID != 54321 {
		t.Fatal("unexpected element id")
	}

	//test error handling
	resp = &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error": "invalid environment"}`)),
		Header:     make(http.Header),
	}
	var pceErrTest ProxyConfigElement
	err = handleJsonResp(resp, http.StatusOK, &pceErrTest)
	if err == nil {
		t.Fatal("expected error not nil")
	}
	if err.Error() != "error calling 3scale system - reason: error - invalid environment - code: 400" {
		t.Fatal("unexpected decoding or error message")
	}

}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func NewTestAdminPortal(t *testing.T) *AdminPortal {
	t.Helper()
	ap, err := NewAdminPortal("https", "www.test.com", 443)
	if err != nil {
		t.Fail()
	}
	return ap

}
