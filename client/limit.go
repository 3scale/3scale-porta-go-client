package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	limitAppPlanCreate           = "/admin/api/application_plans/%s/metrics/%s/limits.xml"
	limitAppPlanList             = "/admin/api/application_plans/%s/limits.xml"
	limitAppPlanUpdateDelete     = "/admin/api/application_plans/%s/metrics/%s/limits/%s.xml "
	limitAppPlanMetricList       = "/admin/api/application_plans/%s/metrics/%s/limits.xml"
	limitEndUserPlanCreateList   = "/admin/api/end_user_plans/%s/metrics/%s/limits.xml"
	limitEndUserPlanUpdateDelete = "/admin/api/end_user_plans/%s/metrics/%s/limits/%s.xml"
)

// CreateLimitAppPlan - Adds a limit to a metric of an application plan.
// All applications with the application plan (application_plan_id) will be constrained by this new limit on the metric (metric_id).
func (c *ThreeScaleClient) CreateLimitAppPlan(credential string, appPlanId string, metricId string, period string, value int) (Limit, error) {
	endpoint := fmt.Sprintf(limitAppPlanCreate, appPlanId, metricId)

	values := url.Values{}
	values.Add("application_plan_id", appPlanId)

	return c.limitCreate(endpoint, credential, metricId, period, value, values)
}

// CreateLimitEndUserPlan - Adds a limit to a metric of an end user plan
// All applications with the application plan (end_user_plan_id) will be constrained by this new limit on the metric (metric_id).
func (c *ThreeScaleClient) CreateLimitEndUserPlan(credential string, endUserPlanId string, metricId string, period string, value int) (Limit, error) {
	endpoint := fmt.Sprintf(limitEndUserPlanCreateList, endUserPlanId, metricId)

	values := url.Values{}
	values.Add("end_user_plan_id", endUserPlanId)

	return c.limitCreate(endpoint, credential, metricId, period, value, values)
}

// UpdateLimitsPerPlan - Updates a limit on a metric of an end user plan
// Valid params keys and their purpose are as follows:
// "period" - Period of the limit
// "value"  - Value of the limit
func (c *ThreeScaleClient) UpdateLimitPerAppPlan(credential string, appPlanId string, metricId string, limitId string, p Params) (Limit, error) {
	endpoint := fmt.Sprintf(limitAppPlanUpdateDelete, appPlanId, metricId, limitId)
	return c.updateLimit(endpoint, credential, p)
}

// UpdateLimitsPerMetric - Updates a limit on a metric of an application plan
// Valid params keys and their purpose are as follows:
// "period" - Period of the limit
// "value"  - Value of the limit
func (c *ThreeScaleClient) UpdateLimitPerEndUserPlan(credential string, userPlanId string, metricId string, limitId string, p Params) (Limit, error) {
	endpoint := fmt.Sprintf(limitEndUserPlanUpdateDelete, userPlanId, metricId, limitId)
	return c.updateLimit(endpoint, credential, p)
}

// DeleteLimitPerAppPlan - Deletes a limit on a metric of an application plan
func (c *ThreeScaleClient) DeleteLimitPerAppPlan(credential string, appPlanId string, metricId string, limitId string) error {
	endpoint := fmt.Sprintf(limitAppPlanUpdateDelete, appPlanId, metricId, limitId)
	return c.deleteLimit(endpoint, credential)
}

// DeleteLimitPerEndUserPlan - Deletes a limit on a metric of an end user plan
func (c *ThreeScaleClient) DeleteLimitPerEndUserPlan(credential string, userPlanId string, metricId string, limitId string) error {
	endpoint := fmt.Sprintf(limitEndUserPlanUpdateDelete, userPlanId, metricId, limitId)
	return c.deleteLimit(endpoint, credential)
}

// ListLimitsPerAppPlan - Returns the list of all limits associated to an application plan.
func (c *ThreeScaleClient) ListLimitsPerAppPlan(credential string, appPlanId string) (LimitList, error) {
	endpoint := fmt.Sprintf(limitAppPlanList, appPlanId)
	return c.listLimits(endpoint, credential)
}

// ListLimitsPerEndUserPlan - Returns the list of all limits associated to an end user plan.
func (c *ThreeScaleClient) ListLimitsPerEndUserPlan(credential string, endUserPlanId string, metricId string) (LimitList, error) {
	endpoint := fmt.Sprintf(limitEndUserPlanCreateList, endUserPlanId, metricId)
	return c.listLimits(endpoint, credential)
}

// ListLimitsPerMetric - Returns the list of all limits associated to a metric of an application plan
func (c *ThreeScaleClient) ListLimitsPerMetric(credential string, appPlanId string, metricId string) (LimitList, error) {
	endpoint := fmt.Sprintf(limitAppPlanMetricList, appPlanId, metricId)
	return c.listLimits(endpoint, credential)
}

func (c *ThreeScaleClient) limitCreate(ep string, credential string, metricId string, period string, value int, values url.Values) (Limit, error) {
	var apiResp Limit

	values.Add("metric_id", metricId)
	values.Add("period", period)
	values.Add("value", strconv.Itoa(value))

	body := strings.NewReader(values.Encode())
	req, err := c.buildPostReq(ep, credential, body)
	if err != nil {
		return apiResp, httpReqError
	}
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return apiResp, err
	}
	defer resp.Body.Close()

	err = handleXMLResp(resp, http.StatusCreated, &apiResp)
	return apiResp, err
}

func (c *ThreeScaleClient) updateLimit(ep string, credential string, p Params) (Limit, error) {
	var l Limit
	values := url.Values{}
	for k, v := range p {
		values.Add(k, v)
	}

	body := strings.NewReader(values.Encode())
	req, err := c.buildUpdateReq(ep, credential, body)
	if err != nil {
		return l, httpReqError
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return l, err
	}
	defer resp.Body.Close()

	err = handleXMLResp(resp, http.StatusOK, &l)
	return l, err
}

func (c *ThreeScaleClient) deleteLimit(ep string, credential string) error {
	values := url.Values{}
	body := strings.NewReader(values.Encode())
	req, err := c.buildDeleteReq(ep, credential, body)
	if err != nil {
		return httpReqError
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return handleXMLResp(resp, http.StatusCreated, nil)
}

// listLimits takes an endpoint and returns a list of limits
func (c *ThreeScaleClient) listLimits(ep string, credential string) (LimitList, error) {
	var ml LimitList

	req, err := c.buildGetReq(ep, credential)
	if err != nil {
		return ml, httpReqError
	}

	values := url.Values{}
	req.URL.RawQuery = values.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ml, err
	}

	defer resp.Body.Close()
	err = handleXMLResp(resp, http.StatusOK, &ml)
	return ml, err
}
