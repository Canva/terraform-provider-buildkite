#!/bin/bash

main() {
    local api_key org_slug module file

    if [[ $# -lt 2 ]] ; then
        echo "You need to provide an api_key and org_slug:

$0 <api-key> <org-slug> [<module>]

" 1>&2
        exit 1
    fi

    api_key="$1"
    org_slug="$2"
    module="${3-buildkite}"
    file="terraform/${module}/org_teams.tf"
    resource_prefix="module.${module}."

    mkdir -p "$(dirname "$file")"

    rm -f "${file}"

    curl -s \
        -H 'content-type: application/json' \
        -H "Authorization: Bearer $api_key" \
        'https://graphql.buildkite.com/v1' \
        -d "
{
  \"operationName\": \"Teams\",
  \"variables\": { \"orgSlug\": \"$org_slug\" },
  \"query\": \"query Teams(\$orgSlug: ID!) { organization(slug: \$orgSlug) { teams(first: 500) { edges { node { slug, name, members(first: 500) { edges { node { id, role, user { email } } } } } } } } }\"
}" | jq -c '.data.organization.teams.edges[].node' | \
        while read -r node; do \
            slug="$(echo ${node} | jq -r '.slug')"
            if [[ "${slug}" == "everyone" ]] ; then
                continue
            fi

            echo "// ${slug}" >> "${file}"
            tf_name="$(echo ${node} | jq -r '.name' | tr '.+-' '___')"
            name="$(echo ${node} | jq -r '.name')"

            echo "terraform import ${resource_prefix}buildkite_team.${tf_name} ${slug}"

            cat <<EOF >> "${file}"
resource "buildkite_team" "${tf_name}" {
  name = "${name}"
}
EOF

            for member in $(echo "${node}" | jq -cr '.members.edges[].node'); do
                member_id="$(echo "${member}" | jq -r '.id')"
                role="$(echo "${member}" | jq -r '.role')"
                user_name="$(echo "${member}" | jq -r '.user.email' | cut -f1 -d@ | tr '.+-' '___')"

                echo "terraform import ${resource_prefix}buildkite_team_member.${user_name}_${tf_name} ${member_id}"

                cat <<EOF >> "${file}"
resource "buildkite_team_member" "${user_name}_${tf_name}" {
  user_id = "\${buildkite_org_member.${user_name}.user_id}"
  team_id = "\${buildkite_team.${tf_name}.team_id}"
  role    = "${role}"
}
EOF

            done
            echo "" >> "${file}"
        done

    terraform fmt "${file}" > /dev/null
}

main "$@"
