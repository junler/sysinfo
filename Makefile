.PHONY: build clean run-cli run-web test install help release tag

# 版本信息
VERSION := 1.0.0
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建标志
LDFLAGS := -ldflags="-s -w -X 'github.com/junler/sysinfo/cmd.Version=$(VERSION)' -X 'github.com/junler/sysinfo/cmd.BuildTime=$(BUILD_TIME)' -X 'github.com/junler/sysinfo/cmd.GitCommit=$(GIT_COMMIT)'"

# 默认目标
help:
	@echo "Available targets:"
	@echo "  build      - Build the sysinfo binary"
	@echo "  clean      - Remove build artifacts"
	@echo "  run-cli    - Run CLI version"
	@echo "  run-web    - Run web server on port 8080"
	@echo "  test       - Run tests"
	@echo "  install    - Install dependencies"
	@echo "  release    - Create a new release tag"
	@echo "  build-all  - Build for all platforms"
	@echo "  help       - Show this help message"

# 构建二进制文件
build:
	@echo "Building sysinfo v$(VERSION)..."
	go build $(LDFLAGS) -o sysinfo .
	@echo "Build complete: ./sysinfo"

# 清理构建产物
clean:
	@echo "Cleaning..."
	rm -f sysinfo
	go clean

# 运行CLI版本
run-cli: build
	@echo "Running sysinfo CLI..."
	./sysinfo info

# 运行Web服务器
run-web: build
	@echo "Starting web server..."
	./sysinfo serve --port 8080

# 运行测试
test:
	@echo "Running tests..."
	go test ./...

# 安装依赖
install:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# 创建发布标签
tag:
	@echo "Current version: $(VERSION)"
	@echo "Creating tag v$(VERSION)..."
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	@echo "Tag created. Push with: git push origin v$(VERSION)"

# 发布流程
release: test build
	@echo "Creating release for version $(VERSION)"
	@echo "1. Run 'make tag' to create the tag"
	@echo "2. Run 'git push origin v$(VERSION)' to trigger GitHub Actions"
	@echo "3. GitHub Actions will build and create the release automatically"

# 构建多平台版本
build-all: clean builds
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o builds/sysinfo-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o builds/sysinfo-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o builds/sysinfo-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o builds/sysinfo-windows-amd64.exe .
	@echo "Multi-platform builds complete in ./builds/"

# 创建builds目录
builds:
	mkdir -p builds
