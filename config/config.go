package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// AppConfig 应用配置结构
type AppConfig struct {
	Port     int      `json:"port"`
	WebRoot  string   `json:"web_root"`
	LogPath  string   `json:"log_path"`
	Services []string `json:"services"`
}

// DefaultConfig 默认配置
var DefaultConfig = AppConfig{
	Port:     8080,
	WebRoot:  "/var/www/html",
	LogPath:  "/var/log",
	Services: []string{"nginx", "mysql", "php-fpm", "php7.4-fpm", "php8.0-fpm"},
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*AppConfig, error) {
	// 如果配置文件不存在，使用默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("配置文件 %s 不存在，使用默认配置", configPath)
		return &DefaultConfig, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// 合并默认配置
	if config.Port == 0 {
		config.Port = DefaultConfig.Port
	}
	if config.WebRoot == "" {
		config.WebRoot = DefaultConfig.WebRoot
	}
	if config.LogPath == "" {
		config.LogPath = DefaultConfig.LogPath
	}
	if len(config.Services) == 0 {
		config.Services = DefaultConfig.Services
	}

	return &config, nil
}

// SaveConfig 保存配置文件
func SaveConfig(config *AppConfig, configPath string) error {
	// 确保目录存在
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetServiceConfigPath 获取服务配置文件路径
func GetServiceConfigPath(serviceName string) string {
	configPaths := map[string]string{
		"nginx":      "/etc/nginx/nginx.conf",
		"mysql":      "/etc/mysql/my.cnf",
		"mariadb":    "/etc/mysql/my.cnf",
		"php-fpm":    "/etc/php/fpm/php-fpm.conf",
		"php7.4-fpm": "/etc/php/7.4/fpm/php-fpm.conf",
		"php8.0-fpm": "/etc/php/8.0/fpm/php-fpm.conf",
		"php8.1-fpm": "/etc/php/8.1/fpm/php-fpm.conf",
		"php8.2-fpm": "/etc/php/8.2/fpm/php-fpm.conf",
		"redis":      "/etc/redis/redis.conf",
		"memcached":  "/etc/memcached.conf",
		"postgresql": "/etc/postgresql/版本/main/postgresql.conf",
		"mongodb":    "/etc/mongod.conf",
	}

	if path, exists := configPaths[serviceName]; exists {
		return path
	}

	// 默认路径
	return "/etc/" + serviceName + "/" + serviceName + ".conf"
}

// ValidateConfig 验证配置
func ValidateConfig(config *AppConfig) error {
	if config.Port < 1 || config.Port > 65535 {
		return &ConfigError{Field: "port", Message: "端口必须在1-65535之间"}
	}

	if config.WebRoot == "" {
		return &ConfigError{Field: "web_root", Message: "Web根目录不能为空"}
	}

	if config.LogPath == "" {
		return &ConfigError{Field: "log_path", Message: "日志路径不能为空"}
	}

	return nil
}

// ConfigError 配置错误
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return e.Field + ": " + e.Message
}
