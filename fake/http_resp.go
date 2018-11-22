package fake

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func CreateAppSuccess(description string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusCreated,
		Body:       ioutil.NopCloser(bytes.NewBufferString(CreateApp(description))),
		Header:     make(http.Header),
	}
}

func CreateAppError() *http.Response {
	return &http.Response{
		StatusCode: http.StatusForbidden,
		Body:       ioutil.NopCloser(bytes.NewBufferString(CreateAppFail())),
		Header:     make(http.Header),
	}
}

func GetProxyConfigLatestSuccess() *http.Response {
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(GetProxyConfigLatestJson())),
		Header:     make(http.Header),
	}
}
