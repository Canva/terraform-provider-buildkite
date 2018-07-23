package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/saymedia/terraform-buildkite/buildkite/version"
)

const (
	defaultBaseURL             = "https://api.buildkite.com/"
	userAgent                  = "Terraform-Buildkite/" + version.Version
	applicationJsonContentType = "application/json"
)

type Client struct {
	client   *http.Client
	baseURL  *url.URL
	orgSlug  string
	apiToken string
}

func NewClient(orgSlug string, apiToken string) *Client {
	var transport http.RoundTripper = NewAuthTransport(apiToken, nil)
	baseURL, _ := url.Parse(defaultBaseURL)

	return &Client{
		client: &http.Client{
			Transport: transport,
		},
		baseURL:  baseURL,
		orgSlug:  orgSlug,
		apiToken: apiToken,
	}
}

func (c *Client) get(relativePath string, responseBody interface{}) error {
	return c.request("GET", relativePath, nil, responseBody)
}

func (c *Client) post(relativePath string, requestBody interface{}, responseBody interface{}) error {
	return c.request("POST", relativePath, requestBody, responseBody)
}

func (c *Client) patch(relativePath string, requestBody interface{}, responseBody interface{}) error {
	return c.request("PATCH", relativePath, requestBody, responseBody)
}

func (c *Client) delete(relativePath string, responseBody interface{}) error {
	return c.request("DELETE", relativePath, nil, responseBody)
}

func (c *Client) request(method string, relativePath string, requestBody interface{}, responseBody interface{}) error {
	log.Printf("[DEBUG] Buildkite Request %s %s\n", method, relativePath)

	req, err := createRequest(method, c.urlPath(relativePath), requestBody)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if responseBody != nil {
		if err = unmarshalResponse(resp, &responseBody); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) urlPath(relativePath string) string {
	return c.baseURL.ResolveReference(&url.URL{
		Path: relativePath,
	}).String()
}

func createRequest(method string, url string, requestBody interface{}) (*http.Request, error) {
	if requestBody == nil {
		return http.NewRequest(method, url, nil)
	}

	body, err := marshalBody(requestBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", applicationJsonContentType)
	return req, nil
}

func marshalBody(body interface{}) (*bytes.Buffer, error) {
	if body == nil {
		return nil, nil
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal body")
	}

	return bytes.NewBuffer(bodyBytes), nil
}

func unmarshalResponse(resp *http.Response, result interface{}) error {
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "could not read response body")
	}
	log.Printf("[TRACE] Buildkite Response body %s\n", string(responseBytes))

	err = json.Unmarshal(responseBytes, &result)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal response body")
	}

	return nil
}
