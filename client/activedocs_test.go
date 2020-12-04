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

func TestListActiveDocs(t *testing.T) {
	var (
		adID1    int64 = 1
		adID2    int64 = 2
		name1          = "ActiveDoc1"
		name2          = "ActiveDoc2"
		body           = "{}"
		endpoint       = activeDocListEndpoint
		list           = ActiveDocList{
			ActiveDocs: []ActiveDoc{
				{
					Element: ActiveDocItem{
						ID:         &adID1,
						Name:       &name1,
						Body:       &body,
						SystemName: &name1,
					},
				},
				{
					Element: ActiveDocItem{
						ID:         &adID2,
						Name:       &name2,
						Body:       &body,
						SystemName: &name2,
					},
				},
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
	resp, err := c.ListActiveDocs()
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("response was nil")
	}

	if !reflect.DeepEqual(*resp, list) {
		got, _ := json.Marshal(*resp)
		expected, _ := json.Marshal(list)
		t.Logf("Expected %s; got %s", string(expected), string(got))
		t.Fatalf("Expected %v; got %v", expected, got)
	}
}

func TestReadActiveDocs(t *testing.T) {
	var (
		adID1     int64 = 1
		name1           = "ActiveDoc1"
		body            = "{}"
		endpoint        = fmt.Sprintf(activeDocEndpoint, adID1)
		activeDoc       = ActiveDoc{
			Element: ActiveDocItem{
				ID:         &adID1,
				Name:       &name1,
				Body:       &body,
				SystemName: &name1,
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

		responseBodyBytes, err := json.Marshal(activeDoc)
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
	resp, err := c.ActiveDoc(adID1)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("response was nil")
	}

	if !reflect.DeepEqual(*resp, activeDoc) {
		got, _ := json.Marshal(*resp)
		expected, _ := json.Marshal(activeDoc)
		t.Logf("Expected %s; got %s", string(expected), string(got))
		t.Fatalf("Expected %v; got %v", expected, got)
	}
}

func TestCreateActiveDocs(t *testing.T) {
	var (
		adID1     int64 = 1
		name1           = "ActiveDoc1"
		body            = "{}"
		endpoint        = activeDocListEndpoint
		activeDoc       = ActiveDoc{
			Element: ActiveDocItem{
				ID:         &adID1,
				Name:       &name1,
				Body:       &body,
				SystemName: &name1,
			},
		}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPost {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPost, req.Method)
		}

		responseBodyBytes, err := json.Marshal(activeDoc)
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
	resp, err := c.CreateActiveDoc(&activeDoc)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("response was nil")
	}

	if !reflect.DeepEqual(*resp, activeDoc) {
		got, _ := json.Marshal(*resp)
		expected, _ := json.Marshal(activeDoc)
		t.Logf("Expected %s; got %s", string(expected), string(got))
		t.Fatalf("Expected %v; got %v", expected, got)
	}
}

func TestUpdateActiveDocs(t *testing.T) {
	var (
		adID1     int64 = 1
		name1           = "ActiveDoc1"
		body            = "{}"
		endpoint        = fmt.Sprintf(activeDocEndpoint, adID1)
		activeDoc       = ActiveDoc{
			Element: ActiveDocItem{
				ID:         &adID1,
				Name:       &name1,
				Body:       &body,
				SystemName: &name1,
			},
		}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPut {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPut, req.Method)
		}

		responseBodyBytes, err := json.Marshal(activeDoc)
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
	resp, err := c.UpdateActiveDoc(&activeDoc)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("response was nil")
	}

	if !reflect.DeepEqual(*resp, activeDoc) {
		got, _ := json.Marshal(*resp)
		expected, _ := json.Marshal(activeDoc)
		t.Logf("Expected %s; got %s", string(expected), string(got))
		t.Fatalf("Expected %v; got %v", expected, got)
	}
}

func TestDeleteActiveDocs(t *testing.T) {
	var (
		adID1    int64 = 1
		endpoint       = fmt.Sprintf(activeDocEndpoint, adID1)
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
	err := c.DeleteActiveDoc(adID1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnbindActiveDocFromProduct(t *testing.T) {
	var (
		adID1     int64 = 1
		name1           = "ActiveDoc1"
		body            = "{}"
		endpoint        = fmt.Sprintf(activeDocEndpoint, adID1)
		activeDoc       = ActiveDoc{
			Element: ActiveDocItem{
				ID:         &adID1,
				Name:       &name1,
				Body:       &body,
				SystemName: &name1,
			},
		}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPut {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPut, req.Method)
		}

		responseBodyBytes, err := json.Marshal(activeDoc)
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
	resp, err := c.UnbindActiveDocFromProduct(adID1)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatal("response was nil")
	}

	if !reflect.DeepEqual(*resp, activeDoc) {
		got, _ := json.Marshal(*resp)
		expected, _ := json.Marshal(activeDoc)
		t.Logf("Expected %s; got %s", string(expected), string(got))
		t.Fatalf("Expected %v; got %v", expected, got)
	}
}
