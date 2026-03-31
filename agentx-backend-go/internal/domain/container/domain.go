package container

import (
	"time"

	"github.com/google/uuid"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/entity"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
	"gorm.io/gorm"
)

// ---- 枚举 ----

// ContainerStatus 容器状态枚举
type ContainerStatus int

const (
	ContainerStatusCreating  ContainerStatus = 1 // 创建中
	ContainerStatusRunning   ContainerStatus = 2 // 运行中
	ContainerStatusStopped   ContainerStatus = 3 // 已停止
	ContainerStatusError     ContainerStatus = 4 // 错误状态
	ContainerStatusDeleting  ContainerStatus = 5 // 删除中
	ContainerStatusDeleted   ContainerStatus = 6 // 已删除
	ContainerStatusSuspended ContainerStatus = 7 // 已暂停
)

// ContainerType 容器类型枚举
type ContainerType string

const (
	ContainerTypeUser   ContainerType = "USER"   // 用户容器
	ContainerTypeReview ContainerType = "REVIEW" // 审核容器
)

// ---- 实体 ----

// ContainerEntity 容器实体
type ContainerEntity struct {
	ID                string          `gorm:"column:id;primaryKey" json:"id"`
	Name              string          `gorm:"column:name" json:"name"`
	UserID            string          `gorm:"column:user_id" json:"userId"`
	Type              ContainerType   `gorm:"column:type" json:"type"`
	Status            ContainerStatus `gorm:"column:status" json:"status"`
	DockerContainerID string          `gorm:"column:docker_container_id" json:"dockerContainerId"`
	Image             string          `gorm:"column:image" json:"image"`
	InternalPort      *int            `gorm:"column:internal_port" json:"internalPort"`
	ExternalPort      *int            `gorm:"column:external_port" json:"externalPort"`
	IPAddress         string          `gorm:"column:ip_address" json:"ipAddress"`
	CPUUsage          *float64        `gorm:"column:cpu_usage" json:"cpuUsage"`
	MemoryUsage       *float64        `gorm:"column:memory_usage" json:"memoryUsage"`
	VolumePath        string          `gorm:"column:volume_path" json:"volumePath"`
	EnvConfig         string          `gorm:"column:env_config" json:"envConfig"`
	ContainerConfig   string          `gorm:"column:container_config" json:"containerConfig"`
	ErrorMessage      string          `gorm:"column:error_message" json:"errorMessage"`
	LastAccessedAt    *time.Time      `gorm:"column:last_accessed_at" json:"lastAccessedAt"`

	entity.BaseEntity
}

func (ContainerEntity) TableName() string {
	return "user_containers"
}

// IsRunning 检查容器是否正在运行
func (e *ContainerEntity) IsRunning() bool {
	return e.Status == ContainerStatusRunning
}

// IsOperatable 检查容器是否可以操作
func (e *ContainerEntity) IsOperatable() bool {
	return e.Status != ContainerStatusDeleting && e.Status != ContainerStatusDeleted
}

// IsSuspended 检查容器是否已暂停
func (e *ContainerEntity) IsSuspended() bool {
	return e.Status == ContainerStatusSuspended
}

// UpdateLastAccessedAt 更新最后访问时间
func (e *ContainerEntity) UpdateLastAccessedAt() {
	now := time.Now()
	e.LastAccessedAt = &now
}

// MarkError 标记容器为错误状态
func (e *ContainerEntity) MarkError(errorMessage string) {
	e.Status = ContainerStatusError
	e.ErrorMessage = errorMessage
}

// UpdateStatus 更新容器状态
func (e *ContainerEntity) UpdateStatus(newStatus ContainerStatus) {
	e.Status = newStatus
	if newStatus == ContainerStatusRunning {
		e.ErrorMessage = ""
	}
}

// UpdateResourceUsage 更新资源使用率
func (e *ContainerEntity) UpdateResourceUsage(cpuUsage, memoryUsage *float64) {
	e.CPUUsage = cpuUsage
	e.MemoryUsage = memoryUsage
}

