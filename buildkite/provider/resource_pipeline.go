package provider

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/saymedia/terraform-buildkite/buildkite/client"
)

func resourcePipeline() *schema.Resource {
	return &schema.Resource{
		Create: CreatePipeline,
		Read:   ReadPipeline,
		Update: UpdatePipeline,
		Delete: DeletePipeline,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"slug": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"web_url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"builds_url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"badge_url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"repository": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"branch_configuration": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_branch": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "master",
			},
			"env": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"provider_settings": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"webhook_url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"step": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"command": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"env": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"timeout_in_minutes": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"agent_query_rules": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"artifact_paths": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"branch_configuration": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"concurrency": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"parallelism": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func CreatePipeline(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] CreatePipeline")

	buildkiteClient := meta.(*client.Client)

	pipeline := preparePipelineRequestPayload(d)

	res, err := buildkiteClient.CreatePipeline(pipeline)
	if err != nil {
		return err
	}

	return updatePipelineFromAPI(d, res)
}

func ReadPipeline(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] ReadPipeline")

	buildkiteClient := meta.(*client.Client)
	slug := d.Id()

	pipeline, err := buildkiteClient.GetPipeline(slug)
	if err != nil {
		if _, ok := err.(*client.NotFound); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	return updatePipelineFromAPI(d, pipeline)
}

func UpdatePipeline(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] UpdatePipeline")

	buildkiteClient := meta.(*client.Client)

	pipeline := preparePipelineRequestPayload(d)

	res, err := buildkiteClient.UpdatePipeline(pipeline)
	if err != nil {
		return err
	}

	return updatePipelineFromAPI(d, res)
}

func DeletePipeline(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] DeletePipeline")

	buildkiteClient := meta.(*client.Client)
	slug := d.Id()

	return buildkiteClient.DeletePipeline(slug)
}

func updatePipelineFromAPI(d *schema.ResourceData, p *client.Pipeline) error {
	d.SetId(p.Slug)
	log.Printf("[INFO] Pipeline ID: %s", d.Id())

	d.Set("env", p.Environment)
	d.Set("name", p.Name)
	d.Set("description", p.Description)
	d.Set("repository", p.Repository)
	d.Set("web_url", p.WebURL)
	d.Set("slug", p.Slug)
	d.Set("builds_url", p.BuildsURL)
	d.Set("branch_configuration", p.BranchConfiguration)
	d.Set("provider_settings", p.Provider.Settings)
	d.Set("webhook_url", p.Provider.WebhookURL)
	d.Set("default_branch", p.DefaultBranch)

	stepMap := make([]interface{}, len(p.Steps))
	for i, element := range p.Steps {
		stepMap[i] = map[string]interface{}{
			"type":                 element.Type,
			"name":                 element.Name,
			"command":              element.Command,
			"env":                  element.Environment,
			"agent_query_rules":    element.AgentQueryRules,
			"branch_configuration": element.BranchConfiguration,
			"artifact_paths":       element.ArtifactPaths,
			"concurrency":          element.Concurrency,
			"parallelism":          element.Parallelism,
			"timeout_in_minutes":   element.TimeoutInMinutes,
		}
	}
	d.Set("step", stepMap)
	return nil
}

func preparePipelineRequestPayload(d *schema.ResourceData) *client.Pipeline {
	req := &client.Pipeline{}

	req.Name = d.Get("name").(string)
	req.DefaultBranch = d.Get("default_branch").(string)
	req.Description = d.Get("description").(string)
	req.Slug = d.Get("slug").(string)
	req.Repository = d.Get("repository").(string)
	req.BranchConfiguration = d.Get("branch_configuration").(string)
	req.Environment = map[string]string{}
	for k, vI := range d.Get("env").(map[string]interface{}) {
		req.Environment[k] = vI.(string)
	}
	req.ProviderSettings = map[string]string{}
	for k, vI := range d.Get("provider_settings").(map[string]interface{}) {
		req.ProviderSettings[k] = vI.(string)
	}

	stepsI := d.Get("step").([]interface{})
	req.Steps = make([]client.Step, len(stepsI))

	for i, stepI := range stepsI {
		stepM := stepI.(map[string]interface{})
		req.Steps[i] = client.Step{
			Type:                stepM["type"].(string),
			Name:                stepM["name"].(string),
			Command:             stepM["command"].(string),
			Environment:         map[string]string{},
			AgentQueryRules:     make([]string, len(stepM["agent_query_rules"].([]interface{}))),
			BranchConfiguration: stepM["branch_configuration"].(string),
			ArtifactPaths:       stepM["artifact_paths"].(string),
			Concurrency:         stepM["concurrency"].(int),
			Parallelism:         stepM["parallelism"].(int),
			TimeoutInMinutes:    stepM["timeout_in_minutes"].(int),
		}

		for k, vI := range stepM["env"].(map[string]interface{}) {
			req.Steps[i].Environment[k] = vI.(string)
		}

		for j, vI := range stepM["agent_query_rules"].([]interface{}) {
			req.Steps[i].AgentQueryRules[j] = vI.(string)
		}
	}

	return req
}
