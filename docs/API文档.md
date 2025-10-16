# LNMP运维面板 API 文档

## 基础信息

- **基础URL**: `http://localhost:8080/api`
- **认证方式**: 当前版本无需认证（生产环境建议配置认证）
- **响应格式**: JSON

## API端点

### 1. 系统信息

#### 获取系统信息
- **端点**: `GET /api/system/info`
- **描述**: 获取系统架构、资源使用等信息
- **响应示例**:
```json
{
  "memory": "内存使用信息",
  "disk": "磁盘使用信息", 
  "cpu": "CPU信息",
  "arch": "arm",
  "os": "linux"
}
```

### 2. 服务管理

#### 获取所有服务状态
- **端点**: `GET /api/services`
- **描述**: 获取所有配置服务的状态信息
- **响应示例**:
```json
[
  {
    "name": "nginx",
    "status": "active",
    "running": true,
    "description": "A high performance web server and a reverse proxy server"
  }
]
```

#### 启动服务
- **端点**: `POST /api/services/{name}/start`
- **描述**: 启动指定名称的服务
- **参数**: `name` - 服务名称
- **响应示例**:
```json
{
  "message": "服务启动成功"
}
```

#### 停止服务
- **端点**: `POST /api/services/{name}/stop`
- **描述**: 停止指定名称的服务
- **参数**: `name` - 服务名称
- **响应示例**:
```json
{
  "message": "服务停止成功"
}
```

#### 重启服务
- **端点**: `POST /api/services/{name}/restart`
- **描述**: 重启指定名称的服务
- **参数**: `name` - 服务名称
- **响应示例**:
```json
{
  "message": "服务重启成功"
}
```

### 3. 配置管理

#### 获取服务配置
- **端点**: `GET /api/config/{name}`
- **描述**: 获取指定服务的配置文件内容
- **参数**: `name` - 服务名称
- **响应示例**:
```json
{
  "config": "配置文件内容",
  "path": "/etc/nginx/nginx.conf"
}
```

#### 更新服务配置
- **端点**: `POST /api/config/{name}`
- **描述**: 更新指定服务的配置文件
- **参数**: 
  - `name` - 服务名称（路径参数）
  - `config` - 新的配置内容（JSON body）
- **请求体示例**:
```json
{
  "config": "新的配置文件内容"
}
```
- **响应示例**:
```json
{
  "message": "配置更新成功"
}
```

### 4. 日志管理

#### 获取服务日志
- **端点**: `GET /api/logs/{name}`
- **描述**: 获取指定服务的日志内容
- **参数**:
  - `name` - 服务名称（路径参数）
  - `lines` - 日志行数（查询参数，默认100）
- **响应示例**:
```json
{
  "logs": "日志内容..."
}
```

### 5. 健康检查

#### 健康检查端点
- **端点**: `GET /health`
- **描述**: 检查面板服务是否正常运行
- **响应示例**:
```json
{
  "status": "ok",
  "arch": "arm"
}
```

## 错误处理

### 错误响应格式
```json
{
  "error": "错误描述信息",
  "code": "错误代码（可选）"
}
```

### 常见错误代码

| 错误代码 | 描述 | HTTP状态码 |
|---------|------|------------|
| SERVICE_NOT_FOUND | 服务不存在 | 404 |
| CONFIG_READ_ERROR | 配置文件读取失败 | 500 |
| CONFIG_WRITE_ERROR | 配置文件写入失败 | 500 |
| PERMISSION_DENIED | 权限不足 | 403 |
| INVALID_REQUEST | 无效的请求参数 | 400 |

## 使用示例

### cURL示例

#### 获取服务状态
```bash
curl http://localhost:8080/api/services
```

#### 启动Nginx服务
```bash
curl -X POST http://localhost:8080/api/services/nginx/start
```

#### 获取Nginx配置
```bash
curl http://localhost:8080/api/config/nginx
```

#### 更新Nginx配置
```bash
curl -X POST http://localhost:8080/api/config/nginx \
  -H "Content-Type: application/json" \
  -d '{"config": "新的配置内容"}'
```

#### 获取Nginx日志（最近50行）
```bash
curl "http://localhost:8080/api/logs/nginx?lines=50"
```

### JavaScript示例

#### 获取所有服务状态
```javascript
fetch('/api/services')
  .then(response => response.json())
  .then(services => {
    console.log('服务状态:', services);
  });
```

#### 重启MySQL服务
```javascript
fetch('/api/services/mysql/restart', {
  method: 'POST'
})
.then(response => response.json())
.then(result => {
  console.log('操作结果:', result);
});
```

#### 更新PHP-FPM配置
```javascript
const newConfig = `[www]
user = www-data
group = www-data
listen = /run/php/php8.1-fpm.sock`;

fetch('/api/config/php-fpm', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({ config: newConfig })
})
.then(response => response.json())
.then(result => {
  console.log('配置更新结果:', result);
});
```

## 安全注意事项

1. **生产环境部署**:
   - 建议配置反向代理（如Nginx）
   - 启用HTTPS加密
   - 配置访问控制列表

2. **API安全**:
   - 限制访问IP范围
   - 考虑添加API密钥认证
   - 记录API访问日志

3. **权限管理**:
   - 确保面板进程有适当的文件权限
   - 避免使用root权限运行
   - 定期审查权限设置

## 性能优化建议

1. **API调用优化**:
   - 批量获取服务状态，避免频繁调用
   - 使用适当的缓存策略
   - 限制日志查询行数

2. **资源使用**:
   - 配置合理的并发连接数
   - 监控内存和CPU使用
   - 定期清理临时文件

## 扩展开发

### 添加新的API端点

要添加新的API端点，修改 `main.go` 文件：

```go
// 添加新的路由
api.GET("/api/custom/endpoint", customHandler)

// 实现处理函数
func customHandler(c *gin.Context) {
    // 处理逻辑
    c.JSON(http.StatusOK, gin.H{
        "message": "自定义端点",
    })
}
```

### 自定义服务支持

要支持新的服务类型，更新配置映射：

```go
// 在 config/config.go 中添加
configPaths := map[string]string{
    "custom-service": "/etc/custom/service.conf",
    // ... 其他服务
}
```

## 版本历史

- v1.0.0: 初始版本，包含基础服务管理功能
- 未来版本计划添加认证、插件系统等高级功能