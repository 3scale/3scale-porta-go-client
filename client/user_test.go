package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestActivateUserOk(t *testing.T) {
	accessToken := "someAccessToken"
	accountID := int64(43)
	userID := int64(86)
	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.Method != http.MethodPut {
			t.Fatalf("wrong http method")
		}

		basicAuthValue, err := fetchBasicAuthHeader(req)
		if err != nil {
			t.Error(err)
		}

		if basicAuthValue != basicAuth("", accessToken) {
			t.Fatalf("Expected access token not found")
		}

		bodyReader := bytes.NewReader(helperLoadBytes(t, "user_response_fixture.json"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bodyReader),
			Header:     make(http.Header),
		}
	})

	c := NewThreeScale(NewTestAdminPortal(t), accessToken, httpClient)

	err := c.ActivateUser(accountID, userID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestActivateUserErrors(t *testing.T) {
	accessToken := "someAccessToken"
	accountID := int64(43)
	userID := int64(86)
	errorTests := []struct {
		Name                string
		ResponseBodyFixture string
		ExpectedErrorMsg    string
		HTTPStatusCode      int
	}{
		{"UnexpectedHTTPStatusCode", "error_response_fixture.json",
			"Test Error", 404},
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
			c := NewThreeScale(NewTestAdminPortal(t), accessToken, httpClient)
			err := c.ActivateUser(accountID, userID)
			if err == nil {
				subTest.Fatalf("activate user did not return error")
			}

			apiError, ok := err.(ApiErr)
			if !ok {
				subTest.Fatalf("expected ApiErr error type")
			}

			if !strings.Contains(apiError.Error(), tt.ExpectedErrorMsg) {
				subTest.Fatalf("got: %s, expected: %s", apiError.Error(), tt.ExpectedErrorMsg)
			}
		})
	}
}

func TestReadUser(t *testing.T) {
	const (
		accessToken = "someAccessToken"
		accountID   = int64(321)
		userID      = int64(74)
	)

	inputs := []struct {
		Name             string
		ExpectErr        bool
		ResponseCode     int
		ResponseBodyFile string
		ExpectedErrorMsg string
	}{
		{
			Name:             "ReadUserOK",
			ExpectErr:        false,
			ResponseCode:     200,
			ResponseBodyFile: "user_response_fixture.json",
			ExpectedErrorMsg: "",
		},
		{
			Name:             "ReadUserErr",
			ExpectErr:        true,
			ResponseCode:     400,
			ResponseBodyFile: "error_response_fixture.json",
			ExpectedErrorMsg: "Test Error",
		},
	}

	for _, input := range inputs {
		httpClient := NewTestClient(func(req *http.Request) *http.Response {
			if req.Method != http.MethodGet {
				t.Fatal("wrong helper called")
			}

			basicAuthValue, err := fetchBasicAuthHeader(req)
			if err != nil {
				t.Error(err)
			}

			if basicAuthValue != basicAuth("", accessToken) {
				t.Fatalf("Expected access token not found")
			}

			if req.URL.Path != fmt.Sprintf(userRead, accountID, userID) {
				t.Fatal("wrong url generated")
			}

			bodyReader := bytes.NewReader(helperLoadBytes(t, input.ResponseBodyFile))
			return &http.Response{
				StatusCode: input.ResponseCode,
				Body:       ioutil.NopCloser(bodyReader),
				Header:     make(http.Header),
			}
		})

		c := NewThreeScale(NewTestAdminPortal(t), accessToken, httpClient)

		t.Run(input.Name, func(subTest *testing.T) {
			user, err := c.ReadUser(accountID, userID)
			if input.ExpectErr {
				if err == nil {
					subTest.Fatalf("client operation did not return error")
				}

				apiError, ok := err.(ApiErr)
				if !ok {
					subTest.Fatalf("expected ApiErr error type")
				}

				if !strings.Contains(apiError.Error(), input.ExpectedErrorMsg) {
					subTest.Fatalf("Expected [%s]: got [%s] ", input.ExpectedErrorMsg, apiError.Error())
				}
			} else {
				if err != nil {
					subTest.Fatal(err)
				}
				if user == nil {
					subTest.Fatalf("user nil")
				}
				if user.ID != userID {
					subTest.Fatalf("user attrs not parsed")
				}
			}
		})
	}
}

