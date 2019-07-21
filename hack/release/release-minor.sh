#!/bin/bash -e

cleanup() {
	set +e
	destroy
}

trap cleanup EXIT

. hack/release/release-prepare.sh

export RELEASE_BRANCH=release-$(echo ${PLUGIN_VERSION_NEXT_MINOR} | cut -d"." -f1)

log "starting minor version ${PLUGIN_VERSION_NEXT_MINOR} release"


