package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"

	"github.com/saymedia/terraform-buildkite/buildkite/client"
)

var (
	providerSettingsExcluded = []string{"repository", "account"}
	pipelineSchema = map[string]*schema.Schema{
		"slug": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"web_url": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"builds_url": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"url": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"badge_url": {
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
		"repository": {
			Type:     schema.TypeString,
			Required: true,
		},
		"branch_configuration": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"default_branch": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "master",
		},
		"env": {
			Type:          schema.TypeMap,
			Optional:      true,
			ConflictsWith: []string{"configuration"},
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"webhook_url": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"configuration": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"step", "env"},
		},
		"team_uuids": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"step": {
			Type:          schema.TypeList,
			Optional:      true,
			ConflictsWith: []string{"configuration"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"command": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"env": {
						Type:     schema.TypeMap,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"timeout_in_minutes": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"agent_query_rules": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"artifact_paths": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"branch_configuration": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"concurrency": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"parallelism": {
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},
		"bitbucket_settings": {
			Type:          schema.TypeList,
			Optional:      true,
			Computed:      true,
			MaxItems:      1,
			ConflictsWith: []string{"github_settings"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"trigger_mode": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  "code",
					},
					"build_pull_requests": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
					"pull_request_branch_filter_enabled": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
					"pull_request_branch_filter_configuration": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"skip_pull_request_builds_for_existing_commits": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
					"prefix_pull_request_fork_branch_names": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
					"build_tags": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
					"publish_commit_status": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
					"publish_commit_status_per_step": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},
		"github_settings": {
			Type:          schema.TypeList,
			Optional:      true,
			Computed:      true,
			MaxItems:      1,
			ConflictsWith: []string{"bitbucket_settings"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"trigger_mode": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  "code",
					},
					"build_pull_requests": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
					"pull_request_branch_filter_enabled": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"pull_request_branch_filter_configuration": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"skip_pull_request_builds_for_existing_commits": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
					"build_pull_request_forks": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"prefix_pull_request_fork_branch_names": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
					"build_tags": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"publish_commit_status": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
					"publish_commit_status_per_step": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"publish_blocked_as_pending": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"separate_pull_request_statuses": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"filter_enabled": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
	}
)

func resourcePipeline() *schema.Resource {
	resource := schema.Resource{
		Create: CreatePipeline,
		Read:   ReadPipeline,
		Update: UpdatePipeline,
		Delete: DeletePipeline,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: pipelineSchema,
	}
	return &resource
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
	log.Printf("[INFO] buildkite: Pipeline ID: %s", d.Id())

	d.Set("env", p.Environment)
	d.Set("name", p.Name)
	d.Set("description", p.Description)
	d.Set("repository", p.Repository)
	d.Set("web_url", p.WebURL)
	d.Set("slug", p.Slug)
	d.Set("builds_url", p.BuildsURL)
	d.Set("branch_configuration", p.BranchConfiguration)
	d.Set("default_branch", p.DefaultBranch)
	d.Set("configuration", p.Configuration)
	d.Set("team_uuids", p.TeamUUIDs)

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
	if err := d.Set("step", stepMap); err != nil {
		return err
	}

	emptySettings := make([]interface{}, 0)
	d.Set("github_settings", emptySettings)
	d.Set("bitbucket_settings", emptySettings)

	log.Printf("[INFO] buildkite: RepositoryProviderId: %s", p.Provider.Id)

	d.Set("webhook_url", p.Provider.WebhookURL)

	switch p.Provider.Id {
	case "github":
		log.Printf("[DEBUG] buildkite: Provider.Settings in github: %+v", p.Provider.Settings)
		if err := d.Set("github_settings", filterProviderSettings("github_settings", p.Provider.Settings)); err != nil {
			return err
		}

	case "bitbucket":
		log.Printf("[DEBUG] buildkite: Provider.Settings in bitbucket: %+v", p.Provider.Settings)
		if err := d.Set("bitbucket_settings", filterProviderSettings("bitbucket_settings", p.Provider.Settings)); err != nil {
			return err
		}

	case "gitlab": // noop
	case "beanstalk": // noop
	default: // unknown, noop
	}

	return nil
}

func filterProviderSettings(
	name string,
	providerSettings map[string]interface{}) []map[string]interface{} {

	result := map[string]interface{}{}
	resultList := []map[string]interface{}{result}

	providerSchema, ok := pipelineSchema[name]
	if !ok{
		log.Printf("[ERROR] could not find provider schema for '%s'", name)
		return resultList
	}

	providerResource, ok := providerSchema.Elem.(*schema.Resource)
	if !ok {
		log.Printf("[ERROR] provider schema for '%s' is not a complex object", name)
		return resultList
	}

	for key, value := range providerSettings {
		if contains(providerSettingsExcluded, key) {
			continue
		}

		if _, keyExistsInSchema := providerResource.Schema[key]; !keyExistsInSchema {
			log.Printf("[DEBUG] '%s.0.%s' does not exist in schema", name, key)
			continue
		}
		result[key] = value
	}

	return []map[string]interface{}{result}
}

func contains(strings []string, value string) bool {
	for _, val := range strings {
		if val == value {
			return true
		}
	}
	return false
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
	teamUUIDs := d.Get("team_uuids").([]interface{})
	req.TeamUUIDs = make([]string, len(teamUUIDs))
	for i, t := range teamUUIDs {
		req.TeamUUIDs[i] = t.(string)
	}

	if val, ok := d.GetOk("configuration"); ok {
		req.Configuration = val.(string)
	} else {
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
	}

	if d.HasChange("github_settings") || d.HasChange("bitbucket_settings") {
		log.Printf("[INFO] buildkite: RepositoryProviderSettings have changed")

		githubSettings := d.Get("github_settings").([]interface{})
		bitbucketSettings := d.Get("bitbucket_settings").([]interface{})
		settings := map[string]interface{}{}

		if len(githubSettings) > 0 {
			s := githubSettings[0].(map[string]interface{})

			for k, vI := range s {
				if _, ok := d.GetOk(fmt.Sprintf("github_settings.0.%s", k)); ok {
					settings[k] = vI
				}
			}
		} else if len(bitbucketSettings) > 0 {
			s := bitbucketSettings[0].(map[string]interface{})

			for k, vI := range s {
				if _, ok := d.GetOk(fmt.Sprintf("bitbucket_settings.0.%s", k)); ok {
					settings[k] = vI
				}
			}
		}
		req.ProviderSettings = settings
	}

	return req
}
