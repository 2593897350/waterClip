#
# 常用开发命令
#
.PHONY: help dev test

help:
	@echo "可用命令："
	@echo "  make dev   - 启动 Docker Compose 开发环境"
	@echo "  make test  - 运行前端、Go API、Python processor 测试"

# 启动本地开发环境
dev:
	docker compose up --build

# 运行项目测试
test:
	pnpm test:web
	cd api && GOCACHE=/Users/liuyafeng/workspace/waterClip/.cache/go-build go test ./...
	cd processor && ../.venv/bin/pytest
