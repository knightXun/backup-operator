#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

echo "$(dirname "${BASH_SOURCE[0]}")/../../../.."

bash "${CODEGEN_PKG}/generate-groups.sh" all \
  github.com/backup-operator/pkg/client \
  github.com/backup-operator/pkg/apis \
  "mydumper:v1alpha1" \
  --output-base "$(dirname "${BASH_SOURCE[0]}")/../../.." \
  --go-header-file "${SCRIPT_ROOT}/hack/boilerplate.go.txt"
