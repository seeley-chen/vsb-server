# 架构与约定

## 模块划分

当前唯一业务模块为 **permission/user**（权限-用户管理），各层路径一一对应：

| 层级 | 路径 |
|------|------|
| Handler | `internal/handler/permission/user/` |
| Service | `internal/service/permission/user/` |
| Repository | `internal/repository/permission/user/` |
| Model | `internal/model/permission/` |

新增业务模块时，建议沿用 `internal/{layer}/{module}/{submodule}/` 的目录结构。

## 路由注册

`internal/router/router.go` 负责组装依赖并注册路由：

- **公开路由**：直接挂载在根 Router（如 `/api/permission/user/login`）
- **受保护路由**：挂载在 `/api` 子路由，并应用 `middleware.Auth`
- **Swagger**：`/swagger/` 前缀，无需鉴权

## 中间件

执行顺序（由外到内）：

```
CORS（main.go 外层包裹）
  → Recover（panic 恢复）
  → MaxBodySize（请求体上限 1MB）
  → Logger（请求日志）
  → Auth（仅 /api 子路由，JWT 校验）
```

| 中间件 | 文件 | 要点 |
|--------|------|------|
| CORS | `middleware/cors.go` | 白名单 Origin，须包裹在 mux 外层以正确处理 OPTIONS 预检 |
| Recover | `middleware/recover.go` | 捕获 panic，返回 500 |
| MaxBodySize | `middleware/bodylimit.go` | 限制请求体大小 |
| Logger | `middleware/logger.go` | 记录请求方法、路径、耗时 |
| Auth | `middleware/auth.go` | 解析 Bearer Token，将 `userId` 写入 context |

## JWT 认证

- 算法：HS256
- Claims 字段：`user_id`（用户 ID）
- 工具包：`pkg/jwt/jwt.go`
- 鉴权后可通过 `middleware.GetUserID(ctx)` 获取当前用户 ID

## 数据存储

- 数据库：MongoDB
- 用户集合：`users`
- 唯一索引：`account`（启动时自动创建）
- 密码：bcrypt 哈希，API 响应中不返回 `passwordHash`
- 用户 ID：毫秒时间戳 + 4 位随机数组成的纯数字字符串

## 优雅关停

`cmd/server/main.go` 监听 `SIGINT` / `SIGTERM`，收到信号后：

1. 停止接受新连接
2. 等待进行中的请求完成（最长 30 秒）
3. 断开 MongoDB 连接
4. 刷新日志缓冲

## HTTP 超时

| 配置项 | 值 |
|--------|-----|
| ReadHeaderTimeout | 5s |
| ReadTimeout | 15s |
| WriteTimeout | 15s |
| IdleTimeout | 60s |

## 扩展建议

添加新 API 模块的典型步骤：

1. 在 `internal/model/` 定义数据模型
2. 在 `internal/repository/` 实现 MongoDB 操作
3. 在 `internal/service/` 编写业务逻辑
4. 在 `internal/handler/` 编写 HTTP 处理器并添加 swag 注解
5. 在 `internal/router/router.go` 注册路由
6. 运行 `swag init` 更新文档
