package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	userActivate = "/admin/api/accounts/%s/users/%s/activate.json"
)

// ActivateUser activates user of a given account from pending state to active
func (c *ThreeScaleClient) ActivateUser(accountID, userID string) error {
	endpoint := fmt.Sprintf(userActivate, accountID, userID)

	values := url.Values{}
	body := strings.NewReader(values.Encode())
	req, err := c.buildUpdateReq(endpoint, body)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handleJsonErrResp(resp)
	}

	return nil
}
