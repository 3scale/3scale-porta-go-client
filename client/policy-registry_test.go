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

func myCustomApicastPolicy1() APIcastPolicy {
	name := "myCustomPolicy01"
	version := "0.0.1"
	var policyID int64 = 1

	configObj := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"someAttr": map[string]interface{}{
				"description": "Some attribute",
				"type":        "integer",
			},
		},
	}

	configBytes, _ := json.Marshal(configObj)
	configuration := json.RawMessage(configBytes)
	schema := "http://json-schema.org/draft-07/schema#"
	summary := "Does cool thing"
	return APIcastPolicy{
		Element: APIcastPolicyItem{
			ID:      &policyID,
			Name:    &name,
			Version: &version,
			Schema: &APIcastPolicySchema{
				Name:          &name,
				Summary:       &summary,
				Schema:        &schema,
				Version:       &version,
				Configuration: &configuration,
			},
		},
	}
}

func myCustomApicastPolicy2() APIcastPolicy {
	name := "myCustomPolicy02"
	version := "0.0.1"
	var policyID int64 = 2
	configObj := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"someAttr": map[string]interface{}{
				"description": "Some attribute",
				"type":        "integer",
			},
		},
	}
	configBytes, _ := json.Marshal(configObj)
	configuration := json.RawMessage(configBytes)

	schema := "http://json-schema.org/draft-07/schema#"
	summary := "Does cool thing"
	return APIcastPolicy{
		Element: APIcastPolicyItem{
			ID:      &policyID,
			Name:    &name,
			Version: &version,
			Schema: &APIcastPolicySchema{
				Name:          &name,
				Summary:       &summary,
				Schema:        &schema,
				Version:       &version,
				Configuration: &configuration,
			},
		},
	}

}

func TestListAPIcastPolicies(t *testing.T) {
	var (
		endpoint = apicastPolicyRegistryEndpoint

		list = APIcastPolicyRegistry{
			Items: []APIcastPolicy{
				myCustomApicastPolicy1(),
				myCustomApicastPolicy2(),
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
	resp, err := c.ListAPIcastPolicies()
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

func TestReadAPIcastPolicy(t *testing.T) {
	var (
		policyID int64 = 1
		endpoint       = fmt.Sprintf(apicastPolicyEndpoint, policyID)
		item           = myCustomApicastPolicy1()
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
	resp, err := c.ReadAPIcastPolicy(policyID)
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

func TestCreateAPIcastPolicy(t *testing.T) {
	var (
		endpoint = apicastPolicyRegistryEndpoint
		item     = myCustomApicastPolicy1()
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
	resp, err := c.CreateAPIcastPolicy(&item)
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

func TestUpdateAPIcastPolicy(t *testing.T) {
	var (
		policyID int64 = 1
		endpoint       = fmt.Sprintf(apicastPolicyEndpoint, policyID)
		item           = myCustomApicastPolicy1()
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
	resp, err := c.UpdateAPIcastPolicy(&item)
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

func TestDeleteAPIcastPolicy(t *testing.T) {
	var (
		policyID int64 = 1
		endpoint       = fmt.Sprintf(apicastPolicyEndpoint, policyID)
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
	err := c.DeleteAPIcastPolicy(policyID)
	if err != nil {
		t.Fatal(err)
	}
}
