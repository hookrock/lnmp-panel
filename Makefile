# LNMP运维面板构建配置

# 变量定义
BINARY_NAME=lnmp-panel
VERSION=1.0.0
BUILD_TIME=$(shell date +%Y%m%d%H%M%S)
GIT_COMMIT=$(shell git rev-parse --short HEAD)

# 构建目标
.PHONY: all build clean test install docker

all: build

# 构建ARMv7l版本
build-arm:
	@echo "构建ARMv7l版本..."
	GOARCH=arm GOOS=linux GOARM=7 go build -ldflags="-s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)" -o bin/$(BINARY_NAME)-armv7l main.go

# 构建AMD64版本
build-amd64:
	@echo "构建AMD64版本..."
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)" -o bin/$(BINARY_NAME)-amd64 main.go

# 构建本地测试版本
build-local:
	@echo "构建本地测试版本..."
	go build -ldflags="-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)" -o bin/$(BINARY_NAME)-local main.go

# 构建所有平台
build: clean build-arm build-amd64 build-local
	@echo "构建完成!"
	@ls -lh bin/

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	mkdir -p bin

# 运行测试
test:
	@echo "运行测试..."
	go test -v ./...

# 安装依赖
deps:
	@echo "安装Go依赖..."
	go mod tidy
	go mod download

# 代码检查
lint:
	@echo "代码检查..."
	golangci-lint run

# 安装到系统
install: build-arm
	@echo "安装到系统..."
	sudo cp bin/$(BINARY_NAME)-armv7l /usr/local/bin/$(BINARY_NAME)
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "安装完成!"

# 创建Docker镜像
docker-build:
	@echo "构建Docker镜像..."
	docker build -t $(BINARY_NAME):$(VERSION) .
	docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest

# 运行Docker容器
docker-run:
	@echo "运行Docker容器..."
	docker run -d \
		--name $(BINARY_NAME) \
		-p 8080:8080 \
		-v /etc/nginx:/etc/nginx \
		-v /etc/mysql:/etc/mysql \
		-v /etc/php:/etc/php \
		-v /var/log:/var/log \
		$(BINARY_NAME):latest

# 发布版本
release: clean deps test build
	@echo "创建发布版本 v$(VERSION)..."
	tar -czf $(BINARY_NAME)-v$(VERSION)-armv7l.tar.gz -C bin $(BINARY_NAME)-armv7l
	tar -czf $(BINARY_NAME)-v$(VERSION)-amd64.tar.gz -C bin $(BINARY_NAME)-amd64
	@echo "发布文件:"
	@ls -lh *.tar.gz

# 开发模式运行
dev:
	@echo "开发模式运行..."
	go run main.go

# 显示帮助信息
help:
	@echo "LNMP运维面板构建系统"
	@echo ""
	@echo "目标:"
	@echo "  build-arm     构建ARMv7l版本"
	@echo "  build-amd64   构建AMD64版本"
	@echo "  build-local   构建本地测试版本"
	@echo "  build         构建所有平台版本"
	@echo "  clean         清理构建文件"
	@echo "  test          运行测试"
	@echo "  deps          安装依赖"
	@echo "  lint          代码检查"
	@echo "  install       安装到系统"
	@echo "  docker-build  构建Docker镜像"
	@echo "  docker-run    运行Docker容器"
	@echo "  release       创建发布版本"
	@echo "  dev           开发模式运行"
	@echo "  help          显示帮助信息"