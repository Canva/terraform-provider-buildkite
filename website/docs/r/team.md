---
layout: "buildkite"
page_title: "Buildkite: buildkite_team resource"
sidebar_current: "docs-buildkite-resource-buildkite-team"
description: |-
  Manages a buildkite team 
---

# buildkite\_team

Buildkite organiation members can be group in teams, where teams can be given certain access to pipelines.

## Example Usage

```hcl
resource "buildkite_team" "backend" {
  name = "backend"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) the team name

* `description` - (Optional) description of the team

* `privacy` - (Optional) the privacy setting of the team. One of: `VISIBLE`, or `SECRET`.
 If the privacy is set to `SECRET` only members will be able to see the team. Defaults to `VISIBLE`.

* `default_member_role` - (Optional) the default role members will get. One of: `MEMBER`, or `MAINTAINER`. 
 Defaults to `MEMBER`.

* `is_default_team` - (Optional) if marked as default team all organization members are added to this team. Defaults to `false`

## Attributes Reference

* `slug` - the slug of the team resource

* `uuid` - the uuid of the team resource

* `team_id` - the id of the team resource

* `created_at` - the time at which the resource was created

## Import

Organization members can be imported using the team slug

```
$ terraform import buildkite_team.backend backend
```
