package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestListApplicationPlansLimits(t *testing.T) {
	var (
		planID   int64 = 97
		endpoint       = fmt.Sprintf(appPlanLimitListResourceEndpoint, planID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		list := ApplicationPlanLimitList{
			Limits: []ApplicationPlanLimit{
				{
					Element: ApplicationPlanLimitItem{
						ID:       1,
						Value:    100,
						Period:   "month",
						PlanID:   34,
						MetricID: 23,
					},
				},
				{
					Element: ApplicationPlanLimitItem{
						ID:       2,
						Value:    101,
						Period:   "year",
						PlanID:   34,
						MetricID: 23,
					},
				},
			},
		}

		responseBodyBytes, err := json.Marshal(list)
		if err != nil {
			t.Fatal(err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer(responseBodyBytes)),
			Header:     make(http.Header),
		}
	})

	credential := "someAccessToken"
	c := NewThreeScale(NewTestAdminPortal(t), credential, httpClient)
	list, err := c.ListApplicationPlansLimits(planID)
	if err != nil {
		t.Fatal(err)
	}

	if list == nil {
		t.Fatal("list returned nil")
	}

	if len(list.Limits) != 2 {
		t.Fatalf("# list items does not match. Expected [%d]; got [%d]", 2, len(list.Limits))
	}
}

func TestCreateApplicationPlanLimit(t *testing.T) {
	var (
		planID   int64 = 97
		metricID int64 = 12
		params         = Params{"value": "123", "period": "month"}
		endpoint       = fmt.Sprintf(appPlanLimitListPerMetricResourceEndpoint, planID, metricID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPost {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPost, req.Method)
		}

		value, _ := strconv.Atoi(params["value"])
		item := ApplicationPlanLimit{
			Element: ApplicationPlanLimitItem{
				ID:       1,
				Value:    value,
				Period:   params["period"],
				PlanID:   planID,
				MetricID: metricID,
			},
		}

		responseBodyBytes, err := json.Marshal(item)
		if err != nil {
			t.Fatal(err)
		}

		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewBuffer(responseBodyBytes)),
			Header:     make(http.Header),
		}
	})

	credential := "someAccessToken"
	c := NewThreeScale(NewTestAdminPortal(t), credential, httpClient)
	obj, err := c.CreateApplicationPlanLimit(planID, metricID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("returned nil")
	}

	if obj.Element.ID != 1 {
		t.Fatalf("method ID does not match. Expected [%d]; got [%d]", 1, obj.Element.ID)
	}

	if obj.Element.Period != params["period"] {
		t.Fatalf("limit period does not match. Expected [%s]; got [%s]", params["period"], obj.Element.Period)
	}
}

func TestDeleteApplicationPlanLimit(t *testing.T) {
	var (
		planID   int64 = 97
		metricID int64 = 12
		limitID  int64 = 16
		endpoint       = fmt.Sprintf(appPlanLimitPerMetricResourceEndpoint, planID, metricID, limitID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodDelete {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodDelete, req.Method)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader("")),
			Header:     make(http.Header),
		}
	})

	credential := "someAccessToken"
	c := NewThreeScale(NewTestAdminPortal(t), credential, httpClient)
	err := c.DeleteApplicationPlanLimit(planID, metricID, limitID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestApplicationPlanLimit(t *testing.T) {
	var (
		planID   int64 = 97
		metricID int64 = 12
		limitID  int64 = 16
		endpoint       = fmt.Sprintf(appPlanLimitPerMetricResourceEndpoint, planID, metricID, limitID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		item := ApplicationPlanLimit{
			Element: ApplicationPlanLimitItem{
				ID:       1,
				Value:    1000,
				Period:   "month",
				PlanID:   planID,
				MetricID: metricID,
			},
		}

		responseBodyBytes, err := json.Marshal(item)
		if err != nil {
			t.Fatal(err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer(responseBodyBytes)),
			Header:     make(http.Header),
		}
	})

	credential := "someAccessToken"
	c := NewThreeScale(NewTestAdminPortal(t), credential, httpClient)
	obj, err := c.ApplicationPlanLimit(planID, metricID, limitID)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("obj returned nil")
	}

	if obj.Element.ID != 1 {
		t.Fatalf("ID does not match. Expected [%d]; got [%d]", 1, obj.Element.ID)
	}
}

func TestUpdateApplicationPlanLimit(t *testing.T) {
	var (
		planID   int64 = 97
		metricID int64 = 12
		limitID  int64 = 16
		params         = Params{"value": "123", "period": "month"}
		endpoint       = fmt.Sprintf(appPlanLimitPerMetricResourceEndpoint, planID, metricID, limitID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPut {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPut, req.Method)
		}

		value, _ := strconv.Atoi(params["value"])
		item := ApplicationPlanLimit{
			Element: ApplicationPlanLimitItem{
				ID:       1,
				Value:    value,
				Period:   params["period"],
				PlanID:   planID,
				MetricID: metricID,
			},
		}

		responseBodyBytes, err := json.Marshal(item)
		if err != nil {
			t.Fatal(err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer(responseBodyBytes)),
			Header:     make(http.Header),
		}
	})

	credential := "someAccessToken"
	c := NewThreeScale(NewTestAdminPortal(t), credential, httpClient)
	obj, err := c.UpdateApplicationPlanLimit(planID, metricID, limitID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("obj returned nil")
	}

	if obj.Element.ID != 1 {
		t.Fatalf("ID does not match. Expected [%d]; got [%d]", 1, obj.Element.ID)
	}
}
