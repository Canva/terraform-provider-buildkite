---
layout: "buildkite"
page_title: "Provider: Buildkite"
sidebar_current: "docs-buildkite-index"
description: |-
  The Buildkite provider allows Terraform to configure Buildkite
---

# Buildkite Provider

The Buildkite provider allows Terraform to configure 
[Buildkite](https://www.buildkite.com/).

## Configuring Buildkite

Terraform can be used by the Buildkite adminstrators to configure Buildkite
organization members, teams, team members, pipelines, and team pipelines.

## Provider Arguments

The provider configuration block accepts the following arguments.
In most cases it is recommended to set them via the indicated environment
variables in order to keep credential information out of the configuration.

* `organization` - (Required) Organization slug which should be managed.
  May be set via the `BUILDKITE_ORGANIZATION` environment variable.

* `api_token` - (Required) Buildkite API token that will be used by Terraform to
  authenticate. May be set via the `BUILDKITE_API_TOKEN` environment variable.
  It needs the `read_pipeline`, `write_pipeline`, and `graphql` privileges.

## Example Usage

```hcl
provider "buildkite" {
  // Get an API token from https://buildkite.com/user/api-access-tokens
  // Needs: read_pipelines, write_pipelines, and GraphQL
  // expose via env variables BUILDKITE_API_TOKEN and BUILDKITE_ORGANIZATION
  version = "0.1.0"
}

resource "buildkite_team" "backend" {
  name = "backend"
}
```
