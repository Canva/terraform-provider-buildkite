package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/saymedia/terraform-buildkite/buildkite/client"
)

func resourceTeamMember() *schema.Resource {
	return &schema.Resource{
		Create: CreateTeamMember,
		Read:   ReadTeamMember,
		Update: UpdateTeamMember,
		Delete: DeleteTeamMember,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      client.TeamMemberRoleMember,
				ValidateFunc: validation.StringInSlice(ValidTeamMemberRole, false),
			},
		},
	}
}

func CreateTeamMember(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] CreateTeamMember")

	buildkiteClient := meta.(*client.Client)

	teamMember := prepareTeamMemberRequestPayload(d)

	res, err := buildkiteClient.CreateTeamMember(teamMember)
	if err != nil {
		return err
	}

	if err = updateTeamMemberFromAPI(d, res); err != nil {
		return err
	}

	// Create does not take role as argument. All team members are created as 'MEMBER'.
	// If that's the desired role we are done, otherwise we need to issue an update API request
	// Handling this here in the resource instead of the client to be able
	// to populate the TF state with some values already
	if teamMember.Role == client.TeamMemberRoleMember {
		return nil
	}

	res.Role = teamMember.Role
	res, err = buildkiteClient.UpdateTeamMember(res)
	if err != nil {
		return err
	}

	return updateTeamMemberFromAPI(d, res)
}

func ReadTeamMember(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] ReadTeamMember")

	buildkiteClient := meta.(*client.Client)
	memberId := d.Id()

	teamMember, err := buildkiteClient.GetTeamMember(memberId)
	if err != nil {
		if _, ok := err.(*client.NotFound); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	return updateTeamMemberFromAPI(d, teamMember)
}

func UpdateTeamMember(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] UpdateTeamMember")

	buildkiteClient := meta.(*client.Client)

	teamMember := prepareTeamMemberRequestPayload(d)

	res, err := buildkiteClient.UpdateTeamMember(teamMember)
	if err != nil {
		return err
	}

	return updateTeamMemberFromAPI(d, res)
}

func DeleteTeamMember(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] DeleteTeamMember")

	buildkiteClient := meta.(*client.Client)
	id := d.Id()

	return buildkiteClient.DeleteTeamMember(id)
}

func updateTeamMemberFromAPI(d *schema.ResourceData, t *client.TeamMember) error {
	d.SetId(t.Id)
	log.Printf("[INFO] buildkite: team member ID: %s", d.Id())

	d.Set("uuid", t.UUID)
	d.Set("role", t.Role)
	d.Set("created_at", t.CreatedAt)
	d.Set("team_id", t.Team.Id)
	d.Set("user_id", t.User.Id)

	return nil
}

func prepareTeamMemberRequestPayload(d *schema.ResourceData) *client.TeamMember {
	req := &client.TeamMember{}

	req.Id = d.Id()
	if val, ok := d.GetOkExists("uuid"); ok {
		req.UUID = val.(string)
	}
	req.Role = d.Get("role").(string)
	req.Team.Id = d.Get("team_id").(string)
	req.User.Id = d.Get("user_id").(string)

	return req
}
