# 切换到 ayhero@gmail.com 的 GitHub 账户
gh-switch-ayhero: ## 切换到 ayhero@gmail.com 的 GitHub 账户
	@echo "🔄 切换到 GitHub 账户 ayhero@gmail.com..."
	gh auth status --hostname github.com | grep "ayhero@gmail.com" >/dev/null 2>&1 && \
	  echo "✅ 已切换到 ayhero@gmail.com" || \
	  (echo "⚠️ 当前未登录 ayhero@gmail.com，请按提示输入 ayhero@gmail.com 进行登录..." && gh auth login -p https -h github.com)
	gh auth status

# InPayOS API Makefile

.PHONY: help build run test clean docker-build docker-run docker-stop migrate dev deps lint format

# 变量定义
APP_NAME=inpayos
BINARY_NAME=build/inpayos
DOCKER_IMAGE=inpayos-api
GO_FILES=$(shell find . -name "*.go" -not -path "./vendor/*")

# 默认目标
help: ## 显示帮助信息
	@echo "InPayOS API 开发工具"
	@echo ""
	@echo "使用方法: make [target]"
	@echo ""
	@echo "目标:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## 下载依赖
	@echo "📦 下载依赖..."
	go mod download
	go mod tidy

build: ## 构建应用
	@echo "🔨 构建应用..."
	@mkdir -p build
	go build -o $(BINARY_NAME) ./main

run: build ## 运行API服务器
	@echo "🚀 启动API服务器..."
	./$(BINARY_NAME) serve

dev: ## 开发模式运行（热重载）
	@echo "🛠️  开发模式启动..."
	go run ./main serve

migrate: build ## 运行数据库迁移
	@echo "🗄️  运行数据库迁移..."
	./$(BINARY_NAME) migrate

test: ## 运行测试
	@echo "🧪 运行测试..."
	go test -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "📊 生成测试覆盖率报告..."
	@mkdir -p build
	go test -coverprofile=build/coverage.out ./...
	go tool cover -html=build/coverage.out -o build/coverage.html
	@echo "📈 覆盖率报告已生成: build/coverage.html"

lint: ## 代码检查
	@echo "🔍 代码检查..."
	golangci-lint run

format: ## 格式化代码
	@echo "✨ 格式化代码..."
	gofmt -s -w $(GO_FILES)
	goimports -w $(GO_FILES)

clean: ## 清理构建文件
	@echo "🧹 清理构建文件..."
	rm -rf build/
	go clean

docker-build: ## 构建Docker镜像
	@echo "🐳 构建Docker镜像..."
	docker build -t $(DOCKER_IMAGE) .

docker-run: docker-build ## 运行Docker容器
	@echo "🐳 启动Docker容器..."
	docker run -d --name inpayos-api -p 8080:8080 -p 8081:8081 -p 8082:8082 -p 8083:8083 -p 8084:8084 $(DOCKER_IMAGE)

docker-stop: ## 停止Docker容器
	@echo "🛑 停止Docker容器..."
	docker stop inpayos-api || true
	docker rm inpayos-api || true

docker-logs: ## 查看Docker日志
	@echo "📋 查看Docker日志..."
	docker logs -f inpayos-api

install-tools: ## 安装开发工具
	@echo "🔧 安装开发工具..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/swaggo/swag/cmd/swag@latest

swagger: ## 生成统一Swagger文档
	@echo "📚 生成统一Swagger文档..."
	make swagger-openapi
	make swagger-merchant
	make swagger-cashier
	make swagger-admin

swagger-openapi: ## 生成OpenAPI Swagger文档
	@echo "📚 生成OpenAPI Swagger文档..."
	swag init -g main/main.go \
		--instanceName openapi \
		--tags "OpenAPI" \
		--parseDependency --parseInternal \
		-o ./docs/openapi

swagger-merchant: ## 生成MerchantAPI Swagger文档
	@echo "📚 生成MerchantAPI Swagger文档..."
	swag init -g main/main.go \
		--instanceName merchant \
		--tags "MerchantAPI" \
		--parseDependency --parseInternal \
		-o ./docs/merchant

