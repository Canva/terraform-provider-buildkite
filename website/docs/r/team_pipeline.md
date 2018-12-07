---
layout: "buildkite"
page_title: "Buildkite: buildkite_team_pipeline resource"
sidebar_current: "docs-buildkite-resource-buildkite-team-pipeline"
description: |-
  Manages a buildkite team pipelines
---

# buildkite\_team\_pipeline

Manages team access to pipelines.

## Example Usage

```hcl
resource "buildkite_team" "backend" {
  name = "backend"
}

resource "buildkite_pipeline" "build_something" {
  name = "Build cool thing"

  default_branch = "master"
  repository     = "git@github.com:my-org/awesome-repo.git"

  github {
  }

  step {
    name    = ":pipeline: Fetch pipeline"
    type    = "script"
    command = "buildkite-agent pipeline upload"

    agent_query_rules = [
    ]
  }
}

resource "buildkite_team_pipeline" "backend_build_something" {
  team_id       = "${buildkite_team.backend.team_id}"
  pipeline_slug = "${buildkite_pipeline.build_something.slug}"
  access_level  = "MANAGE_BUILD_AND_READ"
}
```

## Argument Reference

The following arguments are supported:

* `pipeline_slug` - (Required) the pipeline slug

* `team_id` - (Required) the id of the team

* `access_level` - (Optional) the access level of the team. One of: `READ_ONLY`, `BUILD_AND_READ`,
 or `MANAGE_BUILD_AND_READ`. Defaults to `READ_ONLY`.

## Attributes Reference

* `uuid` - the uuid of the team resource

* `created_at` - the time at which the resource was created

* `pipeline_id` - the id of the pipeline resource

## Import

Team pipelines can be imported using the team pipeline id.

```
$ terraform import buildkite_team_pipeline.backend_build_something <id>
```
