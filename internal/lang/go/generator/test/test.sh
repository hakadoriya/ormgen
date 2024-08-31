#!/usr/bin/env bash

REPO_ROOT=$(git rev-parse --show-toplevel)

if ! cd "${REPO_ROOT:?}"; then
  echo "Failed to cd to REPO_ROOT"
  exit 1
fi

set -x
go run ./cmd/ormgen generate ./internal/lang/go/source/test --go-orm-output-path ./internal/lang/go/generator/test/generated
