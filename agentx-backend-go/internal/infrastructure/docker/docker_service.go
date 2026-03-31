package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	domainContainer "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/container"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/config"
	"go.uber.org/zap"
)

// DockerService Docker服务接口（对应 Java 的 DockerService）
type DockerService interface {
	// CreateContainer 创建容器
	CreateContainer(template *domainContainer.ContainerTemplateEntity, containerName string) (dockerContainerID string, err error)
	// StartContainer 启动容器
	StartContainer(dockerContainerID string) error
	// StopContainer 停止容器
	StopContainer(dockerContainerID string) error
	// RemoveContainer 删除容器
	RemoveContainer(dockerContainerID string) error
	// GetContainerLogs 获取容器日志
	GetContainerLogs(dockerContainerID string, lines int) (string, error)
	// GetContainerStats 获取容器资源使用情况
	GetContainerStats(dockerContainerID string) (*ContainerStats, error)
	// IsContainerRunning 检查容器是否运行中
	IsContainerRunning(dockerContainerID string) (bool, error)
}

// ContainerStats 容器资源统计
type ContainerStats struct {
	CPUUsage    float64 `json:"cpuUsage"`
	MemoryUsage float64 `json:"memoryUsage"`
	MemoryLimit int64   `json:"memoryLimit"`
	NetworkIn   int64   `json:"networkIn"`
	NetworkOut  int64   `json:"networkOut"`
}

// ---- Docker Engine API 客户端实现 ----

// DockerEngineService 基于 Docker Engine API 的实际实现
// 通过 HTTP 直接调用 Docker daemon（无需第三方 SDK）
type DockerEngineService struct {
	httpClient *http.Client
	host       string // Docker daemon 地址
	apiVersion string // API 版本
	logger     *zap.Logger
}

// NewDockerEngineService 创建 Docker Engine 服务
func NewDockerEngineService(dockerCfg *config.DockerConfig, logger *zap.Logger) DockerService {
	host := dockerCfg.Host
	if host == "" {
		host = "unix:///var/run/docker.sock"
	}

	var httpClient *http.Client

	if strings.HasPrefix(host, "unix://") {
		// Unix socket 连接
		socketPath := strings.TrimPrefix(host, "unix://")
		httpClient = &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.DialTimeout("unix", socketPath, 10*time.Second)
				},
			},
			Timeout: 30 * time.Second,
		}
		host = "http://localhost" // Unix socket 使用 localhost 作为 host
	} else if strings.HasPrefix(host, "tcp://") {
		// TCP 连接
		host = "http" + strings.TrimPrefix(host, "tcp")
		httpClient = &http.Client{Timeout: 30 * time.Second}
	} else {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	svc := &DockerEngineService{
		httpClient: httpClient,
		host:       host,
		apiVersion: "v1.41",
		logger:     logger,
	}

	// 测试连接
	if err := svc.ping(); err != nil {
		logger.Warn("Docker daemon 连接失败，容器功能将不可用",
			zap.String("host", dockerCfg.Host),
			zap.Error(err),
		)
	} else {
		logger.Info("Docker daemon 连接成功", zap.String("host", dockerCfg.Host))
	}

	return svc
}

// ping 测试 Docker daemon 连接
func (s *DockerEngineService) ping() error {
	resp, err := s.httpClient.Get(s.host + "/_ping")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Docker daemon ping 失败: %d", resp.StatusCode)
	}
	return nil
}

// apiURL 构建 API URL
func (s *DockerEngineService) apiURL(path string) string {
	return fmt.Sprintf("%s/%s%s", s.host, s.apiVersion, path)
}

