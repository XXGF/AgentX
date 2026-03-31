package portal

import (
	"strconv"

	"github.com/gin-gonic/gin"
	appRag "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/rag"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// ---- FileOperationController 文件操作控制器 ----

type FileOperationController struct {
	appService *appRag.FileOperationAppService
}

func NewFileOperationController(appService *appRag.FileOperationAppService) *FileOperationController {
	return &FileOperationController{appService: appService}
}

// GetFilesByDataSetID 获取数据集下的文件列表
// GET /rag/files/dataset/:dataSetId
func (ctrl *FileOperationController) GetFilesByDataSetID(c *gin.Context) {
	dataSetID := c.Param("dataSetId")
	files, err := ctrl.appService.GetFilesByDataSetID(dataSetID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, files)
}

// GetDocumentUnits 获取文件的文档单元列表
// GET /rag/files/:fileId/document-units
func (ctrl *FileOperationController) GetDocumentUnits(c *gin.Context) {
	fileID := c.Param("fileId")
	units, err := ctrl.appService.GetDocumentUnitsByFileID(fileID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, units)
}

// DeleteFile 删除文件
// DELETE /rag/files/:fileId
func (ctrl *FileOperationController) DeleteFile(c *gin.Context) {
	fileID := c.Param("fileId")
	userID := auth.GetCurrentUserID(c)
	if err := ctrl.appService.DeleteFile(fileID, userID); err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, nil)
}

// ---- RagMarketController RAG市场控制器 ----

type RagMarketController struct {
	appService *appRag.RagMarketAppService
}

func NewRagMarketController(appService *appRag.RagMarketAppService) *RagMarketController {
	return &RagMarketController{appService: appService}
}

// GetMarketList 获取RAG市场列表
// GET /rag/market
func (ctrl *RagMarketController) GetMarketList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "15"))
	userID := auth.GetCurrentUserID(c)
	list, total, err := ctrl.appService.GetMarketList(page, pageSize, userID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, gin.H{"records": list, "total": total})
}

// InstallRag 安装RAG
// POST /rag/market/install
func (ctrl *RagMarketController) InstallRag(c *gin.Context) {
	var req appRag.InstallRagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorJSON(c, 400, "参数错误: "+err.Error())
		return
	}
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.InstallRag(userID, &req)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, result)
}

// UninstallRag 卸载RAG
// DELETE /rag/market/:userRagId
func (ctrl *RagMarketController) UninstallRag(c *gin.Context) {
	userRagID := c.Param("userRagId")
	userID := auth.GetCurrentUserID(c)
	if err := ctrl.appService.UninstallRag(userRagID, userID); err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, nil)
}

// GetUserRags 获取用户安装的RAG列表
// GET /rag/market/my
func (ctrl *RagMarketController) GetUserRags(c *gin.Context) {
	userID := auth.GetCurrentUserID(c)
	rags, err := ctrl.appService.GetUserRags(userID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, rags)
}

// ---- RagPublishController RAG发布控制器 ----

type RagPublishController struct {
	appService *appRag.RagPublishAppService
}

func NewRagPublishController(appService *appRag.RagPublishAppService) *RagPublishController {
	return &RagPublishController{appService: appService}
}

// PublishRag 发布RAG到市场
// POST /rag/publish
func (ctrl *RagPublishController) PublishRag(c *gin.Context) {
	var req appRag.PublishRagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorJSON(c, 400, "参数错误: "+err.Error())
		return
	}
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.PublishRag(userID, &req)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, result)
}

// GetVersions 获取RAG的所有版本
// GET /rag/publish/:originalRagId/versions
func (ctrl *RagPublishController) GetVersions(c *gin.Context) {
	originalRagID := c.Param("originalRagId")
	versions, err := ctrl.appService.GetVersionsByOriginalRagID(originalRagID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, versions)
}

// ---- RagSearchController RAG搜索控制器 ----

type RagSearchController struct {
	appService *appRag.RAGSearchAppService
}

func NewRagSearchController(appService *appRag.RAGSearchAppService) *RagSearchController {
	return &RagSearchController{appService: appService}
}

// RagSearch RAG搜索
// POST /rag/search
func (ctrl *RagSearchController) RagSearch(c *gin.Context) {
	var req appRag.RagSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorJSON(c, 400, "参数错误: "+err.Error())
		return
	}
	userID := auth.GetCurrentUserID(c)
	results, err := ctrl.appService.RagSearch(&req, userID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, results)
}

// RagSearchByUserRag 基于已安装知识库的RAG搜索
// POST /rag/search/user-rag/:userRagId
func (ctrl *RagSearchController) RagSearchByUserRag(c *gin.Context) {
	userRagID := c.Param("userRagId")
	var req appRag.RagSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorJSON(c, 400, "参数错误: "+err.Error())
		return
	}
	userID := auth.GetCurrentUserID(c)
	results, err := ctrl.appService.RagSearchByUserRag(&req, userRagID, userID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, results)
}

// ---- RagQaDatasetController RAG QA数据集控制器 ----

type RagQaDatasetController struct {
	appService *appRag.RagQaDatasetAppService
}

func NewRagQaDatasetController(appService *appRag.RagQaDatasetAppService) *RagQaDatasetController {
	return &RagQaDatasetController{appService: appService}
}

// GetQaDatasets 获取QA数据集列表
// GET /rag/qa/:userRagId
func (ctrl *RagQaDatasetController) GetQaDatasets(c *gin.Context) {
	userRagID := c.Param("userRagId")
	userID := auth.GetCurrentUserID(c)
	datasets, err := ctrl.appService.GetQaDatasets(userRagID, userID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, datasets)
}

// CreateQaDataset 创建QA数据集
// POST /rag/qa/:userRagId
func (ctrl *RagQaDatasetController) CreateQaDataset(c *gin.Context) {
	userRagID := c.Param("userRagId")
	var req struct {
		Question string `json:"question" binding:"required"`
		Answer   string `json:"answer" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorJSON(c, 400, "参数错误: "+err.Error())
		return
	}
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.CreateQaDataset(userRagID, req.Question, req.Answer, userID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, result)
}
