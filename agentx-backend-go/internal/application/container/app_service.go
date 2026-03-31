package container

import (
	"time"

	domainContainer "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/container"
	infraDocker "github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/docker"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// ContainerDTO 容器DTO
type ContainerDTO struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	UserID            string  `json:"userId"`
	UserNickname      string  `json:"userNickname,omitempty"`
	Type              string  `json:"type"`
	Status            int     `json:"status"`
	DockerContainerID string  `json:"dockerContainerId"`
	Image             string  `json:"image"`
	InternalPort      *int    `json:"internalPort"`
	ExternalPort      *int    `json:"externalPort"`
	IPAddress         string  `json:"ipAddress"`
	CPUUsage          *float64 `json:"cpuUsage"`
	MemoryUsage       *float64 `json:"memoryUsage"`
	VolumePath        string  `json:"volumePath"`
	ErrorMessage      string  `json:"errorMessage"`
	LastAccessedAt    *time.Time `json:"lastAccessedAt"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

// ContainerTemplateDTO 容器模板DTO
type ContainerTemplateDTO struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Type            string                 `json:"type"`
	Image           string                 `json:"image"`
	ImageTag        string                 `json:"imageTag"`
	FullImageName   string                 `json:"fullImageName"`
	InternalPort    *int                   `json:"internalPort"`
	CPULimit        *float64               `json:"cpuLimit"`
	MemoryLimit     *int                   `json:"memoryLimit"`
	Environment     map[string]string      `json:"environment"`
	VolumeMountPath string                 `json:"volumeMountPath"`
	Command         []string               `json:"command"`
	NetworkMode     string                 `json:"networkMode"`
	RestartPolicy   string                 `json:"restartPolicy"`
	HealthCheck     map[string]interface{} `json:"healthCheck"`
	ResourceConfig  map[string]interface{} `json:"resourceConfig"`
	Enabled         *bool                  `json:"enabled"`
	IsDefault       *bool                  `json:"isDefault"`
	CreatedBy       string                 `json:"createdBy"`
	SortOrder       *int                   `json:"sortOrder"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`
}

// ContainerHealthStatus 容器健康状态
type ContainerHealthStatus struct {
	Healthy   bool          `json:"healthy"`
	Message   string        `json:"message"`
	Container *ContainerDTO `json:"container,omitempty"`
}

// ContainerStatistics 容器统计信息
type ContainerStatistics struct {
	TotalContainers   int64 `json:"totalContainers"`
	RunningContainers int64 `json:"runningContainers"`
}

func entityToDTO(e *domainContainer.ContainerEntity) *ContainerDTO {
	if e == nil {
		return nil
	}
	return &ContainerDTO{
		ID:                e.ID,
		Name:              e.Name,
		UserID:            e.UserID,
		Type:              string(e.Type),
		Status:            int(e.Status),
		DockerContainerID: e.DockerContainerID,
		Image:             e.Image,
		InternalPort:      e.InternalPort,
		ExternalPort:      e.ExternalPort,
		IPAddress:         e.IPAddress,
		CPUUsage:          e.CPUUsage,
		MemoryUsage:       e.MemoryUsage,
		VolumePath:        e.VolumePath,
		ErrorMessage:      e.ErrorMessage,
		LastAccessedAt:    e.LastAccessedAt,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}

func templateEntityToDTO(e *domainContainer.ContainerTemplateEntity) *ContainerTemplateDTO {
	if e == nil {
		return nil
	}
	return &ContainerTemplateDTO{
		ID:              e.ID,
		Name:            e.Name,
		Description:     e.Description,
		Type:            string(e.Type),
		Image:           e.Image,
		ImageTag:        e.ImageTag,
		FullImageName:   e.GetFullImageName(),
		InternalPort:    e.InternalPort,
		CPULimit:        e.CPULimit,
		MemoryLimit:     e.MemoryLimit,
		Environment:     e.Environment,
		VolumeMountPath: e.VolumeMountPath,
		Command:         e.Command,
		NetworkMode:     e.NetworkMode,
		RestartPolicy:   e.RestartPolicy,
		HealthCheck:     e.HealthCheck,
		ResourceConfig:  e.ResourceConfig,
		Enabled:         e.Enabled,
		IsDefault:       e.IsDefault,
		CreatedBy:       e.CreatedBy,
		SortOrder:       e.SortOrder,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

// AppService 容器应用服务
type AppService struct {
	domainService    *domainContainer.DomainService
	lifecycleService *infraDocker.ContainerLifecycleService
}

func NewAppService(domainService *domainContainer.DomainService, lifecycleService *infraDocker.ContainerLifecycleService) *AppService {
	return &AppService{
		domainService:    domainService,
		lifecycleService: lifecycleService,
	}
}

// GetUserContainer 获取用户容器
func (s *AppService) GetUserContainer(userID string) (*ContainerDTO, error) {
	container, err := s.domainService.GetUserContainer(userID, domainContainer.ContainerTypeUser)
	if err != nil {
		return nil, err
	}
	if container == nil {
		return nil, exception.NewBusinessException("用户容器不存在")
	}
	return entityToDTO(container), nil
}

// GetContainerByID 根据ID获取容器
func (s *AppService) GetContainerByID(containerID string) (*ContainerDTO, error) {
	container, err := s.domainService.GetContainerByID(containerID)
	if err != nil {
		return nil, err
	}
	if container == nil {
		return nil, exception.NewBusinessException("容器不存在")
	}
	return entityToDTO(container), nil
}

// GetContainersPage 分页查询容器
func (s *AppService) GetContainersPage(page, pageSize int, keyword string, status *domainContainer.ContainerStatus, containerType *domainContainer.ContainerType) ([]*ContainerDTO, int64, error) {
	entities, total, err := s.domainService.ContainerRepo().FindPaged(page, pageSize, keyword, status, containerType)
	if err != nil {
		return nil, 0, err
	}
	dtos := make([]*ContainerDTO, len(entities))
	for i, e := range entities {
		dtos[i] = entityToDTO(&e)
	}
	return dtos, total, nil
}

// StartContainer 启动容器
func (s *AppService) StartContainer(containerID string) error {
	return s.domainService.UpdateContainerStatus(containerID, domainContainer.ContainerStatusRunning)
}

// StopContainer 停止容器
func (s *AppService) StopContainer(containerID string) error {
	return s.domainService.UpdateContainerStatus(containerID, domainContainer.ContainerStatusStopped)
}

// DeleteContainer 删除容器
func (s *AppService) DeleteContainer(containerID string) error {
	return s.domainService.DeleteContainer(containerID)
}

// GetContainerLogs 获取容器日志
func (s *AppService) GetContainerLogs(containerID string, lines int) (string, error) {
	if lines <= 0 {
		lines = 100
	}
	return s.lifecycleService.GetContainerLogs(containerID, lines)
}

// ---- 模板应用服务 ----

// TemplateAppService 容器模板应用服务
type TemplateAppService struct {
	domainService *domainContainer.DomainService
}

func NewTemplateAppService(domainService *domainContainer.DomainService) *TemplateAppService {
	return &TemplateAppService{domainService: domainService}
}

// GetTemplate 根据ID获取模板
func (s *TemplateAppService) GetTemplate(templateID string) (*ContainerTemplateDTO, error) {
	template, err := s.domainService.GetTemplate(templateID)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, exception.NewBusinessException("模板不存在")
	}
	return templateEntityToDTO(template), nil
}

// GetDefaultTemplate 获取默认模板
func (s *TemplateAppService) GetDefaultTemplate(containerType domainContainer.ContainerType) (*ContainerTemplateDTO, error) {
	template, err := s.domainService.GetDefaultTemplate(containerType)
	if err != nil {
		return nil, err
	}
	return templateEntityToDTO(template), nil
}

// GetEnabledTemplates 获取所有启用的模板
func (s *TemplateAppService) GetEnabledTemplates() ([]*ContainerTemplateDTO, error) {
	templates, err := s.domainService.GetEnabledTemplates()
	if err != nil {
		return nil, err
	}
	dtos := make([]*ContainerTemplateDTO, len(templates))
	for i, t := range templates {
		dtos[i] = templateEntityToDTO(&t)
	}
	return dtos, nil
}

// CreateTemplate 创建模板
func (s *TemplateAppService) CreateTemplate(template *domainContainer.ContainerTemplateEntity) (*ContainerTemplateDTO, error) {
	created, err := s.domainService.CreateTemplate(template)
	if err != nil {
		return nil, err
	}
	return templateEntityToDTO(created), nil
}

// UpdateTemplate 更新模板
func (s *TemplateAppService) UpdateTemplate(template *domainContainer.ContainerTemplateEntity) (*ContainerTemplateDTO, error) {
	updated, err := s.domainService.UpdateTemplate(template)
	if err != nil {
		return nil, err
	}
	return templateEntityToDTO(updated), nil
}

// DeleteTemplate 删除模板
func (s *TemplateAppService) DeleteTemplate(templateID string) error {
	return s.domainService.DeleteTemplate(templateID)
}
