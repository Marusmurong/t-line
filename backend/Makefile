.PHONY: dev build run test lint migrate-up migrate-down seed

APP_NAME := tline-api
BUILD_DIR := ./bin

dev: ## 热重载开发
	@which air > /dev/null 2>&1 || go install github.com/air-verse/air@latest
	air

build: ## 构建二进制
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

run: build ## 构建并运行
	$(BUILD_DIR)/$(APP_NAME)

test: ## 运行测试
	go test ./... -v -cover -race

lint: ## 代码检查
	golangci-lint run ./...

migrate-up: ## 执行数据库迁移
	bash scripts/migrate.sh up

migrate-down: ## 回滚数据库迁移
	bash scripts/migrate.sh down

seed: ## 导入种子数据
	go run scripts/seed.go

help: ## 显示帮助
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
