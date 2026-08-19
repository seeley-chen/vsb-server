# VSB Server

基于 Go（gorilla/mux + MongoDB）的后端 REST API。

当前模块：`login`、`permission`（department / role / user）、`product`（category）。

Go module：`github.com/Vanselyn/vsb-server`

---

## 技术栈

| 类别      | 选型                         |
| ------- | -------------------------- |
| 语言      | Go 1.26                    |
| HTTP 路由 | gorilla/mux                |
| 数据库     | MongoDB                    |
| 认证      | JWT (HS256)                |
| 日志      | zap                        |
| API 文档  | swaggo/swag + http-swagger |
| 热重载     | air                        |

---

## 项目结构

```
vsb-server/
├── cmd/server/            # 入口（main.go，含 swagger 元信息）
├── config/                # 配置加载与校验
├── docs/                  # swagger 生成结果（不要手改）
├── .githooks/             # git hooks（make setup-hooks 后生效）
├── internal/
│   ├── database/          # MongoDB 连接
│   ├── deps/              # 依赖容器
│   ├── middleware/        # cors / recover / bodylimit / logger / auth / logviewer
│   ├── module/            # 大模块装配，由 module.go 汇总
│   ├── router/            # 路由组装
│   ├── domain/            # 业务：handler / model / repository / service
│   └── tools/
└── pkg/
    ├── idgen/
    ├── jwt/
    ├── logger/
    ├── mongoFilter/
    └── response/
```

---

## 初始化

从别处新 clone 后，按顺序做一次即可。之后日常开发不需要重复。

**前置：** 已安装 Go 1.26+，本机或云端 MongoDB 可连，`$(go env GOPATH)/bin` 已加入 `PATH`。

```bash
git clone git@github.com:Vanselyn/vsb-server.git
cd vsb-server

cp .env.example .env          # 必填：MONGODB_URI、JWT_SECRET
make setup-hooks              # 启用 git hooks（不跑则 commit 检查不会生效）
go mod download               # 拉依赖
```

