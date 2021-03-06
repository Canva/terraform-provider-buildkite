package client

import (
	"fmt"
	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
	"log"
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
	// Buildkite doesn't allow you to create a pipeline if you not an admin or if you a member of more that one team or
	// none of them. So you unable to create a pipeline and attach the "buildkite_team_pipeline" resource to it after it was
	// created in this case.
	TeamIDs []string `json:"team_ids,omitempty"`
	Steps   []Step   `json:"steps,omitempty"`

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
		pipeline.Environment = nil
	}

	pipeline.TeamIDs, err = c.getTeamIDs(slug)
	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (c *Client) CreatePipeline(pipeline *Pipeline) (*Pipeline, error) {
	// Create via the GraphQL API if the YAML based configuration is used
	if len(pipeline.Configuration) > 0 {
		return c.createPipelineGraphQl(pipeline)
	}

	result := Pipeline{}
	relativePath := fmt.Sprintf("/v2/organizations/%s/pipelines", c.orgSlug)
	err := c.post(relativePath, pipeline, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// createPipelineGraphQl will create the pipeline but only set the required fields
// after creation, UpdatePipeline will be used to set the rest of the fields via
// the REST API
func (c *Client) createPipelineGraphQl(pipeline *Pipeline) (*Pipeline, error) {
	req := graphql.NewRequest(`
mutation PipelineCreateRequest($pipelineCreateInput: PipelineCreateInput!) {
  pipelineCreate(input: $pipelineCreateInput) {
    pipeline {
      slug
    }
  }
}`)

	orgID, err := c.GetOrganizationId(c.orgSlug)
	if err != nil {
		return nil, err
	}

	pci := map[string]interface{}{
		"organizationId": orgID,
		"name":           pipeline.Name,
		"repository": map[string]string{
			"url": pipeline.Repository,
		},
		"steps": map[string]string{
			"yaml": pipeline.Configuration,
		},
	}

	if len(pipeline.TeamIDs) != 0 {
		var teamIDs []map[string]string
		// Converting a slice of team UUIDs into the slice of maps since GraphQL API expects this data in this shape.
		// We grant "MANAGE_BUILD_AND_READ" access level to _initial_ teams-owners.
		for _, ID := range pipeline.TeamIDs {
			teamIDs = append(teamIDs, map[string]string{
				"id":          ID,
				"accessLevel": TeamPipelineAccessManage,
			})
		}
		pci["teams"] = teamIDs
	}

	req.Var("pipelineCreateInput", pci)

	var createPipelineResponse struct {
		PipelineCreate struct {
			Pipeline struct {
				Slug string `json:"slug"`
			} `json:"pipeline"`
		} `json:"pipelineCreate"`
	}

	if err := c.graphQLRequest(req, &createPipelineResponse); err != nil {
		return nil, errors.Wrapf(err, "failed to create pipeline %s", pipeline.Slug)
	}

	pipeline.Slug = createPipelineResponse.PipelineCreate.Pipeline.Slug

	// set all other options with the rest api
	return c.UpdatePipeline(pipeline)
}

func (c *Client) UpdatePipeline(pipeline *Pipeline) (*Pipeline, error) {
	// Save other parameters via the REST API
	result := Pipeline{TeamIDs: pipeline.TeamIDs} // Save TeamIDs as long as REST API doesn't provide them in response
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

func (c *Client) getTeamIDs(slug string) ([]string, error) {
	req := graphql.NewRequest(`
query Pipeline($slug: ID!) {
  pipeline(slug: $slug) {
    teams(first: 100) {
      edges {
        node {
          team {
            id
          }
        }
      }
    }
  }
}`)

	req.Var("slug", c.createOrgSlug(slug))
	var resp struct {
		Pipeline struct {
			Teams struct {
				Edges []struct {
					Node struct {
						Team struct {
							ID string `json:"id"`
						} `json:"team"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"teams"`
		} `json:"pipeline"`
	}
	if err := c.graphQLRequest(req, &resp); err != nil {
		return nil, err
	}

	teamIDs := make([]string, len(resp.Pipeline.Teams.Edges))
	for i, edge := range resp.Pipeline.Teams.Edges {
		teamIDs[i] = edge.Node.Team.ID
	}
	log.Printf("[TRACE] got team ids: %v", teamIDs)
	return teamIDs, nil
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
