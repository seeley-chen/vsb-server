# Vsb Server

基于 Go 的轻量 RESTful API 服务，模块化分层架构，方便长期维护。

## 目录结构
vsb-server/
├── cmd/server/main.go          # 程序入口
├── config/config.go            # 配置加载
├── internal/
│   ├── database/mongo.go      # MongoDB 连接
│   ├── handler/user/           # HTTP 处理器（接口层）
│   ├── service/user/           # 业务逻辑层
│   ├── repository/user/        # 数据访问层
│   ├── model/user.go          # 数据模型
│   ├── middleware/             # 中间件（CORS / JWT / 日志）
│   └── router/router.go      # 路由总装
├── pkg/
│   ├── jwt/jwt.go            # JWT 签发/验证
│   ├── response/response.go  # 统一响应封装
│   └── logger/logger.go      # Zap 日志封装
├── readme/                    # 各模块说明文档
├── .env / .env.example        # 环境变量
├── .gitignore
└── go.mod

## 快速开始
1. 克隆项目 git clone git@github.com:seeley-chen/vsb-server.git
2. 复制环境变量 cp .env.example .env
3. 安装依赖 go mod tidy
4. 运行 go run cmd/server/main.go

服务启动后默认监听 `:8080`。

## 鉴权方式

登录后获取 JWT Token，后续请求在 Header 中携带：
Authorization: Bearer <your-token>

## 技术栈

| 组件 | 选型 | 说明 |
|------|------|------|
| 路由 | gorilla/mux | 轻量、成熟 |
| MongoDB | mongo-go-driver | 官方驱动 |
| JWT | golang-jwt/jwt/v5 | HS256 签名 |
| 日志 | uber/zap | 高性能结构化日志 |
| 配置 | joho/godotenv | .env 加载 |
| 密码加密 | golang.org/x/crypto/bcrypt | 加盐哈希 |
| ID 生成 | google/uuid | UUID v4 |

## 开发约定

1. **分层清晰**：handler → service → repository，禁止跨层调用
2. **统一响应**：所有接口返回 `{ code, message, data }` 格式
3. **配置外置**：敏感信息一律走环境变量，不硬编码
4. **模块独立**：新增业务域复制 user 模块结构即可