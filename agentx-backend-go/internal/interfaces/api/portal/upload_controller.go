package portal

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/config"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// UploadController 文件上传控制器（对应 Java 的 UploadController）
type UploadController struct {
	uploadDir string // 本地上传目录
	platform  string // 存储平台：local / s3
	s3Config  *config.S3Config
}

func NewUploadController(uploadCfg *config.UploadConfig) *UploadController {
	uploadDir := uploadCfg.Local.BasePath
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	// 确保上传目录存在
	_ = os.MkdirAll(uploadDir, 0755)

	return &UploadController{
		uploadDir: uploadDir,
		platform:  uploadCfg.Platform,
		s3Config:  &uploadCfg.S3,
	}
}

// GetUploadCredential 获取上传凭证
// GET /upload/credential
func (ctrl *UploadController) GetUploadCredential(c *gin.Context) {
	userID := auth.GetCurrentUserID(c)

	if ctrl.platform == "s3" && ctrl.s3Config.AccessKey != "" {
		// S3 模式：返回预签名 URL 或 STS 临时凭证
		common.SuccessJSON(c, gin.H{
			"platform":  "s3",
			"uploadUrl": ctrl.s3Config.Endpoint,
			"bucket":    ctrl.s3Config.BucketName,
			"region":    ctrl.s3Config.Region,
			"basePath":  ctrl.s3Config.BasePath,
			"domain":    ctrl.s3Config.Domain,
			"userId":    userID,
			"expiresAt": time.Now().Add(30 * time.Minute).Format(time.RFC3339),
		})
		return
	}

	// 本地模式
	uploadToken := uuid.New().String()
	common.SuccessJSON(c, gin.H{
		"platform":    "local",
		"uploadToken": uploadToken,
		"uploadUrl":   "/api/upload/file",
		"userId":      userID,
		"expiresAt":   time.Now().Add(30 * time.Minute).Format(time.RFC3339),
	})
}

// UploadFile 上传文件
// POST /upload/file
func (ctrl *UploadController) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		common.ErrorJSON(c, 400, "请选择要上传的文件")
		return
	}

	// 限制文件大小（50MB）
	if file.Size > 50*1024*1024 {
		common.ErrorJSON(c, 400, "文件大小不能超过50MB")
		return
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// 按日期分目录
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(ctrl.uploadDir, dateDir)
	_ = os.MkdirAll(fullDir, 0755)

	// 保存文件
	dst := filepath.Join(fullDir, newFilename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		common.ErrorJSON(c, 500, "文件保存失败: "+err.Error())
		return
	}

	// 返回文件URL
	var fileURL string
	if ctrl.platform == "s3" && ctrl.s3Config.Domain != "" {
		fileURL = ctrl.s3Config.Domain + ctrl.s3Config.BasePath + dateDir + "/" + newFilename
	} else {
		fileURL = fmt.Sprintf("/uploads/%s/%s", dateDir, newFilename)
	}

	common.SuccessJSON(c, gin.H{
		"url":      fileURL,
		"filename": file.Filename,
		"size":     file.Size,
		"type":     file.Header.Get("Content-Type"),
	})
}