func (s *DockerEngineService) CreateContainer(template *domainContainer.ContainerTemplateEntity, containerName string) (string, error) {
	image := template.GetFullImageName()

	// 构建创建容器的请求体
	createReq := map[string]interface{}{
		"Image": image,
	}

	hostConfig := map[string]interface{}{}
	if template.MemoryLimit != nil && *template.MemoryLimit > 0 {
		hostConfig["Memory"] = int64(*template.MemoryLimit) * 1024 * 1024 // MB -> Bytes
	}
	if template.CPULimit != nil && *template.CPULimit > 0 {
		hostConfig["NanoCpus"] = int64(*template.CPULimit * 1e9) // CPU 核数 -> 纳秒
	}
	if len(hostConfig) > 0 {
		createReq["HostConfig"] = hostConfig
	}

	body, err := json.Marshal(createReq)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	url := s.apiURL(fmt.Sprintf("/containers/create?name=%s", containerName))
	resp, err := s.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("创建容器请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("创建容器失败 (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		ID       string   `json:"Id"`
		Warnings []string `json:"Warnings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析创建容器响应失败: %w", err)
	}

	s.logger.Info("Docker容器创建成功",
		zap.String("containerId", result.ID),
		zap.String("name", containerName),
		zap.String("image", image),
	)

	return result.ID, nil
}

func (s *DockerEngineService) StartContainer(dockerContainerID string) error {
	url := s.apiURL(fmt.Sprintf("/containers/%s/start", dockerContainerID))
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("启动容器请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 204=成功, 304=已在运行
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotModified {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("启动容器失败 (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	s.logger.Info("Docker容器启动成功", zap.String("containerId", dockerContainerID))
	return nil
}

func (s *DockerEngineService) StopContainer(dockerContainerID string) error {
	url := s.apiURL(fmt.Sprintf("/containers/%s/stop?t=10", dockerContainerID))
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("停止容器请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 204=成功, 304=已停止
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotModified {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("停止容器失败 (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	s.logger.Info("Docker容器停止成功", zap.String("containerId", dockerContainerID))
	return nil
}

func (s *DockerEngineService) RemoveContainer(dockerContainerID string) error {
	url := s.apiURL(fmt.Sprintf("/containers/%s?force=true&v=true", dockerContainerID))
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("删除容器请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("删除容器失败 (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	s.logger.Info("Docker容器删除成功", zap.String("containerId", dockerContainerID))
	return nil
}

func (s *DockerEngineService) GetContainerLogs(dockerContainerID string, lines int) (string, error) {
	url := s.apiURL(fmt.Sprintf("/containers/%s/logs?stdout=true&stderr=true&tail=%d", dockerContainerID, lines))
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("获取容器日志请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("获取容器日志失败 (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	logBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取容器日志失败: %w", err)
	}

	// Docker 日志流有 8 字节头部，需要清理
	return cleanDockerLogs(logBytes), nil
}

func (s *DockerEngineService) GetContainerStats(dockerContainerID string) (*ContainerStats, error) {
	url := s.apiURL(fmt.Sprintf("/containers/%s/stats?stream=false", dockerContainerID))
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("获取容器统计请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取容器统计失败 (HTTP %d)", resp.StatusCode)
	}

	var stats struct {
		CPUStats struct {
			CPUUsage struct {
				TotalUsage int64 `json:"total_usage"`
			} `json:"cpu_usage"`
			SystemCPUUsage int64 `json:"system_cpu_usage"`
		} `json:"cpu_stats"`
		PreCPUStats struct {
			CPUUsage struct {
				TotalUsage int64 `json:"total_usage"`
			} `json:"cpu_usage"`
			SystemCPUUsage int64 `json:"system_cpu_usage"`
		} `json:"precpu_stats"`
		MemoryStats struct {
			Usage int64 `json:"usage"`
			Limit int64 `json:"limit"`
		} `json:"memory_stats"`
		Networks map[string]struct {
			RxBytes int64 `json:"rx_bytes"`
			TxBytes int64 `json:"tx_bytes"`
		} `json:"networks"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("解析容器统计失败: %w", err)
	}

	// 计算 CPU 使用率
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemCPUUsage - stats.PreCPUStats.SystemCPUUsage)
	cpuUsage := 0.0
	if systemDelta > 0 {
		cpuUsage = (cpuDelta / systemDelta) * 100.0
	}

	// 计算网络流量
	var networkIn, networkOut int64
	for _, net := range stats.Networks {
		networkIn += net.RxBytes
		networkOut += net.TxBytes
	}

	return &ContainerStats{
		CPUUsage:    cpuUsage,
		MemoryUsage: float64(stats.MemoryStats.Usage) / 1024 / 1024, // Bytes -> MB
		MemoryLimit: stats.MemoryStats.Limit / 1024 / 1024,          // Bytes -> MB
		NetworkIn:   networkIn,
		NetworkOut:  networkOut,
	}, nil
}

