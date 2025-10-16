package service

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

// ServiceManager 服务管理器
type ServiceManager struct {
	UseSudo bool
}

// ServiceStatus 服务状态
type ServiceStatus struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Running     bool   `json:"running"`
	Description string `json:"description,omitempty"`
}

// NewServiceManager 创建服务管理器
func NewServiceManager(useSudo bool) *ServiceManager {
	return &ServiceManager{UseSudo: useSudo}
}

// CheckServiceStatus 检查服务状态
func (sm *ServiceManager) CheckServiceStatus(serviceName string) (*ServiceStatus, error) {
	// 检查系统类型
	if runtime.GOOS != "linux" {
		return &ServiceStatus{
			Name:    serviceName,
			Status:  "unknown",
			Running: false,
		}, fmt.Errorf("非Linux系统不支持服务管理")
	}

	// 检查服务是否存在
	if !sm.serviceExists(serviceName) {
		return &ServiceStatus{
			Name:    serviceName,
			Status:  "not-found",
			Running: false,
		}, fmt.Errorf("服务不存在: %s", serviceName)
	}

	// 检查服务状态
	status, err := sm.getServiceStatus(serviceName)
	if err != nil {
		return &ServiceStatus{
			Name:    serviceName,
			Status:  "error",
			Running: false,
		}, err
	}

	running := status == "active"

	// 获取服务描述
	description := sm.getServiceDescription(serviceName)

	return &ServiceStatus{
		Name:        serviceName,
		Status:      status,
		Running:     running,
		Description: description,
	}, nil
}

// StartService 启动服务
func (sm *ServiceManager) StartService(serviceName string) error {
	if !sm.serviceExists(serviceName) {
		return fmt.Errorf("服务不存在: %s", serviceName)
	}

	var cmd *exec.Cmd
	if sm.UseSudo {
		cmd = exec.Command("sudo", "systemctl", "start", serviceName)
	} else {
		cmd = exec.Command("systemctl", "start", serviceName)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("启动服务失败: %v", err)
	}

	return nil
}

// StopService 停止服务
func (sm *ServiceManager) StopService(serviceName string) error {
	if !sm.serviceExists(serviceName) {
		return fmt.Errorf("服务不存在: %s", serviceName)
	}

	var cmd *exec.Cmd
	if sm.UseSudo {
		cmd = exec.Command("sudo", "systemctl", "stop", serviceName)
	} else {
		cmd = exec.Command("systemctl", "stop", serviceName)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("停止服务失败: %v", err)
	}

	return nil
}

// RestartService 重启服务
func (sm *ServiceManager) RestartService(serviceName string) error {
	if !sm.serviceExists(serviceName) {
		return fmt.Errorf("服务不存在: %s", serviceName)
	}

	var cmd *exec.Cmd
	if sm.UseSudo {
		cmd = exec.Command("sudo", "systemctl", "restart", serviceName)
	} else {
		cmd = exec.Command("systemctl", "restart", serviceName)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("重启服务失败: %v", err)
	}

	return nil
}

// ReloadService 重载服务配置
func (sm *ServiceManager) ReloadService(serviceName string) error {
	if !sm.serviceExists(serviceName) {
		return fmt.Errorf("服务不存在: %s", serviceName)
	}

	var cmd *exec.Cmd
	if sm.UseSudo {
		cmd = exec.Command("sudo", "systemctl", "reload", serviceName)
	} else {
		cmd = exec.Command("systemctl", "reload", serviceName)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("重载服务失败: %v", err)
	}

	return nil
}

// EnableService 启用服务开机自启
func (sm *ServiceManager) EnableService(serviceName string) error {
	if !sm.serviceExists(serviceName) {
		return fmt.Errorf("服务不存在: %s", serviceName)
	}

	var cmd *exec.Cmd
	if sm.UseSudo {
		cmd = exec.Command("sudo", "systemctl", "enable", serviceName)
	} else {
		cmd = exec.Command("systemctl", "enable", serviceName)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("启用服务失败: %v", err)
	}

	return nil
}

// DisableService 禁用服务开机自启
func (sm *ServiceManager) DisableService(serviceName string) error {
	if !sm.serviceExists(serviceName) {
		return fmt.Errorf("服务不存在: %s", serviceName)
	}

	var cmd *exec.Cmd
	if sm.UseSudo {
		cmd = exec.Command("sudo", "systemctl", "disable", serviceName)
	} else {
		cmd = exec.Command("systemctl", "disable", serviceName)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("禁用服务失败: %v", err)
	}

	return nil
}

// GetServiceLogs 获取服务日志
func (sm *ServiceManager) GetServiceLogs(serviceName string, lines int) (string, error) {
	if !sm.serviceExists(serviceName) {
		return "", fmt.Errorf("服务不存在: %s", serviceName)
	}

	var cmd *exec.Cmd
	lineArg := fmt.Sprintf("%d", lines)

	if sm.UseSudo {
		cmd = exec.Command("sudo", "journalctl", "-u", serviceName, "-n", lineArg, "--no-pager")
	} else {
		cmd = exec.Command("journalctl", "-u", serviceName, "-n", lineArg, "--no-pager")
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("获取日志失败: %v", err)
	}

	return string(output), nil
}

// 检查服务是否存在
func (sm *ServiceManager) serviceExists(serviceName string) bool {
	var cmd *exec.Cmd
	if sm.UseSudo {
		cmd = exec.Command("sudo", "systemctl", "list-unit-files", serviceName+".service")
	} else {
		cmd = exec.Command("systemctl", "list-unit-files", serviceName+".service")
	}

	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), serviceName+".service")
}

// 获取服务状态
func (sm *ServiceManager) getServiceStatus(serviceName string) (string, error) {
	var cmd *exec.Cmd
	if sm.UseSudo {
		cmd = exec.Command("sudo", "systemctl", "is-active", serviceName)
	} else {
		cmd = exec.Command("systemctl", "is-active", serviceName)
	}

	output, err := cmd.Output()
	if err != nil {
		// 如果命令执行失败，可能是服务不存在或权限问题
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 3 { // systemctl is-active 返回3表示服务不存在
					return "not-found", nil
				}
			}
		}
		return "unknown", err
	}

	status := strings.TrimSpace(string(output))
	if status == "active" || status == "inactive" || status == "failed" {
		return status, nil
	}

	return "unknown", nil
}

// 获取服务描述
func (sm *ServiceManager) getServiceDescription(serviceName string) string {
	var cmd *exec.Cmd
	if sm.UseSudo {
		cmd = exec.Command("sudo", "systemctl", "show", serviceName, "--property=Description")
	} else {
		cmd = exec.Command("systemctl", "show", serviceName, "--property=Description")
	}

	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(output), "=")
	if len(lines) >= 2 {
		return strings.TrimSpace(lines[1])
	}

	return ""
}
