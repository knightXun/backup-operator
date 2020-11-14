#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

ROOT=$(unset CDPATH && cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd)
cd $ROOT

source "${ROOT}/hack/lib.sh"

backup_target="$ROOT/manifests/backup.yaml"
restore_target="$ROOT/manifests/restore.yaml"

# Ensure that we find the binaries we build before anything else.
export GOBIN="${OUTPUT_BIN}"
PATH="${GOBIN}:${PATH}"

# Enable go modules explicitly.
export GO111MODULE=on
go install github.com/backup-operator/cmd/to-crdgen

to-crdgen generate restore > $restore_target
to-crdgen generate backup > $backup_target