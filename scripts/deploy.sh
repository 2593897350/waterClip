#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
ENV_FILE="${ENV_FILE:-$ROOT_DIR/.env.production}"
COMMAND="${1:-deploy}"

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

ensure_env_file() {
  if [[ -f "$ENV_FILE" ]]; then
    return
  fi

  if [[ -f "$ROOT_DIR/.env.production.example" ]]; then
    cp "$ROOT_DIR/.env.production.example" "$ENV_FILE"
    echo "Created $ENV_FILE from template. Review it before re-running." >&2
    exit 1
  fi

  echo "Missing env file: $ENV_FILE" >&2
  exit 1
}

ensure_runtime_dirs() {
  mkdir -p "$ROOT_DIR/var/uploads" "$ROOT_DIR/var/masks" "$ROOT_DIR/var/results"
}

compose() {
  docker compose --env-file "$ENV_FILE" "$@"
}

require_command docker
ensure_env_file
ensure_runtime_dirs

case "$COMMAND" in
  deploy|up)
    compose up -d --build
    compose ps
    ;;
  down)
    compose down
    ;;
  restart)
    compose down
    compose up -d --build
    compose ps
    ;;
  logs)
    shift || true
    compose logs -f "$@"
    ;;
  status|ps)
    compose ps
    ;;
  *)
    echo "Usage: scripts/deploy.sh [deploy|up|down|restart|logs|status|ps]" >&2
    exit 1
    ;;
esac
