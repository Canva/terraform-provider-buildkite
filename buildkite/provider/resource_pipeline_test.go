package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	buildkiteClient "github.com/saymedia/terraform-buildkite/buildkite/client"
)

func TestAccPipeline_basic_unknown(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildkitePipelineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPipeline_basicBitbucket,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildkitePipelineBasicAttributesFactory("bitbucket"),
					resource.TestCheckResourceAttrSet("buildkite_pipeline.test_bitbucket", "webhook_url"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_bitbucket", "bitbucket_settings.#", "1"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_bitbucket", "bitbucket_settings.0.build_pull_requests", "true"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_bitbucket", "bitbucket_settings.0.build_tags", "false"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_bitbucket", "bitbucket_settings.0.publish_commit_status", "true"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_bitbucket", "bitbucket_settings.0.publish_commit_status_per_step", "false"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_bitbucket", "bitbucket_settings.0.pull_request_branch_filter_configuration", ""),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_bitbucket", "bitbucket_settings.0.pull_request_branch_filter_enabled", "false"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_bitbucket", "bitbucket_settings.0.skip_pull_request_builds_for_existing_commits", "true"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_bitbucket", "github_settings.#", "0"),
				),
			},
		},
	})
}

func TestAccPipeline_basic_beanstalk(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildkitePipelineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPipeline_basicGitlab,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBuildkitePipelineBasicAttributesFactory("gitlab"),
					resource.TestCheckResourceAttrSet("buildkite_pipeline.test_gitlab", "webhook_url"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_gitlab", "github_settings.#", "0"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_gitlab", "bitbucket_settings.#", "0"),
				),
			},
		},
	})
}

func TestAccPipeline_basic_github(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildkitePipelineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPipeline_githubSettingsTriggerModeDeployment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "github_settings.#", "1"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "github_settings.0.build_pull_request_forks", "false"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "github_settings.0.build_pull_requests", "true"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "github_settings.0.build_tags", "false"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "github_settings.0.prefix_pull_request_fork_branch_names", "true"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "github_settings.0.publish_blocked_as_pending", "false"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "github_settings.0.publish_commit_status", "true"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "github_settings.0.publish_commit_status_per_step", "false"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "github_settings.0.pull_request_branch_filter_configuration", ""),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "github_settings.0.pull_request_branch_filter_enabled", "false"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "github_settings.0.skip_pull_request_builds_for_existing_commits", "true"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "github_settings.0.trigger_mode", "deployment"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "bitbucket_settings.#", "0"),
				),
			},
		},
	})
}

func TestAccPipeline_basic_bitbucket(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildkitePipelineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPipeline_githubSettingsBuildTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("buildkite_pipeline.test_foo", "webhook_url"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "github_settings.#", "1"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "github_settings.0.build_tags", "true"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "bitbucket_settings.#", "0"),
				),
			},
		},
	})
}

func TestAccPipeline_basic_gitlab(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBuildkitePipelineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPipeline_bitbucketSettingsBuildTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "bitbucket_settings.#", "1"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "bitbucket_settings.0.build_tags", "true"),
					resource.TestCheckResourceAttr("buildkite_pipeline.test_foo", "github_settings.#", "0"),
				),
			},
		},
	})
}

func testAccCheckBuildkitePipelineExists(id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*buildkiteClient.Client)

		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Pipeline ID is set")
		}

		res, err := client.GetPipeline(rs.Primary.ID)

		if err != nil {
			return err
		}

		if res.Slug != rs.Primary.ID {
			return fmt.Errorf("Pipeline not found")
		}

		return nil
	}
}

func testAccCheckBuildkitePipelineDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*buildkiteClient.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buildkite_pipeline" {
			continue
		}
		if !strings.HasPrefix(rs.Primary.Attributes["name"], "tf-acc-") {
			continue
		}

		res, err := client.GetPipeline(rs.Primary.ID)
		if err == nil {
			if res.Slug == rs.Primary.ID {
				return fmt.Errorf("Pipeline still exists")
			}
		}

		// Verify the error
		if _, ok := err.(*buildkiteClient.NotFound); !ok {
			return err
		}
	}

	return nil
}

