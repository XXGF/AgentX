package rag

import (
	"time"

	"github.com/google/uuid"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/entity"
	"gorm.io/gorm"
)

// ---- 核心实体 ----

// UserRagEntity 用户安装的RAG实体
type UserRagEntity struct {
	ID            string      `gorm:"column:id;primaryKey" json:"id"`
	UserID        string      `gorm:"column:user_id" json:"userId"`
	RagVersionID  string      `gorm:"column:rag_version_id" json:"ragVersionId"`
	Name          string      `gorm:"column:name" json:"name"`
	Description   string      `gorm:"column:description" json:"description"`
	Icon          string      `gorm:"column:icon" json:"icon"`
	Version       string      `gorm:"column:version" json:"version"`
	InstalledAt   *time.Time  `gorm:"column:installed_at" json:"installedAt"`
	OriginalRagID string      `gorm:"column:original_rag_id" json:"originalRagId"`
	InstallType   InstallType `gorm:"column:install_type" json:"installType"`

	entity.BaseEntity
}

func (UserRagEntity) TableName() string {
	return "user_rags"
}

// IsReferenceType 是否为引用类型安装
func (e *UserRagEntity) IsReferenceType() bool {
	return e.InstallType.IsReference()
}

// IsSnapshotType 是否为快照类型安装
func (e *UserRagEntity) IsSnapshotType() bool {
	return e.InstallType.IsSnapshot()
}

// RagVersionEntity RAG版本实体（完整快照）
type RagVersionEntity struct {
	ID              string     `gorm:"column:id;primaryKey" json:"id"`
	Name            string     `gorm:"column:name" json:"name"`
	Icon            string     `gorm:"column:icon" json:"icon"`
	Description     string     `gorm:"column:description" json:"description"`
	UserID          string     `gorm:"column:user_id" json:"userId"`
	Version         string     `gorm:"column:version" json:"version"`
	ChangeLog       string     `gorm:"column:change_log" json:"changeLog"`
	Labels          []string   `gorm:"column:labels;serializer:json" json:"labels"`
	OriginalRagID   string     `gorm:"column:original_rag_id" json:"originalRagId"`
	OriginalRagName string     `gorm:"column:original_rag_name" json:"originalRagName"`
	FileCount       *int       `gorm:"column:file_count" json:"fileCount"`
	TotalSize       *int64     `gorm:"column:total_size" json:"totalSize"`
	DocumentCount   *int       `gorm:"column:document_count" json:"documentCount"`
	PublishStatus   *int       `gorm:"column:publish_status" json:"publishStatus"`
	RejectReason    string     `gorm:"column:reject_reason" json:"rejectReason"`
	ReviewTime      *time.Time `gorm:"column:review_time" json:"reviewTime"`
	PublishedAt     *time.Time `gorm:"column:published_at" json:"publishedAt"`

	entity.BaseEntity
}

func (RagVersionEntity) TableName() string {
	return "rag_versions"
}

// FileDetailEntity 文件详情实体
type FileDetailEntity struct {
	ID                        string   `gorm:"column:id;primaryKey" json:"id"`
	URL                       string   `gorm:"column:url" json:"url"`
	Size                      *int64   `gorm:"column:size" json:"size"`
	Filename                  string   `gorm:"column:filename" json:"filename"`
	OriginalFilename          string   `gorm:"column:original_filename" json:"originalFilename"`
	BasePath                  string   `gorm:"column:base_path" json:"basePath"`
	Path                      string   `gorm:"column:path" json:"path"`
	Ext                       string   `gorm:"column:ext" json:"ext"`
	ContentType               string   `gorm:"column:content_type" json:"contentType"`
	Platform                  string   `gorm:"column:platform" json:"platform"`
	ObjectID                  string   `gorm:"column:object_id" json:"objectId"`
	ObjectType                string   `gorm:"column:object_type" json:"objectType"`
	Metadata                  string   `gorm:"column:metadata" json:"metadata"`
	UserMetadata              string   `gorm:"column:user_metadata" json:"userMetadata"`
	HashInfo                  string   `gorm:"column:hash_info" json:"hashInfo"`
	UploadID                  string   `gorm:"column:upload_id" json:"uploadId"`
	UploadStatus              *int     `gorm:"column:upload_status" json:"uploadStatus"`
	UserID                    string   `gorm:"column:user_id" json:"userId"`
	DataSetID                 string   `gorm:"column:data_set_id" json:"dataSetId"`
	FilePageSize              *int     `gorm:"column:file_page_size" json:"filePageSize"`
	ProcessingStatus          *int     `gorm:"column:processing_status" json:"processingStatus"`
	CurrentOCRPageNumber      *int     `gorm:"column:current_ocr_page_number" json:"currentOcrPageNumber"`
	CurrentEmbeddingPageNumber *int    `gorm:"column:current_embedding_page_number" json:"currentEmbeddingPageNumber"`
	OCRProcessProgress        *float64 `gorm:"column:ocr_process_progress" json:"ocrProcessProgress"`
	EmbeddingProcessProgress  *float64 `gorm:"column:embedding_process_progress" json:"embeddingProcessProgress"`

	entity.BaseEntity
}

