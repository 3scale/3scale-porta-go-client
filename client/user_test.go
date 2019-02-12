package client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestActivateUserOk(t *testing.T) {
	credential := "someAccessToken"
	accountID := "someAccountID"
	userID := "someUserID"
	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method != http.MethodPut {
			t.Fatalf("wrong http method")
		}

		err := req.ParseForm()
		if err != nil {
			t.Fatal(err)
		}

		auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			t.Fatalf("Basic auth header missing or not valid")
		}

		expectedAuth := basicAuth("", credential)
		if auth[1] != expectedAuth {
			t.Fatalf("Invalid authorization header value, expected %s got %s", expectedAuth, auth[1])
		}

		bodyReader := bytes.NewReader(helperLoadBytes(t, "user_response_fixture.xml"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bodyReader),
			Header:     make(http.Header),
		}
	})

	c := NewThreeScale(NewTestAdminPortal(t), credential, httpClient)

	err := c.ActivateUser(accountID, userID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestActivateUserErrors(t *testing.T) {
	credential := "someAccessToken"
	accountID := "someAccountID"
	userID := "someUserID"
	errorTests := []struct {
		Name                string
		ResponseBodyFixture string
		ExpectedErrorMsg    string
		ErrorMsg            string
		HTTPStatusCode      int
	}{
		{"UnexpectedHTTPStatusCode", "error_response_fixture.xml",
			"Test Error", "expected error type is HTTP status error", 404},
	}

	for _, tt := range errorTests {
		t.Run(tt.Name, func(subTest *testing.T) {
			httpClient := NewTestClient(func(req *http.Request) *http.Response {
				bodyReader := bytes.NewReader(helperLoadBytes(subTest, tt.ResponseBodyFixture))

				return &http.Response{
					StatusCode: tt.HTTPStatusCode,
					Body:       ioutil.NopCloser(bodyReader),
					Header:     make(http.Header),
				}
			})
			c := NewThreeScale(NewTestAdminPortal(t), credential, httpClient)
			err := c.ActivateUser(accountID, userID)
			if err == nil {
				subTest.Fatalf("activate user did not return error")
			}

			apiError, ok := err.(ApiErr)
			if !ok {
				subTest.Fatalf("expected ApiErr error type")
			}

			if !strings.Contains(apiError.Error(), tt.ExpectedErrorMsg) {
				subTest.Fatalf("got: %s, expected: %s", apiError.Error(), tt.ErrorMsg)
			}
		})
	}
}
