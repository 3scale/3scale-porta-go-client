package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type clientTenantOperation func(*ThreeScaleClient) error

func helperClientError(t *testing.T, op clientTenantOperation) {
	errorTests := []struct {
		Name                string
		ResponseBodyFixture string
		ExpectedErrorMsg    string
		HTTPStatusCode      int
	}{
		{"CorruptedJsonErrorTest", "accounts_wrong_format_fixture.json", "decoding error - unexpected EOF", 200},
		{"WrongFormatErrorTest", "xml_format_fixture.xml", "decoding error - invalid character", 200},
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
			credential := "someAccessToken"
			c := NewThreeScale(NewTestAdminPortal(t), credential, httpClient)
			err := op(c)
			if err == nil {
				subTest.Fatalf("client operation did not return error")
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
	op := func(c *ThreeScaleClient) error {
		orgName := "someOrgName"
		userName := "someUserName"
		email := "someEmail@example.com"
		password := "somePassword"
		_, err := c.CreateTenant(orgName, userName, email, password)
		return err
	}
	helperClientError(t, op)
}

func TestShowTenantOk(t *testing.T) {
	accessToken := "someAccessToken"
	tenantID := int64(42)
	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method != http.MethodGet {
			t.Fatalf("wrong http method called for create tenant endpoint")
		}

		q := req.URL.Query()
		if q.Get("access_token") != accessToken {
			t.Fatalf("Expected access token not found")
		}

		p := req.URL.Path
		if p != fmt.Sprintf(tenantRead, tenantID) {
			t.Fatalf("Path: expected (%s) found (%s)", fmt.Sprintf(tenantRead, tenantID), p)
		}

		bodyReader := bytes.NewReader(helperLoadBytes(t, "show_tenant_response.json"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bodyReader),
			Header:     make(http.Header),
		}
	})

	c := NewThreeScale(NewTestAdminPortal(t), httpClient)

	tenant, err := c.ShowTenant(accessToken, tenantID)
	if err != nil {
		t.Fatal(err)
	}

	if tenant == nil {
		t.Fatalf("tenant not parsed")
	}

	if tenant.Signup.Account.SupportEmail != "admin@corp24.com" {
		t.Fatalf("Email: expected (%s) found (%s)", "admin@corp24.com", tenant.Signup.Account.SupportEmail)
	}

	if tenant.Signup.Account.OrgName != "Corp24" {
		t.Fatalf("OrgName: expected (%s) found (%s)", "Corp24", tenant.Signup.Account.OrgName)
	}
}

func TestShowTenantErrors(t *testing.T) {
	op := func(c *ThreeScaleClient) error {
		tenantID := int64(42)
		_, err := c.ShowTenant(tenantID)
		return err
	}
	helperClientError(t, op)
}

func TestUpdateTenantOk(t *testing.T) {
	accessToken := "someAccessToken"
	tenantID := int64(42)
	params := Params{
		"support_email": "admin@corp24.com",
		"org_name":      "Corp24",
	}
	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method != http.MethodPut {
			t.Fatalf("wrong http method called for create tenant endpoint")
		}

		p := req.URL.Path
		if p != fmt.Sprintf(tenantRead, tenantID) {
			t.Fatalf("Path: expected (%s) found (%s)", fmt.Sprintf(tenantRead, tenantID), p)
		}

		err := req.ParseForm()
		if err != nil {
			t.Fatal(err)
		}

		if req.Form.Get("access_token") != accessToken {
			t.Fatalf("Expected access token not found")
		}

		// request includes access_token
		if len(req.Form) != len(params)+1 {
			t.Fatalf("Form num params differ: expected (%d) found (%d)", len(params), len(req.Form))
		}

		for paramKey, paramValue := range params {
			if req.Form.Get(paramKey) != paramValue {
				t.Fatalf("field %s: expected (%s) found (%s)", paramKey, paramValue, req.Form.Get(paramKey))
			}
		}

		bodyReader := bytes.NewReader(helperLoadBytes(t, "show_tenant_response.json"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bodyReader),
			Header:     make(http.Header),
		}
	})

	c := NewThreeScale(NewTestAdminPortal(t), accessToken, httpClient)

	tenant, err := c.UpdateTenant(tenantID, params)
	if err != nil {
		t.Fatal(err)
	}

	if tenant == nil {
		t.Fatalf("tenant not parsed")
	}

	if tenant.Signup.Account.SupportEmail != params["support_email"] {
		t.Fatalf("Email: expected (%s) found (%s)", params["support_email"], tenant.Signup.Account.SupportEmail)
	}

	if tenant.Signup.Account.OrgName != params["org_name"] {
		t.Fatalf("OrgName: expected (%s) found (%s)", params["org_name"], tenant.Signup.Account.OrgName)
	}
}

func TestUpdateTenantErrors(t *testing.T) {
	op := func(c *ThreeScaleClient) error {
		accessToken := "someAccessToken"
		tenantID := int64(42)
		params := Params{
			"support_email": "admin@corp24.com",
			"org_name":      "Corp24",
		}
		_, err := c.UpdateTenant(accessToken, tenantID, params)
		return err
	}
	helperClientError(t, op)
}
