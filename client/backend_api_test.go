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

func TestListBackendApi(t *testing.T) {
	backendGenerator := func(startingIndex, n int) BackendApiList {
		pList := BackendApiList{
			Backends: make([]BackendApi, 0, n),
		}

		for idx := 0; idx < n; idx++ {
			pList.Backends = append(pList.Backends, BackendApi{
				Element: BackendApiItem{ID: int64(idx + startingIndex)},
			})
		}

		return pList
	}

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		// Will serve: 3 pages
		// page 1 => BACKENDS_PER_PAGE
		// page 2 => BACKENDS_PER_PAGE
		// page 3 => 51

		if req.URL.Path != backendListResourceEndpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", backendListResourceEndpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		if req.URL.Query().Get("per_page") != strconv.Itoa(BACKENDS_PER_PAGE) {
			t.Fatalf("per_page param does not match. Expected [%d]; got [%s]", BACKENDS_PER_PAGE, req.URL.Query().Get("per_page"))
		}

		var list BackendApiList

		if req.URL.Query().Get("page") == "1" {
			list = backendGenerator(BACKENDS_PER_PAGE*0, BACKENDS_PER_PAGE)
		} else if req.URL.Query().Get("page") == "2" {
			list = backendGenerator(BACKENDS_PER_PAGE*1, BACKENDS_PER_PAGE)
		} else if req.URL.Query().Get("page") == "3" {
			list = backendGenerator(BACKENDS_PER_PAGE*2, 51)
		} else {
			t.Fatalf("page param unexpected value; got [%s]", req.URL.Query().Get("page"))
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
	backendApiList, err := c.ListBackendApis()
	if err != nil {
		t.Fatal(err)
	}

	if backendApiList == nil {
		t.Fatal("backend list returned nil")
	}

	if len(backendApiList.Backends) != 2*BACKENDS_PER_PAGE+51 {
		t.Fatalf("Then number of backend_api's does not match. Expected [%d]; got [%d]", 2*BACKENDS_PER_PAGE+51, len(backendApiList.Backends))
	}
}

func TestListBackendApisPerPage(t *testing.T) {
	t.Run("page and per_page params used", func(subT *testing.T) {
		var (
			pageNum int = 4
			perPage int = 2
		)
		httpClient := NewTestClient(func(req *http.Request) *http.Response {
			if req.URL.Path != backendListResourceEndpoint {
				subT.Fatalf("Path does not match. Expected [%s]; got [%s]", backendListResourceEndpoint, req.URL.Path)
			}

			if req.Method != http.MethodGet {
				subT.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
			}

			if req.URL.Query().Get("page") != strconv.Itoa(pageNum) {
				subT.Fatalf("page param does not match. Expected [%d]; got [%s]", pageNum, req.URL.Query().Get("page"))
			}

			if req.URL.Query().Get("per_page") != strconv.Itoa(perPage) {
				subT.Fatalf("page param does not match. Expected [%d]; got [%s]", perPage, req.URL.Query().Get("per_page"))
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(helperLoadBytes(subT, "backend_api_list_fixture.json"))),
				Header:     make(http.Header),
			}
		})

		credential := "someAccessToken"
		c := NewThreeScale(NewTestAdminPortal(subT), credential, httpClient)
		backendList, err := c.ListBackendApisPerPage(pageNum, perPage)
		if err != nil {
			subT.Fatal(err)
		}

		if backendList == nil {
			subT.Fatal("backend list returned nil")
		}

		if len(backendList.Backends) != 3 {
			subT.Fatalf("Then number of backends does not match. Expected [%d]; got [%d]", 3, len(backendList.Backends))
		}
	})

	t.Run("page and per_page params not used", func(subT *testing.T) {
		httpClient := NewTestClient(func(req *http.Request) *http.Response {
			if req.URL.Path != productListResourceEndpoint {
				subT.Fatalf("Path does not match. Expected [%s]; got [%s]", productListResourceEndpoint, req.URL.Path)
			}

			if req.Method != http.MethodGet {
				subT.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
			}

			if req.URL.Query().Get("page") != "" {
				subT.Fatalf("Query param page does not match. Expected empty; got [%s]", req.URL.Query().Get("page"))
			}

			if req.URL.Query().Get("per_page") != "" {
				subT.Fatalf("page param does not match. Expected empty; got [%s]", req.URL.Query().Get("per_page"))
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(helperLoadBytes(subT, "product_list_fixture.json"))),
				Header:     make(http.Header),
			}
		})

		credential := "someAccessToken"
		c := NewThreeScale(NewTestAdminPortal(subT), credential, httpClient)
		productList, err := c.ListProductsPerPage()
		if err != nil {
			subT.Fatal(err)
		}

		if productList == nil {
			subT.Fatal("product list returned nil")
		}

		if len(productList.Products) != 2 {
			subT.Fatalf("Then number of products does not match. Expected [%d]; got [%d]", 2, len(productList.Products))
		}
	})
}

