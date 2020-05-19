package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestListBackendApi(t *testing.T) {
	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != backendListResourceEndpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", backendListResourceEndpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(helperLoadBytes(t, "backend_api_list_fixture.json"))),
			Header:     make(http.Header),
		}
	})

	credential := "someAccessToken"
	c := NewThreeScale(NewTestAdminPortal(t), credential, httpClient)
	backendApiList, err := c.ListBackendApis()
	if err != nil {
		t.Fatal(err)
	}

	if backendApiList == nil {
		t.Fatal("backend list returned nil")
	}

	if len(backendApiList.Backends) != 3 {
		t.Fatalf("Then number of backend_api's does not match. Expected [%d]; got [%d]", 3, len(backendApiList.Backends))
	}
}

func TestCreateBackendApi(t *testing.T) {
	params := Params{
		"name":             "backendapitest1",
		"private_endpoint": "https://echo-api.3scale.net:443",
	}

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != backendListResourceEndpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", backendListResourceEndpoint, req.URL.Path)
		}

		if req.Method != http.MethodPost {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPost, req.Method)
		}

		backendApi := &BackendApi{
			Element: BackendApiItem{
				ID:              45498,
				Name:            params["name"],
				PrivateEndpoint: params["private_endpoint"],
			},
		}

		responseBodyBytes, err := json.Marshal(backendApi)
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
	backendAPI, err := c.CreateBackendApi(params)
	if err != nil {
		t.Fatal(err)
	}

	if backendAPI == nil {
		t.Fatal("backendapi returned nil")
	}

	if backendAPI.Element.ID != 45498 {
		t.Fatalf("backend_api ID does not match. Expected [%d]; got [%d]", 45498, backendAPI.Element.ID)
	}

	if backendAPI.Element.Name != params["name"] {
		t.Fatalf("backend_api privateEndpoint does not match. Expected [%s]; got [%s]", params["name"], backendAPI.Element.Name)
	}

	if backendAPI.Element.PrivateEndpoint != params["private_endpoint"] {
		t.Fatalf("backend_api privateEndpoint does not match. Expected [%s]; got [%s]", params["private_endpoint"], backendAPI.Element.PrivateEndpoint)
	}
}

func TestDeleteBackendApi(t *testing.T) {
	var backendAAPIID int64 = 12345
	endpoint := fmt.Sprintf(backendResourceEndpoint, backendAAPIID)

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
	err := c.DeleteBackendApi(backendAAPIID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadBackendApi(t *testing.T) {
	var (
		backendapiID int64 = 98765
		endpoint           = fmt.Sprintf(backendResourceEndpoint, backendapiID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		backendApi := &BackendApi{
			Element: BackendApiItem{
				ID:              backendapiID,
				Name:            "someName",
				PrivateEndpoint: "https://example.com:81",
			},
		}

		responseBodyBytes, err := json.Marshal(backendApi)
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
	backendAPI, err := c.BackendApi(backendapiID)
	if err != nil {
		t.Fatal(err)
	}

	if backendAPI == nil {
		t.Fatal("backendapi returned nil")
	}

	if backendAPI.Element.ID != backendapiID {
		t.Fatalf("backend_api ID does not match. Expected [%d]; got [%d]", backendapiID, backendAPI.Element.ID)
	}

	if backendAPI.Element.Name != "someName" {
		t.Fatalf("backend_api privateEndpoint does not match. Expected [%s]; got [%s]", "someName", backendAPI.Element.Name)
	}

	if backendAPI.Element.PrivateEndpoint != "https://example.com:81" {
		t.Fatalf("backend_api privateEndpoint does not match. Expected [%s]; got [%s]", "https://example.com:81", backendAPI.Element.PrivateEndpoint)
	}
}

func TestUpdateBackendApi(t *testing.T) {
	var (
		backendapiID int64 = 98765
		endpoint           = fmt.Sprintf(backendResourceEndpoint, backendapiID)
		params             = Params{"name": "newName"}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPut {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		backendApi := &BackendApi{
			Element: BackendApiItem{
				ID:              backendapiID,
				Name:            params["name"],
				PrivateEndpoint: "https://example.com:81",
			},
		}

		responseBodyBytes, err := json.Marshal(backendApi)
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
	backendAPI, err := c.UpdateBackendApi(backendapiID, params)
	if err != nil {
		t.Fatal(err)
	}

	if backendAPI == nil {
		t.Fatal("backendapi returned nil")
	}

	if backendAPI.Element.ID != backendapiID {
		t.Fatalf("backend_api ID does not match. Expected [%d]; got [%d]", backendapiID, backendAPI.Element.ID)
	}

	if backendAPI.Element.Name != params["name"] {
		t.Fatalf("backend_api privateEndpoint does not match. Expected [%s]; got [%s]", params["name"], backendAPI.Element.Name)
	}
}