func (FileDetailEntity) TableName() string {
	return "file_detail"
}

// GetProcessingStatusEnum 获取处理状态枚举
func (e *FileDetailEntity) GetProcessingStatusEnum() FileProcessingStatus {
	if e.ProcessingStatus == nil {
		return FileProcessingStatusUploaded
	}
	return FileProcessingStatus(*e.ProcessingStatus)
}

// DocumentUnitEntity 文档单元实体
type DocumentUnitEntity struct {
	ID        string `gorm:"column:id;primaryKey" json:"id"`
	FileID    string `gorm:"column:file_id" json:"fileId"`
	Content   string `gorm:"column:content" json:"content"`
	PageIndex *int   `gorm:"column:page_index" json:"pageIndex"`
	ChunkType string `gorm:"column:chunk_type" json:"chunkType"`

	entity.BaseEntity
}

func (DocumentUnitEntity) TableName() string {
	return "document_units"
}

// UserRagFileEntity 用户RAG文件关联实体
type UserRagFileEntity struct {
	ID       string `gorm:"column:id;primaryKey" json:"id"`
	UserRagID string `gorm:"column:user_rag_id" json:"userRagId"`
	FileID   string `gorm:"column:file_id" json:"fileId"`

	entity.BaseEntity
}

func (UserRagFileEntity) TableName() string {
	return "user_rag_files"
}

// UserRagDocumentEntity 用户RAG文档关联实体
type UserRagDocumentEntity struct {
	ID       string `gorm:"column:id;primaryKey" json:"id"`
	UserRagID string `gorm:"column:user_rag_id" json:"userRagId"`
	DocumentID string `gorm:"column:document_id" json:"documentId"`

	entity.BaseEntity
}

func (UserRagDocumentEntity) TableName() string {
	return "user_rag_documents"
}

// RagVersionFileEntity RAG版本文件关联实体
type RagVersionFileEntity struct {
	ID           string `gorm:"column:id;primaryKey" json:"id"`
	RagVersionID string `gorm:"column:rag_version_id" json:"ragVersionId"`
	FileID       string `gorm:"column:file_id" json:"fileId"`

	entity.BaseEntity
}

func (RagVersionFileEntity) TableName() string {
	return "rag_version_files"
}

// RagVersionDocumentEntity RAG版本文档关联实体
type RagVersionDocumentEntity struct {
	ID           string `gorm:"column:id;primaryKey" json:"id"`
	RagVersionID string `gorm:"column:rag_version_id" json:"ragVersionId"`
	DocumentID   string `gorm:"column:document_id" json:"documentId"`

	entity.BaseEntity
}

func (RagVersionDocumentEntity) TableName() string {
	return "rag_version_documents"
}

// RagQaDatasetEntity RAG QA数据集实体
type RagQaDatasetEntity struct {
	ID       string `gorm:"column:id;primaryKey" json:"id"`
	UserRagID string `gorm:"column:user_rag_id" json:"userRagId"`
	Question string `gorm:"column:question" json:"question"`
	Answer   string `gorm:"column:answer" json:"answer"`

	entity.BaseEntity
}

func (RagQaDatasetEntity) TableName() string {
	return "rag_qa_datasets"
}