func TestCreateBackendApi(t *testing.T) {
	params := Params{
		"name":             "backendapitest1",
		"private_endpoint": "https://echo-api.3scale.net:443",
	}

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != backendListResourceEndpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", backendListResourceEndpoint, req.URL.Path)
		}

		if req.Method != http.MethodPost {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPost, req.Method)
		}

		backendApi := &BackendApi{
			Element: BackendApiItem{
				ID:              45498,
				Name:            params["name"],
				PrivateEndpoint: params["private_endpoint"],
			},
		}

		responseBodyBytes, err := json.Marshal(backendApi)
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
	backendAPI, err := c.CreateBackendApi(params)
	if err != nil {
		t.Fatal(err)
	}

	if backendAPI == nil {
		t.Fatal("backendapi returned nil")
	}

	if backendAPI.Element.ID != 45498 {
		t.Fatalf("backend_api ID does not match. Expected [%d]; got [%d]", 45498, backendAPI.Element.ID)
	}

	if backendAPI.Element.Name != params["name"] {
		t.Fatalf("backend_api privateEndpoint does not match. Expected [%s]; got [%s]", params["name"], backendAPI.Element.Name)
	}

	if backendAPI.Element.PrivateEndpoint != params["private_endpoint"] {
		t.Fatalf("backend_api privateEndpoint does not match. Expected [%s]; got [%s]", params["private_endpoint"], backendAPI.Element.PrivateEndpoint)
	}
}

