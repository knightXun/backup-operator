#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

ROOT=$(unset CDPATH && cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd)
cd $ROOT

source "${ROOT}/hack/lib.sh"

crd_target="$ROOT/manifests/crd.yaml"

# Ensure that we find the binaries we build before anything else.
export GOBIN="${OUTPUT_BIN}"
PATH="${GOBIN}:${PATH}"

# Enable go modules explicitly.
export GO111MODULE=on
go install github.com/backup-operator/cmd/to-crdgen

to-crdgen generate restore > $crd_target
to-crdgen generate backup >> $crd_target

hack::ensure_gen_crd_api_references_docs

DOCS_PATH="$ROOT/docs/api-references"

${DOCS_BIN} \
-config "$DOCS_PATH/config.json" \
-template-dir "$DOCS_PATH/template" \
-api-dir "github.com/backup-operator/apis/mydumper/v1alpha1" \
-out-file "$DOCS_PATH/docs.md"
