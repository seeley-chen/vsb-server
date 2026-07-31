# 配置说明

项目使用 `.env` 文件管理环境变量，基于 `godotenv` 加载。

## 环境变量一览

| 变量名 | 必填 | 默认值 | 说明 |
|--------|------|--------|------|
| SERVER_PORT | 否 | 8080 | HTTP 服务监听端口 |
| MONGODB_URI | ✅ | — | MongoDB 连接字符串（支持 SRV 格式） |
| MONGODB_DB | 否 | vsb | 数据库名称 |
| JWT_SECRET | ✅ | — | JWT 签名密钥，建议 32 位以上随机字符串 |
| JWT_EXPIRATION | 否 | 24h | Token 有效期（Go duration 格式） |

## 修改配置

1. 复制模板：`cp .env.example .env`
2. 编辑 `.env` 填入真实值
3. 重启服务生效

> ⚠️ `.env` 已加入 `.gitignore`，切勿提交真实密钥到 Git 仓库。

## 生产环境建议

- JWT_SECRET 使用 `openssl rand -hex 32` 生成
- MongoDB 使用专用读写账号
- 通过反向代理（Nginx/Caddy）终止 TLS