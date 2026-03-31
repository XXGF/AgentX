package file

import "time"

// FileType 文件类型枚举
type FileType string

const (
	FileTypeRAG     FileType = "rag"     // RAG文档文件
	FileTypeAvatar  FileType = "avatar"  // 用户头像文件
	FileTypeGeneral FileType = "general" // 通用文件（默认）
)

// FileTypeFromCode 根据代码获取文件类型
func FileTypeFromCode(code string) FileType {
	switch code {
	case "rag":
		return FileTypeRAG
	case "avatar":
		return FileTypeAvatar
	default:
		return FileTypeGeneral
	}
}

// GetDescription 获取文件类型描述
func (ft FileType) GetDescription() string {
	switch ft {
	case FileTypeRAG:
		return "RAG文档"
	case FileTypeAvatar:
		return "用户头像"
	case FileTypeGeneral:
		return "通用文件"
	default:
		return "通用文件"
	}
}

// FileInfo 文件信息领域模型
type FileInfo struct {
	FileID       string    `json:"fileId"`       // 文件唯一标识
	OriginalName string    `json:"originalName"` // 原始文件名
	StorageName  string    `json:"storageName"`  // 存储文件名
	FileSize     int64     `json:"fileSize"`     // 文件大小（字节）
	ContentType  string    `json:"contentType"`  // 文件类型
	BucketName   string    `json:"bucketName"`   // 存储桶名称
	FilePath     string    `json:"filePath"`     // 文件路径
	AccessURL    string    `json:"accessUrl"`    // 文件访问URL
	CreatedAt    time.Time `json:"createdAt"`    // 创建时间
	MD5Hash      string    `json:"md5Hash"`      // 文件MD5值
}

// NewFileInfo 创建文件信息
func NewFileInfo(fileID, originalName, storageName string, fileSize int64, contentType, bucketName, filePath, accessURL, md5Hash string) *FileInfo {
	return &FileInfo{
		FileID:       fileID,
		OriginalName: originalName,
		StorageName:  storageName,
		FileSize:     fileSize,
		ContentType:  contentType,
		BucketName:   bucketName,
		FilePath:     filePath,
		AccessURL:    accessURL,
		MD5Hash:      md5Hash,
		CreatedAt:    time.Now(),
	}
}

// StorageService 存储服务接口（基础设施层实现）
type StorageService interface {
	// Upload 上传文件
	Upload(fileType FileType, fileName string, data []byte) (*FileInfo, error)
	// Download 下载文件
	Download(filePath string) ([]byte, error)
	// Delete 删除文件
	Delete(filePath string) error
	// GetAccessURL 获取文件访问URL
	GetAccessURL(filePath string) (string, error)
}
