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

func TestListApplicationPlansByProduct(t *testing.T) {
	var (
		productID int64 = 97
		endpoint        = fmt.Sprintf(appPlanListResourceEndpoint, productID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		list := ApplicationPlanJSONList{
			Plans: []ApplicationPlan{
				{
					Element: ApplicationPlanItem{
						ID:         1,
						Name:       "plan01",
						SystemName: "plan01",
					},
				},
				{
					Element: ApplicationPlanItem{
						ID:         2,
						Name:       "plan02",
						SystemName: "plan02",
					},
				},
			},
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
	list, err := c.ListApplicationPlansByProduct(productID)
	if err != nil {
		t.Fatal(err)
	}

	if list == nil {
		t.Fatal("list returned nil")
	}

	if len(list.Plans) != 2 {
		t.Fatalf("# list items does not match. Expected [%d]; got [%d]", 2, len(list.Plans))
	}
}

func TestCreateApplicationPlan(t *testing.T) {
	var (
		productID int64 = 97
		params          = Params{"name": "plan01"}
		endpoint        = fmt.Sprintf(appPlanListResourceEndpoint, productID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPost {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPost, req.Method)
		}

		item := ApplicationPlan{
			Element: ApplicationPlanItem{
				ID:   1,
				Name: params["name"],
			},
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
	obj, err := c.CreateApplicationPlan(productID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("returned nil")
	}

	if obj.Element.Name != params["name"] {
		t.Fatalf("Name does not match. Expected [%s]; got [%s]", params["name"], obj.Element.Name)
	}
}

func TestDeleteApplicationPlan(t *testing.T) {
	var (
		productID int64 = 97
		id        int64 = 3
		endpoint        = fmt.Sprintf(appPlanResourceEndpoint, productID, id)
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
	err := c.DeleteApplicationPlan(productID, id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestApplicationPlan(t *testing.T) {
	var (
		productID int64 = 97
		id        int64 = 3
		endpoint        = fmt.Sprintf(appPlanResourceEndpoint, productID, id)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		item := ApplicationPlan{
			Element: ApplicationPlanItem{
				ID:   id,
				Name: "someName",
			},
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
	obj, err := c.ApplicationPlan(productID, id)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("returned nil")
	}

	if obj.Element.ID != id {
		t.Fatalf("ID does not match. Expected [%d]; got [%d]", id, obj.Element.ID)
	}
}

func TestUpdateApplicationPlan(t *testing.T) {
	var (
		productID int64 = 97
		id        int64 = 3
		params          = Params{"name": "newName"}
		endpoint        = fmt.Sprintf(appPlanResourceEndpoint, productID, id)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPut {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPut, req.Method)
		}

		item := ApplicationPlan{
			Element: ApplicationPlanItem{
				ID:   id,
				Name: params["name"],
			},
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
	obj, err := c.UpdateApplicationPlan(productID, id, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("returned nil")
	}

	if obj.Element.Name != params["name"] {
		t.Fatalf("Name does not match. Expected [%s]; got [%s]", params["name"], obj.Element.Name)
	}
}
