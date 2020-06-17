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

func TestListApplicationPlansPricingRules(t *testing.T) {
	var (
		planID   int64 = 97
		endpoint       = fmt.Sprintf(appPlanRuleListResourceEndpoint, planID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		list := ApplicationPlanPricingRuleList{
			Rules: []ApplicationPlanPricingRule{
				{
					Element: ApplicationPlanPricingRuleItem{
						ID:          1,
						Min:         1,
						Max:         10,
						CostPerUnit: "1.92",
						MetricID:    23,
					},
				},
				{
					Element: ApplicationPlanPricingRuleItem{
						ID:          2,
						Min:         11,
						Max:         20,
						CostPerUnit: "1.22",
						MetricID:    23,
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
	list, err := c.ListApplicationPlansPricingRules(planID)
	if err != nil {
		t.Fatal(err)
	}

	if list == nil {
		t.Fatal("list returned nil")
	}

	if len(list.Rules) != 2 {
		t.Fatalf("# list items does not match. Expected [%d]; got [%d]", 2, len(list.Rules))
	}
}

func TestCreateApplicationPlanPricingRule(t *testing.T) {
	var (
		planID   int64 = 97
		metricID int64 = 12
		params         = Params{"min": "1", "max": "2"}
		endpoint       = fmt.Sprintf(appPlanRuleListPerMetricResourceEndpoint, planID, metricID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPost {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPost, req.Method)
		}

		min, _ := strconv.Atoi(params["min"])
		max, _ := strconv.Atoi(params["max"])
		item := ApplicationPlanPricingRule{
			Element: ApplicationPlanPricingRuleItem{
				ID:          1,
				Min:         min,
				Max:         max,
				CostPerUnit: "1.92",
				MetricID:    metricID,
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
	obj, err := c.CreateApplicationPlanPricingRule(planID, metricID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("returned nil")
	}

	if obj.Element.ID != 1 {
		t.Fatalf("ID does not match. Expected [%d]; got [%d]", 1, obj.Element.ID)
	}
}

func TestDeleteApplicationPlanPricingRule(t *testing.T) {
	var (
		planID   int64 = 97
		metricID int64 = 12
		ruleID   int64 = 16
		endpoint       = fmt.Sprintf(appPlanRulePerMetricResourceEndpoint, planID, metricID, ruleID)
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
	err := c.DeleteApplicationPlanPricingRule(planID, metricID, ruleID)
	if err != nil {
		t.Fatal(err)
	}
}
