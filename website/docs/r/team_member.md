team.md---
layout: "buildkite"
page_title: "Buildkite: buildkite_team_member resource"
sidebar_current: "docs-buildkite-resource-buildkite-team-member"
description: |-
  Manages a buildkite team memberships
---

# buildkite\_team

Manages memberships of organization members into teams.

## Example Usage

```hcl
resource "buildkite_team" "backend" {
  name = "backend"
}

resource "buildkite_org_member" "user1" {
  role = "MEMBER"
}

resource "buildkite_org_member" "user2" {
  role = "MEMBER"
}

resource "buildkite_team_member" "user1_backend" {
  user_id = "${buildkite_org_member.user1.user_id}"
  team_id = "${buildkite_team.backend.team_id}"
  role    = "MEMBER"
}

resource "buildkite_team_member" "user2_backend" {
  user_id = "${buildkite_org_member.user2.user_id}"
  team_id = "${buildkite_team.backend.team_id}"
  role    = "MAINTAINER"
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required) the id of the organization user

* `team_id` - (Requireed) the id of the team

* `role` - (Optional) the role of the team member. One of: `MEMBER`, or `MAINTAINER`. Defaults to `MEMBER`.

## Attributes Reference

* `uuid` - the uuid of the team membership resource

* `created_at` - the time at which the resource was created

## Import

Team members can be imported using the team membership id

```
$ terraform import buildkite_team_member.user1_backend <membership-id>
```

You can get the id via Buildkite's GraphQL API.
