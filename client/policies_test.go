package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestPolicies(t *testing.T) {
	var (
		productID int64 = 97
		endpoint        = fmt.Sprintf(policiesResourceEndpoint, productID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		policies := &PoliciesConfigList{
			Policies: []PolicyConfig{
				{
					Name:          "apicast",
					Version:       "builtin",
					Enabled:       true,
					Configuration: map[string]interface{}{},
				},
			},
		}

		responseBodyBytes, err := json.Marshal(policies)
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
	policies, err := c.Policies(productID)
	if err != nil {
		t.Fatal(err)
	}

	if policies == nil {
		t.Fatal("policies returned nil")
	}

	if len(policies.Policies) != 1 {
		t.Fatalf("Then number of list items does not match. Expected [%d]; got [%d]", 1, len(policies.Policies))
	}
}

func TestUpdatePolicies(t *testing.T) {
	var (
		productID int64 = 98765
		endpoint        = fmt.Sprintf(policiesResourceEndpoint, productID)
		policies        = &PoliciesConfigList{
			Policies: []PolicyConfig{
				{
					Name:          "apicast",
					Version:       "builtin",
					Enabled:       true,
					Configuration: map[string]interface{}{},
				},
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

		respPolicies := &PoliciesConfigList{
			Policies: []PolicyConfig{
				{
					Name:          "apicast",
					Version:       "builtin",
					Enabled:       true,
					Configuration: map[string]interface{}{},
				},
			},
		}
		responseBodyBytes, err := json.Marshal(respPolicies)
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
	obj, err := c.UpdatePolicies(productID, policies)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("resp returned nil")
	}

	if !reflect.DeepEqual(obj, policies) {
		sObj, _ := json.Marshal(obj)
		sPol, _ := json.Marshal(policies)
		t.Logf("Expected %s; got %s", string(sPol), string(sObj))
		t.Fatalf("Expected %v; got %v", policies, obj)
	}
}

func TestPoliciesType(t *testing.T) {
	policiesA := []byte(`{"policies_config": [
						{
						"name": "apicast",
						"version": "builtin",
						"configuration": { "prop1": "value1", "prop2": "value2"},
						"enabled": true
						}
	]}`)

	// Same as policiesA but with slightly different configuration serialization
	policiesB := []byte(`{"policies_config": [
						{
						"name": "apicast",
						"version": "builtin",
						"configuration": {
						"prop1":                         "value1",
						"prop2": "value2"
						},
						"enabled": true
						}
	]}`)

	policiesAObj := &PoliciesConfigList{}
	err := json.Unmarshal(policiesA, policiesAObj)
	if err != nil {
		t.Fatal(err)
	}

	policiesBObj := &PoliciesConfigList{}
	err = json.Unmarshal(policiesB, policiesBObj)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(policiesAObj, policiesBObj) {
		t.Logf("%v | %v", policiesAObj, policiesBObj)
		pa, _ := json.Marshal(policiesAObj)
		pb, _ := json.Marshal(policiesBObj)
		t.Fatalf("objects not equal: %s | %s", string(pa), string(pb))
	}
}
