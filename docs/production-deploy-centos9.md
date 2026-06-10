# CentOS Stream 9 云服务器完整部署指南

本文档基于当前仓库的现状整理，目标是让你在一台 `CentOS Stream 9` 云服务器上把 `WaterClip` 以 `Docker Compose + Nginx + HTTPS` 的方式跑起来。

## 1. 部署结构

当前项目的生产部署结构如下：

- `web`：`Next.js` 前端，对外提供页面访问
- `api`：`Go` API 服务，处理任务创建、状态查询和文件管理
- `processor`：`Python` 图像处理服务，负责检测与修复
- `Nginx`：宿主机反向代理，统一对外暴露 `80/443`

请求路径：

`浏览器 -> Nginx -> web -> /api/* 反代到 api -> api 调用 processor`

## 2. 服务器要求

建议最低配置：

- `2 vCPU`
- `4 GB RAM`
- `40 GB SSD`

建议提前准备：

- 已绑定公网 IP 的云服务器
- 一个已经解析到该公网 IP 的域名
- 已经加入 GitHub 账号的 SSH key

## 3. 安全组 / 防火墙

云厂商安全组至少放行：

- `22/tcp`
- `80/tcp`
- `443/tcp`

如果只是临时联调，也可以额外放行：

- `3000/tcp`

## 4. 使用初始化脚本

仓库中已经提供了 CentOS Stream 9 初始化脚本：

- [scripts/init-server-centos9.sh](/Users/liuyafeng/workspace/waterClip/scripts/init-server-centos9.sh)

登录服务器后执行：

```bash
bash scripts/init-server-centos9.sh \
  --repo-url git@github.com:2593897350/waterClip.git \
  --install-dir /opt/waterClip \
  --web-port 3000
```

这个脚本会完成：

- 校验系统版本是否为 `CentOS Stream 9`
- 安装 `git / curl / docker / docker compose`
- 启动 Docker 并设置开机自启
- 放行 Web 端口
- 克隆或更新项目代码
- 生成 `.env.production`
- 创建运行时目录
- 执行 `bash scripts/deploy.sh deploy`

## 5. 环境文件

部署前会生成：

- `.env.production`

当前默认内容来自：

- [.env.production.example](/Users/liuyafeng/workspace/waterClip/.env.production.example)

默认值：

```env
WEB_PORT=3000
API_ADDRESS=:8080
INTERNAL_API_PROXY_TARGET=http://api:8080
PROCESSOR_BASE_URL=http://processor:8000
NPM_REGISTRY=https://registry.npmmirror.com
NPM_STRICT_SSL=false
GOPROXY=https://goproxy.cn,direct
GOSUMDB=sum.golang.google.cn
PIP_INDEX_URL=https://pypi.org/simple
PIP_TRUSTED_HOSTS=pypi.org files.pythonhosted.org
```

说明：

- `WEB_PORT`：宿主机暴露的前端端口
- `API_ADDRESS`：`api` 容器内部监听地址
- `INTERNAL_API_PROXY_TARGET`：`web` 容器内访问 `api` 的地址
- `PROCESSOR_BASE_URL`：`api` 容器内访问 `processor` 的地址
- `NPM_REGISTRY`：前端构建时使用的 npm / pnpm 镜像源
- `NPM_STRICT_SSL`：前端依赖安装时是否严格校验证书
- `GOPROXY` / `GOSUMDB`：Go 依赖下载与校验源
- `PIP_INDEX_URL` / `PIP_TRUSTED_HOSTS`：Python 依赖下载索引与信任主机列表

说明：

- 当前 `docker-compose.yml` 已经内置了上述镜像源默认值
- 当前 `docker-compose.yml` 里 `NPM_STRICT_SSL` 默认是 `false`，用于规避部分云环境下镜像站证书链不完整的问题
- 当前 `processor` 默认走官方 `PyPI`，并信任 `pypi.org` 与 `files.pythonhosted.org`
- 即使你服务器上的 `.env.production` 是旧文件，只要里面没有手动写这些变量，也会自动使用镜像默认值
- 如果你之前手动配置过官方源，请改成上面的值再重新部署

## 6. Docker Compose 部署

项目当前部署脚本为：

- [scripts/deploy.sh](/Users/liuyafeng/workspace/waterClip/scripts/deploy.sh)

常用命令：

