#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT=$(unset CDPATH && cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd)
cd $ROOT

source "${ROOT}/hack/lib.sh"

# Ensure that we find the binaries we build before anything else.
export GOBIN="${OUTPUT_BIN}"
PATH="${GOBIN}:${PATH}"

# Enable go modules explicitly.
export GO111MODULE=on
go install k8s.io/code-generator/cmd/openapi-gen

openapi-gen --go-header-file=./hack/boilerplate.go.txt \
    -i github.com/vesoft-inc-private/nebula-operator/pkg/apis/nebula/v1alpha1,k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/api/core/v1 \
    -p apis/nebula/v1alpha1 -O openapi_generated -o ./pkg
