#!/usr/bin/env bash
set -euo pipefail

OUTPUT="$(bash scripts/init-server-centos9.sh --help)"

[[ "$OUTPUT" == *"CentOS Stream 9"* ]]
[[ "$OUTPUT" == *"--repo-url"* ]]
[[ "$OUTPUT" == *"--install-dir"* ]]
[[ "$OUTPUT" == *"--web-port"* ]]