// ContainerTemplate 容器模板配置
type ContainerTemplate struct {
	Image           string            `json:"image"`
	InternalPort    *int              `json:"internalPort"`
	CPULimit        *float64          `json:"cpuLimit"`
	MemoryLimit     *int              `json:"memoryLimit"`
	Environment     map[string]string `json:"environment"`
	VolumeMountPath string            `json:"volumeMountPath"`
	Command         []string          `json:"command"`
	NetworkMode     string            `json:"networkMode"`
	RestartPolicy   string            `json:"restartPolicy"`
}

// GetDefaultMcpGatewayTemplate 获取默认的MCP网关容器模板
func GetDefaultMcpGatewayTemplate() *ContainerTemplate {
	port := 8080
	cpu := 1.0
	mem := 512
	return &ContainerTemplate{
		Image:           "ghcr.io/lucky-aeon/mcp-gateway:latest",
		InternalPort:    &port,
		CPULimit:        &cpu,
		MemoryLimit:     &mem,
		VolumeMountPath: "/app/data",
		NetworkMode:     "bridge",
		RestartPolicy:   "unless-stopped",
	}
}

// ContainerTemplateEntity 容器模板实体
type ContainerTemplateEntity struct {
	ID              string                 `gorm:"column:id;primaryKey" json:"id"`
	Name            string                 `gorm:"column:name" json:"name"`
	Description     string                 `gorm:"column:description" json:"description"`
	Type            ContainerType          `gorm:"column:type" json:"type"`
	Image           string                 `gorm:"column:image" json:"image"`
	ImageTag        string                 `gorm:"column:image_tag" json:"imageTag"`
	InternalPort    *int                   `gorm:"column:internal_port" json:"internalPort"`
	CPULimit        *float64               `gorm:"column:cpu_limit" json:"cpuLimit"`
	MemoryLimit     *int                   `gorm:"column:memory_limit" json:"memoryLimit"`
	Environment     map[string]string      `gorm:"column:environment;serializer:json" json:"environment"`
	VolumeMountPath string                 `gorm:"column:volume_mount_path" json:"volumeMountPath"`
	Command         []string               `gorm:"column:command;serializer:json" json:"command"`
	NetworkMode     string                 `gorm:"column:network_mode" json:"networkMode"`
	RestartPolicy   string                 `gorm:"column:restart_policy" json:"restartPolicy"`
	HealthCheck     map[string]interface{} `gorm:"column:health_check;serializer:json" json:"healthCheck"`
	ResourceConfig  map[string]interface{} `gorm:"column:resource_config;serializer:json" json:"resourceConfig"`
	Enabled         *bool                  `gorm:"column:enabled" json:"enabled"`
	IsDefault       *bool                  `gorm:"column:is_default" json:"isDefault"`
	CreatedBy       string                 `gorm:"column:created_by" json:"createdBy"`
	SortOrder       *int                   `gorm:"column:sort_order" json:"sortOrder"`

	entity.BaseEntity
}

func (ContainerTemplateEntity) TableName() string {
	return "container_templates"
}

// GetFullImageName 获取完整的镜像名称
func (e *ContainerTemplateEntity) GetFullImageName() string {
	if e.ImageTag == "" {
		return e.Image
	}
	return e.Image + ":" + e.ImageTag
}

// ToContainerTemplate 转换为容器模板配置
func (e *ContainerTemplateEntity) ToContainerTemplate() *ContainerTemplate {
	return &ContainerTemplate{
		Image:           e.GetFullImageName(),
		InternalPort:    e.InternalPort,
		CPULimit:        e.CPULimit,
		MemoryLimit:     e.MemoryLimit,
		Environment:     e.Environment,
		VolumeMountPath: e.VolumeMountPath,
		Command:         e.Command,
		NetworkMode:     e.NetworkMode,
		RestartPolicy:   e.RestartPolicy,
	}
}

// IsValid 验证模板配置是否有效
func (e *ContainerTemplateEntity) IsValid() bool {
	return e.Name != "" && e.Image != "" &&
		e.InternalPort != nil && *e.InternalPort > 0 && *e.InternalPort <= 65535 &&
		e.CPULimit != nil && *e.CPULimit > 0 &&
		e.MemoryLimit != nil && *e.MemoryLimit > 0
}

// ---- 仓储接口 ----

