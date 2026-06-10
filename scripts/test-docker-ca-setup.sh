#!/usr/bin/env bash
set -euo pipefail

grep -F 'FROM node:20-bookworm-slim' web/Dockerfile >/dev/null
grep -F 'FROM golang:1.22-bookworm' api/Dockerfile >/dev/null
grep -F 'apt-get install -y --no-install-recommends ca-certificates' web/Dockerfile >/dev/null
grep -F 'apt-get install -y --no-install-recommends ca-certificates' api/Dockerfile >/dev/null
if grep -F 'pip install' processor/Dockerfile >/dev/null; then
  echo "processor Dockerfile should not install Python packages during image build" >&2
  exit 1
fi
grep -F 'CMD ["python", "-m", "app.main"]' processor/Dockerfile >/dev/null
grep -F 'ARG NPM_REGISTRY=' web/Dockerfile >/dev/null
grep -F 'ARG NPM_STRICT_SSL=' web/Dockerfile >/dev/null
grep -F 'npm config set strict-ssl "$NPM_STRICT_SSL"' web/Dockerfile >/dev/null
if grep -F 'pnpm' web/Dockerfile >/dev/null; then
  echo "web Dockerfile should not depend on pnpm" >&2
  exit 1
fi
grep -F 'npm install' web/Dockerfile >/dev/null
grep -F 'RUN npm run build' web/Dockerfile >/dev/null
grep -F 'CMD ["npm", "run", "start", "--", "--hostname", "0.0.0.0", "--port", "3000"]' web/Dockerfile >/dev/null
grep -F 'ARG GOPROXY=' api/Dockerfile >/dev/null
grep -F 'go env -w GOPROXY="$GOPROXY"' api/Dockerfile >/dev/null
grep -F 'NPM_REGISTRY: ${NPM_REGISTRY:-https://registry.npmmirror.com}' docker-compose.yml >/dev/null
grep -F 'NPM_STRICT_SSL: ${NPM_STRICT_SSL:-false}' docker-compose.yml >/dev/null
grep -F 'GOPROXY: ${GOPROXY:-https://goproxy.cn,direct}' docker-compose.yml >/dev/null
grep -F 'NPM_REGISTRY=https://registry.npmmirror.com' .env.production.example >/dev/null
grep -F 'NPM_STRICT_SSL=false' .env.production.example >/dev/null
grep -F 'GOPROXY=https://goproxy.cn,direct' .env.production.example >/dev/null
grep -F 'GOSUMDB=sum.golang.google.cn' .env.production.example >/dev/null
if grep -F 'fastapi' processor/pyproject.toml >/dev/null; then
  echo "processor should not depend on fastapi" >&2
  exit 1
fi
if grep -F 'pydantic' processor/pyproject.toml >/dev/null; then
  echo "processor should not depend on pydantic" >&2
  exit 1
fi
