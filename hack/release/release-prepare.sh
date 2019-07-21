#!/bin/bash -e

export SANDBOX=${SANDBOX:-$(mktemp -d)}
export GITHUB_ORGANIZATION=${GITHUB_ORGANIZATION:-"openshift"}
export PLUGIN_REPO=${PLUGIN_REPO:-"${SANDBOX}/src/github.com/${GITHUB_ORGANIZATION}/openshift-azure"}
export TESTS_REPO=${TESTS_REPO:-"${SANDBOX}/src/github.com/${GITHUB_ORGANIZATION}/release"}

export GIT_AUTHOR_NAME=openshift-azure-robot
export GIT_AUTHOR_EMAIL=aos-azure@redhat.com
export GIT_COMMITTER_NAME=openshift-azure-robot
export GIT_COMMITTER_EMAIL=aos-azure@redhat.com

export IMAGE_RESOURCENAME="${IMAGE_RESOURCENAME:-rhel7-3.11-$(TZ=Etc/UTC date +%Y%m%d%H%M)}"

export LOG_FILE=${PWD}/aro-release.log

# git log -n 1 --pretty=format:"%H" origin/release-v6 get commit for branch
# git log -n 1 --tags --pretty=format:"%H" v5.2.1 get commit for tag
log() {
	echo "$@"
	logger -t "aro-release" "$@"
}

err() {
	echo "$@" >&2
	logger -t "aro-release" "$@"
	exit 1
}

plugin_env() {
	export PLUGIN_COMMIT=${PLUGIN_COMMIT:-$(git rev-parse --short HEAD)}
	export PLUGIN_VERSION_CURRENT=${PLUGIN_VERSION_CURRENT:-$(git tag --list | tail -1)}
	export PLUGIN_VERSION_NEXT_PATCH=${PLUGIN_VERSION_NEXT_PATCH:-$(next_version "patch" ${PLUGIN_VERSION_CURRENT})}
	export PLUGIN_VERSION_NEXT_MINOR=${PLUGIN_VERSION_NEXT_MINOR:-$(next_version "minor" ${PLUGIN_VERSION_CURRENT})}
	export PLUGIN_VERSION_NEXT_MAJOR=${PLUGIN_VERSION_NEXT_MAJOR:-$(next_version "major" ${PLUGIN_VERSION_CURRENT})}
}

next_version() {
	local KIND=${1}
	local CURRENT_VERSION=$(echo ${2} | cut -d'v' -f2)
	IFS='.' read -r -a version_array <<< ${CURRENT_VERSION}
	local next_patch=$(echo $(( ${version_array[2]} + 1 )))
	local next_minor=$(echo $(( ${version_array[1]} + 1 )))
	local next_major=$(echo $(( ${version_array[0]} + 1 )))
	[[ ${KIND} == "patch" ]] && (echo "v${version_array[0]}.${version_array[1]}.${next_patch}" && return)
	[[ ${KIND} == "minor" ]] && (echo "v${version_array[0]}.${next_minor}" && return)
	[[ ${KIND} == "major" ]] && (echo "v${next_major}.0" && return)
}

fetch_source() {
	mkdir -p ${PLUGIN_REPO} ${TESTS_REPO}

	log "cloning github.com/openshift/openshift-azure to ${PLUGIN_REPO}"
	git clone https://openshift-azure-robot:${GITHUB_TOKEN}@github.com/${GITHUB_ORGANIZATION}/openshift-azure.git ${PLUGIN_REPO} &>> ${LOG_FILE} || (err "could not fetch github.com/openshift/openshift-azure")
	log "cloning github.com/openshift/release to ${TESTS_REPO}"
	git clone https://openshift-azure-robot:${GITHUB_TOKEN}@github.com/${GITHUB_ORGANIZATION}/release.git ${TESTS_REPO} &>> ${LOG_FILE} || (err "could not fetch github.com/openshift/release")


	export GOPATH=${SANDBOX}
}

destroy() {
	log "cleaning up release resources from filesystem"
#	rm -rf ${SANDBOX}
}

log "preparing ARO plugin release"
fetch_source
ln -sf ${PWD}/secrets ${PLUGIN_REPO}
cd ${PLUGIN_REPO}
plugin_env

log "GOPATH set to ${GOPATH}"
log "current plugin version in production is ${PLUGIN_VERSION_CURRENT}"

log "checking out commit ${PLUGIN_COMMIT} on github.com/openshift/openshift-azure"
git remote rename origin upstream
git remote add origin https://openshift-azure-robot:${GITHUB_TOKEN}@github.com/${GIT_AUTHOR_NAME}/openshift-azure.git
git checkout -q ${PLUGIN_COMMIT}

#log "running make unit on commit ${PLUGIN_COMMIT}"
#make unit &>> ${LOG_FILE} || err "plugin commit ${PLUGIN_COMMIT} failed make unit"
#
#log "running make verify on commit ${PLUGIN_COMMIT}"
#make verify &>> ${LOG_FILE} || err "plugin commit ${PLUGIN_COMMIT} failed make verify"

log "sourcing secrets for azure authorization"
. ./secrets/secret

#log "running make vmimage for ${IMAGE_RESOURCENAME} on commit ${PLUGIN_COMMIT}"
#make vmimage &>> ${LOG_FILE} || err "vm image build ${IMAGE_RESOURCENAME} failed"
#
#log "generating sas url for new vm image ${IMAGE_RESOURCENAME}"
#hack/vmimage-cloudpartner.sh ${IMAGE_RESOURCENAME} || err "could not fetch sas url for vm image ${IMAGE_RESOURCENAME}"
#
#log "publishing vm image ${IMAGE_RESOURCENAME} on azure partner portal (not implemented)"
## TODO(kenny): how can we fix the publish image script so we can plug it in here?
