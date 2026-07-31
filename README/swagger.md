# Swagger 文档使用指南

## 生成文档

在项目根目录执行：

```bash
swag init -g cmd/server/main.go -o docs
```

这会在 `docs/` 目录生成 `docs.go`、`swagger.json`、`swagger.yaml` 三个文件。

## 访问文档

启动服务后访问：http://localhost:8080/swagger/index.html

## 编写注解规范

在每个 Handler 函数上方添加 Swagger 注解：

```go
go

// @Summary 接口简述

// @Description 接口详细描述

// @Tags 分组名称

// @Accept json

// @Produce json

// @Param name query/body/path type true/false "描述"

// @Success 200 {object} response.Response "成功说明"

// @Failure 400 {object} response.Response "失败说明"

// @Security ApiKeyAuth

// @Router /path [method]
```

修改注解后重新运行 `swag init` 即可更新文档。

## 启动 Swagger
如果要在 router.go 中启用 Swagger UI，添加以下导入和路由

```go
import (
	"github.com/swaggo/http-swagger"
	_ "github.com/seeley-chen/vsb-server/docs" // swag 生成的 docs 包
)
```

在 `router.New()` 中添加：

```
r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
```

