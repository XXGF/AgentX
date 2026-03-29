package portal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appLLM "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/llm"
	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/llm"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// LLMController 大模型服务商控制器（对应 Java 的 PortalLLMController）
type LLMController struct {
	llmAppService *appLLM.AppService
}

// NewLLMController 创建LLM控制器
func NewLLMController(llmAppService *appLLM.AppService) *LLMController {
	return &LLMController{llmAppService: llmAppService}
}

// GetProviderDetail 获取服务商详细信息（GET /llms/providers/:providerId）
func (ctrl *LLMController) GetProviderDetail(c *gin.Context) {
	providerID := c.Param("providerId")
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.llmAppService.GetProviderDetail(providerID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetProviders 获取服务商列表（GET /llms/providers）
func (ctrl *LLMController) GetProviders(c *gin.Context) {
	typeStr := c.DefaultQuery("type", "all")
	providerType := domain.ProviderTypeFromCode(typeStr)
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.llmAppService.GetProvidersByType(providerType, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// CreateProvider 创建服务提供商（POST /llms/providers）
func (ctrl *LLMController) CreateProvider(c *gin.Context) {
	var req appLLM.ProviderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.llmAppService.CreateProvider(&req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// UpdateProvider 更新服务提供商（PUT /llms/providers）
func (ctrl *LLMController) UpdateProvider(c *gin.Context) {
	var req appLLM.ProviderUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.llmAppService.UpdateProvider(&req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// UpdateProviderStatus 修改服务商状态（POST /llms/providers/:providerId/status）
func (ctrl *LLMController) UpdateProviderStatus(c *gin.Context) {
	providerID := c.Param("providerId")
	userID := auth.GetCurrentUserID(c)

	if err := ctrl.llmAppService.UpdateProviderStatus(providerID, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}

// DeleteProvider 删除服务提供商（DELETE /llms/providers/:providerId）
func (ctrl *LLMController) DeleteProvider(c *gin.Context) {
	providerID := c.Param("providerId")
	userID := auth.GetCurrentUserID(c)

	if err := ctrl.llmAppService.DeleteProvider(providerID, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}

// GetProviderProtocols 获取服务提供商协议列表（GET /llms/providers/protocols）
func (ctrl *LLMController) GetProviderProtocols(c *gin.Context) {
	c.JSON(http.StatusOK, common.SuccessWithData(ctrl.llmAppService.GetUserProviderProtocols()))
}

// CreateModel 添加模型（POST /llms/models）
func (ctrl *LLMController) CreateModel(c *gin.Context) {
	var req appLLM.ModelCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.llmAppService.CreateModel(&req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// UpdateModel 修改模型（PUT /llms/models）
func (ctrl *LLMController) UpdateModel(c *gin.Context) {
	var req appLLM.ModelUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.llmAppService.UpdateModel(&req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// DeleteModel 删除模型（DELETE /llms/models/:modelId）
func (ctrl *LLMController) DeleteModel(c *gin.Context) {
	modelID := c.Param("modelId")
	userID := auth.GetCurrentUserID(c)

	if err := ctrl.llmAppService.DeleteModel(modelID, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}

// UpdateModelStatus 修改模型状态（PUT /llms/models/:modelId/status）
func (ctrl *LLMController) UpdateModelStatus(c *gin.Context) {
	modelID := c.Param("modelId")
	userID := auth.GetCurrentUserID(c)

	if err := ctrl.llmAppService.UpdateModelStatus(modelID, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}

// GetModelTypes 获取模型类型（GET /llms/models/types）
func (ctrl *LLMController) GetModelTypes(c *gin.Context) {
	c.JSON(http.StatusOK, common.SuccessWithData(domain.AllModelTypes))
}

// GetModels 获取所有激活模型（GET /llms/models）
func (ctrl *LLMController) GetModels(c *gin.Context) {
	userID := auth.GetCurrentUserID(c)
	modelTypeStr := c.Query("modelType")
	officialStr := c.Query("official")

	var modelType *domain.ModelType
	if modelTypeStr != "" {
		mt, ok := domain.ModelTypeFromCode(modelTypeStr)
		if ok {
			modelType = &mt
		}
	}

	providerType := domain.ProviderTypeAll
	if officialStr == "true" {
		providerType = domain.ProviderTypeOfficial
	}

	result, err := ctrl.llmAppService.GetActiveModelsByType(providerType, userID, modelType)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetDefaultModel 获取用户默认模型（GET /llms/models/default）
func (ctrl *LLMController) GetDefaultModel(c *gin.Context) {
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.llmAppService.GetDefaultModel(userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}
