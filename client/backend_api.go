package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	backendListResourceEndpoint       = "/admin/api/backend_apis.json"
	backendResourceEndpoint           = "/admin/api/backend_apis/%d.json"
	backendMethodListResourceEndpoint = "/admin/api/backend_apis/%d/metrics/%d/methods.json"
	backendMethodResourceEndpoint     = "/admin/api/backend_apis/%d/metrics/%d/methods/%d.json"
	backendMetricListResourceEndpoint = "/admin/api/backend_apis/%d/metrics.json"
	backendMetricResourceEndpoint     = "/admin/api/backend_apis/%d/metrics/%d.json"
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

// ListBackendapiMethods List existing backend methods
func (c *ThreeScaleClient) ListBackendapiMethods(backendapiID, hitsID int64) (*MethodList, error) {
	endpoint := fmt.Sprintf(backendMethodListResourceEndpoint, backendapiID, hitsID)
	req, err := c.buildGetReq(endpoint)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	list := &MethodList{}
	err = handleJsonResp(resp, http.StatusOK, list)
	return list, err
}

// CreateBackendApiMethod Create 3scale Backend method
func (c *ThreeScaleClient) CreateBackendApiMethod(backendapiID, hitsID int64, params Params) (*Method, error) {
	endpoint := fmt.Sprintf(backendMethodListResourceEndpoint, backendapiID, hitsID)

	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}

	body := strings.NewReader(values.Encode())
	req, err := c.buildPostReq(endpoint, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	item := &Method{}
	err = handleJsonResp(resp, http.StatusCreated, item)
	return item, err
}

// DeleteBackendApiMethod Delete 3scale Backend method
func (c *ThreeScaleClient) DeleteBackendApiMethod(backendapiID, hitsID, methodID int64) error {
	endpoint := fmt.Sprintf(backendMethodResourceEndpoint, backendapiID, hitsID, methodID)

	req, err := c.buildDeleteReq(endpoint, nil)
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

// BackendApiMethod Read 3scale Backend method
func (c *ThreeScaleClient) BackendApiMethod(backendapiID, hitsID, methodID int64) (*Method, error) {
	endpoint := fmt.Sprintf(backendMethodResourceEndpoint, backendapiID, hitsID, methodID)

	req, err := c.buildGetReq(endpoint)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	item := &Method{}
	err = handleJsonResp(resp, http.StatusOK, item)
	return item, err
}

// UpdateBackendApiMethod Update 3scale Backend method
func (c *ThreeScaleClient) UpdateBackendApiMethod(backendapiID, hitsID, methodID int64, params Params) (*Method, error) {
	endpoint := fmt.Sprintf(backendMethodResourceEndpoint, backendapiID, hitsID, methodID)

	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	body := strings.NewReader(values.Encode())
	req, err := c.buildUpdateReq(endpoint, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	item := &Method{}
	err = handleJsonResp(resp, http.StatusOK, item)
	return item, err
}

// ListBackendapiMetrics List existing backend metric
func (c *ThreeScaleClient) ListBackendapiMetrics(backendapiID int64) (*MetricJSONList, error) {
	endpoint := fmt.Sprintf(backendMetricListResourceEndpoint, backendapiID)
	req, err := c.buildGetReq(endpoint)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	list := &MetricJSONList{}
	err = handleJsonResp(resp, http.StatusOK, list)
	return list, err
}

// CreateBackendApiMetric Create 3scale Backend metric
func (c *ThreeScaleClient) CreateBackendApiMetric(backendapiID int64, params Params) (*MetricJSON, error) {
	endpoint := fmt.Sprintf(backendMetricListResourceEndpoint, backendapiID)

	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}

	body := strings.NewReader(values.Encode())
	req, err := c.buildPostReq(endpoint, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	item := &MetricJSON{}
	err = handleJsonResp(resp, http.StatusCreated, item)
	return item, err
}

// DeleteBackendApiMetric Delete 3scale Backend metric
func (c *ThreeScaleClient) DeleteBackendApiMetric(backendapiID, metricID int64) error {
	endpoint := fmt.Sprintf(backendMetricResourceEndpoint, backendapiID, metricID)

	req, err := c.buildDeleteReq(endpoint, nil)
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

// BackendApiMetric Read 3scale Backend metric
func (c *ThreeScaleClient) BackendApiMetric(backendapiID, metricID int64) (*MetricJSON, error) {
	endpoint := fmt.Sprintf(backendMetricResourceEndpoint, backendapiID, metricID)

	req, err := c.buildGetReq(endpoint)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	item := &MetricJSON{}
	err = handleJsonResp(resp, http.StatusOK, item)
	return item, err
}

// UpdateBackendApiMetric Update 3scale Backend metric
func (c *ThreeScaleClient) UpdateBackendApiMetric(backendapiID, metricID int64, params Params) (*MetricJSON, error) {
	endpoint := fmt.Sprintf(backendMetricResourceEndpoint, backendapiID, metricID)

	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	body := strings.NewReader(values.Encode())
	req, err := c.buildUpdateReq(endpoint, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	item := &MetricJSON{}
	err = handleJsonResp(resp, http.StatusOK, item)
	return item, err
}
