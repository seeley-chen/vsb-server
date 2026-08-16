.PHONY: setup-hooks check fmt vet build tidy run dev install-air status rename-module

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

# 安装 live-reload 工具 air（首次使用时运行一次）
install-air:
	@go install github.com/air-verse/air@latest
	@echo "✅ air 已安装，路径: $$(go env GOPATH)/bin/air（请确保 GOPATH/bin 在 PATH 中）"

# 开发模式：监听文件变化自动 rebuild + 重启（依赖 air）
dev: install-air
	@air

# 检查开发服务状态（air 进程 + 端口监听 + 健康检查），排查"改了代码没生效"时优先运行
status:
	@echo "🔍 服务状态检查..."
	@PORT=$$(grep -E '^SERVER_PORT=' .env 2>/dev/null | cut -d= -f2); PORT=$${PORT:-8080}; \
	AIR_PID=$$(pgrep -x air | head -1); \
	if [ -n "$$AIR_PID" ]; then echo "✅ air 运行中 (pid=$$AIR_PID)"; else echo "❌ air 未运行（用 make dev 启动）"; fi; \
	LISTEN=$$(lsof -i :$$PORT -sTCP:LISTEN 2>/dev/null | tail -n +2); \
	if [ -n "$$LISTEN" ]; then echo "✅ :$$PORT 端口监听中"; echo "   $$LISTEN"; else echo "❌ :$$PORT 端口未监听（服务未启动或编译失败，检查 air 终端输出）"; fi; \
	echo "---"; \
	curl -s -m 2 http://localhost:$$PORT/health >/dev/null 2>&1 && echo "✅ /health 响应正常" || echo "❌ /health 无响应"

# 重命名 Go module path（切换/复制到新仓库时使用）
# 用法: make rename-module NEW=github.com/owner/repo
rename-module:
	@if [ -z "$(NEW)" ]; then echo "❌ 用法: make rename-module NEW=github.com/owner/repo"; exit 1; fi
	@OLD=$$(head -1 go.mod | sed 's/module //'); echo "🔄 $$OLD -> $(NEW)"; go mod edit -module $(NEW); OLD_S=$$(echo "$$OLD" | sed 's/\./_/g; s/\//_/g'); NEW_S=$$(echo "$(NEW)" | sed 's/\./_/g; s/\//_/g'); find . -type f \( -name '*.go' -o -name '*.yaml' -o -name '*.yml' -o -name '*.json' -o -name '*.md' \) -not -path './.git/*' -not -path './.vscode/*' -exec sed -i.bak -e "s|$$OLD|$(NEW)|g" -e "s|$$OLD_S|$$NEW_S|g" {} + ; find . -type f -name '*.bak' -not -path './.git/*' -not -path './.vscode/*' -delete; echo "✅ module path 已更新为 $(NEW)，建议运行 make check 验证"