func TestDeleteBackendApi(t *testing.T) {
	var backendAAPIID int64 = 12345
	endpoint := fmt.Sprintf(backendResourceEndpoint, backendAAPIID)

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
	err := c.DeleteBackendApi(backendAAPIID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadBackendApi(t *testing.T) {
	var (
		backendapiID int64 = 98765
		endpoint           = fmt.Sprintf(backendResourceEndpoint, backendapiID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		backendApi := &BackendApi{
			Element: BackendApiItem{
				ID:              backendapiID,
				Name:            "someName",
				PrivateEndpoint: "https://example.com:81",
			},
		}

		responseBodyBytes, err := json.Marshal(backendApi)
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
	backendAPI, err := c.BackendApi(backendapiID)
	if err != nil {
		t.Fatal(err)
	}

	if backendAPI == nil {
		t.Fatal("backendapi returned nil")
	}

	if backendAPI.Element.ID != backendapiID {
		t.Fatalf("backend_api ID does not match. Expected [%d]; got [%d]", backendapiID, backendAPI.Element.ID)
	}

	if backendAPI.Element.Name != "someName" {
		t.Fatalf("backend_api privateEndpoint does not match. Expected [%s]; got [%s]", "someName", backendAPI.Element.Name)
	}

	if backendAPI.Element.PrivateEndpoint != "https://example.com:81" {
		t.Fatalf("backend_api privateEndpoint does not match. Expected [%s]; got [%s]", "https://example.com:81", backendAPI.Element.PrivateEndpoint)
	}
}

func TestUpdateBackendApi(t *testing.T) {
	var (
		backendapiID int64 = 98765
		endpoint           = fmt.Sprintf(backendResourceEndpoint, backendapiID)
		params             = Params{"name": "newName"}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPut {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		backendApi := &BackendApi{
			Element: BackendApiItem{
				ID:              backendapiID,
				Name:            params["name"],
				PrivateEndpoint: "https://example.com:81",
			},
		}

		responseBodyBytes, err := json.Marshal(backendApi)
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
	backendAPI, err := c.UpdateBackendApi(backendapiID, params)
	if err != nil {
		t.Fatal(err)
	}

	if backendAPI == nil {
		t.Fatal("backendapi returned nil")
	}

	if backendAPI.Element.ID != backendapiID {
		t.Fatalf("backend_api ID does not match. Expected [%d]; got [%d]", backendapiID, backendAPI.Element.ID)
	}

	if backendAPI.Element.Name != params["name"] {
		t.Fatalf("backend_api privateEndpoint does not match. Expected [%s]; got [%s]", params["name"], backendAPI.Element.Name)
	}
}

func TestListBackendApiMethods(t *testing.T) {
	var (
		backendapiID int64 = 98765
		hitsID       int64 = 1
		endpoint           = fmt.Sprintf(backendMethodListResourceEndpoint, backendapiID, hitsID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		methodList := &MethodList{
			Methods: []Method{
				{
					Element: MethodItem{
						ID:         1,
						Name:       "method01",
						SystemName: "method01",
					},
				},
				{
					Element: MethodItem{
						ID:         2,
						Name:       "method02",
						SystemName: "method02",
					},
				},
			},
		}

		responseBodyBytes, err := json.Marshal(methodList)
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
	list, err := c.ListBackendapiMethods(backendapiID, hitsID)
	if err != nil {
		t.Fatal(err)
	}

	if list == nil {
		t.Fatal("backend method list returned nil")
	}

	if len(list.Methods) != 2 {
		t.Fatalf("Then number of backend_api method's does not match. Expected [%d]; got [%d]", 2, len(list.Methods))
	}
}

func TestCreateBackendApiMethod(t *testing.T) {
	var (
		backendapiID int64 = 98765
		hitsID       int64 = 1
		endpoint           = fmt.Sprintf(backendMethodListResourceEndpoint, backendapiID, hitsID)
		params             = Params{"friendly_name": "method5"}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPost {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPost, req.Method)
		}

		method := &Method{
			Element: MethodItem{
				ID:         10,
				Name:       params["friendly_name"],
				SystemName: params["friendly_name"],
			},
		}

		responseBodyBytes, err := json.Marshal(method)
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
	obj, err := c.CreateBackendApiMethod(backendapiID, hitsID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("backend method create returned nil")
	}

	if obj.Element.ID != 10 {
		t.Fatalf("backend_api method ID does not match. Expected [%d]; got [%d]", 10, obj.Element.ID)
	}

	if obj.Element.Name != params["friendly_name"] {
		t.Fatalf("backend_api method name does not match. Expected [%s]; got [%s]", params["friendly_name"], obj.Element.Name)
	}

}

func TestDeleteBackendApiMethod(t *testing.T) {
	var (
		backendapiID int64 = 98765
		hitsID       int64 = 1
		methodID     int64 = 123325
		endpoint           = fmt.Sprintf(backendMethodResourceEndpoint, backendapiID, hitsID, methodID)
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
	err := c.DeleteBackendApiMethod(backendapiID, hitsID, methodID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadBackendApiMethod(t *testing.T) {
	var (
		backendapiID int64 = 98765
		hitsID       int64 = 1
		methodID     int64 = 123325
		endpoint           = fmt.Sprintf(backendMethodResourceEndpoint, backendapiID, hitsID, methodID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		method := &Method{
			Element: MethodItem{
				ID:         methodID,
				Name:       "someName",
				SystemName: "someName2",
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
	obj, err := c.BackendApiMethod(backendapiID, hitsID, methodID)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("backendapi method returned nil")
	}

	if obj.Element.ID != methodID {
		t.Fatalf("backend_api ID does not match. Expected [%d]; got [%d]", methodID, obj.Element.ID)
	}
}

func TestUpdateBackendApiMethod(t *testing.T) {
	var (
		backendapiID int64 = 98765
		hitsID       int64 = 1
		methodID     int64 = 123325
		endpoint           = fmt.Sprintf(backendMethodResourceEndpoint, backendapiID, hitsID, methodID)
		params             = Params{"description": "newDescr"}
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
	obj, err := c.UpdateBackendApiMethod(backendapiID, hitsID, methodID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("backendapi method returned nil")
	}

	if obj.Element.ID != methodID {
		t.Fatalf("backend_api method ID does not match. Expected [%d]; got [%d]", methodID, obj.Element.ID)
	}

	if obj.Element.Description != params["description"] {
		t.Fatalf("backend_api method description does not match. Expected [%s]; got [%s]", params["description"], obj.Element.Description)
	}
}

func TestListBackendApiMetrics(t *testing.T) {
	var (
		backendapiID int64 = 98765
		endpoint           = fmt.Sprintf(backendMetricListResourceEndpoint, backendapiID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		list := &MetricJSONList{
			Metrics: []MetricJSON{
				{
					Element: MetricItem{
						ID:         1,
						Name:       "metric01",
						SystemName: "metric01",
					},
				},
				{
					Element: MetricItem{
						ID:         2,
						Name:       "metric02",
						SystemName: "metric02",
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
	list, err := c.ListBackendapiMetrics(backendapiID)
	if err != nil {
		t.Fatal(err)
	}

	if list == nil {
		t.Fatal("backend metric list returned nil")
	}

	if len(list.Metrics) != 2 {
		t.Fatalf("Then number of backend_api metric's does not match. Expected [%d]; got [%d]", 2, len(list.Metrics))
	}
}

func TestCreateBackendApiMetric(t *testing.T) {
	var (
		backendapiID int64 = 98765
		endpoint           = fmt.Sprintf(backendMetricListResourceEndpoint, backendapiID)
		params             = Params{"friendly_name": "metric05", "unit": "1"}
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
				Unit:       params["unit"],
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
	obj, err := c.CreateBackendApiMetric(backendapiID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("backend metric create returned nil")
	}

	if obj.Element.ID != 10 {
		t.Fatalf("backend_api metric ID does not match. Expected [%d]; got [%d]", 10, obj.Element.ID)
	}

	if obj.Element.Name != params["friendly_name"] {
		t.Fatalf("backend_api metric name does not match. Expected [%s]; got [%s]", params["friendly_name"], obj.Element.Name)
	}
}

func TestDeleteBackendApiMetric(t *testing.T) {
	var (
		backendapiID int64 = 98765
		metricID     int64 = 123325
		endpoint           = fmt.Sprintf(backendMetricResourceEndpoint, backendapiID, metricID)
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
	err := c.DeleteBackendApiMetric(backendapiID, metricID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadBackendApiMetric(t *testing.T) {
	var (
		backendapiID int64 = 98765
		metricID     int64 = 123325
		endpoint           = fmt.Sprintf(backendMetricResourceEndpoint, backendapiID, metricID)
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
				ID:         metricID,
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
	obj, err := c.BackendApiMetric(backendapiID, metricID)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("backendapi metric returned nil")
	}

	if obj.Element.ID != metricID {
		t.Fatalf("backend_api metric ID does not match. Expected [%d]; got [%d]", metricID, obj.Element.ID)
	}
}

func TestUpdateBackendApiMetric(t *testing.T) {
	var (
		backendapiID int64 = 98765
		metricID     int64 = 123325
		endpoint           = fmt.Sprintf(backendMetricResourceEndpoint, backendapiID, metricID)
		params             = Params{"description": "newDescr"}
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
				ID:          metricID,
				Name:        "someName",
				SystemName:  "someName2",
				Description: params["description"],
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
	obj, err := c.UpdateBackendApiMetric(backendapiID, metricID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("backendapi metric returned nil")
	}

	if obj.Element.ID != metricID {
		t.Fatalf("backend_api metric ID does not match. Expected [%d]; got [%d]", metricID, obj.Element.ID)
	}

	if obj.Element.Description != params["description"] {
		t.Fatalf("backend_api metric description does not match. Expected [%s]; got [%s]", params["description"], obj.Element.Description)
	}
}

func TestListBackendApiMappingRules(t *testing.T) {
	var (
		backendapiID int64 = 98765
		endpoint           = fmt.Sprintf(backendMRListResourceEndpoint, backendapiID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		list := &MappingRuleJSONList{
			MappingRules: []MappingRuleJSON{
				{
					Element: MappingRuleItem{
						ID:      1,
						Pattern: "/v1",
					},
				},
				{
					Element: MappingRuleItem{
						ID:      2,
						Pattern: "/v2",
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
	list, err := c.ListBackendapiMappingRules(backendapiID)
	if err != nil {
		t.Fatal(err)
	}

	if list == nil {
		t.Fatal("backend mapping rule list returned nil")
	}

	if len(list.MappingRules) != 2 {
		t.Fatalf("Then number of backend_api mapping rule's does not match. Expected [%d]; got [%d]", 2, len(list.MappingRules))
	}
}

func TestCreateBackendApiMappingRule(t *testing.T) {
	var (
		backendapiID int64 = 98765
		endpoint           = fmt.Sprintf(backendMRListResourceEndpoint, backendapiID)
		params             = Params{"pattern": "/v1", "http_method": "GET", "metric_id": "12"}
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
				Pattern:    params["pattern"],
				HTTPMethod: params["http_method"],
				MetricID:   12,
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
	obj, err := c.CreateBackendapiMappingRule(backendapiID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("backend mapping rule create returned nil")
	}

	if obj.Element.ID != 10 {
		t.Fatalf("backend_api mapping rule ID does not match. Expected [%d]; got [%d]", 10, obj.Element.ID)
	}

	if obj.Element.Pattern != params["pattern"] {
		t.Fatalf("backend_api mapping rule pattern does not match. Expected [%s]; got [%s]", params["pattern"], obj.Element.Pattern)
	}
}

func TestDeleteBackendApiMappingRule(t *testing.T) {
	var (
		backendapiID int64 = 98765
		mrID         int64 = 123325
		endpoint           = fmt.Sprintf(backendMRResourceEndpoint, backendapiID, mrID)
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
	err := c.DeleteBackendapiMappingRule(backendapiID, mrID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadBackendApiMappingRule(t *testing.T) {
	var (
		backendapiID int64 = 98765
		mrID         int64 = 123325
		endpoint           = fmt.Sprintf(backendMRResourceEndpoint, backendapiID, mrID)
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
				ID:         mrID,
				Pattern:    "/v1",
				HTTPMethod: "GET",
				MetricID:   123454,
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
	obj, err := c.BackendapiMappingRule(backendapiID, mrID)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("backendapi mapping rule returned nil")
	}

	if obj.Element.ID != mrID {
		t.Fatalf("backend_api mapping rule ID does not match. Expected [%d]; got [%d]", mrID, obj.Element.ID)
	}
}

func TestUpdateBackendApiMappingRule(t *testing.T) {
	var (
		backendapiID int64 = 98765
		mrID         int64 = 123325
		endpoint           = fmt.Sprintf(backendMRResourceEndpoint, backendapiID, mrID)
		params             = Params{"pattern": "/v2"}
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
				ID:         mrID,
				Pattern:    params["pattern"],
				HTTPMethod: "GET",
				MetricID:   123454,
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
	obj, err := c.UpdateBackendapiMappingRule(backendapiID, mrID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("backendapi mapping rule returned nil")
	}

	if obj.Element.ID != mrID {
		t.Fatalf("backend_api mapping rule ID does not match. Expected [%d]; got [%d]", mrID, obj.Element.ID)
	}

	if obj.Element.Pattern != params["pattern"] {
		t.Fatalf("backend_api mapping rule pattern does not match. Expected [%s]; got [%s]", params["pattern"], obj.Element.Pattern)
	}
}

func TestListBackendApiUsages(t *testing.T) {
	var (
		productID int64 = 98765
		endpoint        = fmt.Sprintf(backendUsageListResourceEndpoint, productID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		list := BackendAPIUsageList{
			{
				Element: BackendAPIUsageItem{
					ID:   1,
					Path: "/v1",
				},
			},
			{
				Element: BackendAPIUsageItem{
					ID:   2,
					Path: "/v2",
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
	list, err := c.ListBackendapiUsages(productID)
	if err != nil {
		t.Fatal(err)
	}

	if list == nil {
		t.Fatal("backend usage list returned nil")
	}

	if len(list) != 2 {
		t.Fatalf("Then number of backend_api usage's does not match. Expected [%d]; got [%d]", 2, len(list))
	}
}

func TestCreateBackendApiUsage(t *testing.T) {
	var (
		productID int64 = 98765
		endpoint        = fmt.Sprintf(backendUsageListResourceEndpoint, productID)
		params          = Params{"path": "/v1", "backend_api_id": "12345"}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPost {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPost, req.Method)
		}

		item := &BackendAPIUsage{
			Element: BackendAPIUsageItem{
				ID:           10,
				Path:         params["path"],
				BackendAPIID: 12345,
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
	obj, err := c.CreateBackendapiUsage(productID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("backend usage create returned nil")
	}

	if obj.Element.ID != 10 {
		t.Fatalf("backend_api usage ID does not match. Expected [%d]; got [%d]", 10, obj.Element.ID)
	}

	if obj.Element.Path != params["path"] {
		t.Fatalf("backend_api usage path does not match. Expected [%s]; got [%s]", params["path"], obj.Element.Path)
	}
}

func TestDeleteBackendApiUsage(t *testing.T) {
	var (
		productID      int64 = 98765
		backendUsageID int64 = 123325
		endpoint             = fmt.Sprintf(backendUsageResourceEndpoint, productID, backendUsageID)
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
	err := c.DeleteBackendapiUsage(productID, backendUsageID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadBackendApiUsage(t *testing.T) {
	var (
		productID      int64 = 98765
		backendUsageID int64 = 123325
		endpoint             = fmt.Sprintf(backendUsageResourceEndpoint, productID, backendUsageID)
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		item := &BackendAPIUsage{
			Element: BackendAPIUsageItem{
				ID:           backendUsageID,
				Path:         "/v1",
				BackendAPIID: 12345,
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
	obj, err := c.BackendapiUsage(productID, backendUsageID)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("backendapi usage returned nil")
	}

	if obj.Element.ID != backendUsageID {
		t.Fatalf("backend_api usage ID does not match. Expected [%d]; got [%d]", backendUsageID, obj.Element.ID)
	}
}

func TestUpdateBackendApiUsage(t *testing.T) {
	var (
		productID      int64 = 98765
		backendUsageID int64 = 123325
		endpoint             = fmt.Sprintf(backendUsageResourceEndpoint, productID, backendUsageID)
		params               = Params{"path": "/v2"}
	)

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.Path != endpoint {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
		}

		if req.Method != http.MethodPut {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodPut, req.Method)
		}

		item := &BackendAPIUsage{
			Element: BackendAPIUsageItem{
				ID:           backendUsageID,
				Path:         "/v2",
				BackendAPIID: 12345,
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
	obj, err := c.UpdateBackendapiUsage(productID, backendUsageID, params)
	if err != nil {
		t.Fatal(err)
	}

	if obj == nil {
		t.Fatal("backendapi usage returned nil")
	}

	if obj.Element.ID != backendUsageID {
		t.Fatalf("backend_api usage ID does not match. Expected [%d]; got [%d]", backendUsageID, obj.Element.ID)
	}

	if obj.Element.Path != params["path"] {
		t.Fatalf("backend_api usage path does not match. Expected [%s]; got [%s]", params["path"], obj.Element.Path)
	}
}
