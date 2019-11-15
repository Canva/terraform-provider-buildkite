---
layout: "buildkite"
page_title: "Buildkite: buildkite_pipeline resource"
sidebar_current: "docs-buildkite-resource-buildkite-pipeline"
description: |-
  Manages a buildkite pipeline 
---

# buildkite\_pipeline

Creates and manages Buildkite pipelines. Have a look at the 
[Pipelines API](https://buildkite.com/docs/apis/rest-api/pipelines), if any parameters are unclear.

## Example Usage

```hcl

resource "buildkite_pipeline" "build_something" {
  name = "Build cool thing"

  default_branch = "master"
  repository     = "git@github.com:my-org/awesome-repo.git"
  
  github {
    build_pull_request_forks = true
  }

  step {
    name    = ":pipeline: Fetch pipeline"
    type    = "script"
    command = "buildkite-agent pipeline upload"

    agent_query_rules = [
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) the pipeline name

* `description` - (Optional) description of the pipeline

* `repository` - (Required) the repository of the code to build

* `branch_configuration` - (Optional) A branch filter pattern to limit which pushed branches trigger builds on this pipeline

* `default_branch` - (Optional) the default branch to build. Defaults to `master`

* `env` - (Optional) pipeline environment variables

* `step` - (Required) nested block list configuring the steps to run. Must provide at least one.

* `bitbucket_settings` - (Optional)

* `github_settings` - (Optional)

Only one of repository settings blocks `bitbucket_settings`, `github_settings` may be present.

### Step Options

For more information about steps, take a look at the [official documentation](https://buildkite.com/docs/pipelines/command-step)

* `type` - (Required) the step type, one of `script`, `waiter`, or `manual`

* `name` - (Optional) the name of the step

* `command` - (Optional) the command to execute

* `env` - (Optional) step environment variables

* `timeout_in_minutes` - (Optional) The number of minutes a job created from this step is allowed to run. If the job does not finish within this limit, it will be automatically cancelled and the build will fail.

* `agent_query_rules` - (Optional) Key-value map to query agents based on queues or other labels.

* `artifact_paths` - (Optional) The glob path or paths of artifacts to upload from this step.

* `branch_configuration` - (Optional) A branch filter pattern to limit for which branches to run this step

* `concurrency` - (Optional) The maximum number of jobs created from this step that are allowed to run at the same time. If you use this attribute, you must also define a label for it with the concurrency_group attribute. Read more about [controlling concurrency](https://buildkite.com/docs/pipelines/controlling-concurrency) in the official docs.

* `parallelism` - (Optional) A unique name for the concurrency group that you are creating with the concurrency attribute.


### Bitbucket Options

* `trigger_mode` - (Optional) The trigger mode for builds. Defaults to `"code"`

* `build_pull_requests` - (Optional) Whether to build pull requests. Defaults to `true`

* `pull_request_branch_filter_configuration` - (Optional) Branch filter for pull request builds

* `pull_request_branch_filter_enabled` - (Optional) Enable branch filtering for pull request builds. Defaults to `false`

* `skip_pull_request_builds_for_existing_commits` - (Optional) Do not rebuild existing commits in pull requests. Defaults to `true`

* `prefix_pull_request_fork_branch_names` - (Optional) Defaults to `true`

* `build_tags` - (Optional) Build git tags. Defaults to `false`

* `publish_commit_status` - (Optional) Publish build status as commit status in Bitbucket. Defaults to `true`

* `publish_commit_status_per_step` - (Optional) Publish a commit status for every step of the pipeline. Defaults to `false`


### GitHub Options

* `trigger_mode` - (Optional) The trigger mode for builds. Defaults to `"code"`

* `build_pull_requests` - (Optional) Whether to build pull requests. Defaults to `true`

* `pull_request_branch_filter_configuration` - (Optional) Branch filter for pull request builds

* `pull_request_branch_filter_enabled` - (Optional) Enable branch filtering for pull request builds

* `skip_pull_request_builds_for_existing_commits` - (Optional) Do not rebuild existing commits in pull requests. Defaults to `true`

* `build_pull_request_forks` - (Optional) Whether to build pull requests from forks.

* `prefix_pull_request_fork_branch_names` - (Optional) Defaults to `true`

* `build_tags` - (Optional) Build git tags. Defaults to `false`

* `publish_commit_status` - (Optional) Publish build status as commit status in GitHub. Defaults to `true`

* `publish_commit_status_per_step` - (Optional) Publish a commit status for every step of the pipeline. Defaults to `false`

* `publish_blocked_as_pending` - (Optional) Whether blocked steps should be published as "pending"

* `separate_pull_request_statuses` - (Optional) Publish separate status for the pull request itself.


## Attributes Reference

* `slug` - the slug of the pipeline

* `created_at` - the time at which the resource was created

* `web_url` - the web url of the pipeline

* `builds_url` - the builds url of the pipeline

* `url` - the url of the pipeline

* `badge_url` - the badge web url of the pipeline

* `webhook_url` - the webhook url of the pipeline
			
## Import

Pipelines can be imported using the pipeline slug

```
$ terraform import buildkite_pipeline.build_something build-cool-thing
```
