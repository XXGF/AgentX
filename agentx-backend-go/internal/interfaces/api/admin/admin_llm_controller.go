package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	appLLM "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/llm"
	domain "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/llm"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// AdminLLMController 管理员LLM控制器（对应 Java 的 AdminLLMController）
type AdminLLMController struct {
	adminLLMAppService *appLLM.AdminAppService
}

// NewAdminLLMController 创建管理员LLM控制器
func NewAdminLLMController(adminLLMAppService *appLLM.AdminAppService) *AdminLLMController {
	return &AdminLLMController{adminLLMAppService: adminLLMAppService}
}

// GetProviders 获取服务商列表（GET /admin/llms/providers）
func (ctrl *AdminLLMController) GetProviders(c *gin.Context) {
	userID := auth.GetCurrentUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	result, err := ctrl.adminLLMAppService.GetOfficialProviders(userID, page, pageSize)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// GetProviderDetail 获取服务商详情（GET /admin/llms/providers/:providerId）
func (ctrl *AdminLLMController) GetProviderDetail(c *gin.Context) {
	providerID := c.Param("providerId")
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.adminLLMAppService.GetProviderDetail(providerID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// CreateProvider 创建服务商（POST /admin/llms/providers）
func (ctrl *AdminLLMController) CreateProvider(c *gin.Context) {
	var req appLLM.ProviderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.adminLLMAppService.CreateProvider(&req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// UpdateProvider 更新服务商（PUT /admin/llms/providers/:id）
func (ctrl *AdminLLMController) UpdateProvider(c *gin.Context) {
	id := c.Param("id")
	var req appLLM.ProviderUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}
	req.ID = id
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.adminLLMAppService.UpdateProvider(&req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// ToggleProviderStatus 切换服务商状态（POST /admin/llms/providers/:id/status）
func (ctrl *AdminLLMController) ToggleProviderStatus(c *gin.Context) {
	id := c.Param("id")
	userID := auth.GetCurrentUserID(c)

	if err := ctrl.adminLLMAppService.ToggleProviderStatus(id, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}

// DeleteProvider 删除服务商（DELETE /admin/llms/providers/:id）
func (ctrl *AdminLLMController) DeleteProvider(c *gin.Context) {
	id := c.Param("id")
	userID := auth.GetCurrentUserID(c)

	if err := ctrl.adminLLMAppService.DeleteProvider(id, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}

// GetProviderProtocols 获取支持的协议列表（GET /admin/llms/providers/protocols）
func (ctrl *AdminLLMController) GetProviderProtocols(c *gin.Context) {
	c.JSON(http.StatusOK, common.SuccessWithData(ctrl.adminLLMAppService.GetProviderProtocols()))
}

// GetModels 获取模型列表（GET /admin/llms/models）
func (ctrl *AdminLLMController) GetModels(c *gin.Context) {
	userID := auth.GetCurrentUserID(c)
	providerIDStr := c.Query("providerId")
	modelTypeStr := c.Query("modelType")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	var providerID *string
	if providerIDStr != "" {
		providerID = &providerIDStr
	}

	var modelType *domain.ModelType
	if modelTypeStr != "" {
		mt, ok := domain.ModelTypeFromCode(modelTypeStr)
		if ok {
			modelType = &mt
		}
	}

	result, err := ctrl.adminLLMAppService.GetOfficialModels(userID, providerID, modelType, page, pageSize)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// CreateModel 创建模型（POST /admin/llms/models）
func (ctrl *AdminLLMController) CreateModel(c *gin.Context) {
	var req appLLM.ModelCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.adminLLMAppService.CreateModel(&req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// UpdateModel 更新模型（PUT /admin/llms/models/:id）
func (ctrl *AdminLLMController) UpdateModel(c *gin.Context) {
	id := c.Param("id")
	var req appLLM.ModelUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequest("参数错误: "+err.Error()))
		return
	}
	req.ID = id
	userID := auth.GetCurrentUserID(c)

	result, err := ctrl.adminLLMAppService.UpdateModel(&req, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessWithData(result))
}

// ToggleModelStatus 切换模型状态（POST /admin/llms/models/:id/status）
func (ctrl *AdminLLMController) ToggleModelStatus(c *gin.Context) {
	id := c.Param("id")
	userID := auth.GetCurrentUserID(c)

	if err := ctrl.adminLLMAppService.ToggleModelStatus(id, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}

// DeleteModel 删除模型（DELETE /admin/llms/models/:id）
func (ctrl *AdminLLMController) DeleteModel(c *gin.Context) {
	id := c.Param("id")
	userID := auth.GetCurrentUserID(c)

	if err := ctrl.adminLLMAppService.DeleteModel(id, userID); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, common.Success())
}

// GetModelTypes 获取模型类型列表（GET /admin/llms/models/types）
func (ctrl *AdminLLMController) GetModelTypes(c *gin.Context) {
	c.JSON(http.StatusOK, common.SuccessWithData(ctrl.adminLLMAppService.GetModelTypes()))
}
