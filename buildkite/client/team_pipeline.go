package client

import (
	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
	"log"
)

const (
	TeamPipelineAccessReadOnly     = "READ_ONLY"
	TeamPipelineAccessBuildAndRead = "BUILD_AND_READ"
	TeamPipelineAccessManage       = "MANAGE_BUILD_AND_READ"
)

type teamPipelineResponse struct {
	TeamPipeline TeamPipeline `json:"teamPipeline"`
}

type TeamPipeline struct {
	Id          string `json:"id,omitempty"`
	UUID        string `json:"uuid,omitempty"`
	AccessLevel string `json:"accessLevel,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	Team        Node   `json:"team,omitempty"`
	Pipeline    Node   `json:"pipeline,omitempty"`
}

type teamPipelineCreateResponse struct {
	TeamPipelineCreate struct {
		TeamPipelineEdge struct {
			Node TeamPipeline
		}
	}
}

type teamPipelineUpdateResponse struct {
	TeamPipelineUpdate struct {
		TeamPipeline TeamPipeline
	}
}

type teamPipelineDeleteResponse struct {
	DeletedTeamPipelineID string `json:"deletedTeamPipelineID"`
}

func (c *Client) GetTeamPipeline(teamPipelineId string) (*TeamPipeline, error) {
	log.Printf("[TRACE] Buildkite client GetTeamPipeline %s", teamPipelineId)

	req := graphql.NewRequest(`
query GetTeamPipeline($teamPipelineId: ID!) {
  teamPipeline: node(id: $teamPipelineId) {
    ... on TeamPipeline {
      id
      uuid
      accessLevel
      createdAt
      pipeline {
        id
		slug
      }
      team {
        id
      }
    }
  }
}
`)
	req.Var("teamPipelineId", teamPipelineId)

	response := teamPipelineResponse{}
	if err := c.graphQLRequest(req, &response); err != nil {
		return nil, errors.Wrapf(err, "failed to get team member %s", teamPipelineId)
	}

	return &response.TeamPipeline, nil
}

func (c *Client) CreateTeamPipeline(teamPipeline *TeamPipeline) (*TeamPipeline, error) {
	log.Printf("[TRACE] Buildkite client CreateTeamPipeline %s", teamPipeline.UUID)

	pipelineId, err := c.GetPipelineNodeId(teamPipeline.Pipeline.Slug)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get pipeline id for slug %s", teamPipeline.Pipeline.Slug)
	}

	req := graphql.NewRequest(`
mutation TeamPipelineNewMutation($teamPipelineCreateInput: TeamPipelineCreateInput!) {
  teamPipelineCreate(input: $teamPipelineCreateInput) {
    teamPipelineEdge {
      node {
        id
        uuid
        accessLevel
        createdAt
        pipeline {
          id
          slug
        }
        team {
          id
        }
      }
    }
  }
}
`)

	req.Var("teamPipelineCreateInput", map[string]interface{}{
		"teamID":     teamPipeline.Team.Id,
		"pipelineID": pipelineId,
	})

	response := teamPipelineCreateResponse{}
	if err := c.graphQLRequest(req, &response); err != nil {
		return nil, errors.Wrapf(err, "failed to create team %s", teamPipeline.Team.Id)
	}

	return &response.TeamPipelineCreate.TeamPipelineEdge.Node, nil
}

func (c *Client) UpdateTeamPipeline(teamPipeline *TeamPipeline) (*TeamPipeline, error) {

	req := graphql.NewRequest(`
mutation TeamPipelineUpdateMutation($teamPipelineUpdateInput: TeamPipelineUpdateInput!) {
  teamPipelineUpdate(input: $teamPipelineUpdateInput) {
    teamPipeline {
      id
      uuid
      accessLevel
      createdAt
      pipeline {
        id
        slug
      }
      team {
        id
      }
    }
  }
}
`)

	req.Var("teamPipelineUpdateInput", map[string]interface{}{
		"id":          teamPipeline.Id,
		"accessLevel": teamPipeline.AccessLevel,
	})

	response := teamPipelineUpdateResponse{}
	if err := c.graphQLRequest(req, &response); err != nil {
		return nil, errors.Wrapf(err, "failed to update teamPipeline %s", teamPipeline.Id)
	}

	return &response.TeamPipelineUpdate.TeamPipeline, nil
}

func (c *Client) DeleteTeamPipeline(teamPipelineId string) error {
	req := graphql.NewRequest(`
mutation TeamPipelineDeleteMutation($teamPipelineDeleteInput: TeamPipelineDeleteInput!) {
  teamPipelineDelete(input: $teamPipelineDeleteInput) {
    deletedTeamPipelineID
  }
}
`)

	req.Var("teamPipelineDeleteInput", map[string]interface{}{
		"id": teamPipelineId,
	})

	response := teamPipelineDeleteResponse{}
	if err := c.graphQLRequest(req, &response); err != nil {
		return errors.Wrapf(err, "failed to delete team member %s", teamPipelineId)
	}

	return nil
}
