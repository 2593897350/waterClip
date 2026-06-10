#!/usr/bin/env bash
set -euo pipefail

grep -F 'FROM node:20-bookworm-slim' web/Dockerfile >/dev/null
grep -F 'FROM golang:1.22-bookworm' api/Dockerfile >/dev/null
grep -F 'apt-get install -y --no-install-recommends ca-certificates' processor/Dockerfile >/dev/null
grep -F 'apt-get install -y --no-install-recommends ca-certificates' web/Dockerfile >/dev/null
grep -F 'apt-get install -y --no-install-recommends ca-certificates' api/Dockerfile >/dev/null
grep -F 'ARG NPM_REGISTRY=' web/Dockerfile >/dev/null
grep -F 'pnpm config set registry "$NPM_REGISTRY"' web/Dockerfile >/dev/null
grep -F 'ARG GOPROXY=' api/Dockerfile >/dev/null
grep -F 'go env -w GOPROXY="$GOPROXY"' api/Dockerfile >/dev/null
grep -F 'ARG PIP_INDEX_URL=' processor/Dockerfile >/dev/null
grep -F 'ARG PIP_TRUSTED_HOST=' processor/Dockerfile >/dev/null
grep -F 'pip install --index-url "$PIP_INDEX_URL"' processor/Dockerfile >/dev/null
grep -F 'NPM_REGISTRY: ${NPM_REGISTRY:-https://registry.npmmirror.com}' docker-compose.yml >/dev/null
grep -F 'GOPROXY: ${GOPROXY:-https://goproxy.cn,direct}' docker-compose.yml >/dev/null
grep -F 'PIP_INDEX_URL: ${PIP_INDEX_URL:-https://pypi.tuna.tsinghua.edu.cn/simple}' docker-compose.yml >/dev/null
grep -F 'PIP_TRUSTED_HOST: ${PIP_TRUSTED_HOST:-pypi.tuna.tsinghua.edu.cn}' docker-compose.yml >/dev/null
grep -F 'NPM_REGISTRY=https://registry.npmmirror.com' .env.production.example >/dev/null
grep -F 'GOPROXY=https://goproxy.cn,direct' .env.production.example >/dev/null
grep -F 'GOSUMDB=sum.golang.google.cn' .env.production.example >/dev/null
grep -F 'PIP_INDEX_URL=https://pypi.tuna.tsinghua.edu.cn/simple' .env.production.example >/dev/null
grep -F 'PIP_TRUSTED_HOST=pypi.tuna.tsinghua.edu.cn' .env.production.example >/dev/null
