# WaterClip

一个面向照片和社交媒体图片的去水印网站 `MVP`。

## 目录结构

- `web/`：`Next.js` 前端，负责上传、编辑和结果展示
- `api/`：`Go` API 服务，负责任务编排、文件落盘和状态查询
- `processor/`：`Python` 图像处理服务，负责检测和去水印修复

## 本地开发

1. 创建 Python 虚拟环境：
   `python3 -m venv .venv`
2. 安装 Python 依赖：
   `./.venv/bin/pip install fastapi httpx pytest uvicorn`
3. 安装前端依赖：
   `pnpm install --dir web`
4. 启动图像处理服务：
   `cd processor && ../.venv/bin/uvicorn app.main:app --reload --port 8000`
5. 启动 API 服务：
   `cd api && GOCACHE=/Users/liuyafeng/workspace/waterClip/.cache/go-build go run ./cmd/server`
6. 启动前端：
   `cd web && pnpm dev`

## 测试命令

- 前端测试：
  `cd web && pnpm test`
- API 测试：
  `cd api && GOCACHE=/Users/liuyafeng/workspace/waterClip/.cache/go-build go test ./...`
- 图像处理服务测试：
  `cd processor && ../.venv/bin/pytest`

## 云服务器部署

1. 复制生产环境模板：
   `cp .env.production.example .env.production`
2. 按需修改 `.env.production`，例如 `WEB_PORT`
3. 执行部署：
   `bash scripts/deploy.sh deploy`
4. 查看服务状态：
   `bash scripts/deploy.sh status`
5. 查看日志：
   `bash scripts/deploy.sh logs web`

### CentOS Stream 9 一键初始化

如果你的服务器是 `CentOS Stream 9`，可以直接使用：

```bash
bash scripts/init-server-centos9.sh --repo-url git@github.com:2593897350/waterClip.git
```

常用参数：

- `--install-dir /opt/waterClip`：指定部署目录
- `--web-port 3000`：指定对外暴露端口
- `--skip-firewall`：跳过 `firewalld` 端口放行
- `--skip-deploy`：只安装环境和拉代码，不执行部署
- `--dry-run`：只打印命令，不真正执行

示例：

```bash
bash scripts/init-server-centos9.sh \
  --repo-url git@github.com:2593897350/waterClip.git \
  --install-dir /opt/waterClip \
  --web-port 3000
```

完整部署说明见：

- [production-deploy-centos9.md](/Users/liuyafeng/workspace/waterClip/docs/production-deploy-centos9.md)

## 生产环境说明

- 浏览器只访问 `web` 容器暴露出来的端口
- `web` 会把 `/api/*` 代理到内部 `api` 容器
- `api` 再通过 `PROCESSOR_BASE_URL` 调用 `processor`
- 公网只需要暴露 Web 端口，不需要把 `api` 和 `processor` 直接暴露出去
- 在 `CentOS Stream 9` / `SELinux Enforcing` 环境下，`docker-compose.yml` 已为 `var/` 挂载增加 `:Z`

## 说明

- 仓库里的设计文档已经是中文
- 从现在开始，新增的 `Markdown` 文档和 `Makefile` 注释默认优先使用中文
