package rag

import (
	"time"

	domainRag "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/rag"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// ---- DTO ----

// UserRagDTO 用户RAG DTO
type UserRagDTO struct {
	ID            string     `json:"id"`
	UserID        string     `json:"userId"`
	RagVersionID  string     `json:"ragVersionId"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Icon          string     `json:"icon"`
	Version       string     `json:"version"`
	InstalledAt   *time.Time `json:"installedAt"`
	OriginalRagID string     `json:"originalRagId"`
	InstallType   string     `json:"installType"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// RagVersionDTO RAG版本DTO
type RagVersionDTO struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Icon            string     `json:"icon"`
	Description     string     `json:"description"`
	UserID          string     `json:"userId"`
	Version         string     `json:"version"`
	ChangeLog       string     `json:"changeLog"`
	Labels          []string   `json:"labels"`
	OriginalRagID   string     `json:"originalRagId"`
	OriginalRagName string     `json:"originalRagName"`
	FileCount       *int       `json:"fileCount"`
	TotalSize       *int64     `json:"totalSize"`
	DocumentCount   *int       `json:"documentCount"`
	PublishStatus   *int       `json:"publishStatus"`
	RejectReason    string     `json:"rejectReason"`
	ReviewTime      *time.Time `json:"reviewTime"`
	PublishedAt     *time.Time `json:"publishedAt"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

// RagMarketDTO RAG市场DTO
type RagMarketDTO struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Icon          string   `json:"icon"`
	Description   string   `json:"description"`
	Version       string   `json:"version"`
	Labels        []string `json:"labels"`
	FileCount     *int     `json:"fileCount"`
	DocumentCount *int     `json:"documentCount"`
	AuthorName    string   `json:"authorName"`
	Installed     bool     `json:"installed"`
}

// DocumentUnitDTO 文档单元DTO
type DocumentUnitDTO struct {
	ID        string `json:"id"`
	FileID    string `json:"fileId"`
	Content   string `json:"content"`
	PageIndex *int   `json:"pageIndex"`
	ChunkType string `json:"chunkType"`
}

// RagSearchRequest RAG搜索请求
type RagSearchRequest struct {
	Query     string `json:"query" binding:"required"`
	TopK      int    `json:"topK"`
	MinScore  float64 `json:"minScore"`
	UserRagID string `json:"userRagId"`
}

// RagStreamChatRequest RAG流式问答请求
type RagStreamChatRequest struct {
	Query     string `json:"query" binding:"required"`
	UserRagID string `json:"userRagId"`
	SessionID string `json:"sessionId"`
}

// InstallRagRequest 安装RAG请求
type InstallRagRequest struct {
	RagVersionID string `json:"ragVersionId" binding:"required"`
	InstallType  string `json:"installType"`
}

// PublishRagRequest 发布RAG请求
type PublishRagRequest struct {
	UserRagID string `json:"userRagId" binding:"required"`
	Version   string `json:"version" binding:"required"`
	ChangeLog string `json:"changeLog"`
}

// CreateDatasetRequest 创建数据集请求
type CreateDatasetRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// RagQaDatasetDTO QA数据集DTO
type RagQaDatasetDTO struct {
	ID        string `json:"id"`
	UserRagID string `json:"userRagId"`
	Question  string `json:"question"`
	Answer    string `json:"answer"`
}

// ---- 转换函数 ----

func userRagEntityToDTO(e *domainRag.UserRagEntity) *UserRagDTO {
	if e == nil {
		return nil
	}
	return &UserRagDTO{
		ID:            e.ID,
		UserID:        e.UserID,
		RagVersionID:  e.RagVersionID,
		Name:          e.Name,
		Description:   e.Description,
		Icon:          e.Icon,
		Version:       e.Version,
		InstalledAt:   e.InstalledAt,
		OriginalRagID: e.OriginalRagID,
		InstallType:   string(e.InstallType),
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

func ragVersionEntityToDTO(e *domainRag.RagVersionEntity) *RagVersionDTO {
	if e == nil {
		return nil
	}
	return &RagVersionDTO{
		ID:              e.ID,
		Name:            e.Name,
		Icon:            e.Icon,
		Description:     e.Description,
		UserID:          e.UserID,
		Version:         e.Version,
		ChangeLog:       e.ChangeLog,
		Labels:          e.Labels,
		OriginalRagID:   e.OriginalRagID,
		OriginalRagName: e.OriginalRagName,
		FileCount:       e.FileCount,
		TotalSize:       e.TotalSize,
		DocumentCount:   e.DocumentCount,
		PublishStatus:   e.PublishStatus,
		RejectReason:    e.RejectReason,
		ReviewTime:      e.ReviewTime,
		PublishedAt:     e.PublishedAt,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

// ---- 应用服务 ----

// FileOperationAppService 文件操作应用服务
type FileOperationAppService struct {
	userRagRepo    domainRag.UserRagRepository
	fileDetailRepo domainRag.FileDetailRepository
	docUnitRepo    domainRag.DocumentUnitRepository
}

func NewFileOperationAppService(
	userRagRepo domainRag.UserRagRepository,
	fileDetailRepo domainRag.FileDetailRepository,
	docUnitRepo domainRag.DocumentUnitRepository,
) *FileOperationAppService {
	return &FileOperationAppService{
		userRagRepo:    userRagRepo,
		fileDetailRepo: fileDetailRepo,
		docUnitRepo:    docUnitRepo,
	}
}

// GetFilesByDataSetID 获取数据集下的文件列表
func (s *FileOperationAppService) GetFilesByDataSetID(dataSetID string) ([]domainRag.FileDetailEntity, error) {
	return s.fileDetailRepo.FindByDataSetID(dataSetID)
}

// GetDocumentUnitsByFileID 获取文件的文档单元列表
func (s *FileOperationAppService) GetDocumentUnitsByFileID(fileID string) ([]DocumentUnitDTO, error) {
	entities, err := s.docUnitRepo.FindByFileID(fileID)
	if err != nil {
		return nil, err
	}
	dtos := make([]DocumentUnitDTO, len(entities))
	for i, e := range entities {
		dtos[i] = DocumentUnitDTO{
			ID:        e.ID,
			FileID:    e.FileID,
			Content:   e.Content,
			PageIndex: e.PageIndex,
			ChunkType: e.ChunkType,
		}
	}
	return dtos, nil
}

// DeleteFile 删除文件
func (s *FileOperationAppService) DeleteFile(fileID, userID string) error {
	file, err := s.fileDetailRepo.FindByID(fileID)
	if err != nil {
		return err
	}
	if file == nil {
		return exception.NewBusinessException("文件不存在")
	}
	if file.UserID != userID {
		return exception.NewBusinessException("无权删除此文件")
	}
	// 删除文档单元
	if err := s.docUnitRepo.DeleteByFileID(fileID); err != nil {
		return err
	}
	// 删除文件记录
	return s.fileDetailRepo.Delete(fileID)
}

// RagMarketAppService RAG市场应用服务
type RagMarketAppService struct {
	ragVersionRepo domainRag.RagVersionRepository
	userRagRepo    domainRag.UserRagRepository
}

func NewRagMarketAppService(
	ragVersionRepo domainRag.RagVersionRepository,
	userRagRepo domainRag.UserRagRepository,
) *RagMarketAppService {
	return &RagMarketAppService{
		ragVersionRepo: ragVersionRepo,
		userRagRepo:    userRagRepo,
	}
}

// GetMarketList 获取RAG市场列表
func (s *RagMarketAppService) GetMarketList(page, pageSize int, userID string) ([]*RagMarketDTO, int64, error) {
	versions, total, err := s.ragVersionRepo.FindPublishedPaged(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	dtos := make([]*RagMarketDTO, len(versions))
	for i, v := range versions {
		installed := false
		if userID != "" {
			existing, _ := s.userRagRepo.FindByUserIDAndOriginalRagID(userID, v.OriginalRagID)
			installed = existing != nil
		}
		dtos[i] = &RagMarketDTO{
			ID:            v.ID,
			Name:          v.Name,
			Icon:          v.Icon,
			Description:   v.Description,
			Version:       v.Version,
			Labels:        v.Labels,
			FileCount:     v.FileCount,
			DocumentCount: v.DocumentCount,
			Installed:     installed,
		}
	}
	return dtos, total, nil
}

// InstallRag 安装RAG
func (s *RagMarketAppService) InstallRag(userID string, req *InstallRagRequest) (*UserRagDTO, error) {
	version, err := s.ragVersionRepo.FindByID(req.RagVersionID)
	if err != nil {
		return nil, err
	}
	if version == nil {
		return nil, exception.NewBusinessException("RAG版本不存在")
	}
	// 检查是否已安装
	existing, _ := s.userRagRepo.FindByUserIDAndOriginalRagID(userID, version.OriginalRagID)
	if existing != nil {
		return nil, exception.NewBusinessException("已安装此知识库")
	}
	installType := domainRag.InstallTypeReference
	if req.InstallType == string(domainRag.InstallTypeSnapshot) {
		installType = domainRag.InstallTypeSnapshot
	}
	now := time.Now()
	userRag := &domainRag.UserRagEntity{
		UserID:        userID,
		RagVersionID:  version.ID,
		Name:          version.Name,
		Description:   version.Description,
		Icon:          version.Icon,
		Version:       version.Version,
		InstalledAt:   &now,
		OriginalRagID: version.OriginalRagID,
		InstallType:   installType,
	}
	if err := s.userRagRepo.Create(userRag); err != nil {
		return nil, err
	}
	return userRagEntityToDTO(userRag), nil
}

// UninstallRag 卸载RAG
func (s *RagMarketAppService) UninstallRag(userRagID, userID string) error {
	userRag, err := s.userRagRepo.FindByID(userRagID)
	if err != nil {
		return err
	}
	if userRag == nil {
		return exception.NewBusinessException("知识库不存在")
	}
	if userRag.UserID != userID {
		return exception.NewBusinessException("无权卸载此知识库")
	}
	return s.userRagRepo.Delete(userRagID)
}

// GetUserRags 获取用户安装的RAG列表
func (s *RagMarketAppService) GetUserRags(userID string) ([]*UserRagDTO, error) {
	entities, err := s.userRagRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	dtos := make([]*UserRagDTO, len(entities))
	for i, e := range entities {
		dtos[i] = userRagEntityToDTO(&e)
	}
	return dtos, nil
}

// RagPublishAppService RAG发布应用服务
type RagPublishAppService struct {
	ragVersionRepo domainRag.RagVersionRepository
	userRagRepo    domainRag.UserRagRepository
}

func NewRagPublishAppService(
	ragVersionRepo domainRag.RagVersionRepository,
	userRagRepo domainRag.UserRagRepository,
) *RagPublishAppService {
	return &RagPublishAppService{
		ragVersionRepo: ragVersionRepo,
		userRagRepo:    userRagRepo,
	}
}

// PublishRag 发布RAG到市场
func (s *RagPublishAppService) PublishRag(userID string, req *PublishRagRequest) (*RagVersionDTO, error) {
	userRag, err := s.userRagRepo.FindByID(req.UserRagID)
	if err != nil {
		return nil, err
	}
	if userRag == nil {
		return nil, exception.NewBusinessException("知识库不存在")
	}
	if userRag.UserID != userID {
		return nil, exception.NewBusinessException("无权发布此知识库")
	}
	pendingStatus := int(domainRag.RagPublishStatusReviewing)
	ragVersion := &domainRag.RagVersionEntity{
		Name:            userRag.Name,
		Icon:            userRag.Icon,
		Description:     userRag.Description,
		UserID:          userID,
		Version:         req.Version,
		ChangeLog:       req.ChangeLog,
		OriginalRagID:   userRag.ID,
		OriginalRagName: userRag.Name,
		PublishStatus:   &pendingStatus,
	}
	if err := s.ragVersionRepo.Create(ragVersion); err != nil {
		return nil, err
	}
	return ragVersionEntityToDTO(ragVersion), nil
}

// GetVersionsByOriginalRagID 获取RAG的所有版本
func (s *RagPublishAppService) GetVersionsByOriginalRagID(originalRagID string) ([]*RagVersionDTO, error) {
	entities, err := s.ragVersionRepo.FindByOriginalRagID(originalRagID)
	if err != nil {
		return nil, err
	}
	dtos := make([]*RagVersionDTO, len(entities))
	for i, e := range entities {
		dtos[i] = ragVersionEntityToDTO(&e)
	}
	return dtos, nil
}

// RAGSearchAppService RAG搜索应用服务
type RAGSearchAppService struct {
	userRagRepo domainRag.UserRagRepository
	// TODO: 集成向量存储服务
}

func NewRAGSearchAppService(userRagRepo domainRag.UserRagRepository) *RAGSearchAppService {
	return &RAGSearchAppService{userRagRepo: userRagRepo}
}

// RagSearch RAG搜索
func (s *RAGSearchAppService) RagSearch(req *RagSearchRequest, userID string) ([]DocumentUnitDTO, error) {
	// TODO: 集成向量存储服务进行实际搜索
	return []DocumentUnitDTO{}, nil
}

// RagSearchByUserRag 基于已安装知识库的RAG搜索
func (s *RAGSearchAppService) RagSearchByUserRag(req *RagSearchRequest, userRagID, userID string) ([]DocumentUnitDTO, error) {
	userRag, err := s.userRagRepo.FindByID(userRagID)
	if err != nil {
		return nil, err
	}
	if userRag == nil {
		return nil, exception.NewBusinessException("知识库不存在")
	}
	if userRag.UserID != userID {
		return nil, exception.NewBusinessException("无权访问此知识库")
	}
	// TODO: 集成向量存储服务进行实际搜索
	return []DocumentUnitDTO{}, nil
}

// RagQaDatasetAppService RAG QA数据集应用服务
type RagQaDatasetAppService struct {
	qaRepo      domainRag.RagQaDatasetRepository
	userRagRepo domainRag.UserRagRepository
}

func NewRagQaDatasetAppService(
	qaRepo domainRag.RagQaDatasetRepository,
	userRagRepo domainRag.UserRagRepository,
) *RagQaDatasetAppService {
	return &RagQaDatasetAppService{qaRepo: qaRepo, userRagRepo: userRagRepo}
}

// GetQaDatasets 获取QA数据集列表
func (s *RagQaDatasetAppService) GetQaDatasets(userRagID, userID string) ([]RagQaDatasetDTO, error) {
	userRag, err := s.userRagRepo.FindByID(userRagID)
	if err != nil {
		return nil, err
	}
	if userRag == nil {
		return nil, exception.NewBusinessException("知识库不存在")
	}
	if userRag.UserID != userID {
		return nil, exception.NewBusinessException("无权访问此知识库")
	}
	entities, err := s.qaRepo.FindByUserRagID(userRagID)
	if err != nil {
		return nil, err
	}
	dtos := make([]RagQaDatasetDTO, len(entities))
	for i, e := range entities {
		dtos[i] = RagQaDatasetDTO{
			ID:        e.ID,
			UserRagID: e.UserRagID,
			Question:  e.Question,
			Answer:    e.Answer,
		}
	}
	return dtos, nil
}

// CreateQaDataset 创建QA数据集
func (s *RagQaDatasetAppService) CreateQaDataset(userRagID, question, answer, userID string) (*RagQaDatasetDTO, error) {
	userRag, err := s.userRagRepo.FindByID(userRagID)
	if err != nil {
		return nil, err
	}
	if userRag == nil {
		return nil, exception.NewBusinessException("知识库不存在")
	}
	if userRag.UserID != userID {
		return nil, exception.NewBusinessException("无权操作此知识库")
	}
	entity := &domainRag.RagQaDatasetEntity{
		UserRagID: userRagID,
		Question:  question,
		Answer:    answer,
	}
	if err := s.qaRepo.Create(entity); err != nil {
		return nil, err
	}
	return &RagQaDatasetDTO{
		ID:        entity.ID,
		UserRagID: entity.UserRagID,
		Question:  entity.Question,
		Answer:    entity.Answer,
	}, nil
}
