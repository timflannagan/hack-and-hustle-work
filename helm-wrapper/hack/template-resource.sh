#! /bin/bash

set -eou pipefail

VALUES_ARGS="${VALUES_ARGS:-""}"
DEBUG=${DEBUG:=false}

if [[ $DEBUG == "true" ]]; then
    set -x
fi

${HELM_BIN} template "$CHART" \
    $VALUES_ARGS \
    -s "$TEMPLATE_RESOURCE_NAME" |
    sed -f "$ROOT_DIR/hack/remove-helm-template-header.sed" \
        > "$OUTPUT_DIR"
