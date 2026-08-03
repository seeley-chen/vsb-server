# 配置说明

配置通过 `.env` 文件或系统环境变量加载，由 `config/config.go` 统一读取并在启动时校验。

## 环境变量

| 变量 | 必填 | 默认值 | 说明 |
|------|------|--------|------|
| `SERVER_PORT` | 否 | `8080` | HTTP 服务端口 |
| `MONGODB_URI` | **是** | — | MongoDB 连接 URI |
| `MONGODB_DB` | 否 | `vsb` | 数据库名称 |
| `JWT_SECRET` | **是** | — | JWT 签名密钥 |
| `JWT_EXPIRATION` | 否 | `24h` | Token 有效期，Go duration 格式（如 `24h`、`30m`） |
| `LOG_LEVEL` | 否 | `info` | 日志级别（zap） |
| `CORS_ALLOWED_ORIGINS` | 否 | 空 | 允许跨域的前端 Origin，逗号分隔 |

## 示例

参考项目根目录 `.env.example`：

```env
SERVER_PORT=8080
MONGODB_URI=mongodb+srv://<user>:<password>@cluster.mongodb.net
MONGODB_DB=vsb
JWT_SECRET=your-jwt-secret
JWT_EXPIRATION=24h
LOG_LEVEL=info
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

## 启动校验

`config.Validate()` 在启动时检查：

- `MONGODB_URI` 不能为空
- `JWT_SECRET` 不能为空
- `JWT_EXPIRATION` 必须是合法的 duration 字符串

校验失败会直接 `log.Fatalf` 退出。

## 加载机制

1. 调用 `godotenv.Load()` 尝试读取项目根目录 `.env`
2. 若 `.env` 不存在或未设置某变量，则回退到系统环境变量
3. 仍无值时使用代码中的默认值（仅非必填项）

## 注意事项

- `.env` 含敏感信息，**不要提交到 Git**
- 生产环境 `JWT_SECRET` 应使用足够长的随机字符串
- `CORS_ALLOWED_ORIGINS` 生产环境需设置为实际前端域名，留空则不允许任何跨域 Origin
