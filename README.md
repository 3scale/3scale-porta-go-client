# 3scale-porta-go-client
[![CircleCI](https://circleci.com/gh/3scale/3scale-porta-go-client.svg?style=svg)](https://circleci.com/gh/3scale/3scale-porta-go-client)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)
[![GitHub release](https://img.shields.io/github/v/release/3scale/3scale-porta-go-client.svg)](https://github.com/3scale/3scale-porta-go-client/releases/latest)

3scale Account Management API Client

## Installation

To install, run:

```bash
$ go get github.com/3scale/3scale-porta-go-client/client
```

And import using:

```go
import "github.com/3scale/3scale-porta-go-client/client"
```

## Usage

Basic usage

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/3scale/3scale-porta-go-client/client"
)

func jsonPrint(obj interface{}) {
	jsonData, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonData))
}

func main() {
	adminPortalURL := "https://3scale-admin.example.com"
	threescaleAccessToken := "************************"

	adminPortal, err := client.NewAdminPortalFromStr(adminPortalURL)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	threescaleClient := client.NewThreeScale(adminPortal, threescaleAccessToken, nil)

	backendList, err := threescaleClient.ListBackendApis()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	jsonPrint(backendList)
}
```

The `client.NewThreeScale` constructor allows injecting you own customized [\*http.Client](https://golang.org/src/net/http/client.go).

In this example, the customized transport will log all http requests and responses

```go
type DebuggerTransport struct {
	Transport http.RoundTripper
}

// RoundTrip is the core part of this module and implements http.RoundTripper.
// Executes HTTP request with request/response logging.
func (t *DebuggerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.logRequest(req)

	resp, err := t.transport().RoundTrip(req)
	if err != nil {
		return resp, err
	}

	t.logResponse(resp)

	return resp, err
}

func (t *DebuggerTransport) logRequest(req *http.Request) {
	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Println(string(dump))
}

func (t *DebuggerTransport) logResponse(resp *http.Response) {
	dump, _ := httputil.DumpResponse(resp, true)
	fmt.Println(string(dump))
}

func (t *DebuggerTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}

	return http.DefaultTransport
}

var transport http.RoundTripper = &DebuggerTransport{}

threescaleClient := client.NewThreeScale(adminPortal, threescaleAccessToken, &http.Client{Transport: transport})
```

In this example, the customized transport allow insecure TLS connections (**warning**: this use is not recommended)

```go
var transport http.RoundTripper = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

threescaleClient := client.NewThreeScale(adminPortal, threescaleAccessToken, &http.Client{Transport: transport})
```

## Development

### Testing

```
$ make test
```

Optionally, add `TEST_NAME` makefile variable to run specific test

```sh
make test TEST_NAME=TestActivateUserErrors
```

or even subtest

```sh
make test TEST_NAME=TestActivateUserErrors/UnexpectedHTTPStatusCode
```

## Contributing

Bug reports and pull requests are welcome on [GitHub](https://github.com/3scale/3scale-porta-go-client)
