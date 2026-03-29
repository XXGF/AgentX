package portal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appApiKey "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/apikey"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// ApiKeyController API密钥控制器
type ApiKeyController struct {
	appService *appApiKey.AppService
}

func NewApiKeyController(appService *appApiKey.AppService) *ApiKeyController {
	return &ApiKeyController{appService: appService}
}

// CreateApiKey 创建API密钥（POST /api-keys）
func (ctrl *ApiKeyController) CreateApiKey(c *gin.Context) {
	var req appApiKey.CreateApiKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest(err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.CreateApiKey(req.AgentID, req.Name, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetUserApiKeys 获取用户的API密钥列表（GET /api-keys）
func (ctrl *ApiKeyController) GetUserApiKeys(c *gin.Context) {
	var req appApiKey.QueryApiKeyRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest(err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.GetUserApiKeys(userID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetAgentApiKeys 获取Agent的API密钥列表（GET /api-keys/agent/:agentId）
func (ctrl *ApiKeyController) GetAgentApiKeys(c *gin.Context) {
	agentID := c.Param("agentId")
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.GetAgentApiKeys(agentID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetApiKey 获取API密钥详情（GET /api-keys/:apiKeyId）
func (ctrl *ApiKeyController) GetApiKey(c *gin.Context) {
	apiKeyID := c.Param("apiKeyId")
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.GetApiKey(apiKeyID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// UpdateApiKeyStatus 更新API密钥状态（PUT /api-keys/:apiKeyId/status）
func (ctrl *ApiKeyController) UpdateApiKeyStatus(c *gin.Context) {
	apiKeyID := c.Param("apiKeyId")
	var req appApiKey.UpdateApiKeyStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest(err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)
	if err := ctrl.appService.UpdateApiKeyStatus(apiKeyID, *req.Status, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}

// DeleteApiKey 删除API密钥（DELETE /api-keys/:apiKeyId）
func (ctrl *ApiKeyController) DeleteApiKey(c *gin.Context) {
	apiKeyID := c.Param("apiKeyId")
	userID := auth.GetCurrentUserID(c)
	if err := ctrl.appService.DeleteApiKey(apiKeyID, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}

// ResetApiKey 重置API密钥（POST /api-keys/:apiKeyId/reset）
func (ctrl *ApiKeyController) ResetApiKey(c *gin.Context) {
	apiKeyID := c.Param("apiKeyId")
	userID := auth.GetCurrentUserID(c)
	result, err := ctrl.appService.ResetApiKey(apiKeyID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}
