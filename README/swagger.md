# Swagger / API 文档

## 在线访问

服务启动后访问：

```
http://localhost:8080/swagger/index.html
```

## 统一响应格式

所有接口返回 JSON，结构如下：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

| code | 含义 |
|------|------|
| `0` | 成功 |
| `400` | 参数错误 |
| `401` | 未授权 / Token 无效 |
| `404` | 资源不存在 |
| `500` | 服务器内部错误 |

分页列表的 `data` 结构：

```json
{
  "list": [],
  "total": 100,
  "pageIndex": 1,
  "pageSize": 20
}
```

## 认证方式

除登录接口外，受保护接口需在 Header 中携带 JWT：

```
Authorization: Bearer <token>
```

Token 由 `POST /api/permission/user/login` 返回，默认有效期 24h（可通过 `JWT_EXPIRATION` 配置）。

## 接口列表

### 用户登录（公开）

```
POST /api/permission/user/login
```

**请求体：**

```json
{
  "account": "admin",
  "password": "123456"
}
```

**成功响应 data：**

```json
{
  "token": "eyJhbG...",
  "user": { "userId": "...", "account": "...", "username": "..." }
}
```

---

### 创建用户（需鉴权）

```
POST /api/permission/user/create
```

**请求体：**

```json
{
  "account": "user001",
  "password": "123456",
  "username": "张三",
  "email": "user@example.com",
  "phone": "13800138000"
}
```

`account` 和 `password` 必填；`username` 为空时默认使用 `account`。

---

### 用户列表（需鉴权）

```
GET /api/permission/user/list?pageIndex=1&pageSize=20
```

| 参数 | 默认 | 说明 |
|------|------|------|
| `pageIndex` | `1` | 页码 |
| `pageSize` | `20` | 每页条数，最大 100 |

---

### 编辑用户（需鉴权）

```
PUT /api/permission/user/edit/{userId}
```

**请求体（字段均可选）：**

```json
{
  "username": "新名称",
  "account": "new_account",
  "email": "new@example.com",
  "phone": "13900139000",
  "password": "new_password"
}
```

---

### 删除用户（需鉴权）

```
DELETE /api/permission/user/delete/{userId}
```

## 重新生成 Swagger 文档

Handler 中使用 swag 注解（如 `@Summary`、`@Router`），修改后需重新生成：

```bash
# 安装 swag CLI（首次）
go install github.com/swaggo/swag/cmd/swag@latest

# 在项目根目录生成
swag init -g cmd/server/main.go -o docs
```

生成文件位于 `docs/` 目录：

- `docs.go` — Go 包，供 router 导入
- `swagger.json` / `swagger.yaml` — OpenAPI 规范

## 注解位置

Swagger 元信息定义在 `cmd/server/main.go` 文件头：

```go
// @title VSB Server API
// @version 1.0
// @host localhost:8080
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
```

各接口注解写在对应 handler 方法上方，例如 `internal/handler/permission/user/login.go`。
