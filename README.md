# 肿瘤医生 Go 项目

基于 Gin 的模块化单体后端骨架，面向可商用项目演进。前期功能包含注册、登录、资讯列表、资讯详情，以及阿里云 OSS 签名上传。

## 技术栈

- Go + Gin
- MySQL + GORM
- Redis
- JWT + bcrypt
- 阿里云 OSS
- Nginx + SSL + CDN
- Zap 日志
- Docker Compose 本地依赖

## 项目结构

```text
cmd/server                 程序入口
configs                    配置样例
internal/app               Gin 服务初始化和路由装配
internal/modules/auth      注册、登录
internal/modules/admin     后台管理员登录和鉴权
internal/modules/user      用户模型
internal/modules/news      资讯列表、详情
internal/modules/upload    OSS 签名上传
internal/middleware        鉴权、日志、跨域、恢复、请求 ID
internal/infrastructure    MySQL、Redis、OSS、Logger
internal/pkg               通用响应、错误码、JWT、密码、分页
migrations                 MySQL 建表脚本
deployments                Nginx 和 systemd 样例
docs                       接口文档
scripts                    运维脚本
```

## 环境变量

当前代码读取环境变量，示例见 `.env.example`：

```bash
APP_ENV=local
APP_ADDR=:8080
MYSQL_DSN='root:root@tcp(127.0.0.1:3306)/doctor_go?charset=utf8mb4&parseTime=True&loc=Local'
REDIS_ADDR=127.0.0.1:6379
JWT_SECRET=change-me
OSS_ENDPOINT=
OSS_ACCESS_KEY_ID=
OSS_ACCESS_KEY_SECRET=
OSS_BUCKET=
OSS_BASE_URL=
```

## 启动

```bash
cp .env.example .env
docker compose up -d
go mod tidy
go run ./cmd/server
```

Docker 构建部署：

```bash
docker build -t doctor-go:latest .
docker compose up -d --build
```

数据库迁移：

```bash
APP_ENV=local go run ./cmd/migrate -action up
APP_ENV=local go run ./cmd/migrate -action version
APP_ENV=local go run ./cmd/migrate -action force -version 1
```

当前迁移文件只包含 `up`，需要回滚时再为对应版本补充 `.down.sql`。
本地启动服务前请先执行迁移，RBAC 表依赖 `006_create_admin_rbac.up.sql`。

Docker Compose 启动 `app` 时会自动执行：

```bash
/app/doctor-go-migrate -action up
```

部署说明：

```text
docs/deploy.md
```

生产容器访问：

```text
http://localhost:8080/health
http://localhost:8080/swagger/index.html
```

配置文件：

```text
configs/config.local.yaml
configs/config.test.yaml
configs/config.prod.yaml
```

通过 `APP_ENV` 选择环境：

```bash
APP_ENV=local go run ./cmd/server
APP_ENV=test go run ./cmd/server
APP_ENV=prod ./doctor-go
```

配置优先级：

```text
环境变量 > configs/config.{APP_ENV}.yaml > 默认值
```

Swagger:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g ./cmd/server/main.go --parseInternal -o ./docs/swagger
go run ./cmd/server
```

访问：

```text
http://localhost:8080/swagger/index.html
```

健康检查：

```bash
curl http://localhost:8080/health
```

返回：

```json
{"code":0,"message":"success","data":{"status":"ok"}}
```

默认后台管理员：

```text
username: admin
password: admin123456
```

Redis 缓存覆盖：

```text
GET /api/v1/news
GET /api/v1/news/:id
GET /api/v1/news/categories
```

后台创建、编辑、发布/下架、删除资讯或分类时，会自动清理相关缓存。

登录安全策略：

```text
POST /api/v1/auth/register       每 IP 每分钟 5 次
POST /api/v1/auth/login          每 IP 每分钟 10 次
POST /api/v1/admin/auth/login    每 IP 每分钟 5 次
```

登录态：

```text
登录返回 access_token + refresh_token
refresh_token 用于刷新登录态
logout 会拉黑当前 access token，并删除 refresh token
Redis key: auth:refresh:* / auth:blacklist:*
```

如果你之前已经手动创建过 `doctor-go-mysql` 容器，首次使用 Compose 前需要先停掉旧容器：

```bash
docker stop doctor-go-mysql
docker rm doctor-go-mysql
docker compose up -d
```
