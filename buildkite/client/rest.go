package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type NotFound struct {
}

func (err *NotFound) Error() string {
	return "404 Not Found"
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

	if resp.StatusCode == http.StatusNotFound {
		return &NotFound{}
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("%s\nResponse body:\n\n%s\n", resp.Status, body)
	}

	if responseBody != nil {
		if err = unmarshalResponse(resp.Body, &responseBody); err != nil {
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
