.PHONY: setup-hooks check fmt vet build tidy run

# 一键配置 git hooks（clone 项目后首先运行）
setup-hooks:
	@echo "Configuring git hooks path..."
	@git config core.hooksPath .githooks
	@chmod +x .githooks/pre-commit
	@echo "✅ Git hooks configured. Pre-commit checks are now active."

# 运行所有检查（与 pre-commit hook 一致）
check:
	@echo "🔍 Running all checks..."
	@UNFORMATTED=$$(gofmt -l . 2>&1); \
	if [ -n "$$UNFORMATTED" ]; then \
		echo "❌ 代码格式检查未通过："; echo "$$UNFORMATTED"; exit 1; \
	fi
	@echo "✅ gofmt 检查通过"
	@go vet ./... && echo "✅ go vet 检查通过"
	@go build ./... && echo "✅ go build 编译通过"
	@echo "🎉 所有检查通过"

# 格式化代码
fmt:
	@gofmt -w .
	@echo "✅ 代码格式化完成"

# 静态分析
vet:
	@go vet ./...

# 编译
build:
	@go build ./...

# 整理依赖
tidy:
	@go mod tidy

# 启动服务
run:
	@go run cmd/server/main.go
