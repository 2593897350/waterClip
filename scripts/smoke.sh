#!/usr/bin/env bash
set -euo pipefail

test -f package.json
test -f Makefile
test -f docker-compose.yml
test -d web
test -d api
test -d processor
