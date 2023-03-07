package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func TestListAccountProxyConfigsParams(t *testing.T) {
	var (
		env        string = "production"
		endpoint   string = fmt.Sprintf(accountProxyConfigGet, env)
		credential string = "someAccessToken"
	)

	errorTests := []struct {
		Name             string
		version          *string
		host             *string
		expectedRawQuery string
	}{
		{"no version or host", nil, nil, ""},
		{"only version", &[]string{"latest"}[0], nil, "version=latest"},
		{"only host", nil, &[]string{"example.com"}[0], "host=example.com"},
		{"version and host set", &[]string{"latest"}[0], &[]string{"example.com"}[0], "host=example.com&version=latest"},
	}

	for _, tt := range errorTests {
		t.Run(tt.Name, func(subT *testing.T) {
			httpClient := NewTestClient(func(req *http.Request) *http.Response {
				if req.URL.Path != endpoint {
					subT.Fatalf("Path does not match. Expected [%s]; got [%s]", endpoint, req.URL.Path)
				}

				if req.URL.RawQuery != tt.expectedRawQuery {
					subT.Fatalf("RawPath does not match. Expected [%s]; got [%s]", tt.expectedRawQuery, req.URL.RawQuery)
				}

				if req.Method != http.MethodGet {
					subT.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(helperLoadBytes(subT, "account_proxy_fixture.json"))),
					Header:     make(http.Header),
				}
			})

			c := NewThreeScale(NewTestAdminPortal(subT), credential, httpClient)
			proxyConfigList, err := c.ListAccountProxyConfigsPerPage(env, tt.version, tt.host)
			if err != nil {
				subT.Fatal(err)
			}

			if proxyConfigList.ProxyConfigs == nil {
				subT.Fatal("backend list returned nil")
			}

			if len(proxyConfigList.ProxyConfigs) != 1 {
				subT.Fatalf("Then number of proxy_configs does not match. Expected [%d]; got [%d]", 1, len(proxyConfigList.ProxyConfigs))
			}
		})
	}
}

func TestListAccountProxyConfigsPagination(t *testing.T) {
	configGenerator := func(startingIndex, n int) ProxyConfigList {
		pList := ProxyConfigList{
			ProxyConfigs: make([]ProxyConfigElement, 0, n),
		}

		for idx := 0; idx < n; idx++ {
			pList.ProxyConfigs = append(pList.ProxyConfigs, ProxyConfigElement{
				ProxyConfig: ProxyConfig{ID: idx + startingIndex},
			})
		}

		return pList
	}

	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		// Will serve: 3 pages
		// page 1 => PROXYCONFIGS_PER_PAGE
		// page 2 => PROXYCONFIGS_PER_PAGE
		// page 3 => 51
		if req.URL.Path != fmt.Sprintf(accountProxyConfigGet, "production") {
			t.Fatalf("Path does not match. Expected [%s]; got [%s]", accountProxyConfigGet, req.URL.Path)
		}

		if req.Method != http.MethodGet {
			t.Fatalf("Method does not match. Expected [%s]; got [%s]", http.MethodGet, req.Method)
		}

		if req.URL.Query().Get("per_page") != strconv.Itoa(PROXYCONFIGS_PER_PAGE) {
			t.Fatalf("per_page param does not match. Expected [%d]; got [%s]", PROXYCONFIGS_PER_PAGE, req.URL.Query().Get("per_page"))
		}

		var list ProxyConfigList

		if req.URL.Query().Get("page") == "1" {
			list = configGenerator(PROXYCONFIGS_PER_PAGE*0, PROXYCONFIGS_PER_PAGE)
		} else if req.URL.Query().Get("page") == "2" {
			list = configGenerator(PROXYCONFIGS_PER_PAGE*1, PROXYCONFIGS_PER_PAGE)
		} else if req.URL.Query().Get("page") == "3" {
			list = configGenerator(PROXYCONFIGS_PER_PAGE*2, 51)
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
	pList, err := c.ListAccountProxyConfigs("production", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	if pList == nil {
		t.Fatal("configlist returned nil")
	}

	if len(pList.ProxyConfigs) != 2*PROXYCONFIGS_PER_PAGE+51 {
		t.Fatalf("Then number of proxy configs does not match. Expected [%d]; got [%d]", 2*PROXYCONFIGS_PER_PAGE+51, len(pList.ProxyConfigs))
	}
}

func TestListAccountProxyConfigsPerPage(t *testing.T) {
	var (
		env string = "production"
	)
	t.Run("page and per_page params used", func(subT *testing.T) {
		var (
			pageNum int = 4
			perPage int = 1
		)

		httpClient := NewTestClient(func(req *http.Request) *http.Response {
			if req.URL.Path != fmt.Sprintf(accountProxyConfigGet, env) {
				subT.Fatalf("Path does not match. Expected [%s]; got [%s]", fmt.Sprintf(accountProxyConfigGet, env), req.URL.Path)
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
				Body:       ioutil.NopCloser(bytes.NewReader(helperLoadBytes(subT, "account_proxy_fixture.json"))),
				Header:     make(http.Header),
			}
		})

		credential := "someAccessToken"
		c := NewThreeScale(NewTestAdminPortal(subT), credential, httpClient)
		pList, err := c.ListAccountProxyConfigsPerPage(env, nil, nil, pageNum, perPage)
		if err != nil {
			subT.Fatal(err)
		}

		if pList == nil {
			subT.Fatal("config list returned nil")
		}

		if len(pList.ProxyConfigs) != 1 {
			subT.Fatalf("Then number of proxy configs does not match. Expected [%d]; got [%d]", 2, len(pList.ProxyConfigs))
		}
	})

	t.Run("page and per_page params not used", func(subT *testing.T) {
		httpClient := NewTestClient(func(req *http.Request) *http.Response {
			if req.URL.Path != fmt.Sprintf(accountProxyConfigGet, "production") {
				subT.Fatalf("Path does not match. Expected [%s]; got [%s]", fmt.Sprintf(accountProxyConfigGet, "production"), req.URL.Path)
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
				Body:       ioutil.NopCloser(bytes.NewReader(helperLoadBytes(subT, "account_proxy_fixture.json"))),
				Header:     make(http.Header),
			}
		})

		credential := "someAccessToken"
		c := NewThreeScale(NewTestAdminPortal(subT), credential, httpClient)
		pList, err := c.ListAccountProxyConfigsPerPage(env, nil, nil)
		if err != nil {
			subT.Fatal(err)
		}

		if pList == nil {
			subT.Fatal("configslist returned nil")
		}

		if len(pList.ProxyConfigs) != 1 {
			subT.Fatalf("Then number of proxy configs does not match. Expected [%d]; got [%d]", 1, len(pList.ProxyConfigs))
		}
	})
}
