#!/usr/bin/env bash
set -euo pipefail

grep -F 'context: .' docker-compose.yml >/dev/null
grep -F 'dockerfile: web/Dockerfile' docker-compose.yml >/dev/null
grep -F 'COPY web/package.json package.json' web/Dockerfile >/dev/null
grep -F 'COPY pnpm-lock.yaml pnpm-lock.yaml' web/Dockerfile >/dev/null
