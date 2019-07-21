#!/bin/bash -e

cleanup() {
	set +e
	destroy
}

trap cleanup EXIT

. hack/release/release-prepare.sh

export RELEASE_BRANCH=release-$(echo ${PLUGIN_VERSION_NEXT_MAJOR} | cut -d"." -f1)

log "starting major version ${PLUGIN_VERSION_NEXT_MAJOR} release"
cd ${PLUGIN_REPO}

log "creating and pushing ${RELEASE_BRANCH} branch to upstream"
git checkout -b ${RELEASE_BRANCH} &>> ${LOG_FILE}
#git push upstream ${RELEASE_BRANCH} &>> ${LOG_FILE} || err "could not push ${RELEASE_BRANCH} to github.com/${GITHUB_ORGANIZATION}/openshift-azure"
git checkout master &>> ${LOG_FILE}

log "prepare branch ${RELEASE_BRANCH}-pluginconfig for major release"
git fetch --all &>> ${LOG_FILE}
git checkout -b ${RELEASE_BRANCH}-pluginconfig &>> ${LOG_FILE}

log "updating pluginconfig with new image ${IMAGE_RESOURCENAME} and major version ${PLUGIN_VERSION_NEXT_MAJOR}"
# TODO: replace imageVersion in pluginconfig/pluginconfig-311.yaml with ${IMAGE_RESOURCENAME} (sed)

log "generating changelog between ${PLUGIN_VERSION_CURRENT} and master (${PLUGIN_COMMIT})"
go run cmd/releasenotes/releasenotes -start-sha=${PLUGIN_VERSION_CURRENT} -end-sha=master -release-version=${PLUGIN_VERSION_NEXT_MAJOR} -output-file=CHANGELOG.md || err "could not generate release notes"
## TODO: update image version to point to specific version instead of latest tag.

log "creating pull request for major release ${PLUGIN_VERSION_NEXT_MAJOR}"
git add --all &>> ${LOG_FILE}
git commit -m "Release ARO plugin config ${PLUGIN_VERSION_NEXT_MAJOR}" &>> ${LOG_FILE}
#git push -u origin ${RELEASE_BRANCH}-pluginconfig &>> ${LOG_FILE} || err "could not push ${RELEASE_BRANCH}-pluginconfig to github.com/${GIT_AUTHOR_NAME}/openshift-azure"
# TODO: create pr to plugin repo with robot token