// ContainerRepository 容器仓储接口
type ContainerRepository interface {
	FindByID(id string) (*ContainerEntity, error)
	FindByUserIDAndType(userID string, containerType ContainerType) (*ContainerEntity, error)
	FindByDockerContainerID(dockerContainerID string) (*ContainerEntity, error)
	FindByExternalPort(externalPort int) (*ContainerEntity, error)
	FindByUserID(userID string) ([]ContainerEntity, error)
	FindByStatus(status ContainerStatus) ([]ContainerEntity, error)
	FindPaged(page, pageSize int, keyword string, status *ContainerStatus, containerType *ContainerType) ([]ContainerEntity, int64, error)
	IsPortOccupied(port int) (bool, error)
	CountRunningContainers() (int64, error)
	CountByUserID(userID string) (int64, error)
	Create(entity *ContainerEntity) error
	Update(entity *ContainerEntity) error
	Delete(id string) error
}

// ContainerTemplateRepository 容器模板仓储接口
type ContainerTemplateRepository interface {
	FindByID(id string) (*ContainerTemplateEntity, error)
	FindByType(containerType string) ([]ContainerTemplateEntity, error)
	FindDefaultByType(containerType ContainerType) (*ContainerTemplateEntity, error)
	FindEnabledTemplates() ([]ContainerTemplateEntity, error)
	FindByName(name string) (*ContainerTemplateEntity, error)
	ExistsByName(name, excludeID string) (bool, error)
	FindPaged(page, pageSize int, keyword, containerType string, enabled *bool) ([]ContainerTemplateEntity, int64, error)
	FindByCreatedBy(createdBy string) ([]ContainerTemplateEntity, error)
	CountTemplates() (int64, error)
	CountEnabledTemplates() (int64, error)
	Create(entity *ContainerTemplateEntity) error
	Update(entity *ContainerTemplateEntity) error
	Delete(id string) error
}

// ---- GORM 实现 ----

type ContainerRepositoryImpl struct {
	DB *gorm.DB
}

func NewContainerRepository(db *gorm.DB) ContainerRepository {
	return &ContainerRepositoryImpl{DB: db}
}

func (r *ContainerRepositoryImpl) FindByID(id string) (*ContainerEntity, error) {
	var e ContainerEntity
	result := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *ContainerRepositoryImpl) FindByUserIDAndType(userID string, containerType ContainerType) (*ContainerEntity, error) {
	var e ContainerEntity
	result := r.DB.Where("user_id = ? AND type = ? AND deleted_at IS NULL", userID, containerType).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *ContainerRepositoryImpl) FindByDockerContainerID(dockerContainerID string) (*ContainerEntity, error) {
	var e ContainerEntity
	result := r.DB.Where("docker_container_id = ? AND deleted_at IS NULL", dockerContainerID).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *ContainerRepositoryImpl) FindByExternalPort(externalPort int) (*ContainerEntity, error) {
	var e ContainerEntity
	result := r.DB.Where("external_port = ? AND deleted_at IS NULL", externalPort).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *ContainerRepositoryImpl) FindByUserID(userID string) ([]ContainerEntity, error) {
	var entities []ContainerEntity
	result := r.DB.Where("user_id = ? AND deleted_at IS NULL", userID).Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *ContainerRepositoryImpl) FindByStatus(status ContainerStatus) ([]ContainerEntity, error) {
	var entities []ContainerEntity
	result := r.DB.Where("status = ? AND deleted_at IS NULL", status).Find(&entities)
	return entities, result.Error
}

func (r *ContainerRepositoryImpl) FindPaged(page, pageSize int, keyword string, status *ContainerStatus, containerType *ContainerType) ([]ContainerEntity, int64, error) {
	var entities []ContainerEntity
	var total int64
	query := r.DB.Model(&ContainerEntity{}).Where("deleted_at IS NULL")
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if containerType != nil {
		query = query.Where("type = ?", *containerType)
	}
	query.Count(&total)
	result := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities)
	return entities, total, result.Error
}

func (r *ContainerRepositoryImpl) IsPortOccupied(port int) (bool, error) {
	var count int64
	result := r.DB.Model(&ContainerEntity{}).Where("external_port = ? AND deleted_at IS NULL", port).Count(&count)
	return count > 0, result.Error
}

func (r *ContainerRepositoryImpl) CountRunningContainers() (int64, error) {
	var count int64
	result := r.DB.Model(&ContainerEntity{}).Where("status = ? AND deleted_at IS NULL", ContainerStatusRunning).Count(&count)
	return count, result.Error
}

