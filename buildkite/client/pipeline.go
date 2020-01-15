package client

import (
	"fmt"
	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
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
	TeamUUIDs           []string               `json:"team_uuids,omitempty"`
	Steps               []Step                 `json:"steps,omitempty"`

	// Configuration is the "new" YAML based pipeline setup
	// This value can only be set via the GraphQL API
	Configuration string `json:"configuration,omitempty"`
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

	// If the yaml configuration is used, both Configuration and Steps will be set
	// ignore the value of Steps
	if len(pipeline.Configuration) > 0 {
		pipeline.Steps = nil
	}

	return &pipeline, nil
}

func (c *Client) CreatePipeline(pipeline *Pipeline) (*Pipeline, error) {
	// Create via the REST API if the YAML based configuration is used
	if len(pipeline.Configuration) > 0 {
		return c.createYAMLPipeline(pipeline)
	}

	result := Pipeline{}
	relativePath := fmt.Sprintf("/v2/organizations/%s/pipelines", c.orgSlug)
	err := c.post(relativePath, pipeline, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// createYAMLPipeline will create the pipeline but only set the required fields
// filled with stubs. After creation, UpdatePipeline will be used to set the rest
// of the fields via the REST API
func (c *Client) createYAMLPipeline(pipeline *Pipeline) (*Pipeline, error) {
	created, err := c.CreatePipeline(&Pipeline{
		Name:       pipeline.Name,
		Repository: pipeline.Repository,
		TeamUUIDs:  pipeline.TeamUUIDs,
		Steps: []Step{
			{
				Type:    "script",
				Name:    "Script",
				Command: "command.sh",
			},
		},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create an empty pipeline %s", pipeline.Name)
	}
	pipeline.Slug = created.Slug
	return c.UpdatePipeline(pipeline)
}

func (c *Client) UpdatePipeline(pipeline *Pipeline) (*Pipeline, error) {
	// Save other parameters via the REST API
	result := Pipeline{}
	relativePath := fmt.Sprintf("/v2/organizations/%s/pipelines/%s", c.orgSlug, pipeline.Slug)
	err := c.patch(relativePath, pipeline, &result)
	if err != nil {
		return nil, err
	}

	// Set YAML steps via the GraphQL API
	if len(pipeline.Configuration) > 0 {
		err := c.savePipelineYaml(pipeline)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (c *Client) savePipelineYaml(pipeline *Pipeline) error {
	req := graphql.NewRequest(`
mutation PipelineUpdateMutation($pipelineUpdateInput: PipelineUpdateInput!) {
  pipelineUpdate(input: $pipelineUpdateInput) {
    pipeline {
      steps {
        yaml
      }
    }
  }
}`)

	nodeID, err := c.GetPipelineNodeId(pipeline.Slug)
	if err != nil {
		return errors.Wrapf(err, "failed to get GraphQL node id for %s", pipeline.Slug)
	}

	req.Var("pipelineUpdateInput", map[string]interface{}{
		"id": nodeID,
		"steps": map[string]interface{}{
			"yaml": pipeline.Configuration,
		},
	})

	var gres interface{}
	if err := c.graphQLRequest(req, &gres); err != nil {
		return errors.Wrapf(err, "failed to update pipeline %s", pipeline.Slug)
	}

	return nil
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