```bash
bash scripts/deploy.sh deploy
bash scripts/deploy.sh status
bash scripts/deploy.sh logs web
bash scripts/deploy.sh logs api
bash scripts/deploy.sh logs processor
bash scripts/deploy.sh restart
bash scripts/deploy.sh down
```

如果你怀疑 Docker 仍然复用了旧缓存，可以强制重建：

```bash
docker compose --env-file .env.production build --no-cache
bash scripts/deploy.sh deploy
```

## 7. SELinux 注意事项

`CentOS Stream 9` 默认大概率启用了 `SELinux`。

当前 [docker-compose.yml](/Users/liuyafeng/workspace/waterClip/docker-compose.yml) 已经把运行时目录挂载改成了：

```yml
- ./var:/app/var:Z
```

这样可以避免容器无法写入挂载目录的问题。

检查状态：

```bash
getenforce
```

如果输出是 `Enforcing`，当前 compose 文件已经按推荐方式处理好了，不需要再临时关闭 `SELinux`。

## 8. 加 Nginx 作为公网入口

正式上线建议不要直接把 `3000` 暴露给公网，而是用 `Nginx` 做反向代理。

安装：

```bash
dnf install -y nginx
systemctl enable --now nginx
```

仓库里已提供一个基础模板：

- [deploy/nginx/waterclip.conf](/Users/liuyafeng/workspace/waterClip/deploy/nginx/waterclip.conf)

复制到服务器：

```bash
cp /opt/waterClip/deploy/nginx/waterclip.conf /etc/nginx/conf.d/waterclip.conf
```

如果你已经有域名，把配置里的：

```nginx
server_name _;
```

改成你的域名，例如：

```nginx
server_name waterclip.example.com;
```

检查并重载：

```bash
nginx -t
systemctl reload nginx
```

## 9. 配 HTTPS

安装 `certbot`：

```bash
dnf install -y epel-release
dnf install -y certbot python3-certbot-nginx
```

申请证书：

```bash
certbot --nginx -d 你的域名
```

例如：

```bash
certbot --nginx -d waterclip.example.com
```

成功后验证：

```bash
systemctl status certbot-renew.timer
```

## 10. 首次访问验证

如果你还没接域名，先直接访问：

```text
http://服务器公网IP:3000
```

如果已经接了 Nginx 和域名：

```text
http://你的域名
https://你的域名
```

## 11. 更新代码流程

以后更新项目，推荐这样做：

```bash
cd /opt/waterClip
git pull
bash scripts/deploy.sh deploy
```

如果只想重启：

```bash
bash scripts/deploy.sh restart
```

## 12. 排错清单

### 12.1 页面打不开

检查：

```bash
bash scripts/deploy.sh status
ss -lntp | grep 3000
firewall-cmd --list-ports
```

### 12.2 容器启动失败

检查：

```bash
bash scripts/deploy.sh logs web
bash scripts/deploy.sh logs api
bash scripts/deploy.sh logs processor
```

如果报错集中在 `pip install`、`pnpm install`、`go mod download` 这类依赖下载阶段，优先检查：

```bash
grep -E 'NPM_REGISTRY|GOPROXY|GOSUMDB|PIP_INDEX_URL|PIP_TRUSTED_HOSTS' .env.production
docker compose --env-file .env.production config | sed -n '1,120p'
```

如果镜像源已经生效但仍然出现证书错误，说明这台服务器所在网络可能还存在自定义 CA 或 TLS 代理，需要额外导入云厂商证书链。

### 12.3 图片处理后没有结果

重点看：

```bash
bash scripts/deploy.sh logs api
bash scripts/deploy.sh logs processor
ls -R /opt/waterClip/var
```

### 12.4 Nginx 配置不生效

检查：

```bash
nginx -t
systemctl status nginx
cat /etc/nginx/conf.d/waterclip.conf
```

## 13. 当前生产限制

当前项目已经可以部署和运行，但还属于 `MVP` 阶段，正式运营前需要明确这些限制：

- 运行时文件仍然保存在服务器本地 `var/`
- 没有对象存储
- 没有数据库持久化任务状态
- 没有登录、额度、支付、风控
- 没有统一监控、告警、日志采集

## 14. 建议的下一步

如果你准备真正运营，建议接下来按这个顺序升级：

1. 对象存储
2. 数据库
3. 登录与额度控制
4. 支付
5. 限流与安全策略
6. 监控与日志系统
