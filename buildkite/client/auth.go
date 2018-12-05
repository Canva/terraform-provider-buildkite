package client

import (
	"fmt"
	"net/http"
)

type authTransport struct {
	APIToken  string
	Transport http.RoundTripper
	userAgent string
}

func NewAuthTransport(apiToken string, userAgent string, transport *http.RoundTripper) *authTransport {
	if transport == nil {
		transport = &http.DefaultTransport
	}

	return &authTransport{
		APIToken:  apiToken,
		Transport: *transport,
		userAgent: userAgent,
	}
}

func (t authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", t.userAgent)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.APIToken))

	return t.Transport.RoundTrip(req)
}
