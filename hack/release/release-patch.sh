#!/bin/bash -e

cleanup() {
	set +e
	destroy
}

trap cleanup EXIT

. hack/release/release-prepare.sh

export RELEASE_BRANCH=release-$(echo ${PLUGIN_VERSION_NEXT_PATCH} | cut -d"." -f1)

log "starting patch version ${PLUGIN_VERSION_NEXT_PATCH} release"
