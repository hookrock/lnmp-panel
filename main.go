package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"lnmp-panel/config"
	"lnmp-panel/service"

	"github.com/gin-gonic/gin"
)

var (
	appConfig  *config.AppConfig
	svcManager *service.ServiceManager
	version    = "1.0.0"
	buildTime  = "unknown"
	gitCommit  = "unknown"
)

func main() {
	// 显示版本信息
	log.Printf("LNMP运维面板 v%s (构建时间: %s, Git提交: %s)", version, buildTime, gitCommit)

	// 检查是否为armv7l架构
	if runtime.GOARCH != "arm" && runtime.GOOS != "linux" {
		log.Println("警告：当前系统不是armv7l架构，部分功能可能无法正常工作")
	}

	// 初始化配置
	configPath := getConfigPath()
	var err error
	appConfig, err = config.LoadConfig(configPath)
	if err != nil {
		log.Printf("加载配置文件失败: %v，使用默认配置", err)
		appConfig = &config.DefaultConfig
	}

	// 验证配置
	if err := config.ValidateConfig(appConfig); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}

	// 初始化服务管理器（使用sudo权限）
	svcManager = service.NewServiceManager(true)

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// 静态文件服务
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	// API路由组
	api := router.Group("/api")
	{
		api.GET("/services", getServicesStatus)
		api.POST("/services/:name/start", startService)
		api.POST("/services/:name/stop", stopService)
		api.POST("/services/:name/restart", restartService)
		api.POST("/services/:name/reload", reloadService)
		api.POST("/services/:name/enable", enableService)
		api.POST("/services/:name/disable", disableService)
		api.GET("/config/:service", getServiceConfig)
		api.POST("/config/:service", updateServiceConfig)
		api.GET("/logs/:service", getServiceLogs)
		api.GET("/system/info", getSystemInfo)
		api.GET("/version", getVersionInfo)
	}

	// 页面路由
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   "LNMP运维面板",
			"version": version,
			"arch":    runtime.GOARCH,
		})
	})

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"arch":    runtime.GOARCH,
			"version": version,
		})
	})

	log.Printf("LNMP运维面板启动在端口 %d，架构: %s\n", appConfig.Port, runtime.GOARCH)

	// 启动服务器
	addr := fmt.Sprintf(":%d", appConfig.Port)
	if err := router.Run(addr); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}

// 获取配置路径
func getConfigPath() string {
	// 优先使用环境变量
	if path := os.Getenv("LNMP_PANEL_CONFIG"); path != "" {
		return path
	}

	// 其次使用当前目录的配置文件
	if _, err := os.Stat("config.json"); err == nil {
		return "config.json"
	}

	// 使用默认路径
	return "/etc/lnmp-panel/config.json"
}

// 获取服务状态
func getServicesStatus(c *gin.Context) {
	var services []service.ServiceStatus

	for _, serviceName := range appConfig.Services {
		status, err := svcManager.CheckServiceStatus(serviceName)
		if err != nil {
			log.Printf("检查服务状态失败 %s: %v", serviceName, err)
			services = append(services, service.ServiceStatus{
				Name:    serviceName,
				Status:  "error",
				Running: false,
			})
			continue
		}
		services = append(services, *status)
	}

	c.JSON(http.StatusOK, services)
}

// 启动服务
func startService(c *gin.Context) {
	serviceName := c.Param("name")

	if err := svcManager.StartService(serviceName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("启动服务失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "服务启动成功"})
}

// 停止服务
func stopService(c *gin.Context) {
	serviceName := c.Param("name")

	if err := svcManager.StopService(serviceName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("停止服务失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "服务停止成功"})
}

// 重启服务
func restartService(c *gin.Context) {
	serviceName := c.Param("name")

	if err := svcManager.RestartService(serviceName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("重启服务失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "服务重启成功"})
}

// 重载服务配置
func reloadService(c *gin.Context) {
	serviceName := c.Param("name")

	if err := svcManager.ReloadService(serviceName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("重载服务失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "服务重载成功"})
}

// 启用服务开机自启
func enableService(c *gin.Context) {
	serviceName := c.Param("name")

	if err := svcManager.EnableService(serviceName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("启用服务失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "服务启用成功"})
}

// 禁用服务开机自启
func disableService(c *gin.Context) {
	serviceName := c.Param("name")

	if err := svcManager.DisableService(serviceName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("禁用服务失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "服务禁用成功"})
}

// 获取服务配置
func getServiceConfig(c *gin.Context) {
	serviceName := c.Param("name")
	configPath := config.GetServiceConfigPath(serviceName)

	content, err := os.ReadFile(configPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("读取配置文件失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"config": string(content),
		"path":   configPath,
	})
}

// 更新服务配置
func updateServiceConfig(c *gin.Context) {
	serviceName := c.Param("name")
	configPath := config.GetServiceConfigPath(serviceName)

	var request struct {
		Config string `json:"config"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	err := os.WriteFile(configPath, []byte(request.Config), 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("更新配置文件失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "配置更新成功"})
}

// 获取服务日志
func getServiceLogs(c *gin.Context) {
	serviceName := c.Param("name")
	lines := c.DefaultQuery("lines", "100")

	lineCount := 100
	if l, err := fmt.Sscanf(lines, "%d", &lineCount); err == nil && l == 1 {
		// 解析成功
	} else {
		lineCount = 100
	}

	logs, err := svcManager.GetServiceLogs(serviceName, lineCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取日志失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
	})
}

// 获取系统信息
func getSystemInfo(c *gin.Context) {
	// 获取内存信息
	memCmd := exec.Command("free", "-h")
	memOutput, _ := memCmd.Output()

	// 获取磁盘信息
	diskCmd := exec.Command("df", "-h")
	diskOutput, _ := diskCmd.Output()

	// 获取CPU信息
	cpuCmd := exec.Command("lscpu")
	cpuOutput, _ := cpuCmd.Output()

	c.JSON(http.StatusOK, gin.H{
		"memory": strings.TrimSpace(string(memOutput)),
		"disk":   strings.TrimSpace(string(diskOutput)),
		"cpu":    strings.TrimSpace(string(cpuOutput)),
		"arch":   runtime.GOARCH,
		"os":     runtime.GOOS,
	})
}

// 获取版本信息
func getVersionInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":    version,
		"build_time": buildTime,
		"git_commit": gitCommit,
		"go_version": runtime.Version(),
		"arch":       runtime.GOARCH,
		"os":         runtime.GOOS,
	})
}
