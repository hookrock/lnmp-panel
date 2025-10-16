#!/bin/bash

# LNMP运维面板完整安装脚本 - ARMv7l架构优化版

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示横幅
show_banner() {
    echo -e "${BLUE}"
    echo "================================================"
    echo "    LNMP运维面板安装程序 - ARMv7l优化版"
    echo "================================================"
    echo -e "${NC}"
}

# 检查系统要求
check_requirements() {
    log_info "检查系统要求..."
    
    # 检查操作系统
    if [[ "$(uname -s)" != "Linux" ]]; then
        log_error "只支持Linux系统"
        exit 1
    fi
    
    # 检查架构
    local arch=$(uname -m)
    if [[ "$arch" != "armv7l" ]]; then
        log_warn "检测到系统架构: $arch"
        log_warn "本面板专为ARMv7l架构优化，其他架构可能无法达到最佳性能"
        read -p "是否继续安装? (y/n): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        log_success "检测到ARMv7l架构，继续安装..."
    fi
    
    # 检查内存
    local mem_kb=$(grep MemTotal /proc/meminfo | awk '{print $2}')
    local mem_mb=$((mem_kb / 1024))
    
    if [[ $mem_mb -lt 512 ]]; then
        log_warn "系统内存较低（${mem_mb}MB），建议至少512MB内存"
    else
        log_success "系统内存: ${mem_mb}MB"
    fi
    
    # 检查磁盘空间
    local disk_free=$(df / | awk 'NR==2 {print $4}')
    if [[ $disk_free -lt 1048576 ]]; then # 小于1GB
        log_warn "磁盘空间较低，建议确保有足够空间"
    fi
}

# 安装依赖
install_dependencies() {
    log_info "安装系统依赖..."
    
    # 检测包管理器
    if command -v apt &> /dev/null; then
        # Debian/Ubuntu
        sudo apt update
        sudo apt install -y wget curl systemd
    elif command -v yum &> /dev/null; then
        # CentOS/RHEL
        sudo yum install -y wget curl systemd
    elif command -v dnf &> /dev/null; then
        # Fedora
        sudo dnf install -y wget curl systemd
    elif command -v apk &> /dev/null; then
        # Alpine
        sudo apk add wget curl
    else
        log_error "不支持的包管理器"
        exit 1
    fi
    
    log_success "系统依赖安装完成"
}

# 创建系统用户和目录
setup_environment() {
    log_info "设置运行环境..."
    
    # 创建系统用户
    if ! id "lnmp-panel" &>/dev/null; then
        sudo useradd -r -s /bin/false -d /opt/lnmp-panel lnmp-panel
        log_success "创建系统用户: lnmp-panel"
    else
        log_info "用户lnmp-panel已存在"
    fi
    
    # 创建安装目录
    sudo mkdir -p /opt/lnmp-panel
    sudo mkdir -p /etc/lnmp-panel
    sudo mkdir -p /var/log/lnmp-panel
    
    log_success "环境设置完成"
}

# 下载并安装面板
install_panel() {
    log_info "下载LNMP运维面板..."
    
    local version="1.0.0"
    local download_url="https://github.com/your-repo/lnmp-panel/releases/download/v${version}/lnmp-panel-armv7l"
    
    cd /opt/lnmp-panel
    
    # 下载面板二进制文件
    if sudo wget -q "$download_url" -O lnmp-panel; then
        sudo chmod +x lnmp-panel
        log_success "面板下载成功"
    else
        log_error "面板下载失败，尝试从源码构建..."
        build_from_source
    fi
}

# 从源码构建
build_from_source() {
    log_info "从源码构建面板..."
    
    # 安装Go环境
    if ! command -v go &> /dev/null; then
        log_info "安装Go语言环境..."
        install_go
    fi
    
    # 下载源码
    cd /tmp
    git clone https://github.com/your-repo/lnmp-panel.git
    cd lnmp-panel
    
    # 构建ARMv7l版本
    GOARCH=arm GOOS=linux GOARM=7 go build -o lnmp-panel main.go
    
    # 复制到安装目录
    sudo cp lnmp-panel /opt/lnmp-panel/
    sudo chmod +x /opt/lnmp-panel/lnmp-panel
    
    log_success "从源码构建完成"
}

# 安装Go环境
install_go() {
    local go_version="1.19"
    local arch=$(uname -m)
    
    if [[ "$arch" == "armv7l" ]]; then
        arch="armv6l" # Go的ARM版本使用armv6l
    fi
    
    local go_package="go${go_version}.linux-${arch}.tar.gz"
    local download_url="https://golang.org/dl/${go_package}"
    
    cd /tmp
    wget -q "$download_url"
    sudo tar -C /usr/local -xzf "$go_package"
    
    # 设置环境变量
    echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
    echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /home/$(whoami)/.profile
    
    source /etc/profile
    log_success "Go环境安装完成"
}

