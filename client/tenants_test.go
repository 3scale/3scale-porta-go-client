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

		bodyReader := bytes.NewReader(helperLoadBytes(t, "create_tenant_response.json"))
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bodyReader),
			Header:     make(http.Header),
		}
	})

	c := NewThreeScale(NewTestAdminPortal(t), credential, httpClient)

	signup, err := c.CreateTenant(orgName, userName, email, password)
	if err != nil {
		t.Fatal(err)
	}

	if signup == nil {
		t.Fatalf("signup not parsed")
	}

	if signup.Signup.Account.SupportEmail != email {
		t.Fatalf("Email: expected (%s) found (%s)", email, signup.Signup.Account.SupportEmail)
	}

	if signup.Signup.Account.OrgName != orgName {
		t.Fatalf("OrgName: expected (%s) found (%s)", orgName, signup.Signup.Account.OrgName)
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
		HTTPStatusCode      int
	}{
		{"CorruptedJsonErrorTest", "accounts_wrong_format_fixture.json", "decoding error - unexpected EOF", 201},
		{"WrongFormatErrorTest", "xml_format_fixture.xml", "decoding error - invalid character", 201},
		{"UnexpectedHTTPStatusCode", "error_response_fixture.json", "Test Error", 400},
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
			c := NewThreeScale(NewTestAdminPortal(t),credential, httpClient)
			_, err := c.CreateTenant(orgName, userName, email, password)
			if err == nil {
				subTest.Fatalf("create tenant did not return error")
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
