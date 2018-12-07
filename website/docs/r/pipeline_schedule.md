---
layout: "buildkite"
page_title: "Buildkite: buildkite_pipeline_schedule resource"
sidebar_current: "docs-buildkite-resource-buildkite-pipeline-schedule"
description: |-
  Manages a buildkite pipeline schedules
---

# buildkite\_pipeline\_schedule

Creates and manages Buildkite pipeline schedules.

## Example Usage

```hcl

resource "buildkite_pipeline" "build_something" {
  name = "Build cool thing"

  default_branch = "master"
  repository     = "git@github.com:my-org/awesome-repo.git"

  step {
    name    = ":pipeline: Fetch pipeline"
    type    = "script"
    command = "buildkite-agent pipeline upload"
  }
}

resource "buildkite_pipeline_schedule" "build_something_weekly" {
  pipeline_slug = "${buildkite_pipeline.build_something.slug}"
  label         = "Wow, scheduled!"
  cron_schedule = "0 5 * * 1"

  env {
    SPECIAL_VAR_USED_FOR_SCHEDULES = "WOW"
  }
}
```

## Argument Reference

The following arguments are supported:

* `pipeline_slug` - (Required) the pipeline slug

* `label` - (Required) label of the schedule build

* `message` - (Optional) Message to display for scheduled builds. Defaults to `"Scheduled build"`.

* `cron_schedule` - (Required) Cron schedule for frequency of the scheduled build

* `commit` - (Optional) Commit for which to run scheduled builds. Defaults to `"HEAD"`.

* `branch` - (Optional) Branch for which to run scheduled builds. Defaults to `"master"`.
 
* `env` - (Optional) Environment parameters for scheduled builds.
 
* `enabled` - (Optional). Defaults to `true`. 
			
## Attributes Reference

* `pipeline_id` - the GraphQL node id of the pipeline

* `schedule_id` - the GraphQL node id of the pipeline schedule

* `created_at` - the time at which the resource was created
			
## Import

Pipelines can be imported using the pipeline slug and schedule UUID:

```
$ terraform import buildkite_pipeline_schedule.build_something_weekly build-cool-thing/<UUID>
```