func (s *DockerEngineService) IsContainerRunning(dockerContainerID string) (bool, error) {
	url := s.apiURL(fmt.Sprintf("/containers/%s/json", dockerContainerID))
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return false, fmt.Errorf("检查容器状态请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("检查容器状态失败 (HTTP %d)", resp.StatusCode)
	}

	var info struct {
		State struct {
			Running bool `json:"Running"`
		} `json:"State"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return false, err
	}

	return info.State.Running, nil
}

// cleanDockerLogs 清理 Docker 日志流的 8 字节头部
func cleanDockerLogs(raw []byte) string {
	var result strings.Builder
	for len(raw) > 0 {
		if len(raw) < 8 {
			result.Write(raw)
			break
		}
		// Docker 日志头: [stream_type(1), 0, 0, 0, size(4)]
		size := int(raw[4])<<24 | int(raw[5])<<16 | int(raw[6])<<8 | int(raw[7])
		raw = raw[8:]
		if size > len(raw) {
			size = len(raw)
		}
		result.Write(raw[:size])
		raw = raw[size:]
	}
	return result.String()
}

// ---- ContainerLifecycleService 容器生命周期管理服务 ----

// ContainerLifecycleService 容器生命周期管理（对应 Java 的 ContainerLifecycleService）
type ContainerLifecycleService struct {
	containerDomainService *domainContainer.DomainService
	dockerService          DockerService
	logger                 *zap.Logger
}

func NewContainerLifecycleService(
	containerDomainService *domainContainer.DomainService,
	dockerService DockerService,
	logger *zap.Logger,
) *ContainerLifecycleService {
	return &ContainerLifecycleService{
		containerDomainService: containerDomainService,
		dockerService:          dockerService,
		logger:                 logger,
	}
}

// CreateAndStartContainer 创建并启动容器
func (s *ContainerLifecycleService) CreateAndStartContainer(
	userID string,
	containerType domainContainer.ContainerType,
) (*domainContainer.ContainerEntity, error) {
	// 1. 获取默认模板
	template, err := s.containerDomainService.GetDefaultTemplate(containerType)
	if err != nil {
		return nil, fmt.Errorf("获取默认模板失败: %w", err)
	}
	if template == nil {
		return nil, fmt.Errorf("未找到类型为 %s 的默认模板", containerType)
	}

	// 2. 创建容器记录
	containerName := fmt.Sprintf("%s-%s-%d", containerType, userID[:8], time.Now().Unix())
	container := &domainContainer.ContainerEntity{
		Name:   containerName,
		UserID: userID,
		Type:   containerType,
		Status: domainContainer.ContainerStatusCreating,
		Image:  template.GetFullImageName(),
	}

	created, err := s.containerDomainService.CreateContainer(container)
	if err != nil {
		return nil, fmt.Errorf("创建容器记录失败: %w", err)
	}

	// 3. 调用Docker创建容器
	dockerID, err := s.dockerService.CreateContainer(template, containerName)
	if err != nil {
		_ = s.containerDomainService.UpdateContainerStatus(created.ID, domainContainer.ContainerStatusError)
		return nil, fmt.Errorf("创建Docker容器失败: %w", err)
	}

	created.DockerContainerID = dockerID

	// 4. 启动容器
	if err := s.dockerService.StartContainer(dockerID); err != nil {
		_ = s.containerDomainService.UpdateContainerStatus(created.ID, domainContainer.ContainerStatusError)
		return nil, fmt.Errorf("启动Docker容器失败: %w", err)
	}

	// 5. 更新状态为运行中
	_ = s.containerDomainService.UpdateContainerStatus(created.ID, domainContainer.ContainerStatusRunning)
	created.Status = domainContainer.ContainerStatusRunning

	s.logger.Info("容器创建并启动成功",
		zap.String("containerId", created.ID),
		zap.String("dockerId", dockerID),
		zap.String("userId", userID),
	)

	return created, nil
}

// StopAndRemoveContainer 停止并删除容器
func (s *ContainerLifecycleService) StopAndRemoveContainer(containerID string) error {
	container, err := s.containerDomainService.GetContainerByID(containerID)
	if err != nil {
		return err
	}
	if container == nil {
		return fmt.Errorf("容器不存在")
	}

	if container.DockerContainerID != "" {
		if err := s.dockerService.StopContainer(container.DockerContainerID); err != nil {
			s.logger.Warn("停止Docker容器失败", zap.Error(err))
		}
		if err := s.dockerService.RemoveContainer(container.DockerContainerID); err != nil {
			s.logger.Warn("删除Docker容器失败", zap.Error(err))
		}
	}

	return s.containerDomainService.DeleteContainer(containerID)
}

// GetContainerLogs 获取容器日志
func (s *ContainerLifecycleService) GetContainerLogs(containerID string, lines int) (string, error) {
	container, err := s.containerDomainService.GetContainerByID(containerID)
	if err != nil {
		return "", err
	}
	if container == nil {
		return "", fmt.Errorf("容器不存在")
	}
	if container.DockerContainerID == "" {
		return "", fmt.Errorf("容器未关联Docker实例")
	}
	return s.dockerService.GetContainerLogs(container.DockerContainerID, lines)
}

// HealthCheck 容器健康检查
func (s *ContainerLifecycleService) HealthCheck(containerID string) (bool, string, error) {
	container, err := s.containerDomainService.GetContainerByID(containerID)
	if err != nil {
		return false, "获取容器信息失败", err
	}
	if container == nil {
		return false, "容器不存在", nil
	}
	if container.DockerContainerID == "" {
		return false, "容器未关联Docker实例", nil
	}

	running, err := s.dockerService.IsContainerRunning(container.DockerContainerID)
	if err != nil {
		return false, "检查容器状态失败", err
	}
	if !running {
		return false, "容器未运行", nil
	}

	return true, "容器运行正常", nil
}