#!/bin/bash

# LNMP运维面板安装脚本 - ARMv7l架构

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查系统架构
check_architecture() {
    local arch=$(uname -m)
    if [[ "$arch" != "armv7l" ]]; then
        log_warn "检测到系统架构: $arch"
        log_warn "本面板专为ARMv7l架构优化，其他架构可能无法正常工作"
        read -p "是否继续安装? (y/n): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        log_info "检测到ARMv7l架构，继续安装..."
    fi
}

# 检查依赖
check_dependencies() {
    log_info "检查系统依赖..."
    
    local missing_deps=()
    
    # 检查必要命令
    for cmd in systemctl wget curl; do
        if ! command -v $cmd &> /dev/null; then
            missing_deps+=($cmd)
        fi
    done
    
    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        log_error "缺少必要依赖: ${missing_deps[*]}"
        exit 1
    fi
    
    log_info "系统依赖检查通过"
}

# 创建系统用户
create_user() {
    if ! id "lnmp-panel" &>/dev/null; then
        log_info "创建系统用户: lnmp-panel"
        useradd -r -s /bin/false -d /opt/lnmp-panel lnmp-panel
    else
        log_info "用户lnmp-panel已存在"
    fi
}