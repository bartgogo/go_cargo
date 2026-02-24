.PHONY: build run clean test dev

# 应用名称
APP_NAME=go-cargo
# 构建目录
BUILD_DIR=bin

# 构建
build:
	@echo "Building $(APP_NAME)..."
	@go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME).exe ./cmd/server/
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME).exe"

# 开发模式运行
run:
	@go run ./cmd/server/

# 开发模式 (带热重载, 需要 air)
dev:
	@air

# 清理构建产物
clean:
	@if exist $(BUILD_DIR) rd /s /q $(BUILD_DIR)
	@echo "Cleaned."

# 运行测试
test:
	@go test -v ./...

# 测试覆盖率
test-cover:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# 格式化代码
fmt:
	@go fmt ./...

# 静态检查
vet:
	@go vet ./...

# 安装依赖
deps:
	@go mod tidy
	@go mod download
