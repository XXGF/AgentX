package rag

// ---- 枚举 ----

// RagPublishStatus RAG版本发布状态枚举
type RagPublishStatus int

const (
	RagPublishStatusReviewing RagPublishStatus = 1 // 审核中
	RagPublishStatusPublished RagPublishStatus = 2 // 已发布
	RagPublishStatusRejected  RagPublishStatus = 3 // 拒绝
	RagPublishStatusRemoved   RagPublishStatus = 4 // 已下架
)

// FileProcessingStatus 文件处理状态枚举
type FileProcessingStatus int

const (
	FileProcessingStatusUploaded            FileProcessingStatus = 0 // 已上传，待开始处理
	FileProcessingStatusOCRProcessing       FileProcessingStatus = 1 // OCR处理中
	FileProcessingStatusOCRCompleted        FileProcessingStatus = 2 // OCR处理完成，待向量化
	FileProcessingStatusEmbeddingProcessing FileProcessingStatus = 3 // 向量化处理中
	FileProcessingStatusCompleted           FileProcessingStatus = 4 // 全部处理完成
	FileProcessingStatusOCRFailed           FileProcessingStatus = 5 // OCR处理失败
	FileProcessingStatusEmbeddingFailed     FileProcessingStatus = 6 // 向量化处理失败
)

// IsProcessing 是否为处理中状态
func (s FileProcessingStatus) IsProcessing() bool {
	return s == FileProcessingStatusOCRProcessing || s == FileProcessingStatusEmbeddingProcessing
}

// IsFailed 是否为失败状态
func (s FileProcessingStatus) IsFailed() bool {
	return s == FileProcessingStatusOCRFailed || s == FileProcessingStatusEmbeddingFailed
}

// IsCompleted 是否为完成状态
func (s FileProcessingStatus) IsCompleted() bool {
	return s == FileProcessingStatusCompleted
}

// CanStartOCR 是否可以开始OCR处理
func (s FileProcessingStatus) CanStartOCR() bool {
	return s == FileProcessingStatusUploaded
}

// CanStartEmbedding 是否可以开始向量化处理
func (s FileProcessingStatus) CanStartEmbedding() bool {
	return s == FileProcessingStatusOCRCompleted
}

// EmbeddingStatus 向量化状态枚举
type EmbeddingStatus int

const (
	EmbeddingStatusPending    EmbeddingStatus = 0 // 待处理
	EmbeddingStatusProcessing EmbeddingStatus = 1 // 处理中
	EmbeddingStatusCompleted  EmbeddingStatus = 2 // 已完成
	EmbeddingStatusFailed     EmbeddingStatus = 3 // 失败
)

// InstallType 安装类型枚举
type InstallType string

const (
	InstallTypeSnapshot  InstallType = "SNAPSHOT"  // 快照安装
	InstallTypeReference InstallType = "REFERENCE" // 引用安装
)

// IsReference 是否为引用类型
func (t InstallType) IsReference() bool {
	return t == InstallTypeReference
}

// IsSnapshot 是否为快照类型
func (t InstallType) IsSnapshot() bool {
	return t == InstallTypeSnapshot
}

// SearchType 搜索类型枚举
type SearchType string

const (
	SearchTypeVector  SearchType = "VECTOR"  // 向量搜索
	SearchTypeKeyword SearchType = "KEYWORD" // 关键词搜索
	SearchTypeHybrid  SearchType = "HYBRID"  // 混合搜索
)

// DocumentProcessingType 文档处理类型枚举
type DocumentProcessingType string

const (
	DocumentProcessingTypeMarkdown DocumentProcessingType = "MARKDOWN" // Markdown文档
	DocumentProcessingTypePDF      DocumentProcessingType = "PDF"      // PDF文档
	DocumentProcessingTypeTXT      DocumentProcessingType = "TXT"      // 纯文本文档
	DocumentProcessingTypeWORD     DocumentProcessingType = "WORD"     // Word文档
)

// SegmentType 分段类型枚举
type SegmentType string

const (
	SegmentTypeTitle     SegmentType = "TITLE"     // 标题
	SegmentTypeParagraph SegmentType = "PARAGRAPH" // 段落
	SegmentTypeCode      SegmentType = "CODE"      // 代码块
	SegmentTypeTable     SegmentType = "TABLE"     // 表格
	SegmentTypeList      SegmentType = "LIST"      // 列表
	SegmentTypeImage     SegmentType = "IMAGE"     // 图片
)
