package oauth_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type httpInterceptor struct {
	Transport http.RoundTripper
	Response  *http.Response
}

func (i httpInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	return i.Response, nil
}

func intercept(statusCode int, body json.RawMessage) *http.Client {
	jsonPayload, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	res := &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBuffer(jsonPayload)),
	}

	return &http.Client{
		Transport: httpInterceptor{
			Response: res,
		},
	}
}
