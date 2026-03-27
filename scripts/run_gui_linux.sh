#!/usr/bin/env bash
set -euo pipefail

TAGS="dev"
if pkg-config --exists webkit2gtk-4.1; then
  TAGS="dev webkit2_41"
fi

echo "Running GUI with tags: ${TAGS}"
go run -tags "${TAGS}" ./cmd/slay_the_spire_tool