func TestListUser(t *testing.T) {
	const (
		accessToken = "someAccessToken"
		accountID   = int64(321)
	)

	inputs := []struct {
		Name             string
		FilterParams     Params
		ExpectErr        bool
		ResponseCode     int
		ResponseBodyFile string
		ExpectedErrorMsg string
	}{
		{
			Name:             "ListUserOK",
			FilterParams:     Params{},
			ExpectErr:        false,
			ResponseCode:     200,
			ResponseBodyFile: "user_list_response_fixture.json",
			ExpectedErrorMsg: "",
		},
		{
			Name:             "ListUserNilParamsOK",
			FilterParams:     nil,
			ExpectErr:        false,
			ResponseCode:     200,
			ResponseBodyFile: "user_list_response_fixture.json",
			ExpectedErrorMsg: "",
		},
		{
			Name:             "ListUserErr",
			ExpectErr:        true,
			ResponseCode:     400,
			ResponseBodyFile: "error_response_fixture.json",
			ExpectedErrorMsg: "Test Error",
		},
	}

	for _, input := range inputs {
		httpClient := NewTestClient(func(req *http.Request) *http.Response {
			if req.Method != http.MethodGet {
				t.Fatal("wrong helper called")
			}

			basicAuthValue, err := fetchBasicAuthHeader(req)
			if err != nil {
				t.Error(err)
			}

			if basicAuthValue != basicAuth("", accessToken) {
				t.Fatalf("Expected access token not found")
			}

			if req.URL.Path != fmt.Sprintf(userList, accountID) {
				t.Fatal("wrong url generated")
			}

			bodyReader := bytes.NewReader(helperLoadBytes(t, input.ResponseBodyFile))
			return &http.Response{
				StatusCode: input.ResponseCode,
				Body:       ioutil.NopCloser(bodyReader),
				Header:     make(http.Header),
			}
		})

		c := NewThreeScale(NewTestAdminPortal(t), accessToken, httpClient)

		t.Run(input.Name, func(subTest *testing.T) {
			userList, err := c.ListUsers(accountID, input.FilterParams)
			if input.ExpectErr {
				if err == nil {
					subTest.Fatalf("client operation did not return error")
				}

				apiError, ok := err.(ApiErr)
				if !ok {
					subTest.Fatalf("expected ApiErr error type")
				}

				if !strings.Contains(apiError.Error(), input.ExpectedErrorMsg) {
					subTest.Fatalf("Expected [%s]: got [%s] ", input.ExpectedErrorMsg, apiError.Error())
				}
			} else {
				if err != nil {
					subTest.Fatal(err)
				}
				if userList == nil {
					subTest.Fatalf("user list nil")
				}
				if len(userList.Users) == 0 {
					subTest.Fatalf("user list attrs not parsed")
				}
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	const (
		accessToken = "someAccessToken"
		accountID   = int64(321)
		userID      = int64(74)
	)

	inputs := []struct {
		Name             string
		FilterParams     Params
		ExpectErr        bool
		ResponseCode     int
		ResponseBodyFile string
		ExpectedErrorMsg string
	}{
		{
			Name:             "UpdateUserOK",
			FilterParams:     Params{},
			ExpectErr:        false,
			ResponseCode:     200,
			ResponseBodyFile: "user_response_fixture.json",
			ExpectedErrorMsg: "",
		},
		{
			Name:             "UpdateUserNilParamsOK",
			FilterParams:     nil,
			ExpectErr:        false,
			ResponseCode:     200,
			ResponseBodyFile: "user_response_fixture.json",
			ExpectedErrorMsg: "",
		},
		{
			Name:             "UpdateUserErr",
			ExpectErr:        true,
			ResponseCode:     400,
			ResponseBodyFile: "error_response_fixture.json",
			ExpectedErrorMsg: "Test Error",
		},
	}

	for _, input := range inputs {
		httpClient := NewTestClient(func(req *http.Request) *http.Response {
			if req.Method != http.MethodPut {
				t.Fatal("wrong helper called")
			}

			basicAuthValue, err := fetchBasicAuthHeader(req)
			if err != nil {
				t.Error(err)
			}

			if basicAuthValue != basicAuth("", accessToken) {
				t.Fatalf("Expected access token not found")
			}

			if req.URL.Path != fmt.Sprintf(userUpdate, accountID, userID) {
				t.Fatal("wrong url generated")
			}

			bodyReader := bytes.NewReader(helperLoadBytes(t, input.ResponseBodyFile))
			return &http.Response{
				StatusCode: input.ResponseCode,
				Body:       ioutil.NopCloser(bodyReader),
				Header:     make(http.Header),
			}
		})

		c := NewThreeScale(NewTestAdminPortal(t), accessToken, httpClient)

		t.Run(input.Name, func(subTest *testing.T) {
			user, err := c.UpdateUser(accountID, userID, input.FilterParams)
			if input.ExpectErr {
				if err == nil {
					subTest.Fatalf("client operation did not return error")
				}

				apiError, ok := err.(ApiErr)
				if !ok {
					subTest.Fatalf("expected ApiErr error type")
				}

				if !strings.Contains(apiError.Error(), input.ExpectedErrorMsg) {
					subTest.Fatalf("Expected [%s]: got [%s] ", input.ExpectedErrorMsg, apiError.Error())
				}
			} else {
				if err != nil {
					subTest.Fatal(err)
				}
				if user == nil {
					subTest.Fatalf("user nil")
				}
				if user.ID != userID {
					subTest.Fatalf("user attrs not parsed")
				}
			}
		})
	}
}
