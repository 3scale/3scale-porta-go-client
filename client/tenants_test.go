package client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestCreateTenantOk(t *testing.T) {
	credential := "someAccessToken"
	orgName := "someOrgName"
	userName := "someUserName"
	email := "someEmail@example.com"
	password := "somePassword"
	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method != http.MethodPost {
			t.Fatalf("wrong http method called for create tenant endpoint")
		}

		err := req.ParseForm()
		if err != nil {
			t.Fatal(err)
		}

		params := []struct {
			ParamName          string
			ParamExpectedValue string
		}{
			{"org_name", orgName},
			{"username", userName},
			{"email", email},
			{"password", password},
		}

		for _, param := range params {
			if req.Form.Get(param.ParamName) != param.ParamExpectedValue {
				t.Fatalf("field %s: expected (%s) found (%s)", param.ParamName, param.ParamExpectedValue, req.Form.Get(param.ParamName))
			}
		}

		bodyReader := bytes.NewReader(helperLoadBytes(t, "create_tenant_response.xml"))
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bodyReader),
			Header:     make(http.Header),
		}
	})

	c := NewThreeScale(NewTestAdminPortal(t), httpClient)

	signup, err := c.CreateTenant(credential, orgName, userName, email, password)
	if err != nil {
		t.Fatal(err)
	}

	if signup == nil {
		t.Fatalf("signup not parsed")
	}

	if signup.Account.SupportEmail != email {
		t.Fatalf("Email: expected (%s) found (%s)", email, signup.Account.SupportEmail)
	}

	if signup.Account.OrgName != orgName {
		t.Fatalf("OrgName: expected (%s) found (%s)", orgName, signup.Account.OrgName)
	}
}

func TestCreateTenantErrors(t *testing.T) {
	credential := "someAccessToken"
	orgName := "someOrgName"
	userName := "someUserName"
	email := "someEmail@example.com"
	password := "somePassword"
	errorTests := []struct {
		Name                string
		ResponseBodyFixture string
		ExpectedErrorMsg    string
		ErrorMsg            string
		HTTPStatusCode      int
	}{
		{"BadSyntaxTest", "accounts_wrong_format_fixture.xml",
			"XML syntax error", "expected error type is XML syntax error", 201},
		{"DecodingErrorTest", "accounts_json_format_fixture.json",
			"decoding error", "expected error type is decoding error", 201},
		{"UnexpectedHTTPStatusCode", "error_response_fixture.xml",
			"Test Error", "expected error type is HTTP status error", 200},
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
			c := NewThreeScale(NewTestAdminPortal(t), httpClient)
			_, err := c.CreateTenant(credential, orgName, userName, email, password)
			if err == nil {
				subTest.Fatalf("create tenant did not return error")
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
