package client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestListAccountsOk(t *testing.T) {
	access_token := "someAccessToken"
	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method != http.MethodGet {
			t.Fatalf("wrong helper called for account list api")
		}

		q := req.URL.Query()
		if q.Get("access_token") != access_token {
			t.Fatalf("Expected access token not found")
		}

		bodyReader := bytes.NewReader(helperLoadBytes(t, "accounts_fixture.xml"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bodyReader),
			Header:     make(http.Header),
		}
	})

	c := NewThreeScale(NewTestAdminPortal(t), httpClient)

	accountList, err := c.ListAccounts(access_token)
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
	access_token := "someAccessToken"
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
			c := NewThreeScale(NewTestAdminPortal(t), httpClient)
			_, err := c.ListAccounts(access_token)
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
