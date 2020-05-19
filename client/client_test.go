package client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/3scale/3scale-porta-go-client/fake"
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
	expectedErr := `error calling 3scale system - reason: {"error": "invalid environment"} - code: 400`
	if err.Error() != expectedErr {
		t.Fatalf("Expected error: [%s]; got [%s]", expectedErr, err.Error())
	}

}

func TestHandleJsonRespNilBody(t *testing.T) {
	emptyResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader("")),
		Header:     make(http.Header),
	}
	err := handleJsonResp(emptyResp, http.StatusOK, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandleJsonErrResp(mainT *testing.T) {
	inputs := []struct {
		name             string
		responseErr      *http.Response
		expectedCode     int
		expectedErrorMsg string
	}{
		{"unprocessableEntityError", fake.CreateStatusUnprocessableEntityError(),
			http.StatusUnprocessableEntity,
			`error calling 3scale system - reason: {"system_name":["has already been taken"]} - code: 422`,
		},
		{"ForbidenError", fake.CreateAppError(),
			http.StatusForbidden,
			`error calling 3scale system - reason: { "error": "Your access token does not have the correct permissions" } - code: 403`,
		},
	}

	for _, input := range inputs {
		mainT.Run(input.name, func(t *testing.T) {
			err := handleJsonErrResp(input.responseErr)
			if err == nil {
				t.Fatal("error expected")
			}

			apiErr, ok := err.(ApiErr)
			if !ok {
				t.Fatalf("error is not ApiErr type: %T", err)
			}

			if apiErr.Code() != input.expectedCode {
				t.Fatalf("Expected code: [%d]; got [%d]", input.expectedCode, apiErr.Code())
			}

			if err.Error() != input.expectedErrorMsg {
				t.Fatalf("Expected error: [%s]; got [%s]", input.expectedErrorMsg, err.Error())
			}
		})
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
