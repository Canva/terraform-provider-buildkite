package client

import (
	"fmt"
	"net/http"
)

type buildkiteTransport struct {
	APIToken  string
	Transport http.RoundTripper
}

func NewAuthTransport(apiToken string, transport *http.RoundTripper) *buildkiteTransport {
	if transport == nil {
		transport = &http.DefaultTransport
	}

	return &buildkiteTransport{
		APIToken:  apiToken,
		Transport: *transport,
	}
}

func (t buildkiteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", userAgent)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.APIToken))

	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, &NotFound{}
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	return resp, nil
}

type NotFound struct {
}

func (err *NotFound) Error() string {
	return "404 Not Found"
}
