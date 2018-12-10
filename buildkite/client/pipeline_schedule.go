package client

import (
	"log"
	"strings"

	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
)

type pipelineScheduleResponse struct {
	PipelineSchedule PipelineSchedule `json:"pipelineSchedule"`
}

type PipelineSchedule struct {
	Id           string   `json:"id,omitempty"`
	UUID         string   `json:"uuid,omitempty"`
	Pipeline     Node     `json:"pipeline,omitempty"`
	CreatedAt    string   `json:"createdAt,omitempty"`
	Label        string   `json:"label,omitempty"`
	CronSchedule string   `json:"cronline,omitempty"`
	Message      string   `json:"message,omitempty"`
	Commit       string   `json:"commit,omitempty"`
	Branch       string   `json:"Branch,omitempty"`
	Environment  []string `json:"env,omitempty"`
	Enabled      bool     `json:"enabled"`
}

type pipelineScheduleCreateResponse struct {
	PipelineScheduleCreate struct {
		PipelineScheduleEdge struct {
			Node PipelineSchedule
		}
	}
}

type pipelineScheduleUpdateResponse struct {
	PipelineScheduleUpdate struct {
		PipelineSchedule PipelineSchedule
	}
}

type pipelineScheduleDeleteResponse struct {
	DeletedPipelineScheduleID string `json:"deletedPipelineScheduleID"`
}

func (c *Client) GetPipelineSchedule(pipelineScheduleSlug string) (*PipelineSchedule, error) {
	log.Printf("[TRACE] Buildkite client GetPipelineSchedule %s", pipelineScheduleSlug)

	req := graphql.NewRequest(`
query GetPipelineSchedule($pipelineScheduleSlug: ID!) {
  pipelineSchedule(slug: $pipelineScheduleSlug) {
    id
    uuid
    label
    cronline
    message
    commit
    branch
    env
    enabled
    createdAt
    pipeline {
      id
      slug
    }
  }
}
`)
	req.Var("pipelineScheduleSlug", c.createOrgSlug(pipelineScheduleSlug))

	response := pipelineScheduleResponse{}
	if err := c.graphQLRequest(req, &response); err != nil {
		return nil, errors.Wrapf(err, "failed to get pipeline schedule %s", pipelineScheduleSlug)
	}

	return &response.PipelineSchedule, nil
}

func (c *Client) CreatePipelineSchedule(pipelineSchedule *PipelineSchedule) (*PipelineSchedule, error) {
	log.Printf("[TRACE] Buildkite client CreatePipelineSchedule %s", pipelineSchedule.UUID)

	pipelineId, err := c.GetPipelineNodeId(pipelineSchedule.Pipeline.Slug)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get schedule id for slug %s", pipelineSchedule.Pipeline.Slug)
	}

	req := graphql.NewRequest(`
mutation PipelineScheduleNewMutation($pipelineScheduleCreateInput: PipelineScheduleCreateInput!) {
  pipelineScheduleCreate(input: $pipelineScheduleCreateInput) {
    pipelineScheduleEdge {
      node {
        id
        uuid
        label
        cronline
        message
        commit
        branch
        env
        enabled
        createdAt
        pipeline {
          id
          slug
        }
      }
    }
  }
}
`)

	req.Var("pipelineScheduleCreateInput", map[string]interface{}{
		"pipelineID": pipelineId,
		"label":      pipelineSchedule.Label,
		"cronline":   pipelineSchedule.CronSchedule,
		"message":    pipelineSchedule.Message,
		"commit":     pipelineSchedule.Commit,
		"branch":     pipelineSchedule.Branch,
		"env":        listToString(pipelineSchedule.Environment),
		"enabled":    pipelineSchedule.Enabled,
	})

	response := pipelineScheduleCreateResponse{}
	if err := c.graphQLRequest(req, &response); err != nil {
		return nil, errors.Wrapf(err, "failed to create pipeline schedule for pipeline %s", pipelineSchedule.Pipeline.Slug)
	}

	return &response.PipelineScheduleCreate.PipelineScheduleEdge.Node, nil
}

func (c *Client) UpdatePipelineSchedule(pipelineSchedule *PipelineSchedule) (*PipelineSchedule, error) {

	req := graphql.NewRequest(`
mutation PipelineScheduleUpdateMutation($pipelineScheduleUpdateInput: PipelineScheduleUpdateInput!) {
  pipelineScheduleUpdate(input: $pipelineScheduleUpdateInput) {
    pipelineSchedule {
      id
      uuid
      label
      cronline
      message
      commit
      branch
      env
      enabled
      createdAt
      pipeline {
        id
        slug
      }
    }
  }
}
`)

	req.Var("pipelineScheduleUpdateInput", map[string]interface{}{
		"id":       pipelineSchedule.Id,
		"label":    pipelineSchedule.Label,
		"cronline": pipelineSchedule.CronSchedule,
		"message":  pipelineSchedule.Message,
		"commit":   pipelineSchedule.Commit,
		"branch":   pipelineSchedule.Branch,
		"env":      listToString(pipelineSchedule.Environment),
		"enabled":  pipelineSchedule.Enabled,
	})

	response := pipelineScheduleUpdateResponse{}
	if err := c.graphQLRequest(req, &response); err != nil {
		return nil, errors.Wrapf(err, "failed to update pipeline schedule %s", pipelineSchedule.Id)
	}

	return &response.PipelineScheduleUpdate.PipelineSchedule, nil
}

func (c *Client) DeletePipelineSchedule(pipelineScheduleId string) error {
	req := graphql.NewRequest(`
mutation PipelineScheduleDeleteMutation($pipelineScheduleDeleteInput: PipelineScheduleDeleteInput!) {
  pipelineScheduleDelete(input: $pipelineScheduleDeleteInput) {
    deletedPipelineScheduleID
  }
}
`)

	req.Var("pipelineScheduleDeleteInput", map[string]interface{}{
		"id": pipelineScheduleId,
	})

	response := pipelineScheduleDeleteResponse{}
	if err := c.graphQLRequest(req, &response); err != nil {
		return errors.Wrapf(err, "failed to delete pipeline schedule %s", pipelineScheduleId)
	}

	return nil
}

func listToString(list []string) string {
	return strings.Join(list, "\n")
}
