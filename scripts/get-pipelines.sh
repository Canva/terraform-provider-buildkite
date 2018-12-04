#!/bin/bash

FILE=pipelines.tf

print_attr() {
    value="$(echo $1 | jq ".$2")"
    if [[ "$?" -ne 0 ]] ; then
      return
    fi
    if [[ "$value" != "null" ]] && [[ "$value" != '""' ]] && [[ "$value" != '{}' ]] ; then
      echo "    $2 = $value"
    fi
}

print_step() {
  echo "  step {"
  print_attr "$1" name
  print_attr "$1" type
  print_attr "$1" command
  print_attr "$1" env
  print_attr "$1" timeout_in_minutes
  print_attr "$1" agent_query_rules
  print_attr "$1" artifact_paths
  print_attr "$1" branch_configuration
  print_attr "$1" concurrency
  print_attr "$1" parallelism
  echo "  }"
}

print_repo_attr() {
    value="$(echo $1 | jq ".provider.settings.$2")"
    default="$3"
    if [[ "$?" -ne 0 ]] ; then
      return
    fi
    if [[ "$value" != "null" ]] && [[ "$value" != '""' ]] && [[ "$value" != '{}' ]] && [[ "$value" != "$default" ]] ; then
      echo "    $2 = $value"
    fi
}

print_github_settings() {
  provider="$(echo "$1" | jq -r '.provider.id')"
  if [[ $? -gt 0 ]] || [[ "$provider" != "github" ]] ; then
    return
  fi
  echo "  github_settings {"
  print_repo_attr "$1" trigger_mode '"code"'
  print_repo_attr "$1" build_pull_requests "true"
  print_repo_attr "$1" pull_request_branch_filter_enabled "false"
  print_repo_attr "$1" skip_pull_request_builds_for_existing_commits "true"
  print_repo_attr "$1" build_pull_request_forks "false"
  print_repo_attr "$1" prefix_pull_request_fork_branch_names "true"
  print_repo_attr "$1" build_tags "false"
  print_repo_attr "$1" publish_commit_status "true"
  print_repo_attr "$1" publish_commit_status_per_step "false"
  print_repo_attr "$1" separate_pull_request_statuses "false"
  print_repo_attr "$1" publish_blocked_as_pending "false"
  echo "  }"
}

get_team_pipelines() {
    local api_key org_slug pipeline
    api_key="$1"
    org_slug="$2"
    pipeline="$3"

    curl -s \
        -H 'content-type: application/json' \
        -H "Authorization: Bearer $api_key" \
        'https://graphql.buildkite.com/v1' \
        -d "
{
  \"operationName\": \"Teams\",
  \"variables\": { \"pipelineSlug\": \"$org_slug/$pipeline\" },
  \"query\": \"query Teams(\$pipelineSlug: ID!) { pipeline(slug: \$pipelineSlug) { teams(first: 500) { edges { node { id, accessLevel, team { slug } } } } } }\"
}" | jq -c '.data.pipeline.teams.edges[].node' | \
        while read -r node; do \
            team_slug="$(echo ${node} | jq -r '.team.slug')"
            if [[ "${team_slug}" == "everyone" ]] ; then
                continue
            fi
            id="$(echo ${node} | jq -r '.id')"
            access_level="$(echo ${node} | jq -r '.accessLevel')"
            team_name="$(echo ${team_slug} | tr '.+-' '___')"
            pipeline_name="$(echo ${pipeline} | tr '.+-' '___')"

            echo "terraform import buildkite_team_pipeline.${team_name}_${pipeline_name} ${id}"

            cat <<EOF >> "${FILE}"
resource "buildkite_team_pipeline" "${team_name}_${pipeline_name}" {
  team_id       = "\${buildkite_team.${team_name}.team_id}"
  pipeline_slug = "\${buildkite_pipeline.${pipeline_name}.slug}"
  access_level  = "${access_level}"
}
EOF
        done
}

main() {
    local api_key org_slug

    if [[ $# -lt 2 ]] ; then
        echo "You need to provide an api_key and org_slug:
$0 <api-key> <org-slug>
" 1>&2
        exit 1
    fi

    api_key="$1"
    org_slug="$2"

    rm -f "${FILE}"

    curl -s \
        -H 'content-type: application/json' \
        -H "Authorization: Bearer $api_key" \
        "https://api.buildkite.com/v2/organizations/${org_slug}/pipelines?per_page=100" \
         | jq -cr '.[]' | \
        while read -r node; do \

            slug="$(echo ${node} | jq -r '.slug')"
            tf_name="$(echo ${node} | jq -r '.slug' | tr '.+-' '___')"

            echo "terraform import buildkite_pipeline.${tf_name} ${slug}"

            cat <<EOF >> "${FILE}"
resource "buildkite_pipeline" "${tf_name}" {
  name            = "$(echo ${node} | jq -r '.name')"
$(print_attr "${node}" description)
  default_branch  = "$(echo ${node} | jq -r '.default_branch')"
  repository      = "$(echo ${node} | jq -r '.repository')"
$(print_github_settings "${node}")

$(echo ${node} | jq -c '.steps[]' | while read -r step; do print_step "$step"; done)
}
EOF

            get_team_pipelines "$api_key" "$org_slug" "$slug"

            echo "" >> "${FILE}"
        done
}

main "$@"
