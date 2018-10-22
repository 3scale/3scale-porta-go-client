package fake

import (
        "bytes"
        "io/ioutil"
        "net/http"
)

func GetProxyConfigLatestSuccess() *http.Response {
        return &http.Response{
                StatusCode: 200,
                Body:       ioutil.NopCloser(bytes.NewBufferString(GetProxyConfigLatestJson())),
                Header:     make(http.Header),
        }
}
