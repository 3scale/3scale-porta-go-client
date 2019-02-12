package client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestListAccountsOk(t *testing.T) {
	credential := "someAccessToken"
	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method != http.MethodGet {
			t.Fatalf("wrong helper called for account list api")
		}

		auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			t.Fatalf("Basic auth header missing or not valid")
		}

		expectedAuth := basicAuth("", credential)
		if auth[1] != expectedAuth {
			t.Fatalf("Invalid authorization header value, expected %s got %s", expectedAuth, auth[1])
		}

		bodyReader := bytes.NewReader(helperLoadBytes(t, "accounts_fixture.xml"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bodyReader),
			Header:     make(http.Header),
		}
	})

	c := NewThreeScale(NewTestAdminPortal(t), credential, httpClient)

	accountList, err := c.ListAccounts()
	if err != nil {
		t.Fatal(err)
	}

	if accountList == nil {
		t.Fatalf("account list not parsed")
	}

	expectedNumAccounts := 7
	if len(accountList.Accounts) != expectedNumAccounts {
		t.Fatalf("expected number of accounts was %d", expectedNumAccounts)
	}
}

func TestListAccountsErrors(t *testing.T) {
	credential := "someAccessToken"
	errorTests := []struct {
		Name                string
		ResponseBodyFixture string
		ExpectedErrorMsg    string
		ErrorMsg            string
	}{
		{"BadSyntaxTest", "accounts_wrong_format_fixture.xml",
			"XML syntax error", "expected error type is XML syntax error"},
		{"DecodingErrorTest", "accounts_json_format_fixture.json",
			"decoding error", "expected error type is decoding error"},
	}

	for _, tt := range errorTests {
		t.Run(tt.Name, func(subTest *testing.T) {
			httpClient := NewTestClient(func(req *http.Request) *http.Response {
				bodyReader := bytes.NewReader(helperLoadBytes(subTest, tt.ResponseBodyFixture))

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bodyReader),
					Header:     make(http.Header),
				}
			})
			c := NewThreeScale(NewTestAdminPortal(t), credential, httpClient)
			_, err := c.ListAccounts()
			if err == nil {
				t.Fatalf("account list did not return error")
			}

			apiError, ok := err.(ApiErr)
			if !ok {
				t.Fatalf("expected ApiErr error type")
			}

			if !strings.Contains(apiError.Error(), tt.ExpectedErrorMsg) {
				t.Fatalf(tt.ErrorMsg)
			}
		})
	}
}