# 创建配置文件
create_config() {
    log_info "创建配置文件..."
    
    sudo cat > /etc/lnmp-panel/config.json << EOF
{
    "port": 8080,
    "web_root": "/var/www/html",
    "log_path": "/var/log",
    "services": ["nginx", "mysql", "php-fpm", "php7.4-fpm", "php8.0-fpm"]
}
EOF
    
    # 创建示例Nginx配置
    sudo mkdir -p /etc/lnmp-panel/examples
    sudo cat > /etc/lnmp-panel/examples/nginx.conf << 'EOF'
# LNMP优化配置示例
user www-data;
worker_processes 1;  # ARM设备建议设置为1

events {
    worker_connections 512;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;

    # 限制缓冲区大小，减少内存占用
    client_body_buffer_size 8K;
    client_header_buffer_size 1k;
    client_max_body_size 8m;

    gzip on;
    gzip_min_length 1k;
    gzip_types text/plain text/css application/json;

    include /etc/nginx/conf.d/*.conf;
    include /etc/nginx/sites-enabled/*;
}
EOF
    
    sudo chown -R lnmp-panel:lnmp-panel /etc/lnmp-panel
    log_success "配置文件创建完成"
}

# 创建系统服务
create_systemd_service() {
    log_info "创建系统服务..."
    
    sudo cat > /etc/systemd/system/lnmp-panel.service << EOF
[Unit]
Description=LNMP运维面板
Documentation=https://github.com/your-repo/lnmp-panel
After=network.target
Wants=network.target

[Service]
Type=simple
User=lnmp-panel
Group=lnmp-panel
WorkingDirectory=/opt/lnmp-panel
ExecStart=/opt/lnmp-panel/lnmp-panel
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

# ARM设备资源限制
LimitNOFILE=65536
MemoryLimit=256M
CPUQuota=80%

[Install]
WantedBy=multi-user.target
EOF
    
    sudo systemctl daemon-reload
    log_success "系统服务创建完成"
}

# 配置防火墙
setup_firewall() {
    log_info "配置防火墙..."
    
    if command -v ufw &> /dev/null; then
        sudo ufw allow 8080/tcp
        sudo ufw reload
        log_success "UFW防火墙配置完成"
    elif command -v firewall-cmd &> /dev/null; then
        sudo firewall-cmd --permanent --add-port=8080/tcp
        sudo firewall-cmd --reload
        log_success "firewalld配置完成"
    else
        log_warn "未检测到支持的防火墙，请手动开放8080端口"
    fi
}

# 设置权限
setup_permissions() {
    log_info "设置文件权限..."
    
    sudo chown -R lnmp-panel:lnmp-panel /opt/lnmp-panel
    sudo chmod 755 /opt/lnmp-panel
    sudo chmod 600 /etc/lnmp-panel/config.json
    
    # 设置日志目录权限
    sudo chown lnmp-panel:lnmp-panel /var/log/lnmp-panel
    sudo chmod 755 /var/log/lnmp-panel
    
    log_success "权限设置完成"
}

# 启动服务
start_service() {
    log_info "启动LNMP运维面板服务..."
    
    sudo systemctl enable lnmp-panel
    sudo systemctl start lnmp-panel
    
    # 等待服务启动
    sleep 3
    
    # 检查服务状态
    if sudo systemctl is-active --quiet lnmp-panel; then
        log_success "LNMP运维面板启动成功"
    else
        log_error "LNMP运维面板启动失败"
        sudo systemctl status lnmp-panel
        exit 1
    fi
}

# 运行健康检查
health_check() {
    log_info "运行健康检查..."
    
    local max_attempts=10
    local attempt=1
    
    while [[ $attempt -le $max_attempts ]]; do
        if curl -f -s http://localhost:8080/health > /dev/null; then
            log_success "健康检查通过"
            return 0
        fi
        
        log_info "等待服务启动... ($attempt/$max_attempts)"
        sleep 5
        ((attempt++))
    done
    
    log_error "健康检查失败，服务未正常启动"
    return 1
}

# 显示安装完成信息
show_completion() {
    local ip=$(hostname -I | awk '{print $1}')
    
    echo -e "${GREEN}"
    echo "================================================"
    echo "      LNMP运维面板安装完成！"
    echo "================================================"
    echo -e "${NC}"
    echo -e "${BLUE}访问信息:${NC}"
    echo "  面板地址: http://${ip}:8080"
    echo "  本地访问: http://localhost:8080"
    echo ""
    echo -e "${BLUE}管理命令:${NC}"
    echo "  sudo systemctl start lnmp-panel    # 启动服务"
    echo "  sudo systemctl stop lnmp-panel     # 停止服务"
    echo "  sudo systemctl restart lnmp-panel  # 重启服务"
    echo "  sudo systemctl status lnmp-panel   # 查看状态"
    echo "  journalctl -u lnmp-panel -f        # 查看日志"
    echo ""
    echo -e "${BLUE}配置文件:${NC}"
    echo "  主配置: /etc/lnmp-panel/config.json"
    echo "  日志文件: /var/log/lnmp-panel/"
    echo ""
    echo -e "${BLUE}下一步:${NC}"
    echo "  1. 在浏览器中访问面板地址"
    echo "  2. 检查LNMP服务状态"
    echo "  3. 根据需要调整配置"
    echo -e "${GREEN}"
    echo "================================================"
    echo -e "${NC}"
}

# 主安装流程
main() {
    show_banner
    check_requirements
    install_dependencies
    setup_environment
    install_panel
    create_config
    create_systemd_service
    setup_firewall
    setup_permissions
    start_service
    health_check
    show_completion
    
    log_success "安装完成！"
}

# 错误处理
trap 'log_error "安装过程被中断"; exit 1' INT TERM

# 执行主函数
main "$@"