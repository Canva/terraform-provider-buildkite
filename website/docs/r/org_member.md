---
layout: "buildkite"
page_title: "Buildkite: buildkite_org_member resource"
sidebar_current: "docs-buildkite-resource-buildkite-org-member"
description: |-
  Manages a buildkite organization membership 
---

# buildkite\_org\_member

Organization members cannot be create directly. They can only be invited, once a user has singed up to the organization
they membership can be managed via this resource. This only the update, and delete operations are supported.
You need to define the resource, and then import it to be able to manage it.

## Example Usage

```hcl
resource "buildkite_org_member" "test_user" {
  role = "MEMBER"
}

resource "buildkite_org_member" "admin_user" {
  role = "ADMIN"
}
```

## Argument Reference

The following arguments are supported:

* `role` - (Required) the organization role of the member, one of: `MEMBER`, `ADMIN`


## Attributes Reference

* `uuid` - the uuid of the organization member resource

* `member_id` - the id of the organization member resource

* `created_at` - the time at which the resource was created

* `user_id` - the id of the user

* `user_name` - the name of the user

* `user_email` - the email of the user

## Import

Organization members can be imported using the uuid of the membership, e.g.:

```
$ terraform import buildkite_org_member.admin_user <uuid>
```

You can get the uuid via Buildkite's GraphQL API.
