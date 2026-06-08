#!/usr/bin/env bash
set -euo pipefail

grep -F 'FROM node:20-bookworm-slim' web/Dockerfile >/dev/null
grep -F 'FROM golang:1.22-bookworm' api/Dockerfile >/dev/null
grep -F 'apt-get install -y --no-install-recommends ca-certificates' processor/Dockerfile >/dev/null
grep -F 'apt-get install -y --no-install-recommends ca-certificates' web/Dockerfile >/dev/null
grep -F 'apt-get install -y --no-install-recommends ca-certificates' api/Dockerfile >/dev/null
