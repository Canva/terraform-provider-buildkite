package client

import (
	"fmt"

	"github.com/machinebox/graphql"
)

type Pipeline struct {
	Id                  string                 `json:"id,omitempty"`
	Environment         map[string]string      `json:"env"`
	Slug                string                 `json:"slug,omitempty"`
	WebURL              string                 `json:"web_url,omitempty"`
	BuildsURL           string                 `json:"builds_url,omitempty"`
	Url                 string                 `json:"url,omitempty"`
	DefaultBranch       string                 `json:"default_branch,omitempty"`
	BadgeURL            string                 `json:"badge_url,omitempty"`
	CreatedAt           string                 `json:"created_at,omitempty"`
	Repository          string                 `json:"repository,omitempty"`
	Name                string                 `json:"name,omitempty"`
	Description         string                 `json:"description"`
	BranchConfiguration string                 `json:"branch_configuration"`
	Provider            BuildkiteProvider      `json:"provider,omitempty"`
	ProviderSettings    map[string]interface{} `json:"provider_settings,omitempty"`
	Steps               []Step                 `json:"steps"`
}

type BuildkiteProvider struct {
	Id         string                 `json:"id"`
	Settings   map[string]interface{} `json:"settings"`
	WebhookURL string                 `json:"webhook_url"`
}

type Step struct {
	Type                string            `json:"type"`
	Name                string            `json:"name,omitempty"`
	Command             string            `json:"command,omitempty"`
	Environment         map[string]string `json:"env"`
	TimeoutInMinutes    int               `json:"timeout_in_minutes,omitempty"`
	AgentQueryRules     []string          `json:"agent_query_rules"`
	BranchConfiguration string            `json:"branch_configuration"`
	ArtifactPaths       string            `json:"artifact_paths"`
	Concurrency         int               `json:"concurrency,omitempty"`
	Parallelism         int               `json:"parallelism,omitempty"`
}

type pipelineIdResponse struct {
	Pipeline Node `json:"pipeline"`
}

func (c *Client) GetPipeline(slug string) (*Pipeline, error) {
	pipeline := Pipeline{}
	relativePath := fmt.Sprintf("/v2/organizations/%s/pipelines/%s", c.orgSlug, slug)
	err := c.get(relativePath, &pipeline)
	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (c *Client) CreatePipeline(pipeline *Pipeline) (*Pipeline, error) {
	result := Pipeline{}
	relativePath := fmt.Sprintf("/v2/organizations/%s/pipelines", c.orgSlug)
	err := c.post(relativePath, pipeline, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) UpdatePipeline(pipeline *Pipeline) (*Pipeline, error) {
	result := Pipeline{}
	relativePath := fmt.Sprintf("/v2/organizations/%s/pipelines/%s", c.orgSlug, pipeline.Slug)
	err := c.patch(relativePath, pipeline, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) DeletePipeline(slug string) error {
	relativePath := fmt.Sprintf("/v2/organizations/%s/pipelines/%s", c.orgSlug, slug)
	err := c.delete(relativePath, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetPipelineNodeId(slug string) (string, error) {
	req := graphql.NewRequest(`
query GetPipelineId($pipelineSlug: ID!) {
  pipeline(slug: $pipelineSlug) {
    id
  }
}`)
	req.Var("pipelineSlug", c.createOrgSlug(slug))

	idResponse := pipelineIdResponse{}
	if err := c.graphQLRequest(req, &idResponse); err != nil {
		return "", err
	}

	return idResponse.Pipeline.Id, nil
}
