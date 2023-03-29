package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func developerAccount1() DeveloperAccount {
	return DeveloperAccount{Element: DeveloperAccountItem{ID: &[]int64{1}[0]}}
}

func developerAccount2() DeveloperAccount {
	return DeveloperAccount{Element: DeveloperAccountItem{ID: &[]int64{2}[0]}}
}

func TestListDeveloperAccounts(t *testing.T) {
	var (
		endpoint = developerAccountListResourceEndpoint
	)

	develAccountGenerator := func(startingIndex, n int) DeveloperAccountList {
		pList := DeveloperAccountList{
			Items: make([]DeveloperAccount, 0, n),
		}

		for idx := 0; idx < n; idx++ {
			pList.Items = append(pList.Items, DeveloperAccount{
				Element: DeveloperAccountItem{ID: &[]int64{int64(idx + startingIndex)}[0]},
			})
		}

		return pList
	}

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		// Will serve: 3 pages
		// page 1 => DEVELOPERACCOUNTS_PER_PAGE
		// page 2 => DEVELOPERACCOUNTS_PER_PAGE
		// page 3 => 51

		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		if req.URL.Query().Get("per_page") != strconv.Itoa(DEVELOPERACCOUNTS_PER_PAGE) {
			t.Fatalf("per_page param does not match. Expected [%d]; got [%s]", DEVELOPERACCOUNTS_PER_PAGE, req.URL.Query().Get("per_page"))
		}

		var list DeveloperAccountList

		if req.URL.Query().Get("page") == "1" {
			list = develAccountGenerator(DEVELOPERACCOUNTS_PER_PAGE*0, DEVELOPERACCOUNTS_PER_PAGE)
		} else if req.URL.Query().Get("page") == "2" {
			list = develAccountGenerator(DEVELOPERACCOUNTS_PER_PAGE*1, DEVELOPERACCOUNTS_PER_PAGE)
		} else if req.URL.Query().Get("page") == "3" {
			list = develAccountGenerator(DEVELOPERACCOUNTS_PER_PAGE*2, 51)
		} else {
			t.Fatalf("page param unexpected value; got [%s]", req.URL.Query().Get("page"))
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
	list, err := c.ListDeveloperAccounts()
	if err != nil {
		t.Fatal(err)
	}
	if list == nil {
		t.Fatal("developer account list returned nil")
	}

	if len(list.Items) != 2*DEVELOPERACCOUNTS_PER_PAGE+51 {
		t.Fatalf("The number of developer accounts does not match. Expected [%d]; got [%d]", 2*DEVELOPERACCOUNTS_PER_PAGE+51, len(list.Items))
	}
}

func TestListDeveloperAccountsPerPage(t *testing.T) {
	var (
		endpoint = developerAccountListResourceEndpoint
	)

	t.Run("page and per_page params used", func(subT *testing.T) {
		var (
			pageNum int = 4
			perPage int = 2
		)
		httpClient := NewTestClient(func(req *http.Request) *http.Response {
			if req.URL.Path != endpoint {
				subT.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
			}

			if req.Method != http.MethodGet {
				subT.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
			}

			if req.URL.Query().Get("page") != strconv.Itoa(pageNum) {
				subT.Fatalf("page param does not match. Expected [%d]; got [%s]", pageNum, req.URL.Query().Get("page"))
			}

			if req.URL.Query().Get("per_page") != strconv.Itoa(perPage) {
				subT.Fatalf("page param does not match. Expected [%d]; got [%s]", perPage, req.URL.Query().Get("per_page"))
			}

			list := DeveloperAccountList{
				Items: []DeveloperAccount{developerAccount1(), developerAccount2()},
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
		c := NewThreeScale(NewTestAdminPortal(subT), credential, httpClient)
		list, err := c.ListDeveloperAccountsPerPage(pageNum, perPage)
		if err != nil {
			subT.Fatal(err)
		}

		if list == nil {
			subT.Fatal("developer account list returned nil")
		}

		if len(list.Items) != 2 {
			subT.Fatalf("Then number of developer accounts does not match. Expected [%d]; got [%d]", 2, len(list.Items))
		}
	})

	t.Run("page and per_page params not used", func(subT *testing.T) {
		httpClient := NewTestClient(func(req *http.Request) *http.Response {
			if req.URL.Path != endpoint {
				subT.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
			}

			if req.Method != http.MethodGet {
				subT.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
			}

			if req.URL.Query().Get("page") != "" {
				subT.Fatalf("Query param page does not match. Expected empty; got [%s]", req.URL.Query().Get("page"))
			}

			if req.URL.Query().Get("per_page") != "" {
				subT.Fatalf("page param does not match. Expected empty; got [%s]", req.URL.Query().Get("per_page"))
			}

			list := DeveloperAccountList{
				Items: []DeveloperAccount{developerAccount1(), developerAccount2()},
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
		c := NewThreeScale(NewTestAdminPortal(subT), credential, httpClient)
		list, err := c.ListDeveloperAccountsPerPage()
		if err != nil {
			subT.Fatal(err)
		}

		if list == nil {
			subT.Fatal("developer account list returned nil")
		}

		if len(list.Items) != 2 {
			subT.Fatalf("Then number of developer account items does not match. Expected [%d]; got [%d]", 2, len(list.Items))
		}
	})
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
