package client

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	tenantCreate = "/master/api/providers.xml"
)

// CreateTenant creates new tenant using 3scale API
func (c *ThreeScaleClient) CreateTenant(credential, orgName, username, email, password string) (*Signup, error) {
	values := url.Values{}
	values.Add("org_name", orgName)
	values.Add("username", username)
	values.Add("email", email)
	values.Add("password", password)

	body := strings.NewReader(values.Encode())
	req, err := c.buildPostReq(tenantCreate, credential, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	signup := &Signup{}
	err = handleXMLResp(resp, http.StatusCreated, signup)
	return signup, err
}
