package provider

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/saymedia/terraform-buildkite/buildkite/client"
)

func resourcePipelineSchedule() *schema.Resource {
	return &schema.Resource{
		Create: CreatePipelineSchedule,
		Read:   ReadPipelineSchedule,
		Update: UpdatePipelineSchedule,
		Delete: DeletePipelineSchedule,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"pipeline_slug": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"schedule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pipeline_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"message": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Scheduled build",
			},
			"cron_schedule": {
				Type:     schema.TypeString,
				Required: true,
			},
			"commit": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "HEAD",
			},
			"branch": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "master",
			},
			"env": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func CreatePipelineSchedule(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] CreatePipelineSchedule")

	buildkiteClient := meta.(*client.Client)

	pipelineSchedule := preparePipelineScheduleRequestPayload(d)

	res, err := buildkiteClient.CreatePipelineSchedule(pipelineSchedule)
	if err != nil {
		return err
	}

	return updatePipelineScheduleFromAPI(d, res)
}

func ReadPipelineSchedule(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] ReadPipelineSchedule")

	buildkiteClient := meta.(*client.Client)
	memberId := d.Id()

	pipelineSchedule, err := buildkiteClient.GetPipelineSchedule(memberId)
	if err != nil {
		if _, ok := err.(*client.NotFound); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	return updatePipelineScheduleFromAPI(d, pipelineSchedule)
}

func UpdatePipelineSchedule(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] UpdatePipelineSchedule")

	buildkiteClient := meta.(*client.Client)

	pipelineSchedule := preparePipelineScheduleRequestPayload(d)

	res, err := buildkiteClient.UpdatePipelineSchedule(pipelineSchedule)
	if err != nil {
		return err
	}

	return updatePipelineScheduleFromAPI(d, res)
}

func DeletePipelineSchedule(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] DeletePipelineSchedule")

	buildkiteClient := meta.(*client.Client)
	id := d.Get("schedule_id").(string)

	return buildkiteClient.DeletePipelineSchedule(id)
}

func updatePipelineScheduleFromAPI(d *schema.ResourceData, t *client.PipelineSchedule) error {
	d.SetId(fmt.Sprintf("%s/%s", t.Pipeline.Slug, t.UUID))
	log.Printf("[INFO] buildkite: team member ID: %s", d.Id())

	d.Set("schedule_id", t.Id)
	d.Set("created_at", t.CreatedAt)
	d.Set("pipeline_id", t.Pipeline.Id)
	d.Set("pipeline_slug", t.Pipeline.Slug)
	d.Set("label", t.Label)
	d.Set("message", t.Message)
	d.Set("cron_schedule", t.CronSchedule)
	d.Set("commit", t.Commit)
	d.Set("branch", t.Branch)
	d.Set("env", listToMap(t.Environment))
	d.Set("enabled", t.Enabled)

	return nil
}

func preparePipelineScheduleRequestPayload(d *schema.ResourceData) *client.PipelineSchedule {
	req := &client.PipelineSchedule{}

	req.UUID = d.Id()

	if val, ok := d.GetOkExists("pipeline_id"); ok {
		req.Pipeline.Id = val.(string)
	}
	req.Pipeline.Slug = d.Get("pipeline_slug").(string)

	req.Id = d.Get("schedule_id").(string)
	req.Label = d.Get("label").(string)
	req.Message = d.Get("message").(string)
	req.CronSchedule = d.Get("cron_schedule").(string)
	req.Commit = d.Get("commit").(string)
	req.Branch = d.Get("branch").(string)
	req.Environment = mapToList(d.Get("env").(map[string]interface{}))
	req.Enabled = d.Get("enabled").(bool)

	return req
}

func listToMap(list []string) map[string]string {
	result := map[string]string{}
	for _, value := range list {
		keyValue := strings.Split(value, "=")
		result[keyValue[0]] = keyValue[1]
	}
	return result
}

func mapToList(m map[string]interface{}) []string {
	var result []string
	for key, value := range m {
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}
	return result
}
