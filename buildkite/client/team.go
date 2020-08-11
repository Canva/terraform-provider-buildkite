package client

import (
	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
)

const (
	TeamPrivacyVisible = "VISIBLE"
	TeamPrivacySecret  = "SECRET"

	TeamMemberRoleMember     = "MEMBER"
	TeamMemberRoleMaintainer = "MAINTAINER"
)

type teamResponse struct {
	Team Team `json:"team"`
}

type Team struct {
	Id                          string                    `json:"id,omitempty"`
	UUID                        string                    `json:"uuid,omitempty"`
	Slug                        string                    `json:"slug,omitempty"`
	Name                        string                    `json:"name,omitempty"`
	Description                 string                    `json:"description,omitempty"`
	Privacy                     string                    `json:"privacy,omitempty"`
	IsDefaultTeam               bool                      `json:"isDefaultTeam,omitempty"`
	DefaultMemberRole           string                    `json:"defaultMemberRole,omitempty"`
	CreatedAt                   string                    `json:"createdAt,omitempty"`
	MembersCanCreatePipelines   bool                      `json:"membersCanCreatePipelines,omitempty"`
}

type teamCreateResponse struct {
	TeamCreate struct {
		TeamEdge struct {
			Node Team
		}
	}
}

type teamUpdateResponse struct {
	TeamUpdate struct {
		Team Team
	}
}

type teamDeleteResponse struct {
	DeletedTeamID string `json:"deletedTeamID"`
}

func (c *Client) GetTeam(slug string) (*Team, error) {
	req := graphql.NewRequest(`
query GetTeam($teamSlug: ID!) {
  team(slug: $teamSlug) {
    id
    uuid
    slug
    name
    description
    createdAt
    privacy
    isDefaultTeam
    defaultMemberRole
    membersCanCreatePipelines
  }
}`)
	req.Var("teamSlug", c.createOrgSlug(slug))

	teamResponse := teamResponse{}
	if err := c.graphQLRequest(req, &teamResponse); err != nil {
		return nil, errors.Wrapf(err, "failed to get team %s", slug)
	}

	return &teamResponse.Team, nil
}

func (c *Client) CreateTeam(team *Team) (*Team, error) {

	orgId, err := c.GetOrganizationId(c.orgSlug)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch organization id")
	}

	req := graphql.NewRequest(`
mutation TeamNewMutation($teamCreateInput: TeamCreateInput!) {
  teamCreate(input: $teamCreateInput) {
    teamEdge {
      node {
        id
        uuid
        slug
        name
        description
        createdAt
        privacy
        isDefaultTeam
        defaultMemberRole
        membersCanCreatePipelines
      }
    }
  }
}
`)

	req.Var("teamCreateInput", map[string]interface{}{
		"organizationID":               orgId,
		"name":                         team.Name,
		"description":                  team.Description,
		"isDefaultTeam":                team.IsDefaultTeam,
		"defaultMemberRole":            team.DefaultMemberRole,
		"membersCanCreatePipelines":    team.MembersCanCreatePipelines,
		"privacy":                      team.Privacy,
	})

	teamCreateResponse := teamCreateResponse{}
	if err := c.graphQLRequest(req, &teamCreateResponse); err != nil {
		return nil, errors.Wrapf(err, "failed to create team %s", team.Name)
	}

	return &teamCreateResponse.TeamCreate.TeamEdge.Node, nil
}

func (c *Client) UpdateTeam(team *Team) (*Team, error) {

	req := graphql.NewRequest(`
mutation TeamUpdateMutation($teamUpdateInput: TeamUpdateInput!) {
  teamUpdate(input: $teamUpdateInput) {
    team {
      id
      uuid
      slug
      name
      description
      createdAt
      privacy
      isDefaultTeam
      defaultMemberRole
      membersCanCreatePipelines
    }
  }
}
`)

	req.Var("teamUpdateInput", map[string]interface{}{
		"id":                           team.Id,
		"name":                         team.Name,
		"description":                  team.Description,
		"isDefaultTeam":                team.IsDefaultTeam,
		"defaultMemberRole":            team.DefaultMemberRole,
		"membersCanCreatePipelines":    team.MembersCanCreatePipelines,
		"privacy":                      team.Privacy,
	})

	teamUpdateResponse := teamUpdateResponse{}
	if err := c.graphQLRequest(req, &teamUpdateResponse); err != nil {
		return nil, errors.Wrapf(err, "failed to update team %s", team.Id)
	}

	return &teamUpdateResponse.TeamUpdate.Team, nil
}

func (c *Client) DeleteTeam(id string) error {
	req := graphql.NewRequest(`
mutation TeamDeleteMutation($teamDeleteInput: TeamDeleteInput!) {
  teamDelete(input: $teamDeleteInput) {
    deletedTeamID
  }
}
`)

	req.Var("teamDeleteInput", map[string]interface{}{
		"id": id,
	})

	teamDeleteResponse := teamDeleteResponse{}
	if err := c.graphQLRequest(req, &teamDeleteResponse); err != nil {
		return errors.Wrapf(err, "failed to delete team %s", id)
	}

	return nil
}
