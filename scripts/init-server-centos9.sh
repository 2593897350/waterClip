#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="/opt/waterClip"
REPO_URL=""
WEB_PORT="3000"
SKIP_FIREWALL="0"
SKIP_DEPLOY="0"
DRY_RUN="0"

usage() {
  cat <<'EOF'
CentOS Stream 9 服务器初始化脚本

用法：
  bash scripts/init-server-centos9.sh --repo-url <仓库地址> [选项]

选项：
  --repo-url <url>       Git 仓库地址，例如 git@github.com:2593897350/waterClip.git
  --install-dir <path>   项目部署目录，默认 /opt/waterClip
  --web-port <port>      对外暴露的 Web 端口，默认 3000
  --skip-firewall        跳过 firewalld 端口放行
  --skip-deploy          只初始化环境，不执行 deploy.sh
  --dry-run              只打印将要执行的命令，不真正执行
  --help                 显示帮助
EOF
}

run() {
  if [[ "$DRY_RUN" == "1" ]]; then
    printf '[dry-run] %s\n' "$*"
    return
  fi

  "$@"
}

run_shell() {
  if [[ "$DRY_RUN" == "1" ]]; then
    printf '[dry-run] %s\n' "$*"
    return
  fi

  bash -lc "$*"
}

require_root() {
  if [[ "${EUID}" -ne 0 ]]; then
    echo "请使用 root 用户执行该脚本。" >&2
    exit 1
  fi
}

verify_os() {
  if [[ ! -f /etc/os-release ]]; then
    echo "无法识别当前系统：缺少 /etc/os-release" >&2
    exit 1
  fi

  # shellcheck disable=SC1091
  source /etc/os-release

  if [[ "${ID:-}" != "centos" || "${VERSION_ID:-}" != "9" ]]; then
    echo "该脚本仅支持 CentOS Stream 9，当前系统是：${PRETTY_NAME:-unknown}" >&2
    exit 1
  fi
}

install_base_packages() {
  run dnf install -y git curl ca-certificates dnf-plugins-core
}

install_docker() {
  run dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
  run dnf install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
  run systemctl enable --now docker
}

open_firewall_port() {
  if [[ "$SKIP_FIREWALL" == "1" ]]; then
    return
  fi

  if ! command -v firewall-cmd >/dev/null 2>&1; then
    echo "未检测到 firewall-cmd，跳过防火墙端口配置。" >&2
    return
  fi

  run firewall-cmd --permanent --add-port="${WEB_PORT}/tcp"
  run firewall-cmd --reload
}

clone_or_update_repo() {
  if [[ -d "$INSTALL_DIR/.git" ]]; then
    run_shell "cd '$INSTALL_DIR' && git fetch --all --prune"
    run_shell "cd '$INSTALL_DIR' && git checkout main"
    run_shell "cd '$INSTALL_DIR' && git pull --ff-only origin main"
    return
  fi

  if [[ -z "$REPO_URL" ]]; then
    echo "目标目录不存在仓库时，必须通过 --repo-url 提供 Git 地址。" >&2
    exit 1
  fi

  run mkdir -p "$(dirname "$INSTALL_DIR")"
  run git clone "$REPO_URL" "$INSTALL_DIR"
}

prepare_env_file() {
  local env_file="$INSTALL_DIR/.env.production"
  local example_file="$INSTALL_DIR/.env.production.example"

  if [[ ! -f "$env_file" ]]; then
    run cp "$example_file" "$env_file"
  fi

  if [[ "$DRY_RUN" == "1" ]]; then
    printf '[dry-run] 更新 %s 中的 WEB_PORT=%s\n' "$env_file" "$WEB_PORT"
    return
  fi

  python3 - <<PY
from pathlib import Path

env_path = Path("$env_file")
web_port = "$WEB_PORT"
lines = env_path.read_text().splitlines()
updated = []
found = False
for line in lines:
    if line.startswith("WEB_PORT="):
        updated.append(f"WEB_PORT={web_port}")
        found = True
    else:
        updated.append(line)
if not found:
    updated.append(f"WEB_PORT={web_port}")
env_path.write_text("\n".join(updated) + "\n")
PY
}

prepare_runtime_dirs() {
  run mkdir -p "$INSTALL_DIR/var/uploads" "$INSTALL_DIR/var/masks" "$INSTALL_DIR/var/results"
}

deploy_project() {
  if [[ "$SKIP_DEPLOY" == "1" ]]; then
    return
  fi

  run_shell "cd '$INSTALL_DIR' && bash scripts/deploy.sh deploy"
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo-url)
      REPO_URL="${2:-}"
      shift 2
      ;;
    --install-dir)
      INSTALL_DIR="${2:-}"
      shift 2
      ;;
    --web-port)
      WEB_PORT="${2:-}"
      shift 2
      ;;
    --skip-firewall)
      SKIP_FIREWALL="1"
      shift
      ;;
    --skip-deploy)
      SKIP_DEPLOY="1"
      shift
      ;;
    --dry-run)
      DRY_RUN="1"
      shift
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      echo "未知参数：$1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

require_root
verify_os
install_base_packages
install_docker
open_firewall_port
clone_or_update_repo
prepare_env_file
prepare_runtime_dirs
deploy_project

if [[ "$DRY_RUN" == "1" ]]; then
  echo "dry-run 完成。"
else
  echo "初始化完成。"
fi