swagger-cashier: ## 生成CashierAPI Swagger文档
	@echo "📚 生成CashierAPI Swagger文档..."
	swag init -g main/main.go \
		--instanceName cashier \
		--tags "CashierAPI" \
		--parseDependency --parseInternal \
		-o ./docs/cashier

swagger-admin: ## 生成AdminAPI Swagger文档
	@echo "📚 生成AdminAPI Swagger文档..."
	swag init -g main/main.go \
		--instanceName admin \
		--tags "AdminAPI" \
		--parseDependency --parseInternal \
		-o ./docs/admin

security-scan: ## 安全扫描
	@echo "🔒 执行安全扫描..."
	gosec ./...

mod-update: ## 更新所有依赖到最新版本
	@echo "⬆️  更新依赖..."
	go get -u ./...
	go mod tidy

benchmark: ## 运行性能测试
	@echo "⚡ 运行性能测试..."
	go test -bench=. -benchmem ./...

logs: ## 查看应用日志
	@echo "📋 查看应用日志..."
	tail -f logs/app.log

health-check: ## 健康检查
	@echo "💊 执行健康检查..."
	@echo "尝试OpenAPI端口8080..."
	@curl -f http://localhost:8080/health 2>/dev/null && echo "✅ OpenAPI端口8080服务正常" || echo "❌ OpenAPI端口8080无响应"
	@echo "尝试CashierAPI端口8081..."
	@curl -f http://localhost:8081/health 2>/dev/null && echo "✅ CashierAPI端口8081服务正常" || echo "❌ CashierAPI端口8081无响应"
	@echo "尝试MerchantAPI端口8082..."
	@curl -f http://localhost:8082/health 2>/dev/null && echo "✅ MerchantAPI端口8082服务正常" || echo "❌ MerchantAPI端口8082无响应"
	@echo "尝试CashierAdminAPI端口8083..."
	@curl -f http://localhost:8083/health 2>/dev/null && echo "✅ CashierAdminAPI端口8083服务正常" || echo "❌ CashierAdminAPI端口8083无响应"
	@echo "尝试AdminAPI端口8084..."
	@curl -f http://localhost:8084/health 2>/dev/null && echo "✅ AdminAPI端口8084服务正常" || echo "❌ AdminAPI端口8084无响应"

stop: ## 停止当前运行的服务
	@echo "🛑 停止服务..."
	@pkill -f "$(BINARY_NAME)" || echo "ℹ️  没有找到运行中的服务"
	@pkill -f "go run ./main" || echo "ℹ️  没有找到开发模式服务"

stop-all: ## 一键关闭所有五个服务 (OpenAPI、CashierAPI、MerchantAPI、CashierAdminAPI、AdminAPI)
	@echo "🛑 正在关闭所有InPayOS服务..."
	@echo "📍 关闭OpenAPI服务 (端口8080)..."
	@lsof -ti:8080 | xargs kill -9 2>/dev/null || echo "ℹ️  端口8080无进程运行"
	@echo "📍 关闭CashierAPI服务 (端口8081)..."
	@lsof -ti:8081 | xargs kill -9 2>/dev/null || echo "ℹ️  端口8081无进程运行"
	@echo "📍 关闭MerchantAPI服务 (端口8082)..."
	@lsof -ti:8082 | xargs kill -9 2>/dev/null || echo "ℹ️  端口8082无进程运行"
	@echo "📍 关闭CashierAdminAPI服务 (端口8083)..."
	@lsof -ti:8083 | xargs kill -9 2>/dev/null || echo "ℹ️  端口8083无进程运行"
	@echo "📍 关闭AdminAPI服务 (端口8084)..."
	@lsof -ti:8084 | xargs kill -9 2>/dev/null || echo "ℹ️  端口8084无进程运行"
	@echo "🧹 清理相关进程..."
	@pkill -f "inpayos" 2>/dev/null || echo "ℹ️  没有找到inpayos进程"
	@pkill -f "go run ./main" 2>/dev/null || echo "ℹ️  没有找到go run进程"
	@echo "✅ 所有服务已关闭"

