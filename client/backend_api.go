package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	backendListResourceEndpoint = "/admin/api/backend_apis.json"
	backendResourceEndpoint     = "/admin/api/backend_apis/%d.json"
)

// ListBackends List existing backends
func (c *ThreeScaleClient) ListBackendApis() (*BackendApiList, error) {
	req, err := c.buildGetReq(backendListResourceEndpoint)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	backendList := &BackendApiList{}
	err = handleJsonResp(resp, http.StatusOK, backendList)
	return backendList, err
}

// CreateBackendApi Create 3scale Backend
func (c *ThreeScaleClient) CreateBackendApi(params Params) (*BackendApi, error) {
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}

	body := strings.NewReader(values.Encode())
	req, err := c.buildPostReq(backendListResourceEndpoint, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	backendApi := &BackendApi{}
	err = handleJsonResp(resp, http.StatusCreated, backendApi)
	return backendApi, err
}

// DeleteBackendApi Delete existing backend
func (c *ThreeScaleClient) DeleteBackendApi(id int64) error {
	backendEndpoint := fmt.Sprintf(backendResourceEndpoint, id)

	req, err := c.buildDeleteReq(backendEndpoint, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return handleJsonResp(resp, http.StatusOK, nil)
}

// BackendApi Read 3scale Backend
func (c *ThreeScaleClient) BackendApi(id int64) (*BackendApi, error) {
	backendEndpoint := fmt.Sprintf(backendResourceEndpoint, id)

	req, err := c.buildGetReq(backendEndpoint)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	backendAPI := &BackendApi{}
	err = handleJsonResp(resp, http.StatusOK, backendAPI)
	return backendAPI, err
}

// UpdateBackendApi Update 3scale Backend
func (c *ThreeScaleClient) UpdateBackendApi(id int64, params Params) (*BackendApi, error) {
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}

	backendEndpoint := fmt.Sprintf(backendResourceEndpoint, id)

	body := strings.NewReader(values.Encode())
	req, err := c.buildUpdateReq(backendEndpoint, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	backendAPI := &BackendApi{}
	err = handleJsonResp(resp, http.StatusOK, backendAPI)
	return backendAPI, err
}
