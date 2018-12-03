package client

import (
	"github.com/machinebox/graphql"
	"sync"
)

type orgIdResponse struct {
	Organization Node `json:"organization"`
}

type Node struct {
	Id string `json:"id,omitempty"`
}

var (
	orgIds   = map[string]string{}
	orgMutex = &sync.Mutex{}
)

func (c *Client) GetOrganizationId(slug string) (string, error) {
	orgMutex.Lock()
	defer orgMutex.Unlock()

	if val, ok := orgIds[slug]; ok {
		return val, nil
	}

	val, err := c.fetchOrganizationId(slug)
	if err != nil {
		return "", err
	}

	orgIds[slug] = val
	return val, nil
}

func (c *Client) fetchOrganizationId(slug string) (string, error) {
	req := graphql.NewRequest(`
query Organization($orgSlug: ID!) {
  organization(slug: $orgSlug) {
    id
  }
}`)
	req.Var("orgSlug", c.orgSlug)

	idResponse := orgIdResponse{}
	if err := c.graphQLRequest(req, &idResponse); err != nil {
		return "", err
	}

	return idResponse.Organization.Id, nil
}
