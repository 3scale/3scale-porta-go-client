package client

import (
	"bytes"
	"encoding/json"
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

		bodyReader := bytes.NewReader(helperLoadBytes(t, "accounts_fixture.json"))
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

	expectedNumAccounts := 5
	if len(accountList.Accounts) != expectedNumAccounts {
		t.Fatalf("expected number of accounts was %d", expectedNumAccounts)
	}

	attrTests := []struct {
		Name          string
		ValidateParam func(account Account) bool
		ErrorMsg      string
	}{
		{
			Name: "testID",
			ValidateParam: func(account Account) bool {
				return account.ID == 2
			},
			ErrorMsg: "ID did not match",
		},
		{
			Name: "testState",
			ValidateParam: func(account Account) bool {
				return account.State == "approved"
			},
			ErrorMsg: "STATE did not match",
		},
		{
			Name: "testOrgName",
			ValidateParam: func(account Account) bool {
				return account.OrgName == "Provider Name"
			},
			ErrorMsg: "OrgName did not match",
		},
		{
			Name: "testSupportEmail",
			ValidateParam: func(account Account) bool {
				return account.SupportEmail == "admin@3scale.amp24.127.0.0.1.nip.io"
			},
			ErrorMsg: "SupportEmail did not match",
		},
		{
			Name: "testAdminDomain",
			ValidateParam: func(account Account) bool {
				return account.AdminDomain == "3scale-admin.amp24.127.0.0.1.nip.io"
			},
			ErrorMsg: "AdminDomain did not match",
		},
		{
			Name: "testDomain",
			ValidateParam: func(account Account) bool {
				return account.Domain == "3scale.amp24.127.0.0.1.nip.io"
			},
			ErrorMsg: "Domain did not match",
		},
	}

	account := accountList.Accounts[1].Account

	for _, tt := range attrTests {
		t.Run(tt.Name, func(subTest *testing.T) {
			if !tt.ValidateParam(account) {
				subTest.Fatalf(tt.ErrorMsg)
			}
		})
	}
}

func TestListAccountsErrors(t *testing.T) {
	credential := "someAccessToken"
	errorTests := []struct {
		Name                string
		ResponseBodyFixture string
		ExpectedErrorMsg    string
	}{
		{"CorruptedJsonErrorTest", "accounts_wrong_format_fixture.json", "decoding error - unexpected EOF"},
		{"WrongFormatErrorTest", "xml_format_fixture.xml", "decoding error - invalid character"},
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
				subTest.Fatalf("account list did not return error")
			}

			apiError, ok := err.(ApiErr)
			if !ok {
				subTest.Fatalf("expected ApiErr error type")
			}

			if !strings.Contains(apiError.Error(), tt.ExpectedErrorMsg) {
				subTest.Fatalf("Expected [%s]: got [%s] ", tt.ExpectedErrorMsg, apiError.Error())
			}
		})
	}
}

func TestFindAccount(t *testing.T) {
	var (
		accountID int64 = 3
		username        = "John"
		endpoint        = findAccount
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		account := &AccountElem{
			Account{
				ID: accountID,
			},
		}
		responseBodyBytes, err := json.Marshal(account)
		if err != nil {
			t.Fatal(err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer(responseBodyBytes)),
			Header:     make(http.Header),
		}
	})

	credential := "someAccessToken"
	c := NewThreeScale(NewTestAdminPortal(t), credential, httpClient)
	obj, err := c.FindAccount(username)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("application returned nil")
	}

	if obj.ID != accountID {
		t.Fatalf("obj state does not match. Expected [%d]; got [%d]", accountID, obj.ID)
	}
}
