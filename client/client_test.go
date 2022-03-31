package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/3scale/3scale-porta-go-client/fake"
)

func TestNewAdminPortal(t *testing.T) {
	ap, err := NewAdminPortal("https", "www.test.com", 443)
	if err != nil {
		t.Errorf("unexpected error when creating client")
	}
	equals(t, "https://www.test.com:443", ap.rawURL)

	ap, err = NewAdminPortal("https", "www.test.com", 0)
	if err != nil {
		t.Errorf("unexpected error when creating client")
	}
	equals(t, "https://www.test.com", ap.rawURL)

	_, err = NewAdminPortalFromStr("https://www.test.com:443")
	if err != nil {
		t.Errorf("expected nil error when creating with valid url")
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
	ap, err := NewAdminPortalFromStr("https://www.test.com:443")
	if err != nil {
		t.Fail()
	}
	return ap

}

func TestAdminPortalPathIsPreserved(t *testing.T) {
	const scheme = "https"
	const host = "www.test.com"
	const port = "443"
	ap, err := NewAdminPortalFromStr("https://www.test.com:443/example/")
	if err != nil {
		t.Errorf("unexpected err when creating admin portal")
	}
	if ap.rawURL != "https://www.test.com:443/example" {
		t.Errorf("expected trailing slash to be stripped")
	}

	verify := func(req *http.Request, path string) {
		equals(t, host, req.URL.Hostname())
		equals(t, port, req.URL.Port())
		equals(t, scheme, req.URL.Scheme)
		equals(t, path, req.URL.Path)
	}

	// Test path is preserved for GET
	_, err = NewThreeScale(ap, "any", NewTestClient(func(req *http.Request) *http.Response {
		verify(req, "/example/admin/api/services.xml")
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(`<services></services>`))),
		}
	})).ListServices()
	if err != nil {
		t.Fatal(err)
	}

	// Test path is preserved for POST
	_, err = NewThreeScale(ap, "any", NewTestClient(func(req *http.Request) *http.Response {
		verify(req, "/example/admin/api/services.xml")
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(`<service></service>`))),
		}
	})).CreateService("any")
	if err != nil {
		t.Fatal(err)
	}

	// Test path is preserved for DELETE
	err = NewThreeScale(ap, "any", NewTestClient(func(req *http.Request) *http.Response {
		verify(req, "/example/admin/api/services/any.xml")
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(""))),
		}
	})).DeleteService("any")
	if err != nil {
		t.Fatal(err)
	}

	// Test path is preserved for PUT
	_, err = NewThreeScale(ap, "any", NewTestClient(func(req *http.Request) *http.Response {
		verify(req, "/example/admin/api/services/any.xml")
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(`<service></service>`))),
		}
	})).UpdateService("any", NewParams())
	if err != nil {
		t.Fatal(err)
	}
}

func equals(t *testing.T, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		t.FailNow()
	}
}