func (r *ContainerRepositoryImpl) CountByUserID(userID string) (int64, error) {
	var count int64
	result := r.DB.Model(&ContainerEntity{}).Where("user_id = ? AND deleted_at IS NULL", userID).Count(&count)
	return count, result.Error
}

func (r *ContainerRepositoryImpl) Create(e *ContainerEntity) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return r.DB.Create(e).Error
}

func (r *ContainerRepositoryImpl) Update(e *ContainerEntity) error {
	return r.DB.Save(e).Error
}

func (r *ContainerRepositoryImpl) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&ContainerEntity{}).Error
}

// ---- ContainerTemplate GORM 实现 ----

type ContainerTemplateRepositoryImpl struct {
	DB *gorm.DB
}

func NewContainerTemplateRepository(db *gorm.DB) ContainerTemplateRepository {
	return &ContainerTemplateRepositoryImpl{DB: db}
}

func (r *ContainerTemplateRepositoryImpl) FindByID(id string) (*ContainerTemplateEntity, error) {
	var e ContainerTemplateEntity
	result := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *ContainerTemplateRepositoryImpl) FindByType(containerType string) ([]ContainerTemplateEntity, error) {
	var entities []ContainerTemplateEntity
	result := r.DB.Where("type = ? AND deleted_at IS NULL", containerType).Order("sort_order ASC").Find(&entities)
	return entities, result.Error
}

func (r *ContainerTemplateRepositoryImpl) FindDefaultByType(containerType ContainerType) (*ContainerTemplateEntity, error) {
	var e ContainerTemplateEntity
	result := r.DB.Where("type = ? AND is_default = true AND deleted_at IS NULL", containerType).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *ContainerTemplateRepositoryImpl) FindEnabledTemplates() ([]ContainerTemplateEntity, error) {
	var entities []ContainerTemplateEntity
	result := r.DB.Where("enabled = true AND deleted_at IS NULL").Order("sort_order ASC").Find(&entities)
	return entities, result.Error
}

func (r *ContainerTemplateRepositoryImpl) FindByName(name string) (*ContainerTemplateEntity, error) {
	var e ContainerTemplateEntity
	result := r.DB.Where("name = ? AND deleted_at IS NULL", name).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *ContainerTemplateRepositoryImpl) ExistsByName(name, excludeID string) (bool, error) {
	var count int64
	query := r.DB.Model(&ContainerTemplateEntity{}).Where("name = ? AND deleted_at IS NULL", name)
	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}
	result := query.Count(&count)
	return count > 0, result.Error
}

func (r *ContainerTemplateRepositoryImpl) FindPaged(page, pageSize int, keyword, containerType string, enabled *bool) ([]ContainerTemplateEntity, int64, error) {
	var entities []ContainerTemplateEntity
	var total int64
	query := r.DB.Model(&ContainerTemplateEntity{}).Where("deleted_at IS NULL")
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if containerType != "" {
		query = query.Where("type = ?", containerType)
	}
	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
	}
	query.Count(&total)
	result := query.Order("sort_order ASC, created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities)
	return entities, total, result.Error
}

func (r *ContainerTemplateRepositoryImpl) FindByCreatedBy(createdBy string) ([]ContainerTemplateEntity, error) {
	var entities []ContainerTemplateEntity
	result := r.DB.Where("created_by = ? AND deleted_at IS NULL", createdBy).Order("sort_order ASC").Find(&entities)
	return entities, result.Error
}

func (r *ContainerTemplateRepositoryImpl) CountTemplates() (int64, error) {
	var count int64
	result := r.DB.Model(&ContainerTemplateEntity{}).Where("deleted_at IS NULL").Count(&count)
	return count, result.Error
}

func (r *ContainerTemplateRepositoryImpl) CountEnabledTemplates() (int64, error) {
	var count int64
	result := r.DB.Model(&ContainerTemplateEntity{}).Where("enabled = true AND deleted_at IS NULL").Count(&count)
	return count, result.Error
}

func (r *ContainerTemplateRepositoryImpl) Create(e *ContainerTemplateEntity) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return r.DB.Create(e).Error
}

func (r *ContainerTemplateRepositoryImpl) Update(e *ContainerTemplateEntity) error {
	return r.DB.Save(e).Error
}

