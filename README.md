# VSB Server

基于 Go（gorilla/mux + MongoDB）的后端 REST API 服务。当前包含 **login** 与 **permission**（department / role / user）模块。

当前 Go module path：`github.com/Vanselyn/vsb-server`

## 技术栈

| 类别 | 选型 |
|------|------|
| 语言 | Go 1.26 |
| HTTP 路由 | gorilla/mux |
| 数据库 | MongoDB |
| 认证 | JWT (HS256) |
| 日志 | zap |
| API 文档 | swaggo/swag + http-swagger |

## 项目结构

```
vsb-server/
├── cmd/server/            # 程序入口（main.go，含 swag 元信息）
├── config/                # 配置加载与校验
├── docs/                  # Swagger 生成文件（docs.go / swagger.json / swagger.yaml）
├── internal/
│   ├── database/          # MongoDB 连接
│   ├── deps/              # 依赖容器（Deps：DB、JWT、Logger、EnsureIndexes）
│   ├── middleware/        # 中间件：cors / recover / bodylimit / logger / auth
│   ├── module/            # 大模块装配（permission、login），由 module.go 汇总
│   ├── router/            # 路由组装（公开 / 受保护子路由）
│   ├── domain/            # 业务模块，每个模块自包含 handler/model/repository/service
│   └── tools/             # 通用工具（TrimSpace / IsEmpty / Default / MaskString 等）
└── pkg/
    ├── idgen/             # ID 生成
    ├── jwt/               # JWT 工具
    ├── logger/            # 日志初始化
    └── response/          # 统一响应 / 参数绑定 / 分页
```

## 分层架构

每个业务模块自包含四层（login 无独立 repository，复用 user 的 UserRepo）：

```
HTTP Request
    ↓
Middleware（CORS → Recover → MaxBodySize → Logger → Auth）
    ↓
Handler（handler.go：参数解析、校验、调用 Service）
    ↓
Service（service.go：业务逻辑）
    ↓
Repository（repository.go：MongoDB CRUD）
    ↓
Model（model.go：请求 / 响应 / 存储结构）
```

模块装配在 `internal/module/<name>/register.go`，由 `internal/module/module.go` 的 `All` 汇总，`internal/router/router.go` 统一调用。

## 快速开始

```bash
git clone git@github.com:Vanselyn/vsb-server.git
cd vsb-server
cp .env.example .env        # 填写 MONGODB_URI / JWT_SECRET 等必填项
make setup-hooks            # 配置 git hooks（clone 后首次运行）
make run                    # 启动服务，默认监听 http://localhost:8080
```

健康检查：

```bash
curl http://localhost:8080/health
# ok
```

Swagger UI：`http://localhost:8080/swagger/index.html`

## 配置说明

通过 `.env` 或系统环境变量加载，由 `config/config.go` 统一读取并在启动时校验。

| 变量 | 必填 | 默认值 | 说明 |
|------|------|--------|------|
| `SERVER_PORT` | 否 | `8080` | HTTP 服务端口 |
| `MONGODB_URI` | **是** | — | MongoDB 连接 URI |
| `MONGODB_DB` | 否 | `vsb` | 数据库名称 |
| `JWT_SECRET` | **是** | — | JWT 签名密钥 |
| `JWT_EXPIRATION` | 否 | `24h` | Token 有效期，Go duration 格式（如 `24h`、`30m`） |
| `LOG_LEVEL` | 否 | `info` | 日志级别（zap） |
| `CORS_ALLOWED_ORIGINS` | 否 | 空 | 允许跨域的前端 Origin，逗号分隔 |

启动校验：`MONGODB_URI`、`JWT_SECRET` 不能为空，`JWT_EXPIRATION` 必须为合法 duration，否则直接退出。`.env` 含敏感信息，不要提交到 Git。

## Make 命令

| 命令 | 说明 |
| --- | --- |
| `make setup-hooks` | 配置 git hooks（pre-commit / commit-msg），clone 后首次运行 |
| `make check` | 全部检查：gofmt + go vet + go build（与 pre-commit 一致） |
| `make fmt` | gofmt 格式化全部代码 |
| `make vet` | 静态分析 |
| `make build` | 编译 |
| `make tidy` | 整理依赖（go mod tidy） |
| `make run` | 启动服务 |
| `make rename-module NEW=<新地址>` | 重命名 Go module path（见文末） |

