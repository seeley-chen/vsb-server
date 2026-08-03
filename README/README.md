# VSB Server

VSB 服务端，基于 Go + MongoDB 的 REST API 服务，当前已实现权限模块下的用户管理（登录、创建、列表、编辑、删除）。

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
├── cmd/server/          # 程序入口
├── config/              # 配置加载与校验
├── docs/                # Swagger 生成文件（docs.go / swagger.json / swagger.yaml）
├── internal/
│   ├── database/      # MongoDB 连接
│   ├── handler/       # HTTP 处理器（按业务模块划分）
│   ├── middleware/    # 中间件（鉴权、CORS、日志、限流等）
│   ├── model/         # 数据模型
│   ├── repository/    # 数据访问层
│   ├── router/        # 路由注册
│   └── service/       # 业务逻辑层
└── pkg/
    ├── jwt/           # JWT 工具
    ├── logger/        # 日志初始化
    └── response/      # 统一响应格式
```

## 分层架构

```
HTTP Request
    ↓
Middleware（Recover → BodyLimit → Logger → CORS → Auth）
    ↓
Handler（参数解析、校验、调用 Service）
    ↓
Service（业务逻辑、JWT 签发）
    ↓
Repository（MongoDB CRUD）
```

## 快速启动

### 1. 配置环境变量

复制 `.env.example` 为 `.env`，填写必填项（详见 [config.md](./config.md)）：

```bash
cp .env.example .env
```

### 2. 安装依赖并运行

```bash
go mod download
go run cmd/server/main.go
```

服务默认监听 `http://localhost:8080`。

### 3. 健康检查

```bash
curl http://localhost:8080/health
# ok
```

## 文档索引

| 文档 | 说明 |
|------|------|
| [config.md](./config.md) | 环境变量与配置项说明 |
| [swagger.md](./swagger.md) | API 文档与 Swagger 使用方式 |
| [architecture.md](./architecture.md) | 目录约定、中间件与认证机制 |

## 当前 API 概览

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/health` | 否 | 健康检查 |
| POST | `/api/permission/user/login` | 否 | 用户登录 |
| POST | `/api/permission/user/create` | 是 | 创建用户 |
| GET | `/api/permission/user/list` | 是 | 用户列表（分页） |
| PUT | `/api/permission/user/edit/{userId}` | 是 | 编辑用户 |
| DELETE | `/api/permission/user/delete/{userId}` | 是 | 删除用户 |

## VS Code 调试

已配置 `.vscode/launch.json`，可直接使用 **Launch Server** 启动调试。
