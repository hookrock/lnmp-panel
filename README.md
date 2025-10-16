# LNMP运维面板 - ARMv7l架构优化版

专为ARMv7l架构设备优化的LNMP（Linux + Nginx + MySQL + PHP）运维管理面板，使用Go语言开发。

## 🎯 项目特色

- **ARMv7l深度优化**: 专为ARM架构设计，资源占用极低
- **完整运维功能**: 服务管理、配置编辑、日志查看一体化
- **现代化Web界面**: 响应式设计，支持移动端操作
- **一键部署**: 支持多种安装方式，快速上手

## 🚀 快速开始

### 方式一：一键安装（推荐）
```bash
# 下载安装脚本
wget https://raw.githubusercontent.com/your-repo/lnmp-panel/main/deploy/full_install.sh
chmod +x full_install.sh
sudo ./full_install.sh
```

### 方式二：Docker部署
```bash
docker run -d --name lnmp-panel -p 8080:8080 \
  -v /etc/nginx:/etc/nginx -v /etc/mysql:/etc/mysql \
  -v /etc/php:/etc/php -v /var/log:/var/log \
  your-repo/lnmp-panel:armv7l
```

### 方式三：手动安装
```bash
# 下载二进制文件
wget https://github.com/your-repo/lnmp-panel/releases/latest/lnmp-panel-armv7l
sudo mv lnmp-panel-armv7l /usr/local/bin/lnmp-panel
sudo chmod +x /usr/local/bin/lnmp-panel

# 启动服务
sudo lnmp-panel
```

访问地址: `http://你的服务器IP:8080`

## ✨ 核心功能

### 🔧 服务管理
- 实时监控Nginx、MySQL、PHP-FPM等服务状态
- 一键启动、停止、重启操作
- 开机自启管理

### 📝 配置管理  
- 在线编辑服务配置文件
- 语法高亮支持
- 配置验证和备份

### 📊 日志查看
- 实时日志监控
- 日志级别过滤
- 搜索和导出功能

### 📈 系统监控
- CPU、内存、磁盘使用率
- 网络连接状态
- 系统负载监控

## 🛠 技术架构

### 后端技术栈
- **语言**: Go 1.19+
- **框架**: Gin Web Framework
- **服务管理**: systemd集成
- **配置格式**: JSON

### 前端技术栈  
- **UI框架**: Bootstrap 5
- **交互**: 原生JavaScript
- **样式**: CSS3现代化设计
- **兼容性**: 移动端适配

## 📋 系统要求

### 硬件要求
- **架构**: ARMv7l (兼容其他架构)
- **内存**: 最低512MB，推荐1GB+
- **存储**: 100MB可用空间

### 软件要求
- **操作系统**: Linux (Debian/Ubuntu/CentOS/Raspberry Pi OS等)
- **依赖服务**: systemd, curl, wget
- **可选服务**: Nginx, MySQL, PHP-FPM

## 🔧 配置说明

### 主要配置项
```json
{
  "port": 8080,
  "web_root": "/var/www/html", 
  "log_path": "/var/log",
  "services": ["nginx", "mysql", "php-fpm"]
}
```

### 环境变量支持
- `LNMP_PANEL_CONFIG`: 自定义配置文件路径
- `PORT`: 服务端口设置

## 📖 详细文档

- [📚 使用教程](docs/使用教程.md) - 完整的功能使用指南
- [🔧 API文档](docs/API文档.md) - 详细的API接口说明  
- [🛠 故障排除](docs/故障排除指南.md) - 常见问题解决方案
- [🏗 开发指南](项目总结.md) - 项目架构和扩展开发

## 🐛 故障排除

### 常见问题
1. **面板无法访问**: 检查防火墙8080端口
2. **服务状态异常**: 确认对应服务已安装
3. **配置保存失败**: 检查文件权限设置

### 获取帮助
- 查看详细故障排除指南
- 检查系统日志: `journalctl -u lnmp-panel`
- 提交Issue到项目仓库

## 🤝 参与贡献

欢迎提交Issue和Pull Request！

### 开发环境搭建
```bash
git clone https://github.com/your-repo/lnmp-panel.git
cd lnmp-panel
go mod tidy
go run main.go
```

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 🏆 版本信息

**当前版本**: v1.0.0  
**发布日期**: 2024-01-01  
**支持架构**: ARMv7l (主要) / AMD64 / ARM64

---

**项目状态**: ✅ 生产就绪  
**文档完整性**: ✅ 完整中文文档  
**测试覆盖率**: ✅ 核心功能测试通过

访问 [GitHub仓库](https://github.com/your-repo/lnmp-panel) 获取最新版本和更新！