#!/usr/bin/env bash
set -euo pipefail

grep -F 'apt-get install -y --no-install-recommends ca-certificates' processor/Dockerfile >/dev/null
grep -F 'apk add --no-cache ca-certificates' web/Dockerfile >/dev/null
grep -F 'apk add --no-cache ca-certificates' api/Dockerfile >/dev/null
