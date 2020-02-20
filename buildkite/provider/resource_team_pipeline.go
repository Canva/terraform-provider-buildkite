package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/saymedia/terraform-buildkite/buildkite/client"
)

var (
	ValidTeamPipelineAccessLevels = []string{
		client.TeamPipelineAccessReadOnly,
		client.TeamPipelineAccessBuildAndRead,
		client.TeamPipelineAccessManage,
	}
)

func resourceTeamPipeline() *schema.Resource {
	return &schema.Resource{
		Create: CreateTeamPipeline,
		Read:   ReadTeamPipeline,
		Update: UpdateTeamPipeline,
		Delete: DeleteTeamPipeline,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"pipeline_slug": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pipeline_id": {
				Type:     schema.TypeString,
				Computed: true,
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
			"access_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      client.TeamPipelineAccessReadOnly,
				ValidateFunc: validation.StringInSlice(ValidTeamPipelineAccessLevels, false),
			},
		},
	}
}

func CreateTeamPipeline(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] CreateTeamPipeline")

	buildkiteClient := meta.(*client.Client)

	teamPipeline := prepareTeamPipelineRequestPayload(d)

	res, err := buildkiteClient.CreateTeamPipeline(teamPipeline)
	if err != nil {
		return err
	}

	if err = updateTeamPipelineFromAPI(d, res); err != nil {
		return err
	}

	// Create does not take accessLevel as argument. All team pipelines are created as 'READ_ONLY'.
	// If that's the desired role we are done, otherwise we need to issue an update API request
	// Handling this here in the resource instead of the client to be able
	// to populate the TF state with some values already
	if teamPipeline.AccessLevel == client.TeamPipelineAccessReadOnly {
		return nil
	}

	res.AccessLevel = teamPipeline.AccessLevel
	res, err = buildkiteClient.UpdateTeamPipeline(res)
	if err != nil {
		return err
	}

	return updateTeamPipelineFromAPI(d, res)
}

func ReadTeamPipeline(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] ReadTeamPipeline")

	buildkiteClient := meta.(*client.Client)
	memberId := d.Id()

	teamPipeline, err := buildkiteClient.GetTeamPipeline(memberId)
	if err != nil {
		if _, ok := err.(*client.NotFound); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	return updateTeamPipelineFromAPI(d, teamPipeline)
}

func UpdateTeamPipeline(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] UpdateTeamPipeline")

	buildkiteClient := meta.(*client.Client)

	teamPipeline := prepareTeamPipelineRequestPayload(d)

	res, err := buildkiteClient.UpdateTeamPipeline(teamPipeline)
	if err != nil {
		return err
	}

	return updateTeamPipelineFromAPI(d, res)
}

func DeleteTeamPipeline(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] DeleteTeamPipeline")

	buildkiteClient := meta.(*client.Client)
	id := d.Id()

	return buildkiteClient.DeleteTeamPipeline(id)
}

func updateTeamPipelineFromAPI(d *schema.ResourceData, t *client.TeamPipeline) error {
	d.SetId(t.Id)
	log.Printf("[INFO] buildkite: team member ID: %s", d.Id())

	d.Set("uuid", t.UUID)
	d.Set("access_level", t.AccessLevel)
	d.Set("created_at", t.CreatedAt)
	d.Set("team_id", t.Team.Id)
	d.Set("pipeline_id", t.Pipeline.Id)
	d.Set("pipeline_slug", t.Pipeline.Slug)

	return nil
}

func prepareTeamPipelineRequestPayload(d *schema.ResourceData) *client.TeamPipeline {
	req := &client.TeamPipeline{}

	req.Id = d.Id()
	if val, ok := d.GetOkExists("uuid"); ok {
		req.UUID = val.(string)
	}
	req.AccessLevel = d.Get("access_level").(string)
	req.Team.Id = d.Get("team_id").(string)

	if val, ok := d.GetOkExists("pipeline_id"); ok {
		req.Pipeline.Id = val.(string)
	}
	req.Pipeline.Slug = d.Get("pipeline_slug").(string)

	return req
}
