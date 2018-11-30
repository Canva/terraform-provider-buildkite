package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/machinebox/graphql"
	"github.com/pkg/errors"

	"github.com/saymedia/terraform-buildkite/buildkite/version"
)

const (
	defaultBaseURL             = "https://api.buildkite.com/"
	defaultGraphQLUrl          = "https://graphql.buildkite.com/v1"
	userAgent                  = "Terraform-Buildkite/" + version.Version
	applicationJsonContentType = "application/json"
)

type Client struct {
	client   *http.Client
	graphQl  *graphql.Client
	baseURL  *url.URL
	orgSlug  string
	apiToken string
}

func NewClient(orgSlug string, apiToken string) *Client {
	var authTransport http.RoundTripper = NewAuthTransport(apiToken, nil)
	baseURL, _ := url.Parse(defaultBaseURL)

	return &Client{
		client: &http.Client{
			Transport: authTransport,
		},
		graphQl: graphql.NewClient(defaultGraphQLUrl, graphql.WithHTTPClient(&http.Client{
			Transport: authTransport,
		})),
		baseURL:  baseURL,
		orgSlug:  orgSlug,
		apiToken: apiToken,
	}
}

func (c *Client) graphQLRequest(req *graphql.Request, result interface{}) error {
	ctx := context.Background()
	return c.graphQl.Run(ctx, req, &result)
}

func (c *Client) createOrgSlug(slug string) string {
	return fmt.Sprintf("%s/%s", c.orgSlug, slug)
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

func unmarshalResponse(body io.Reader, result interface{}) error {
	responseBytes, err := ioutil.ReadAll(body)
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
