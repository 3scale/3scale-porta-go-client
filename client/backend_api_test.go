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

func TestListBackendApiMethods(t *testing.T) {
	var (
		backendapiID int64 = 98765
		hitsID       int64 = 1
		endpoint           = fmt.Sprintf(backendMethodListResourceEndpoint, backendapiID, hitsID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		methodList := &MethodList{
			Methods: []Method{
				{
					Element: MethodItem{
						ID:         1,
						Name:       "method01",
						SystemName: "method01",
					},
				},
				{
					Element: MethodItem{
						ID:         2,
						Name:       "method02",
						SystemName: "method02",
					},
				},
			},
		}

		responseBodyBytes, err := json.Marshal(methodList)
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
	list, err := c.ListBackendapiMethods(backendapiID, hitsID)
	if err != nil {
		t.Fatal(err)
	}

	if list == nil {
		t.Fatal("backend method list returned nil")
	}

	if len(list.Methods) != 2 {
		t.Fatalf("Then number of backend_api method's does not match. Expected [%d]; got [%d]", 2, len(list.Methods))
	}
}

func TestCreateBackendApiMethod(t *testing.T) {
	var (
		backendapiID int64 = 98765
		hitsID       int64 = 1
		endpoint           = fmt.Sprintf(backendMethodListResourceEndpoint, backendapiID, hitsID)
		params             = Params{"friendly_name": "method5"}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPost {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPost, req.Method)
		}

		method := &Method{
			Element: MethodItem{
				ID:         10,
				Name:       params["friendly_name"],
				SystemName: params["friendly_name"],
			},
		}

		responseBodyBytes, err := json.Marshal(method)
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
	obj, err := c.CreateBackendApiMethod(backendapiID, hitsID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("backend method create returned nil")
	}

	if obj.Element.ID != 10 {
		t.Fatalf("backend_api method ID does not match. Expected [%d]; got [%d]", 10, obj.Element.ID)
	}

	if obj.Element.Name != params["friendly_name"] {
		t.Fatalf("backend_api privateEndpoint does not match. Expected [%s]; got [%s]", params["friendly_name"], obj.Element.Name)
	}

}

func TestDeleteBackendApiMethod(t *testing.T) {
	var (
		backendapiID int64 = 98765
		hitsID       int64 = 1
		methodID     int64 = 123325
		endpoint           = fmt.Sprintf(backendMethodResourceEndpoint, backendapiID, hitsID, methodID)
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
	err := c.DeleteBackendApiMethod(backendapiID, hitsID, methodID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadBackendApiMethod(t *testing.T) {
	var (
		backendapiID int64 = 98765
		hitsID       int64 = 1
		methodID     int64 = 123325
		endpoint           = fmt.Sprintf(backendMethodResourceEndpoint, backendapiID, hitsID, methodID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		method := &Method{
			Element: MethodItem{
				ID:         methodID,
				Name:       "someName",
				SystemName: "someName2",
			},
		}

		responseBodyBytes, err := json.Marshal(method)
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
	obj, err := c.BackendApiMethod(backendapiID, hitsID, methodID)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("backendapi method returned nil")
	}

	if obj.Element.ID != methodID {
		t.Fatalf("backend_api ID does not match. Expected [%d]; got [%d]", methodID, obj.Element.ID)
	}
}

func TestUpdateBackendApiMethod(t *testing.T) {
	var (
		backendapiID int64 = 98765
		hitsID       int64 = 1
		methodID     int64 = 123325
		endpoint           = fmt.Sprintf(backendMethodResourceEndpoint, backendapiID, hitsID, methodID)
		params             = Params{"description": "newDescr"}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPut {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPut, req.Method)
		}

		method := &Method{
			Element: MethodItem{
				ID:          methodID,
				Name:        "someName",
				SystemName:  "someName2",
				Description: params["description"],
			},
		}

		responseBodyBytes, err := json.Marshal(method)
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
	obj, err := c.UpdateBackendApiMethod(backendapiID, hitsID, methodID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("backendapi method returned nil")
	}

	if obj.Element.ID != methodID {
		t.Fatalf("backend_api method ID does not match. Expected [%d]; got [%d]", methodID, obj.Element.ID)
	}

	if obj.Element.Description != params["description"] {
		t.Fatalf("backend_api method description does not match. Expected [%s]; got [%s]", params["description"], obj.Element.Description)
	}
}
