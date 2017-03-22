package share

import (
	"io"
	"net/http"
)

// mockHTTPPost return a mock instead of executing the given POST request.
func mockHTTPPost(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	return resp, err
}

// mockHTTPGet return a mock instead of executing the given GET request.
func mockHTTPGet(url string) (resp *http.Response, err error) {
	return resp, err
}
