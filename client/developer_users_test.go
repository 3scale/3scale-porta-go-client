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

func developerUser1() DeveloperUser {
	var userID int64 = 1

	return DeveloperUser{
		Element: DeveloperUserItem{
			ID: &userID,
		},
	}

}

func developerUser2() DeveloperUser {
	var userID int64 = 2

	return DeveloperUser{
		Element: DeveloperUserItem{
			ID: &userID,
		},
	}

}

func TestListDeveloperUsers(t *testing.T) {
	var (
		accountID int64 = 12

		endpoint = fmt.Sprintf(developerUserListResourceEndpoint, accountID)

		list = DeveloperUserList{
			Items: []DeveloperUser{
				developerUser1(),
				developerUser2(),
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
	resp, err := c.ListDeveloperUsers(accountID, nil)
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

func TestDeveloperUser(t *testing.T) {
	var (
		accountID int64 = 12
		userID    int64 = 1
		endpoint        = fmt.Sprintf(developerUserResourceEndpoint, accountID, userID)
		item            = developerUser1()
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
	resp, err := c.DeveloperUser(accountID, userID)
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

func TestUpdateDeveloperUser(t *testing.T) {
	var (
		accountID int64 = 12
		userID    int64 = 1
		endpoint        = fmt.Sprintf(developerUserResourceEndpoint, accountID, userID)
		item            = developerUser1()
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
	resp, err := c.UpdateDeveloperUser(accountID, &item)
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

func TestDeleteDeveloperUser(t *testing.T) {
	var (
		accountID int64 = 12
		userID    int64 = 1
		endpoint        = fmt.Sprintf(developerUserResourceEndpoint, accountID, userID)
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
	err := c.DeleteDeveloperUser(accountID, userID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestActivateDeveloperUser(t *testing.T) {
	var (
		accountID int64 = 12
		userID    int64 = 1
		endpoint        = fmt.Sprintf(developerUserActivateEndpoint, accountID, userID)
		item            = developerUser1()
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
	resp, err := c.ActivateDeveloperUser(accountID, userID)
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

func TestCreateDeveloperUser(t *testing.T) {
	var (
		accountID int64 = 12
		endpoint        = fmt.Sprintf(developerUserListResourceEndpoint, accountID)
		item            = developerUser1()
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

	credential := "someAccessToken"
	c := NewThreeScale(NewTestAdminPortal(t), credential, httpClient)
	resp, err := c.CreateDeveloperUser(accountID, &item)
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

func TestChangeRoleToMemberDeveloperUser(t *testing.T) {
	var (
		accountID int64 = 12
		userID    int64 = 1
		endpoint        = fmt.Sprintf(developerUserMemberResourceEndpoint, accountID, userID)
		item            = developerUser1()
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
	resp, err := c.ChangeRoleToMemberDeveloperUser(accountID, userID)
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

func TestChangeRoleToAdminDeveloperUser(t *testing.T) {
	var (
		accountID int64 = 12
		userID    int64 = 1
		endpoint        = fmt.Sprintf(developerUserAdminResourceEndpoint, accountID, userID)
		item            = developerUser1()
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
	resp, err := c.ChangeRoleToAdminDeveloperUser(accountID, userID)
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

func TestSuspendDeveloperUser(t *testing.T) {
	var (
		accountID int64 = 12
		userID    int64 = 1
		endpoint        = fmt.Sprintf(developerUserSuspendResourceEndpoint, accountID, userID)
		item            = developerUser1()
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
	resp, err := c.SuspendDeveloperUser(accountID, userID)
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

func TestUnsuspendDeveloperUser(t *testing.T) {
	var (
		accountID int64 = 12
		userID    int64 = 1
		endpoint        = fmt.Sprintf(developerUserUnsuspendResourceEndpoint, accountID, userID)
		item            = developerUser1()
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
	resp, err := c.UnsuspendDeveloperUser(accountID, userID)
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