func (r *ContainerTemplateRepositoryImpl) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&ContainerTemplateEntity{}).Error
}

// ---- 领域服务 ----

// DomainService 容器领域服务
type DomainService struct {
	containerRepo ContainerRepository
	templateRepo  ContainerTemplateRepository
}

func NewDomainService(containerRepo ContainerRepository, templateRepo ContainerTemplateRepository) *DomainService {
	return &DomainService{containerRepo: containerRepo, templateRepo: templateRepo}
}

// ContainerRepo 获取容器仓储（供应用层使用）
func (s *DomainService) ContainerRepo() ContainerRepository {
	return s.containerRepo
}

// TemplateRepo 获取模板仓储（供应用层使用）
func (s *DomainService) TemplateRepo() ContainerTemplateRepository {
	return s.templateRepo
}

// GetContainerByID 根据ID获取容器
func (s *DomainService) GetContainerByID(id string) (*ContainerEntity, error) {
	return s.containerRepo.FindByID(id)
}

// GetUserContainer 获取用户指定类型的容器
func (s *DomainService) GetUserContainer(userID string, containerType ContainerType) (*ContainerEntity, error) {
	return s.containerRepo.FindByUserIDAndType(userID, containerType)
}

// GetUserContainers 获取用户所有容器
func (s *DomainService) GetUserContainers(userID string) ([]ContainerEntity, error) {
	return s.containerRepo.FindByUserID(userID)
}

// CreateContainer 创建容器
func (s *DomainService) CreateContainer(container *ContainerEntity) (*ContainerEntity, error) {
	if container.UserID == "" {
		return nil, exception.NewBusinessException("用户ID不能为空")
	}
	if container.Name == "" {
		return nil, exception.NewBusinessException("容器名称不能为空")
	}
	container.Status = ContainerStatusCreating
	if err := s.containerRepo.Create(container); err != nil {
		return nil, err
	}
	return container, nil
}

// UpdateContainerStatus 更新容器状态
func (s *DomainService) UpdateContainerStatus(id string, status ContainerStatus) error {
	container, err := s.containerRepo.FindByID(id)
	if err != nil {
		return err
	}
	if container == nil {
		return exception.NewBusinessException("容器不存在")
	}
	container.UpdateStatus(status)
	return s.containerRepo.Update(container)
}

// DeleteContainer 删除容器
func (s *DomainService) DeleteContainer(id string) error {
	return s.containerRepo.Delete(id)
}

// IsPortOccupied 检查端口是否被占用
func (s *DomainService) IsPortOccupied(port int) (bool, error) {
	return s.containerRepo.IsPortOccupied(port)
}

// GetTemplate 获取容器模板
func (s *DomainService) GetTemplate(id string) (*ContainerTemplateEntity, error) {
	return s.templateRepo.FindByID(id)
}

// GetDefaultTemplate 获取默认模板
func (s *DomainService) GetDefaultTemplate(containerType ContainerType) (*ContainerTemplateEntity, error) {
	return s.templateRepo.FindDefaultByType(containerType)
}

// GetEnabledTemplates 获取所有启用的模板
func (s *DomainService) GetEnabledTemplates() ([]ContainerTemplateEntity, error) {
	return s.templateRepo.FindEnabledTemplates()
}

// CreateTemplate 创建容器模板
func (s *DomainService) CreateTemplate(template *ContainerTemplateEntity) (*ContainerTemplateEntity, error) {
	if !template.IsValid() {
		return nil, exception.NewBusinessException("模板配置无效")
	}
	exists, err := s.templateRepo.ExistsByName(template.Name, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, exception.NewBusinessException("模板名称已存在")
	}
	if err := s.templateRepo.Create(template); err != nil {
		return nil, err
	}
	return template, nil
}

// UpdateTemplate 更新容器模板
func (s *DomainService) UpdateTemplate(template *ContainerTemplateEntity) (*ContainerTemplateEntity, error) {
	if template.ID == "" {
		return nil, exception.NewBusinessException("模板ID不能为空")
	}
	if err := s.templateRepo.Update(template); err != nil {
		return nil, err
	}
	return template, nil
}

// DeleteTemplate 删除容器模板
func (s *DomainService) DeleteTemplate(id string) error {
	return s.templateRepo.Delete(id)
}