func testAccCheckBuildkitePipelineBasicAttributesFactory(repoProvider string) resource.TestCheckFunc {
	PipelineStateId := fmt.Sprintf("buildkite_pipeline.test_%v", repoProvider)
	PipelineName := fmt.Sprintf("tf-acc-basic-%v", repoProvider)

	return resource.ComposeTestCheckFunc(
		testAccCheckBuildkitePipelineExists(PipelineStateId),
		resource.TestCheckResourceAttr(PipelineStateId, "id", PipelineName),
		resource.TestCheckResourceAttr(PipelineStateId, "slug", PipelineName),
		resource.TestCheckResourceAttr(PipelineStateId, "name", PipelineName),
		resource.TestCheckResourceAttrSet(PipelineStateId, "repository"),
		resource.TestCheckResourceAttr(PipelineStateId, "step.#", "1"),
		resource.TestCheckResourceAttr(PipelineStateId, "step.0.agent_query_rules.#", "0"),
		resource.TestCheckResourceAttr(PipelineStateId, "step.0.artifact_paths", ""),
		resource.TestCheckResourceAttr(PipelineStateId, "step.0.branch_configuration", ""),
		resource.TestCheckResourceAttr(PipelineStateId, "step.0.command", "echo 'Hello World'"),
		resource.TestCheckResourceAttr(PipelineStateId, "step.0.concurrency", "0"),
		resource.TestCheckResourceAttr(PipelineStateId, "step.0.env.%", "0"),
		resource.TestCheckResourceAttr(PipelineStateId, "step.0.name", "test"),
		resource.TestCheckResourceAttr(PipelineStateId, "step.0.parallelism", "0"),
		resource.TestCheckResourceAttr(PipelineStateId, "step.0.timeout_in_minutes", "0"),
		resource.TestCheckResourceAttr(PipelineStateId, "step.0.type", "script"),
		resource.TestCheckResourceAttr(PipelineStateId, "default_branch", "master"),
		resource.TestCheckResourceAttr(PipelineStateId, "branch_configuration", ""),
		resource.TestCheckResourceAttr(PipelineStateId, "description", ""),
		resource.TestCheckResourceAttr(PipelineStateId, "env.%", "0"),
		resource.TestCheckResourceAttrSet(PipelineStateId, "builds_url"),
		resource.TestCheckResourceAttrSet(PipelineStateId, "web_url"),
	)
}

const testAccPipeline_basicUnknown = `
resource "buildkite_pipeline" "test_unknown" {
  name = "tf-acc-basic-unknown"
  repository = "git@example.com:terraform-provider-buildkite/terraform-buildkite.git"

  step {
    type = "script"
    name = "test"
    command = "echo 'Hello World'"
  }
}
`

const testAccPipeline_basicBeanstalk = `
resource "buildkite_pipeline" "test_beanstalk" {
  name = "tf-acc-basic-beanstalk"
  repository = "git@terraform-provider-buildkite.git.beanstalkapp.com:/terraform-provider-buildkite/terraform-buildkite.git"

  step {
    type = "script"
    name = "test"
    command = "echo 'Hello World'"
  }
}
`
const testAccPipeline_basicGithub = `
resource "buildkite_pipeline" "test_github" {
  name = "tf-acc-basic-github"
  repository = "git@github.com:saymedia/terraform-provider-buildkite.git"

  step {
    type = "script"
    name = "test"
    command = "echo 'Hello World'"
  }
}
`

const testAccPipeline_basicBitbucket = `
resource "buildkite_pipeline" "test_bitbucket" {
  name = "tf-acc-basic-bitbucket"
  repository = "git@bitbucket.org:terraform-provider-buildkite/terraform-buildkite.git"

  step {
    type = "script"
    name = "test"
    command = "echo 'Hello World'"
  }
}
`

const testAccPipeline_basicGitlab = `
resource "buildkite_pipeline" "test_gitlab" {
  name = "tf-acc-basic-gitlab"
  repository = "git@gitlab.com:terraform-provider-buildkite/terraform-buildkite.git"

  step {
    type = "script"
    name = "test"
    command = "echo 'Hello World'"
  }
}
`

const testAccPipeline_githubSettingsTriggerModeDeployment = `
resource "buildkite_pipeline" "test_foo" {
  name = "tf-acc-foo"
  repository = "git@github.com:saymedia/terraform-provider-buildkite.git"

  step {
    type = "script"
    name = "test"
    command = "echo 'Hello World'"
  }

  github_settings {
	trigger_mode = "deployment"
  }
}
`

const testAccPipeline_githubSettingsBuildTags = `
resource "buildkite_pipeline" "test_foo" {
  name = "tf-acc-foo"
  repository = "git@github.com:saymedia/terraform-provider-buildkite.git"

  step {
    type = "script"
    name = "test"
    command = "echo 'Hello World'"
  }

  github_settings {
	  build_tags = true
  }
}
`

const testAccPipeline_bitbucketSettingsBuildTags = `
resource "buildkite_pipeline" "test_foo" {
  name = "tf-acc-foo"
  repository = "git@bitbucket.org:terraform-provider-buildkite/terraform-buildkite.git"

  step {
    type = "script"
    name = "test"
    command = "echo 'Hello World'"
  }
  
  bitbucket_settings {
	  build_tags = true
  }
}
`
