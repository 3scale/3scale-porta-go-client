package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestListProductMethods(t *testing.T) {
	var (
		productID int64 = 97
		hitsID    int64 = 98
		endpoint        = fmt.Sprintf(productMethodListResourceEndpoint, productID, hitsID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		list := MethodList{
			Methods: []Method{
				{
					Element: MethodItem{
						ID:         1,
						Name:       "method01",
						SystemName: "desc01",
					},
				},
				{
					Element: MethodItem{
						ID:         2,
						Name:       "method02",
						SystemName: "desc02",
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
	list, err := c.ListProductMethods(productID, hitsID)
	if err != nil {
		t.Fatal(err)
	}

	if list == nil {
		t.Fatal("list returned nil")
	}

	if len(list.Methods) != 2 {
		t.Fatalf("Then number of list items does not match. Expected [%d]; got [%d]", 2, len(list.Methods))
	}
}

func TestCreateProductMethod(t *testing.T) {
	var (
		productID int64 = 97
		hitsID    int64 = 98
		params          = Params{"friendly_name": "method5"}
		endpoint        = fmt.Sprintf(productMethodListResourceEndpoint, productID, hitsID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPost {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPost, req.Method)
		}

		item := &Method{
			Element: MethodItem{
				ID:         10,
				Name:       params["friendly_name"],
				SystemName: params["friendly_name"],
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
	obj, err := c.CreateProductMethod(productID, hitsID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("returned nil")
	}

	if obj.Element.ID != 10 {
		t.Fatalf("method ID does not match. Expected [%d]; got [%d]", 10, obj.Element.ID)
	}

	if obj.Element.Name != params["friendly_name"] {
		t.Fatalf("method name does not match. Expected [%s]; got [%s]", params["friendly_name"], obj.Element.Name)
	}
}

func TestDeleteProductMethod(t *testing.T) {
	var (
		productID int64 = 97
		hitsID    int64 = 98
		methodID  int64 = 123325
		endpoint        = fmt.Sprintf(productMethodResourceEndpoint, productID, hitsID, methodID)
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
	err := c.DeleteProductMethod(productID, hitsID, methodID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadProductMethod(t *testing.T) {
	var (
		productID int64 = 97
		hitsID    int64 = 98
		methodID  int64 = 123325
		endpoint        = fmt.Sprintf(productMethodResourceEndpoint, productID, hitsID, methodID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		item := &Method{
			Element: MethodItem{
				ID:         methodID,
				Name:       "someName",
				SystemName: "someName2",
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
	obj, err := c.ProductMethod(productID, hitsID, methodID)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("method returned nil")
	}

	if obj.Element.ID != methodID {
		t.Fatalf("ID does not match. Expected [%d]; got [%d]", methodID, obj.Element.ID)
	}
}

func TestUpdateProductMethod(t *testing.T) {
	var (
		productID int64 = 98765
		hitsID    int64 = 1
		methodID  int64 = 123325
		endpoint        = fmt.Sprintf(productMethodResourceEndpoint, productID, hitsID, methodID)
		params          = Params{"description": "newDescr"}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPut {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPut, req.Method)
		}

		method := &Method{
			Element: MethodItem{
				ID:          methodID,
				Name:        "someName",
				SystemName:  "someName2",
				Description: params["description"],
			},
		}

		responseBodyBytes, err := json.Marshal(method)
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
	obj, err := c.UpdateProductMethod(productID, hitsID, methodID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("method returned nil")
	}

	if obj.Element.ID != methodID {
		t.Fatalf("method ID does not match. Expected [%d]; got [%d]", methodID, obj.Element.ID)
	}

	if obj.Element.Description != params["description"] {
		t.Fatalf("method description does not match. Expected [%s]; got [%s]", params["description"], obj.Element.Description)
	}
}

func TestListProductMetrics(t *testing.T) {
	var (
		productID int64 = 97
		endpoint        = fmt.Sprintf(productMetricListResourceEndpoint, productID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		list := MetricJSONList{
			Metrics: []MetricJSON{
				{
					Element: MetricItem{
						ID:         1,
						Name:       "hits",
						SystemName: "hits",
						Unit:       "hit",
					},
				},
				{
					Element: MetricItem{
						ID:         2,
						Name:       "method02",
						SystemName: "desc02",
						Unit:       "1",
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
	list, err := c.ListProductMetrics(productID)
	if err != nil {
		t.Fatal(err)
	}

	if list == nil {
		t.Fatal("list returned nil")
	}

	if len(list.Metrics) != 2 {
		t.Fatalf("Then number of list items does not match. Expected [%d]; got [%d]", 2, len(list.Metrics))
	}
}

func TestCreateProductMetric(t *testing.T) {
	var (
		productID int64 = 97
		params          = Params{"friendly_name": "metric02"}
		endpoint        = fmt.Sprintf(productMetricListResourceEndpoint, productID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPost {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPost, req.Method)
		}

		item := &MetricJSON{
			Element: MetricItem{
				ID:         10,
				Name:       params["friendly_name"],
				SystemName: params["friendly_name"],
				Unit:       "1",
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
	obj, err := c.CreateProductMetric(productID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("returned nil")
	}

	if obj.Element.ID != 10 {
		t.Fatalf("ID does not match. Expected [%d]; got [%d]", 10, obj.Element.ID)
	}

	if obj.Element.Name != params["friendly_name"] {
		t.Fatalf("name does not match. Expected [%s]; got [%s]", params["friendly_name"], obj.Element.Name)
	}
}

func TestDeleteProductMetric(t *testing.T) {
	var (
		productID int64 = 97
		itemID    int64 = 123325
		endpoint        = fmt.Sprintf(productMetricResourceEndpoint, productID, itemID)
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
	err := c.DeleteProductMetric(productID, itemID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadProductMetric(t *testing.T) {
	var (
		productID int64 = 97
		itemID    int64 = 123325
		endpoint        = fmt.Sprintf(productMetricResourceEndpoint, productID, itemID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		item := &MetricJSON{
			Element: MetricItem{
				ID:         itemID,
				Name:       "someName",
				SystemName: "someName2",
				Unit:       "1",
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
	obj, err := c.ProductMetric(productID, itemID)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("returned nil")
	}

	if obj.Element.ID != itemID {
		t.Fatalf("ID does not match. Expected [%d]; got [%d]", itemID, obj.Element.ID)
	}
}

func TestUpdateProductMetric(t *testing.T) {
	var (
		productID int64 = 98765
		itemID    int64 = 123325
		endpoint        = fmt.Sprintf(productMetricResourceEndpoint, productID, itemID)
		params          = Params{"description": "newDescr"}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPut {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPut, req.Method)
		}

		item := &MetricJSON{
			Element: MetricItem{
				ID:          itemID,
				Name:        "someName",
				SystemName:  "someName2",
				Description: params["description"],
				Unit:        "1",
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
	obj, err := c.UpdateProductMetric(productID, itemID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("returned nil")
	}

	if obj.Element.ID != itemID {
		t.Fatalf("ID does not match. Expected [%d]; got [%d]", itemID, obj.Element.ID)
	}

	if obj.Element.Description != params["description"] {
		t.Fatalf("description does not match. Expected [%s]; got [%s]", params["description"], obj.Element.Description)
	}
}

func TestListProductMappingRules(t *testing.T) {
	var (
		productID int64 = 97
		endpoint        = fmt.Sprintf(productMappingRuleListResourceEndpoint, productID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		list := MappingRuleJSONList{
			MappingRules: []MappingRuleJSON{
				{
					Element: MappingRuleItem{
						ID:         1,
						MetricID:   2,
						Pattern:    "/v1",
						HTTPMethod: "GET",
					},
				},
				{
					Element: MappingRuleItem{
						ID:         2,
						MetricID:   2,
						Pattern:    "/v2",
						HTTPMethod: "GET",
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
	list, err := c.ListProductMappingRules(productID)
	if err != nil {
		t.Fatal(err)
	}

	if list == nil {
		t.Fatal("list returned nil")
	}

	if len(list.MappingRules) != 2 {
		t.Fatalf("Then number of list items does not match. Expected [%d]; got [%d]", 2, len(list.MappingRules))
	}
}

func TestCreateProductMappingRule(t *testing.T) {
	var (
		productID int64 = 97
		params          = Params{"pattern": "/somePath"}
		endpoint        = fmt.Sprintf(productMappingRuleListResourceEndpoint, productID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPost {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPost, req.Method)
		}

		item := &MappingRuleJSON{
			Element: MappingRuleItem{
				ID:         10,
				MetricID:   2,
				Pattern:    params["pattern"],
				HTTPMethod: "GET",
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
	obj, err := c.CreateProductMappingRule(productID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("returned nil")
	}

	if obj.Element.ID != 10 {
		t.Fatalf("ID does not match. Expected [%d]; got [%d]", 10, obj.Element.ID)
	}

	if obj.Element.Pattern != params["pattern"] {
		t.Fatalf("name does not match. Expected [%s]; got [%s]", params["pattern"], obj.Element.Pattern)
	}
}

func TestDeleteProductMappingRule(t *testing.T) {
	var (
		productID int64 = 97
		itemID    int64 = 123325
		endpoint        = fmt.Sprintf(productMappingRuleResourceEndpoint, productID, itemID)
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
	err := c.DeleteProductMappingRule(productID, itemID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadProductMappingRule(t *testing.T) {
	var (
		productID int64 = 97
		itemID    int64 = 123325
		endpoint        = fmt.Sprintf(productMappingRuleResourceEndpoint, productID, itemID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		item := &MappingRuleJSON{
			Element: MappingRuleItem{
				ID:         itemID,
				MetricID:   2,
				Pattern:    "/v1",
				HTTPMethod: "GET",
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
	obj, err := c.ProductMappingRule(productID, itemID)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("returned nil")
	}

	if obj.Element.ID != itemID {
		t.Fatalf("ID does not match. Expected [%d]; got [%d]", itemID, obj.Element.ID)
	}
}

func TestUpdateProductMappingRule(t *testing.T) {
	var (
		productID int64 = 98765
		itemID    int64 = 123325
		endpoint        = fmt.Sprintf(productMappingRuleResourceEndpoint, productID, itemID)
		params          = Params{"pattern": "/newPath"}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPut {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPut, req.Method)
		}

		item := &MappingRuleJSON{
			Element: MappingRuleItem{
				ID:         itemID,
				MetricID:   2,
				Pattern:    params["pattern"],
				HTTPMethod: "GET",
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
	obj, err := c.UpdateProductMappingRule(productID, itemID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("returned nil")
	}

	if obj.Element.ID != itemID {
		t.Fatalf("ID does not match. Expected [%d]; got [%d]", itemID, obj.Element.ID)
	}

	if obj.Element.Pattern != params["pattern"] {
		t.Fatalf("description does not match. Expected [%s]; got [%s]", params["pattern"], obj.Element.Pattern)
	}
}

func TestReadProductProxy(t *testing.T) {
	var (
		productID          int64 = 97
		productionEndpoint       = "prod.example.com"
		endpoint                 = fmt.Sprintf(productProxyResourceEndpoint, productID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		item := &ProxyJSON{
			Element: ProxyItem{
				Endpoint:        productionEndpoint,
				SandboxEndpoint: "staging.example.com",
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
	obj, err := c.ProductProxy(productID)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("returned nil")
	}

	if obj.Element.Endpoint != productionEndpoint {
		t.Fatalf("Endpoint does not match. Expected [%s]; got [%s]", productionEndpoint, obj.Element.Endpoint)
	}
}

func TestUpdateProductProxy(t *testing.T) {
	var (
		productID          int64 = 97
		productionEndpoint       = "prod.example.com"
		endpoint                 = fmt.Sprintf(productProxyResourceEndpoint, productID)
		params                   = Params{"endpoint": productionEndpoint}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPut {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPut, req.Method)
		}

		item := &ProxyJSON{
			Element: ProxyItem{
				Endpoint:        productionEndpoint,
				SandboxEndpoint: "staging.example.com",
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
	obj, err := c.UpdateProductProxy(productID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("returned nil")
	}

	if obj.Element.Endpoint != productionEndpoint {
		t.Fatalf("Endpoint does not match. Expected [%s]; got [%s]", productionEndpoint, obj.Element.Endpoint)
	}

}

func TestDeployProductProxy(t *testing.T) {
	var (
		productID          int64 = 97
		productionEndpoint       = "prod.example.com"
		endpoint                 = fmt.Sprintf(productProxyDeployResourceEndpoint, productID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPost {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPost, req.Method)
		}

		item := &ProxyJSON{
			Element: ProxyItem{
				Endpoint:        productionEndpoint,
				SandboxEndpoint: "staging.example.com",
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
	obj, err := c.DeployProductProxy(productID)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("returned nil")
	}

	if obj.Element.Endpoint != productionEndpoint {
		t.Fatalf("Endpoint does not match. Expected [%s]; got [%s]", productionEndpoint, obj.Element.Endpoint)
	}
}

func TestReadProduct(t *testing.T) {
	var (
		productID int64 = 98765
		endpoint        = fmt.Sprintf(productResourceEndpoint, productID)
		product         = &Product{
			Element: ProductItem{
				ID:   productID,
				Name: "myProduct",
			},
		}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		responseBodyBytes, err := json.Marshal(*product)
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
	obj, err := c.Product(productID)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("backendapi returned nil")
	}

	if *obj != *product {
		t.Fatalf("Expected %v; got %v", *product, *obj)
	}
}

func TestCreateProduct(t *testing.T) {
	var (
		name      string = "producttest1"
		productID int64  = 98765
		params    Params = Params{}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != productListResourceEndpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", productListResourceEndpoint, req.URL.Path)
		}

		if req.Method != http.MethodPost {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPost, req.Method)
		}

		product := &Product{
			Element: ProductItem{
				ID:   productID,
				Name: name,
			},
		}

		responseBodyBytes, err := json.Marshal(product)
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
	obj, err := c.CreateProduct(name, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("CreateProduct returned nil")
	}

	if obj.Element.ID != productID {
		t.Fatalf("obj ID does not match. Expected [%d]; got [%d]", productID, obj.Element.ID)
	}

	if obj.Element.Name != name {
		t.Fatalf("obj name does not match. Expected [%s]; got [%s]", name, obj.Element.Name)
	}
}

func TestUpdateProduct(t *testing.T) {
	var (
		productID int64 = 98765
		endpoint        = fmt.Sprintf(productResourceEndpoint, productID)
		params          = Params{"name": "newName"}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPut {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		product := &Product{
			Element: ProductItem{
				ID:   productID,
				Name: "newName",
			},
		}

		responseBodyBytes, err := json.Marshal(product)
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
	obj, err := c.UpdateProduct(productID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("product returned nil")
	}

	if obj.Element.ID != productID {
		t.Fatalf("obj ID does not match. Expected [%d]; got [%d]", productID, obj.Element.ID)
	}

	if obj.Element.Name != "newName" {
		t.Fatalf("obj name does not match. Expected [%s]; got [%s]", "newName", obj.Element.Name)
	}
}

func TestDeleteProduct(t *testing.T) {
	var (
		productID int64 = 98765
		endpoint        = fmt.Sprintf(productResourceEndpoint, productID)
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
	err := c.DeleteProduct(productID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListProducts(t *testing.T) {
	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != productListResourceEndpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", backendListResourceEndpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(helperLoadBytes(t, "product_list_fixture.json"))),
			Header:     make(http.Header),
		}
	})

	credential := "someAccessToken"
	c := NewThreeScale(NewTestAdminPortal(t), credential, httpClient)
	productList, err := c.ListProducts()
	if err != nil {
		t.Fatal(err)
	}

	if productList == nil {
		t.Fatal("product list returned nil")
	}

	if len(productList.Products) != 2 {
		t.Fatalf("Then number of products does not match. Expected [%d]; got [%d]", 2, len(productList.Products))
	}
}
