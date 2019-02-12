package client

// This package provides bare minimum functionality for all the endpoints it exposes,
// which is a subset of the Account Management API.

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	mappingRuleEndpoint             = "/admin/api/services/%s/proxy/mapping_rules.xml"
	createListMetricEndpoint        = "/admin/api/services/%s/metrics.xml"
	updateDeleteMetricEndpoint      = "/admin/api/services/%s/metrics/%s.xml"
	updateDeleteMappingRuleEndpoint = "/admin/api/services/%s/proxy/mapping_rules/%s.xml"
)

var httpReqError = errors.New("error building http request")

// Returns a custom AdminPortal which integrates with the users Account Management API.
// Supported schemes are http and https
func NewAdminPortal(scheme string, host string, port int) (*AdminPortal, error) {
	url2, err := verifyUrl(fmt.Sprintf("%s://%s:%d", scheme, host, port))
	if err != nil {
		return nil, err
	}
	return &AdminPortal{scheme, host, port, url2}, nil
}

// Creates a ThreeScaleClient to communicate with Account Management API.
// If http Client is nil, the default http client will be used
func NewThreeScale(backEnd *AdminPortal, credential string, httpClient *http.Client) *ThreeScaleClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &ThreeScaleClient{backEnd, credential, httpClient}
}

func NewParams() Params {
	params := make(map[string]string)
	return params
}

func (e ApiErr) Error() string {
	return fmt.Sprintf("error calling 3scale system - reason: %s - code: %d", e.err, e.code)
}

func (e ApiErr) Code() int {
	return e.code
}

func (p Params) AddParam(key string, value string) {
	p[key] = value
}

// SetCredentials allow the user to set the client credentials
func (c *ThreeScaleClient) SetCredentials(credential string) {
	c.credential = credential
}

// Request builder for GET request to the provided endpoint
func (c *ThreeScaleClient) buildGetReq(ep string) (*http.Request, error) {
	path := &url.URL{Path: ep}
	req, err := http.NewRequest("GET", c.adminPortal.baseUrl.ResolveReference(path).String(), nil)
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Authorization", "Basic "+basicAuth("", c.credential))
	return req, err
}

// Request builder for POST request to the provided endpoint
func (c *ThreeScaleClient) buildPostReq(ep string, body io.Reader) (*http.Request, error) {
	path := &url.URL{Path: ep}
	req, err := http.NewRequest("POST", c.adminPortal.baseUrl.ResolveReference(path).String(), body)
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+basicAuth("", c.credential))
	return req, err
}

// Request builder for PUT request to the provided endpoint
func (c *ThreeScaleClient) buildUpdateReq(ep string, body io.Reader) (*http.Request, error) {
	path := &url.URL{Path: ep}
	req, err := http.NewRequest("PUT", c.adminPortal.baseUrl.ResolveReference(path).String(), body)
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+basicAuth("", c.credential))
	return req, err
}

// Request builder for DELETE request to the provided endpoint
func (c *ThreeScaleClient) buildDeleteReq(ep string, body io.Reader) (*http.Request, error) {
	path := &url.URL{Path: ep}
	req, err := http.NewRequest("DELETE", c.adminPortal.baseUrl.ResolveReference(path).String(), body)
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+basicAuth("", c.credential))
	return req, err
}

// Request builder for PUT request to the provided endpoint
func (c *ThreeScaleClient) buildPutReq(ep string, body io.Reader) (*http.Request, error) {
	path := &url.URL{Path: ep}
	req, err := http.NewRequest("PUT", c.adminPortal.baseUrl.ResolveReference(path).String(), body)
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+basicAuth("", c.credential))
	return req, err
}

// Verifies a custom admin portal is valid
func verifyUrl(urlToCheck string) (*url.URL, error) {
	url2, err := url.ParseRequestURI(urlToCheck)
	if err == nil {
		if url2.Scheme != "http" && url2.Scheme != "https" {
			err = fmt.Errorf("unsupported schema %s passed to adminPortal", url2.Scheme)
		}

	}
	return url2, err
}

// handleXMLResp takes a http response and validates it against an expected status code
// if response code is unexpected or it fails to decode into the interface provided
// by the caller, an error of type ApiErr is returned
func handleXMLResp(resp *http.Response, expectCode int, decodeInto interface{}) error {
	if resp.StatusCode != expectCode {
		return handleXMLErrResp(resp)
	}

	if decodeInto == nil {
		return nil
	}

	if err := xml.NewDecoder(resp.Body).Decode(decodeInto); err != nil {
		return createApiErr(resp.StatusCode, createDecodingErrorMessage(err))

	}
	return nil
}

// handleJsonResp takes a http response and validates it against an expected status code
// if response code is unexpected or it fails to decode into the interface provided
// by the caller, an error of type ApiErr is returned
func handleJsonResp(resp *http.Response, expectCode int, decodeInto interface{}) error {
	if resp.StatusCode != expectCode {
		return handleJsonErrResp(resp)
	}

	if err := json.NewDecoder(resp.Body).Decode(decodeInto); err != nil {
		return createApiErr(resp.StatusCode, createDecodingErrorMessage(err))
	}

	return nil
}

// handleXMLErrResp decodes an XML response from 3scale system
// into an error of type ApiErr
func handleXMLErrResp(resp *http.Response) error {
	var errResp ErrorResp

	if err := xml.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		return createApiErr(resp.StatusCode, createDecodingErrorMessage(err))
	}

	return ApiErr{resp.StatusCode, errResp.Text}
}

// handleJsonErrResp decodes a JSON response from 3scale system
// into an error of type APiErr
func handleJsonErrResp(resp *http.Response) error {
	var errMap map[string]string

	if err := json.NewDecoder(resp.Body).Decode(&errMap); err != nil {
		return createApiErr(resp.StatusCode, createDecodingErrorMessage(err))
	}

	msg := "error"
	for _, v := range errMap {
		msg = fmt.Sprintf("%s - %s ", msg, v)
	}

	return createApiErr(resp.StatusCode, msg)
}

func createApiErr(statusCode int, message string) ApiErr {
	return ApiErr{
		code: statusCode,
		err:  message,
	}
}

func createDecodingErrorMessage(err error) string {
	return fmt.Sprintf("decoding error - %s", err.Error())
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
