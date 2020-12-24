package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func developerAccount1() DeveloperAccount {
	var accountID int64 = 1

	return DeveloperAccount{
		Element: DeveloperAccountItem{
			ID: &accountID,
		},
	}

}

func developerAccount2() DeveloperAccount {
	var accountID int64 = 2

	return DeveloperAccount{
		Element: DeveloperAccountItem{
			ID: &accountID,
		},
	}

}

func TestListDeveloperAccounts(t *testing.T) {
	var (
		endpoint = developerAccountListResourceEndpoint

		list = DeveloperAccountList{
			Items: []DeveloperAccount{
				developerAccount1(),
				developerAccount2(),
			},
		}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		responseBodyBytes, err := json.Marshal(list)
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
	resp, err := c.ListDeveloperAccounts()
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("response was nil")
	}

	if !reflect.DeepEqual(*resp, list) {
		got, _ := json.Marshal(*resp)
		expected, _ := json.Marshal(list)
		t.Fatalf("Expected %s; got %s", string(expected), string(got))
	}
}

func TestReadDeveloperAccount(t *testing.T) {
	var (
		accountID int64 = 1
		endpoint        = fmt.Sprintf(developerAccountResourceEndpoint, accountID)
		item            = developerAccount1()
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		responseBodyBytes, err := json.Marshal(item)
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
	resp, err := c.DeveloperAccount(accountID)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("response was nil")
	}

	if !reflect.DeepEqual(*resp, item) {
		got, _ := json.Marshal(*resp)
		expected, _ := json.Marshal(item)
		t.Fatalf("Expected %s; got %s", string(expected), string(got))
	}
}

func TestDeveloperAccountSignup(t *testing.T) {
	var (
		endpoint = signupResourceEndpoint
		item     = developerAccount1()
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPost {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPost, req.Method)
		}

		responseBodyBytes, err := json.Marshal(item)
		if err != nil {
			t.Fatal(err)
		}

		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewBuffer(responseBodyBytes)),
			Header:     make(http.Header),
		}
	})

	params := Params{
		"org_name": "otherorg",
		"username": "username01",
		"email":    "username01@otherorg.com",
		"password": "1234",
	}
	credential := "someAccessToken"
	c := NewThreeScale(NewTestAdminPortal(t), credential, httpClient)
	resp, err := c.Signup(params)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("response was nil")
	}

	if !reflect.DeepEqual(*resp, item) {
		got, _ := json.Marshal(*resp)
		expected, _ := json.Marshal(item)
		t.Fatalf("Expected %s; got %s", string(expected), string(got))
	}
}

func TestUpdateDeveloperAccount(t *testing.T) {
	var (
		accountID int64 = 1
		endpoint        = fmt.Sprintf(developerAccountResourceEndpoint, accountID)
		item            = developerAccount1()
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPut {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPut, req.Method)
		}

		responseBodyBytes, err := json.Marshal(item)
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
	resp, err := c.UpdateDeveloperAccount(&item)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("response was nil")
	}

	if !reflect.DeepEqual(*resp, item) {
		got, _ := json.Marshal(*resp)
		expected, _ := json.Marshal(item)
		t.Fatalf("Expected %s; got %s", string(expected), string(got))
	}
}

func TestDeleteDeveloperAccount(t *testing.T) {
	var (
		accountID int64 = 1
		endpoint        = fmt.Sprintf(developerAccountResourceEndpoint, accountID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodDelete {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodDelete, req.Method)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader("")),
			Header:     make(http.Header),
		}
	})

	credential := "someAccessToken"
	c := NewThreeScale(NewTestAdminPortal(t), credential, httpClient)
	err := c.DeleteDeveloperAccount(accountID)
	if err != nil {
		t.Fatal(err)
	}
}
