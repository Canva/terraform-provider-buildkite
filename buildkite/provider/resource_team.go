package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/saymedia/terraform-buildkite/buildkite/client"
)

var (
	ValidTeamPrivacy    = []string{client.TeamPrivacySecret, client.TeamPrivacyVisible}
	ValidTeamMemberRole = []string{client.TeamMemberRoleMember, client.TeamMemberRoleMaintainer}
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		Create: CreateTeam,
		Read:   ReadTeam,
		Update: UpdateTeam,
		Delete: DeleteTeam,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"privacy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      client.TeamPrivacyVisible,
				ValidateFunc: validation.StringInSlice(ValidTeamPrivacy, false),
			},
			"default_member_role": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      client.TeamMemberRoleMember,
				ValidateFunc: validation.StringInSlice(ValidTeamMemberRole, false),
			},
			"is_default_team": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
		    },
			"members_can_create_pipelines": {
			    Type:     schema.TypeBool,
			    Optional: true,
			    Default:  false,
			},
		},
	}
}

func CreateTeam(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] CreatePipeline")

	buildkiteClient := meta.(*client.Client)

	team := prepareTeamRequestPayload(d)

	res, err := buildkiteClient.CreateTeam(team)
	if err != nil {
		return err
	}

	return updateTeamFromAPI(d, res)
}

func ReadTeam(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] ReadPipeline")

	buildkiteClient := meta.(*client.Client)
	slug := d.Id()

	team, err := buildkiteClient.GetTeam(slug)
	if err != nil {
		if _, ok := err.(*client.NotFound); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	return updateTeamFromAPI(d, team)
}

func UpdateTeam(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] UpdatePipeline")

	buildkiteClient := meta.(*client.Client)

	team := prepareTeamRequestPayload(d)

	res, err := buildkiteClient.UpdateTeam(team)
	if err != nil {
		return err
	}

	return updateTeamFromAPI(d, res)
}

func DeleteTeam(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] DeletePipeline")

	buildkiteClient := meta.(*client.Client)
	id := d.Get("team_id").(string)

	return buildkiteClient.DeleteTeam(id)
}

func updateTeamFromAPI(d *schema.ResourceData, t *client.Team) error {
	d.SetId(t.Slug)
	log.Printf("[INFO] buildkite: Pipeline ID: %s", d.Id())

	d.Set("team_id", t.Id)
	d.Set("uuid", t.UUID)
	d.Set("slug", t.Slug)
	d.Set("name", t.Name)
	d.Set("description", t.Description)
	d.Set("created_at", t.CreatedAt)
	d.Set("privacy", t.Privacy)
	d.Set("is_default_team", t.IsDefaultTeam)
	d.Set("default_member_role", t.DefaultMemberRole)
	d.Set("members_can_create_pipelines", t.MembersCanCreatePipelines)

	return nil
}

func prepareTeamRequestPayload(d *schema.ResourceData) *client.Team {
	req := &client.Team{}

	if val, ok := d.GetOkExists("team_id"); ok {
		req.Id = val.(string)
	}
	if val, ok := d.GetOkExists("uuid"); ok {
		req.UUID = val.(string)
	}
	req.Slug = d.Get("slug").(string)
	req.Name = d.Get("name").(string)
	req.Description = d.Get("description").(string)
	req.Privacy = d.Get("privacy").(string)
	req.CreatedAt = d.Get("created_at").(string)
	req.IsDefaultTeam = d.Get("is_default_team").(bool)
	req.DefaultMemberRole = d.Get("default_member_role").(string)
	req.MembersCanCreatePipelines = d.Get("members_can_create_pipelines").(bool)

	return req
}
