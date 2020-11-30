package client

import (
	"fmt"
	"net/http"
)

const (
	activeDocListEndpoint = "/admin/api/active_docs.json"
	activeDocEndpoint     = "/admin/api/active_docs/%d.json"
)

// ListActiveDocs List existing activedocs for the client provider account
func (c *ThreeScaleClient) ListActiveDocs() (*ActiveDocList, error) {
	req, err := c.buildGetReq(activeDocListEndpoint)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	activeDocList := &ActiveDocList{}
	err = handleJsonResp(resp, http.StatusOK, activeDocList)
	return activeDocList, err
}

// ActiveDoc Reads 3scale Activedoc
func (c *ThreeScaleClient) ActiveDoc(id int64) (*ActiveDoc, error) {
	endpoint := fmt.Sprintf(activeDocEndpoint, id)

	req, err := c.buildGetJSONReq(endpoint)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	activeDoc := &ActiveDoc{}
	err = handleJsonResp(resp, http.StatusOK, activeDoc)
	return activeDoc, err
}
