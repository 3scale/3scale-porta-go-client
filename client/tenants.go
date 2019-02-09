package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	tenantCreate = "/master/api/providers.json"
	tenantRead   = "/master/api/providers/%d.json"
)

// CreateTenant creates new tenant using 3scale API
func (c *ThreeScaleClient) CreateTenant(orgName, username, email, password string) (*Tenant, error) {
	values := url.Values{}
	values.Add("org_name", orgName)
	values.Add("username", username)
	values.Add("email", email)
	values.Add("password", password)

	body := strings.NewReader(values.Encode())
	req, err := c.buildPostReq(tenantCreate, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tenant := &Tenant{}
	err = handleJsonResp(resp, http.StatusCreated, tenant)
	return tenant, err
}

// ShowTenant - Returns tenant info for the specified ID
func (c *ThreeScaleClient) ShowTenant(accessToken string, tenantID int64) (*Tenant, error) {
	endpoint := fmt.Sprintf(tenantRead, tenantID)
	req, err := c.buildGetReq(endpoint)
	if err != nil {
		return nil, httpReqError
	}

	values := url.Values{}
	values.Add("access_token", accessToken)
	req.URL.RawQuery = values.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tenant := &Tenant{}
	err = handleJsonResp(resp, http.StatusOK, tenant)
	return tenant, err
}
