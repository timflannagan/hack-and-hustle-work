#! /bin/bash

set -e

if [[ -z "$1" ]]; then
    echo "Invalid configuration: you need to provide the quay tag parameter as an argument to this script"
    exit 1
fi

export IMAGE_TAG="$1"
export IMAGE_ORG="${IMAGE_ORG:="tflannag"}"
export IMAGE_REPO="${IMAGE_REPO:="origin-metering-ansible-operator"}"
export DOCKER_RUNTIME="${DOCKER_RUNTIME:="docker"}"
export DOCKER_BUILD_CMD="${DOCKER_BUILD_CMD:="${DOCKER_RUNTIME} build"}"

"$DOCKER_RUNTIME" build -f Dockerfile.metering-ansible-operator.okd -t quay.io/"$IMAGE_ORG"/"$IMAGE_REPO":"$IMAGE_TAG" .
"$DOCKER_RUNTIME" push quay.io/"$IMAGE_ORG"/"$IMAGE_REPO":"$IMAGE_TAG"