然后任选一种方式启动，见 [启动](#启动)。

---

## 配置

通过项目根目录 `.env` 加载（也可用系统环境变量）。启动时由 `config/config.go` 校验，不合法直接退出。

`.env` 含密钥，不要提交到 Git。

| 变量                     | 必填  | 默认       | 说明                                  |
| ---------------------- | --- | -------- | ----------------------------------- |
| `SERVER_PORT`          | 否   | `8080`   | HTTP 端口                             |
| `MONGODB_URI`          | 是   | -        | MongoDB 连接 URI                      |
| `MONGODB_DB`           | 否   | `vsb`    | 数据库名                                |
| `JWT_SECRET`           | 是   | -        | JWT 签名密钥                            |
| `JWT_EXPIRATION`       | 否   | `24h`    | Token 有效期，如 `24h`、`30m`             |
| `LOG_LEVEL`            | 否   | `info`   | `debug` / `info` / `warn` / `error` |
| `LOG_BODY`             | 否   | `masked` | `full` / `masked` / `off`，见「日志」     |
| `CORS_ALLOWED_ORIGINS` | 否   | -        | 允许跨域的 Origin，逗号分隔                   |

---

## Make

所有常用操作都走 Makefile，不需要记底层命令。

| 命令                              | 说明                                         |
| ------------------------------- | ------------------------------------------ |
| `make setup-hooks`              | 启用 git hooks。clone 后只跑一次                   |
| `make check`                    | gofmt + go vet + go build（与 pre-commit 一致） |
| `make fmt`                      | 格式化全部代码                                    |
| `make vet`                      | 静态分析                                       |
| `make build`                    | 编译                                         |
| `make tidy`                     | `go mod tidy`                              |
| `make run`                      | 直接启动（不热重载）                                 |
| `make dev`                      | 开发模式：文件变化自动 rebuild + 重启                   |
| `make install-air`              | 安装 air。`make dev` 会自动装，一般不用单独跑             |
| `make status`                   | 检查 air / 端口 / `/health`                    |
| `make swagger`                  | 按 handler 注解重新生成 docs/                     |
| `make rename-module NEW=<path>` | 更换 Go module path，见「切换仓库」                  |

---

## 启动

**一次性启动：**

```bash
make run
```

**日常开发（推荐）：** 监听 `.go` / `.env` 变化，自动 rebuild + 重启。首次会自动安装 air。

```bash
make dev
```

启动成功后终端会出现 banner：

```
🚀 vsb-server started | pid=12345 | addr=:8080 | log_body=full | time=2026-08-19 20:15:03
```

改代码保存后应再出现一条新 banner。没有新 banner = 编译失败（看终端红色错误）或端口被占。

排查「改了代码没生效」或「前端连不上」：

```bash
make status
```

---

## Swagger

### 查看文档

服务启动后打开：

```
http://localhost:8080/swagger/index.html
```

无需登录。接口列表来自 `docs/`（由 `make swagger` 生成），不是运行时实时扫代码。

### 什么时候要更新

改了 handler 上的 swag 注解（`@Summary`、`@Router`、`@Param` 等），或新增了接口，必须重新生成，否则 UI 里看不到。

### 怎么更新

1. 在对应 `handler.go` 的方法上方写 / 改注解。示例：

```go
// Login 用户登录
// @Summary 用户登录
// @Tags 登录
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=LoginResponse}
// @Router /api/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) { ... }

// 分页列表要把 list 的元素类型写出来，否则 swagger 里 list 没有字段：
// @Success 200 {object} response.Response{data=response.PageData{list=[]Response}}
```

2. 全局标题、host、鉴权方式写在 `cmd/server/main.go` 文件头，一般不用动。

3. 本机第一次生成前，安装 swag（只需一次）：

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

4. 生成：

```bash
make swagger
```

等价于：`swag init -g cmd/server/main.go -o docs --parseInternal`

`--parseInternal` 必须带：handler 在 `internal/` 下，不加的话生成结果只有元信息、没有接口。

5. 重启服务（`make run` 需手动重启；`make dev` 会因 `docs/` 变化自动重启）。刷新浏览器即可。

不要手改 `docs/docs.go`、`docs/swagger.json`、`docs/swagger.yaml`，下次 `make swagger` 会覆盖。

---

## Git Hooks

`make setup-hooks` 会把 `core.hooksPath` 指到 `.githooks/`。

| Hook         | 时机           | 做什么                                  |
| ------------ | ------------ | ------------------------------------ |
| `pre-commit` | commit 前     | gofmt、go vet、go build、go.mod 是否 tidy |
| `commit-msg` | 填写 message 后 | 校验 Conventional Commits              |

提交信息格式：`type: 说明` 或 `type(scope): 说明`

允许的 type：`feat` `fix` `docs` `style` `refactor` `perf` `test` `chore` `build` `ci` `revert`

```
feat: 新增用户导出
fix(role): 修复角色权限合并
docs: 更新 README
```

本地想先自检，不必等 hook：

```bash
make check
make fmt          # 格式不过时先格式化
make tidy         # go.mod 不同步时先整理
```

---

## VS Code 调试

`.vscode/launch.json` 已配好，二者共用 `:8080`，不能同时开。

| 配置                    | 用途                 |
| --------------------- | ------------------ |
| Dev (Air Hot Reload)  | 日常开发，等价 `make dev` |
| Launch Server (Debug) | 断点调试。先停掉上面的 Dev    |

---

## 常用地址

默认端口 `8080`，可被 `.env` 的 `SERVER_PORT` 覆盖。

| 地址                    | 说明                |
| --------------------- | ----------------- |
| `GET /health`         | 健康检查，返回 `ok`      |
| `/swagger/index.html` | API 文档            |
| `/admin/logs`         | 实时请求日志网页（开发用，无鉴权） |

```bash
curl http://localhost:8080/health
```

---

## API 约定

统一响应：

```json
{ "code": 0, "message": "success", "data": {} }
```

| code  | 含义             |
| ----- | -------------- |
| `0`   | 成功             |
| `400` | 参数错误           |
| `401` | 未授权 / Token 无效 |
| `404` | 资源不存在          |
| `500` | 服务器内部错误        |

分页 `data`：

```json
{ "list": [], "total": 100, "pageIndex": 1, "pageSize": 20 }
```

鉴权：受保护接口 Header 带 `Authorization: Bearer <token>`。token 由 `POST /api/login` 返回。

| 类型  | 路径                                                    | 鉴权  |
| --- | ----------------------------------------------------- | --- |
| 公开  | `/health`、`POST /api/login`、`/swagger/`、`/admin/logs` | 否   |
| 受保护 | `/api/permission/*`、`/api/product/*`                  | 是   |

---

## 分层架构

每个业务模块自包含四层（login 无独立 repository，复用 user 的 UserRepo）：

```
HTTP Request
  → Middleware（CORS → Recover → MaxBodySize → Logger → Auth）
  → Handler（参数解析、校验、调 Service）
  → Service（业务逻辑）
  → Repository（MongoDB）
  → Model（请求 / 响应 / 存储结构）
```

模块装配在 `internal/module/<name>/register.go`，由 `internal/module/module.go` 的 `All` 汇总，`internal/router/router.go` 统一调用。

---

## 中间件

执行顺序（由外到内）：

```
CORS（main.go 外层包裹，处理 OPTIONS）
  → Recover
  → MaxBodySize（1MB）
  → Logger
  → Auth（仅 /api 受保护子路由）
```

JWT：HS256，Claims 字段 `user_id`。鉴权后用 `middleware.GetUserID(ctx)` 取当前用户。

---

## 日志

### 请求日志

每个请求一条结构化日志，含 `request_id`、`method`、`path`、`status`、`duration`、`user_id`、`req_body`、`resp_body` 等。

`request_id` 同时写入响应头 `X-Request-ID`。前端报错时拿这个值，就能在终端或日志网页里对上那一次请求。

| 状态码      | 级别    |
| -------- | ----- |
| `>= 500` | ERROR |
| `>= 400` | WARN  |
| 其他       | INFO  |

未知 500 会再打一条带同一 `request_id` 的底层 error（如 Mongo 细节）。

### LOG_BODY

控制日志里请求/响应 body 展示多少。开发建议 `full`。

| 值        | 行为                                             |
| -------- | ---------------------------------------------- |
| `full`   | 完整内容（超 1KB 截断）                                 |
| `masked` | `password` / `secret` / `token` 等替换为 `***`（默认） |
| `off`    | 只记长度，如 `(123 bytes)`                           |

`make dev` 下改 `.env` 会自动重启。

### 日志网页

```
http://localhost:8080/admin/logs
```

SSE 实时表格：级别筛选、按 path / request_id 搜索、展开看 body。缓冲区约 500 条。body 同样受 `LOG_BODY` 控制。

公开路由，仅本地排障用。生产请用反向代理挡掉或加 IP 白名单。

---

## 扩展新模块

1. `internal/domain/<module>/` 下建 `handler.go` / `model.go` / `repository.go` / `service.go`
2. `internal/module/<module>/register.go` 实现 `Register(d, public, protected)`
3. 在 `internal/module/module.go` 的 `All` 里追加
4. 需要索引时：`d.EnsureIndexes(repo.EnsureIndexes)`
5. handler 写 swagger 注解后执行 `make swagger`
6. `make check` 验证

---

## 数据与关停

- 集合：users / roles / departments / product categories 等
- 密码 bcrypt 哈希，响应不返回 `passwordHash`
- department / role / category 索引在启动时自动创建

收到 `SIGINT` / `SIGTERM` 后停止接新连接，最多等 30 秒再断 Mongo、刷日志。

| HTTP 超时           | 值   |
| ----------------- | --- |
| ReadHeaderTimeout | 5s  |
| ReadTimeout       | 15s |
| WriteTimeout      | 15s |
| IdleTimeout       | 60s |

---

## 切换仓库

Go module path 是所有内部 import 的前缀。仓库地址变了或复制到新仓库时：

```bash
make rename-module NEW=github.com/<owner>/<repo>
make check
```

会改 `go.mod` 以及 `.go` / `.yaml` / `.json` / `.md` 里的引用（含 swagger 下划线形式）。自动跳过 `.git`、`.vscode`。
