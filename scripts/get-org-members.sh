#!/bin/bash

FILE=org_members.tf

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
        'https://graphql.buildkite.com/v1' \
        -d "
{
  \"operationName\": \"OrganizationMembers\",
  \"variables\": { \"orgSlug\": \"$org_slug\" },
  \"query\": \"query OrganizationMembers(\$orgSlug: ID!) { organization(slug: \$orgSlug) { members(first: 500) { edges { node { uuid, role, user { id, email } } } } } }\"
}" | jq -cr '.data.organization.members.edges[].node ' | \
        while read -r node; do \
            uuid="$(echo "$node" | jq -r '.uuid')"
            role="$(echo "$node" | jq -r '.role')"
            user_id="$(echo "$node" | jq -r '.user.id')"
            name="$(echo "$node" | jq -r '.user.email' | cut -f1 -d@ | tr '.+-' '__-')"

            echo "terraform import buildkite_org_member.${name} ${uuid}"

            cat <<EOF >> "${FILE}"
resource "buildkite_org_member" "${name}" {
  role = "${role}"
}
EOF
        done
}

main "$@"
