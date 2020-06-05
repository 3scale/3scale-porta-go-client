package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	productListResourceEndpoint       = "/admin/api/services.json"
	productResourceEndpoint           = "/admin/api/services/%d.json"
	productMethodListResourceEndpoint = "/admin/api/services/%d/metrics/%d/methods.json"
	productMethodResourceEndpoint     = "/admin/api/services/%d/metrics/%d/methods/%d.json"
)

// CreateProduct Create 3scale Product
func (c *ThreeScaleClient) CreateProduct(name string, params Params) (*Product, error) {
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	values.Add("name", name)

	body := strings.NewReader(values.Encode())
	req, err := c.buildPostReq(productListResourceEndpoint, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	product := &Product{}
	err = handleJsonResp(resp, http.StatusCreated, product)
	return product, err
}

// UpdateProduct Update existing product
func (c *ThreeScaleClient) UpdateProduct(id int64, params Params) (*Product, error) {
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}

	putProductEndpoint := fmt.Sprintf(productResourceEndpoint, id)

	body := strings.NewReader(values.Encode())
	req, err := c.buildUpdateReq(putProductEndpoint, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	product := &Product{}
	err = handleJsonResp(resp, http.StatusOK, product)
	return product, err
}

// DeleteProduct Delete existing product
func (c *ThreeScaleClient) DeleteProduct(id int64) error {
	productEndpoint := fmt.Sprintf(productResourceEndpoint, id)

	req, err := c.buildDeleteReq(productEndpoint, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var empty struct{}
	return handleJsonResp(resp, http.StatusOK, &empty)
}

// ListProducts List existing products
func (c *ThreeScaleClient) ListProducts() (*ProductList, error) {
	req, err := c.buildGetReq(productListResourceEndpoint)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	productList := &ProductList{}
	err = handleJsonResp(resp, http.StatusOK, productList)
	return productList, err
}

// ListProductMethods List existing product methods
func (c *ThreeScaleClient) ListProductMethods(productID, hitsID int64) (*MethodList, error) {
	endpoint := fmt.Sprintf(productMethodListResourceEndpoint, productID, hitsID)
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

// CreateProductMethod Create 3scale product method
func (c *ThreeScaleClient) CreateProductMethod(productID, hitsID int64, params Params) (*Method, error) {
	endpoint := fmt.Sprintf(productMethodListResourceEndpoint, productID, hitsID)

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

// DeleteProductMethod Delete 3scale product method
func (c *ThreeScaleClient) DeleteProductMethod(productID, hitsID, methodID int64) error {
	endpoint := fmt.Sprintf(productMethodResourceEndpoint, productID, hitsID, methodID)

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

// ProductMethod Read 3scale product method
func (c *ThreeScaleClient) ProductMethod(productID, hitsID, methodID int64) (*Method, error) {
	endpoint := fmt.Sprintf(productMethodResourceEndpoint, productID, hitsID, methodID)

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

// UpdateProductMethod Update 3scale product method
func (c *ThreeScaleClient) UpdateProductMethod(productID, hitsID, methodID int64, params Params) (*Method, error) {
	endpoint := fmt.Sprintf(productMethodResourceEndpoint, productID, hitsID, methodID)

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
