#! /bin/bash

set -eou pipefail

kustomize build "${KUSTOMIZE_OVERLAY_PATH}" > "${OUTPUT_PATH}"