kill-all: ## 强制停止所有相关进程和线程
	@echo "💀 强制停止所有相关进程..."
	@pkill -9 -f "inpayos" || echo "ℹ️  没有找到inpayos进程"
	@pkill -9 -f "go run ./main" || echo "ℹ️  没有找到go run进程"
	@pkill -9 -f "dlv dap" || echo "ℹ️  没有找到调试进程"
	@echo "🧹 清理临时文件..."
	@rm -f /tmp/inpayos-*.pid 2>/dev/null || true

status: ## 检查服务运行状态
	@echo "📊 检查服务状态..."
	@echo "=== InPayOS API 进程 ==="
	@ps aux | grep -E "(inpayos|go run.*main)" | grep -v grep || echo "❌ 没有找到API服务进程"
	@echo ""
	@echo "=== 调试进程 ==="
	@ps aux | grep "dlv dap" | grep -v grep || echo "ℹ️  没有调试进程运行"
	@echo ""
	@echo "=== 端口占用情况 ==="
	@lsof -i :8080 2>/dev/null || echo "ℹ️  端口8080未被占用"
	@lsof -i :8081 2>/dev/null || echo "ℹ️  端口8081未被占用"
	@lsof -i :8082 2>/dev/null || echo "ℹ️  端口8082未被占用"
	@lsof -i :8083 2>/dev/null || echo "ℹ️  端口8083未被占用"
	@lsof -i :8084 2>/dev/null || echo "ℹ️  端口8084未被占用"

restart: stop build run ## 重启服务

all: deps build test ## 执行完整构建流程

# 服务器管理
ssh-dev: ## 连接到AWS开发服务器
	@echo "🌐 连接到AWS开发服务器..."
	ssh aws-dev

ssh-prod: ## 连接到AWS生产服务器
	@echo "🌐 连接到AWS生产服务器..."
	ssh aws-prod

sync-config-dev: ## 同步dev.yaml配置到GitHub DEV环境Secrets
	@echo "🔄 同步dev.yaml配置到GitHub DEV环境..."
	@if [ ! -f dev.yaml ]; then \
		echo "❌ 错误: dev.yaml文件不存在"; \
		exit 1; \
	fi
	@echo "📤 上传dev.yaml内容到GitHub Secret CONFIG (DEV环境)..."
	@cat dev.yaml | gh secret set CONFIG --env DEV
	@echo "✅ 配置同步成功!"

sync-config-prod: ## 同步prod.yaml配置到GitHub PROD环境Secrets
	@echo "🔄 同步prod.yaml配置到GitHub PROD环境..."
	@if [ ! -f prod.yaml ]; then \
		echo "❌ 错误: prod.yaml文件不存在"; \
		exit 1; \
	fi
	@echo "📤 上传prod.yaml内容到GitHub Secret CONFIG (PROD环境)..."
	@cat prod.yaml | gh secret set CONFIG --env PROD
	@echo "✅ 配置同步成功!"

check-github-auth: ## 检查GitHub CLI认证状态
	@echo "🔐 检查GitHub CLI认证状态..."
	@gh auth status || echo "❌ 请先运行: gh auth login"

sync-config: ## 手动同步配置到GitHub Secrets
	@echo "🔄 手动同步配置到GitHub Secrets..."
	@echo "请手动运行以下命令同步配置:"
	@echo "  gh secret set CONFIG --env DEV < dev.yaml"
	@echo "  gh secret set CONFIG --env PROD < prod.yaml"

sync-config-all: ## 同步所有环境配置
	@echo "🔄 同步所有环境配置..."
	@echo "同步DEV环境配置..."
	@gh secret set CONFIG --env DEV < dev.yaml
	@echo "同步PROD环境配置..."
	@gh secret set CONFIG --env PROD < prod.yaml
	@echo "✅ 所有配置同步完成!"