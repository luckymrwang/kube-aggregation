#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Create an environment that manages an multiple clusters of different versions

cases="${1}"

source "$(dirname "${BASH_SOURCE[0]}")/../helper.sh"

releases=(
    v1.23.5
    v1.22.8
    v1.21.11
    v1.20.15
    v1.19.16
    v1.18.20
    v1.17.17
    v1.16.15
    v1.15.12
    v1.14.10
    v1.13.12
    v1.12.10
    v1.11.10
    v1.10.13
)

function cleanup() {
    "${ROOT}/hack/clean-clusterconfigs.sh" >/dev/null 2>&1
    for release in "${releases[@]}"; do
        delete_data_plane "data-${release//./-}" >/dev/null 2>&1
    done
    delete_control_plane control-v1-23 >/dev/null 2>&1
}
trap cleanup EXIT

create_control_plane control-v1-23 v1.23.4 || {
    echo "Failed to create control plane"
    exit 1
}
for release in "${releases[@]}"; do
    create_data_plane "data-${release//./-}" "${release}" || {
        echo "Failed to create data plane"
        exit 1
    }
done
"${ROOT}/hack/gen-clusterconfigs.sh"

"${cases}"
