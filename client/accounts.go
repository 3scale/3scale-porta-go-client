package client

import (
	"net/http"
	"net/url"
)

const (
	accountList = "/admin/api/accounts.xml"
)

func (c *ThreeScaleClient) ListAccounts(accessToken string) (*AccountList, error) {
	req, err := c.buildGetReq(accountList)
	if err != nil {
		return nil, err
	}

	urlValues := url.Values{}
	urlValues.Add("access_token", accessToken)

	req.URL.RawQuery = urlValues.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	accountList := &AccountList{}
	err = handleXMLResp(resp, http.StatusOK, accountList)
	return accountList, err
}