// VectorStoreResult 向量存储搜索结果
type VectorStoreResult struct {
	DocumentID string  `json:"documentId"`
	Content    string  `json:"content"`
	Score      float64 `json:"score"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// ProcessedSegment 处理后的文档分段
type ProcessedSegment struct {
	Content     string      `json:"content"`
	SegmentType SegmentType `json:"segmentType"`
	PageIndex   int         `json:"pageIndex"`
	Order       int         `json:"order"`
}

// ModelConfig RAG模型配置
type ModelConfig struct {
	EmbeddingModelID string `json:"embeddingModelId"`
	OCRModelID       string `json:"ocrModelId"`
	RerankModelID    string `json:"rerankModelId"`
}

// HybridSearchConfig 混合搜索配置
type HybridSearchConfig struct {
	VectorWeight  float64 `json:"vectorWeight"`
	KeywordWeight float64 `json:"keywordWeight"`
	TopK          int     `json:"topK"`
	MinScore      float64 `json:"minScore"`
	UseRerank     bool    `json:"useRerank"`
	RerankTopK    int     `json:"rerankTopK"`
}

// ---- 仓储接口 ----

type UserRagRepository interface {
	FindByID(id string) (*UserRagEntity, error)
	FindByUserID(userID string) ([]UserRagEntity, error)
	FindByUserIDAndOriginalRagID(userID, originalRagID string) (*UserRagEntity, error)
	Create(entity *UserRagEntity) error
	Update(entity *UserRagEntity) error
	Delete(id string) error
}

type RagVersionRepository interface {
	FindByID(id string) (*RagVersionEntity, error)
	FindByOriginalRagID(originalRagID string) ([]RagVersionEntity, error)
	FindPublished() ([]RagVersionEntity, error)
	FindPublishedPaged(page, pageSize int) ([]RagVersionEntity, int64, error)
	Create(entity *RagVersionEntity) error
	Update(entity *RagVersionEntity) error
	Delete(id string) error
}

type FileDetailRepository interface {
	FindByID(id string) (*FileDetailEntity, error)
	FindByDataSetID(dataSetID string) ([]FileDetailEntity, error)
	FindByUserID(userID string) ([]FileDetailEntity, error)
	Create(entity *FileDetailEntity) error
	Update(entity *FileDetailEntity) error
	Delete(id string) error
}

type DocumentUnitRepository interface {
	FindByFileID(fileID string) ([]DocumentUnitEntity, error)
	Create(entity *DocumentUnitEntity) error
	BatchCreate(entities []DocumentUnitEntity) error
	DeleteByFileID(fileID string) error
}

type UserRagFileRepository interface {
	FindByUserRagID(userRagID string) ([]UserRagFileEntity, error)
	Create(entity *UserRagFileEntity) error
	DeleteByUserRagID(userRagID string) error
}

type UserRagDocumentRepository interface {
	FindByUserRagID(userRagID string) ([]UserRagDocumentEntity, error)
	Create(entity *UserRagDocumentEntity) error
	DeleteByUserRagID(userRagID string) error
}

type RagVersionFileRepository interface {
	FindByRagVersionID(ragVersionID string) ([]RagVersionFileEntity, error)
	Create(entity *RagVersionFileEntity) error
	DeleteByRagVersionID(ragVersionID string) error
}

type RagVersionDocumentRepository interface {
	FindByRagVersionID(ragVersionID string) ([]RagVersionDocumentEntity, error)
	Create(entity *RagVersionDocumentEntity) error
	DeleteByRagVersionID(ragVersionID string) error
}

type RagQaDatasetRepository interface {
	FindByUserRagID(userRagID string) ([]RagQaDatasetEntity, error)
	Create(entity *RagQaDatasetEntity) error
	DeleteByUserRagID(userRagID string) error
}

type VectorStoreRepository interface {
	Search(collectionName, query string, topK int) ([]VectorStoreResult, error)
	Upsert(collectionName string, documents []map[string]interface{}) error
	Delete(collectionName string, ids []string) error
	CreateCollection(collectionName string, dimension int) error
	DeleteCollection(collectionName string) error
}

// ---- GORM 实现 ----

type UserRagRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRagRepository(db *gorm.DB) UserRagRepository {
	return &UserRagRepositoryImpl{DB: db}
}

func (r *UserRagRepositoryImpl) FindByID(id string) (*UserRagEntity, error) {
	var e UserRagEntity
	result := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *UserRagRepositoryImpl) FindByUserID(userID string) ([]UserRagEntity, error) {
	var entities []UserRagEntity
	result := r.DB.Where("user_id = ? AND deleted_at IS NULL", userID).Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *UserRagRepositoryImpl) FindByUserIDAndOriginalRagID(userID, originalRagID string) (*UserRagEntity, error) {
	var e UserRagEntity
	result := r.DB.Where("user_id = ? AND original_rag_id = ? AND deleted_at IS NULL", userID, originalRagID).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *UserRagRepositoryImpl) Create(e *UserRagEntity) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return r.DB.Create(e).Error
}

func (r *UserRagRepositoryImpl) Update(e *UserRagEntity) error {
	return r.DB.Save(e).Error
}

func (r *UserRagRepositoryImpl) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&UserRagEntity{}).Error
}

// ---- RagVersion GORM 实现 ----

type RagVersionRepositoryImpl struct {
	DB *gorm.DB
}

func NewRagVersionRepository(db *gorm.DB) RagVersionRepository {
	return &RagVersionRepositoryImpl{DB: db}
}

func (r *RagVersionRepositoryImpl) FindByID(id string) (*RagVersionEntity, error) {
	var e RagVersionEntity
	result := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *RagVersionRepositoryImpl) FindByOriginalRagID(originalRagID string) ([]RagVersionEntity, error) {
	var entities []RagVersionEntity
	result := r.DB.Where("original_rag_id = ? AND deleted_at IS NULL", originalRagID).Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *RagVersionRepositoryImpl) FindPublished() ([]RagVersionEntity, error) {
	var entities []RagVersionEntity
	result := r.DB.Where("publish_status = ? AND deleted_at IS NULL", int(RagPublishStatusPublished)).Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *RagVersionRepositoryImpl) FindPublishedPaged(page, pageSize int) ([]RagVersionEntity, int64, error) {
	var entities []RagVersionEntity
	var total int64
	query := r.DB.Model(&RagVersionEntity{}).Where("publish_status = ? AND deleted_at IS NULL", int(RagPublishStatusPublished))
	query.Count(&total)
	result := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities)
	return entities, total, result.Error
}

func (r *RagVersionRepositoryImpl) Create(e *RagVersionEntity) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return r.DB.Create(e).Error
}

func (r *RagVersionRepositoryImpl) Update(e *RagVersionEntity) error {
	return r.DB.Save(e).Error
}

func (r *RagVersionRepositoryImpl) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&RagVersionEntity{}).Error
}

// ---- FileDetail GORM 实现 ----

type FileDetailRepositoryImpl struct {
	DB *gorm.DB
}

func NewFileDetailRepository(db *gorm.DB) FileDetailRepository {
	return &FileDetailRepositoryImpl{DB: db}
}

func (r *FileDetailRepositoryImpl) FindByID(id string) (*FileDetailEntity, error) {
	var e FileDetailEntity
	result := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&e)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &e, nil
}

func (r *FileDetailRepositoryImpl) FindByDataSetID(dataSetID string) ([]FileDetailEntity, error) {
	var entities []FileDetailEntity
	result := r.DB.Where("data_set_id = ? AND deleted_at IS NULL", dataSetID).Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *FileDetailRepositoryImpl) FindByUserID(userID string) ([]FileDetailEntity, error) {
	var entities []FileDetailEntity
	result := r.DB.Where("user_id = ? AND deleted_at IS NULL", userID).Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *FileDetailRepositoryImpl) Create(e *FileDetailEntity) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return r.DB.Create(e).Error
}

func (r *FileDetailRepositoryImpl) Update(e *FileDetailEntity) error {
	return r.DB.Save(e).Error
}

func (r *FileDetailRepositoryImpl) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&FileDetailEntity{}).Error
}

// ---- DocumentUnit GORM 实现 ----

type DocumentUnitRepositoryImpl struct {
	DB *gorm.DB
}

func NewDocumentUnitRepository(db *gorm.DB) DocumentUnitRepository {
	return &DocumentUnitRepositoryImpl{DB: db}
}

func (r *DocumentUnitRepositoryImpl) FindByFileID(fileID string) ([]DocumentUnitEntity, error) {
	var entities []DocumentUnitEntity
	result := r.DB.Where("file_id = ? AND deleted_at IS NULL", fileID).Order("page_index ASC").Find(&entities)
	return entities, result.Error
}

func (r *DocumentUnitRepositoryImpl) Create(e *DocumentUnitEntity) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return r.DB.Create(e).Error
}

func (r *DocumentUnitRepositoryImpl) BatchCreate(entities []DocumentUnitEntity) error {
	if len(entities) == 0 {
		return nil
	}
	for i := range entities {
		if entities[i].ID == "" {
			entities[i].ID = uuid.New().String()
		}
	}
	return r.DB.Create(&entities).Error
}

func (r *DocumentUnitRepositoryImpl) DeleteByFileID(fileID string) error {
	return r.DB.Where("file_id = ?", fileID).Delete(&DocumentUnitEntity{}).Error
}

// ---- RagQaDataset GORM 实现 ----

type RagQaDatasetRepositoryImpl struct {
	DB *gorm.DB
}

func NewRagQaDatasetRepository(db *gorm.DB) RagQaDatasetRepository {
	return &RagQaDatasetRepositoryImpl{DB: db}
}

func (r *RagQaDatasetRepositoryImpl) FindByUserRagID(userRagID string) ([]RagQaDatasetEntity, error) {
	var entities []RagQaDatasetEntity
	result := r.DB.Where("user_rag_id = ? AND deleted_at IS NULL", userRagID).Order("created_at DESC").Find(&entities)
	return entities, result.Error
}

func (r *RagQaDatasetRepositoryImpl) Create(e *RagQaDatasetEntity) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return r.DB.Create(e).Error
}

func (r *RagQaDatasetRepositoryImpl) DeleteByUserRagID(userRagID string) error {
	return r.DB.Where("user_rag_id = ?", userRagID).Delete(&RagQaDatasetEntity{}).Error
}
