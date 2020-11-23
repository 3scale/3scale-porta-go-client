package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestOIDCConfiguration(t *testing.T) {
	var (
		productID int64 = 97
		endpoint        = fmt.Sprintf(oidcResourceEndpoint, productID)
		oidcConf        = &OIDCConfiguration{
			Element: OIDCConfigurationItem{
				StandardFlowEnabled:       false,
				ImplicitFlowEnabled:       true,
				ServiceAccountsEnabled:    false,
				DirectAccessGrantsEnabled: true,
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

		responseBodyBytes, err := json.Marshal(oidcConf)
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
	obj, err := c.OIDCConfiguration(productID)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("response was nil")
	}

	if *obj != *oidcConf {
		t.Fatalf("Expected %v; got %v", *oidcConf, *obj)
	}
}

func TestUpdateOIDCConfiguration(t *testing.T) {
	var (
		productID int64 = 98765
		endpoint        = fmt.Sprintf(oidcResourceEndpoint, productID)
		oidcConf        = &OIDCConfiguration{
			Element: OIDCConfigurationItem{
				StandardFlowEnabled:       false,
				ImplicitFlowEnabled:       true,
				ServiceAccountsEnabled:    false,
				DirectAccessGrantsEnabled: true,
			},
		}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPatch {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPatch, req.Method)
		}

		responseBodyBytes, err := json.Marshal(*oidcConf)
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
	obj, err := c.UpdateOIDCConfiguration(productID, oidcConf)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("resp returned nil")
	}

	if *obj != *oidcConf {
		t.Fatalf("Expected %v; got %v", *oidcConf, *obj)
	}
}
