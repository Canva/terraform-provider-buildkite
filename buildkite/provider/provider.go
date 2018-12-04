package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/saymedia/terraform-buildkite/buildkite/client"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"buildkite_org_member":    resourceOrgMember(),
			"buildkite_pipeline":      resourcePipeline(),
			"buildkite_team":          resourceTeam(),
			"buildkite_team_member":   resourceTeamMember(),
			"buildkite_team_pipeline": resourceTeamPipeline(),
		},

		Schema: map[string]*schema.Schema{
			"organization": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BUILDKITE_ORGANIZATION", nil),
			},
			"api_token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BUILDKITE_API_TOKEN", nil),
			},
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	orgName := d.Get("organization").(string)
	apiToken := d.Get("api_token").(string)

	return client.NewClient(orgName, apiToken), nil
}
