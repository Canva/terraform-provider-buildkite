package client

import (
	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
	"log"
)

type teamMemberResponse struct {
	TeamMember TeamMember `json:"teamMember"`
}

type TeamMember struct {
	Id        string `json:"id,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	Role      string `json:"role,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	Team      Node   `json:"team,omitempty"`
	User      Node   `json:"user,omitempty"`
}

type teamMemberCreateResponse struct {
	TeamMemberCreate struct {
		TeamMemberEdge struct {
			Node TeamMember
		}
	}
}

type teamMemberUpdateResponse struct {
	TeamMemberUpdate struct {
		TeamMember TeamMember
	}
}

type teamMemberDeleteResponse struct {
	DeletedTeamMemberID string `json:"deletedTeamMemberID"`
}

func (c *Client) GetTeamMember(teamMemberId string) (*TeamMember, error) {
	log.Printf("[TRACE] Buildkite client GetTeamMember %s", teamMemberId)

	req := graphql.NewRequest(`
query GetTeamMember($teamMemberId: ID!) {
  teamMember: node(id: $teamMemberId) {
    ... on TeamMember {
      id
      uuid
      role
      createdAt
      user {
        id
      }
      team {
        id
      }
    }
  }
}
`)
	req.Var("teamMemberId", teamMemberId)

	response := teamMemberResponse{}
	if err := c.graphQLRequest(req, &response); err != nil {
		return nil, errors.Wrapf(err, "failed to get team member %s", teamMemberId)
	}

	return &response.TeamMember, nil
}

func (c *Client) CreateTeamMember(teamMember *TeamMember) (*TeamMember, error) {
	log.Printf("[TRACE] Buildkite client CreateTeamMember %s", teamMember.UUID)

	req := graphql.NewRequest(`
mutation TeamMemberNewMutation($teamMemberCreateInput: TeamMemberCreateInput!) {
  teamMemberCreate(input: $teamMemberCreateInput) {
    teamMemberEdge {
      node {
        id
        uuid
        role
        createdAt
        team {
          id
        }
        user {
          id
        }
      }
    }
  }
}
`)

	req.Var("teamMemberCreateInput", map[string]interface{}{
		"teamID": teamMember.Team.Id,
		"userID": teamMember.User.Id,
	})

	response := teamMemberCreateResponse{}
	if err := c.graphQLRequest(req, &response); err != nil {
		return nil, errors.Wrapf(err, "failed to create team %s", teamMember.Team.Id)
	}

	return &response.TeamMemberCreate.TeamMemberEdge.Node, nil
}

func (c *Client) UpdateTeamMember(teamMember *TeamMember) (*TeamMember, error) {

	req := graphql.NewRequest(`
mutation TeamMemberUpdateMutation($teamMemberUpdateInput: TeamMemberUpdateInput!) {
  teamMemberUpdate(input: $teamMemberUpdateInput) {
    teamMember {
      id
      uuid
      role
      createdAt
      team {
        id
      }
      user {
        id
      }
    }
  }
}
`)

	req.Var("teamMemberUpdateInput", map[string]interface{}{
		"id":   teamMember.Id,
		"role": teamMember.Role,
	})

	response := teamMemberUpdateResponse{}
	if err := c.graphQLRequest(req, &response); err != nil {
		return nil, errors.Wrapf(err, "failed to update teamMember %s", teamMember.Id)
	}

	return &response.TeamMemberUpdate.TeamMember, nil
}

func (c *Client) DeleteTeamMember(teamMemberId string) error {
	req := graphql.NewRequest(`
mutation TeamMemberDeleteMutation($teamMemberDeleteInput: TeamMemberDeleteInput!) {
  teamMemberDelete(input: $teamMemberDeleteInput) {
    deletedTeamMemberID
  }
}
`)

	req.Var("teamMemberDeleteInput", map[string]interface{}{
		"id": teamMemberId,
	})

	response := teamMemberDeleteResponse{}
	if err := c.graphQLRequest(req, &response); err != nil {
		return errors.Wrapf(err, "failed to delete team member %s", teamMemberId)
	}

	return nil
}