## 模块与 API

统一响应格式：

```json
{ "code": 0, "message": "success", "data": {} }
```

| code | 含义 |
|------|------|
| `0` | 成功 |
| `400` | 参数错误 |
| `401` | 未授权 / Token 无效 |
| `404` | 资源不存在 |
| `500` | 服务器内部错误 |

分页列表的 `data`：

```json
{ "list": [], "total": 100, "pageIndex": 1, "pageSize": 20 }
```

鉴权：受保护接口需在 Header 携带 `Authorization: Bearer <token>`，token 由 `POST /api/login` 返回。


## 架构与约定

### 路由注册

`internal/router/router.go` 组装依赖并注册路由：

- **公开路由**：挂在根 Router（如 `/health`、`/api/login`）
- **受保护路由**：挂在 `/api` 子路由并应用 `middleware.Auth`（department / role / user / privileges）
- **Swagger**：`/swagger/` 前缀，无需鉴权

### 中间件

执行顺序（由外到内）：

```
CORS（main.go 外层包裹）
  → Recover（panic 恢复）
  → MaxBodySize（请求体上限 1MB）
  → Logger（请求日志）
  → Auth（仅 /api 受保护子路由，JWT 校验）
```

| 中间件 | 文件 | 要点 |
|--------|------|------|
| CORS | `middleware/cors.go` | 白名单 Origin，须包裹在 mux 外层以处理 OPTIONS 预检 |
| Recover | `middleware/recover.go` | 捕获 panic，返回 500 |
| MaxBodySize | `middleware/bodylimit.go` | 限制请求体大小 |
| Logger | `middleware/logger.go` | 记录请求方法、路径、耗时 |
| Auth | `middleware/auth.go` | 解析 Bearer Token，将 `userId` 写入 context |

### JWT 认证

- 算法：HS256
- Claims 字段：`user_id`
- 工具包：`pkg/jwt/jwt.go`
- 鉴权后通过 `middleware.GetUserID(ctx)` 获取当前用户 ID

### 数据存储

- 数据库：MongoDB（集合 users / roles / departments）
- 密码：bcrypt 哈希，API 响应不返回 `passwordHash`
- 用户 ID：由 `pkg/idgen` 生成
- department / role 的索引在启动时自动创建（role 的 `role_id`、`name` 为唯一索引）

### 优雅关停与超时

`cmd/server/main.go` 监听 `SIGINT` / `SIGTERM`，收到信号后停止接受新连接、等待进行中请求完成（最长 30 秒）、断开 MongoDB 连接、刷新日志。

| 配置项 | 值 |
|--------|-----|
| ReadHeaderTimeout | 5s |
| ReadTimeout | 15s |
| WriteTimeout | 15s |
| IdleTimeout | 60s |

### 扩展新模块

1. 在 `internal/domain/<module>/` 下创建 `handler.go` / `model.go` / `repository.go` / `service.go`
2. 在 `internal/module/<module>/register.go` 实现 `Register(d, public, protected)`，完成依赖注入与路由挂载
3. 在 `internal/module/module.go` 的 `All` 中追加该 `Register`
4. 需要索引时通过 `d.EnsureIndexes(repo.EnsureIndexes)` 注册
5. 运行 `make check` 验证

## Swagger 文档

Handler 中使用 swag 注解（如 `@Summary`、`@Router`），修改后重新生成：

```bash
go install github.com/swaggo/swag/cmd/swag@latest   # 首次安装
swag init -g cmd/server/main.go -o docs
```

元信息定义在 `cmd/server/main.go` 文件头，各接口注解写在对应 `handler.go` 的 handler 方法上方。

## VS Code 调试

已配置 `.vscode/launch.json`，可直接使用 **Launch Server** 启动调试。

## 切换 / 复制到新仓库

Go 的 module path 是所有内部包 import 路径的前缀（定义在 `go.mod` 的 `module` 行）。当仓库地址变更或复制到新仓库时，一条命令同步更新（自动重写 `go.mod` 及所有 `.go` / `.yaml` / `.json` / `.md` 中的引用，含 swagger 下划线形式）：

```bash
make rename-module NEW=github.com/<owner>/<repo>
make check
```

说明：自动排除 `.git`、`.vscode`；使用 `sed -i.bak` + 自动清理 `.bak`，macOS / Linux 通用，无需额外安装软件。
