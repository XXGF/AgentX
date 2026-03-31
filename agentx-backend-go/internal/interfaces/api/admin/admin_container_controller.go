package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	appContainer "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/container"
	domainContainer "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/container"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// AdminContainerController 管理员容器控制器
type AdminContainerController struct {
	appService *appContainer.AppService
}

func NewAdminContainerController(appService *appContainer.AppService) *AdminContainerController {
	return &AdminContainerController{appService: appService}
}

// GetContainersPage 分页查询容器
// GET /admin/containers
func (ctrl *AdminContainerController) GetContainersPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "15"))
	keyword := c.Query("keyword")

	var status *domainContainer.ContainerStatus
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		st := domainContainer.ContainerStatus(v)
		status = &st
	}

	var containerType *domainContainer.ContainerType
	if t := c.Query("type"); t != "" {
		ct := domainContainer.ContainerType(t)
		containerType = &ct
	}

	containers, total, err := ctrl.appService.GetContainersPage(page, pageSize, keyword, status, containerType)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, gin.H{"records": containers, "total": total})
}

// GetContainerByID 根据ID获取容器
// GET /admin/containers/:containerId
func (ctrl *AdminContainerController) GetContainerByID(c *gin.Context) {
	containerID := c.Param("containerId")
	container, err := ctrl.appService.GetContainerByID(containerID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, container)
}

// StartContainer 启动容器
// POST /admin/containers/:containerId/start
func (ctrl *AdminContainerController) StartContainer(c *gin.Context) {
	containerID := c.Param("containerId")
	if err := ctrl.appService.StartContainer(containerID); err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, nil)
}

// StopContainer 停止容器
// POST /admin/containers/:containerId/stop
func (ctrl *AdminContainerController) StopContainer(c *gin.Context) {
	containerID := c.Param("containerId")
	if err := ctrl.appService.StopContainer(containerID); err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, nil)
}

// DeleteContainer 删除容器
// DELETE /admin/containers/:containerId
func (ctrl *AdminContainerController) DeleteContainer(c *gin.Context) {
	containerID := c.Param("containerId")
	if err := ctrl.appService.DeleteContainer(containerID); err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, nil)
}

// GetContainerLogs 获取容器日志
// GET /admin/containers/:containerId/logs
func (ctrl *AdminContainerController) GetContainerLogs(c *gin.Context) {
	containerID := c.Param("containerId")
	lines, _ := strconv.Atoi(c.DefaultQuery("lines", "100"))
	logs, err := ctrl.appService.GetContainerLogs(containerID, lines)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, logs)
}

// AdminContainerTemplateController 管理员容器模板控制器
type AdminContainerTemplateController struct {
	appService *appContainer.TemplateAppService
}

func NewAdminContainerTemplateController(appService *appContainer.TemplateAppService) *AdminContainerTemplateController {
	return &AdminContainerTemplateController{appService: appService}
}

// GetTemplate 获取模板详情
// GET /admin/container-templates/:templateId
func (ctrl *AdminContainerTemplateController) GetTemplate(c *gin.Context) {
	templateID := c.Param("templateId")
	template, err := ctrl.appService.GetTemplate(templateID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, template)
}

// GetEnabledTemplates 获取所有启用的模板
// GET /admin/container-templates
func (ctrl *AdminContainerTemplateController) GetEnabledTemplates(c *gin.Context) {
	templates, err := ctrl.appService.GetEnabledTemplates()
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, templates)
}

// DeleteTemplate 删除模板
// DELETE /admin/container-templates/:templateId
func (ctrl *AdminContainerTemplateController) DeleteTemplate(c *gin.Context) {
	templateID := c.Param("templateId")
	if err := ctrl.appService.DeleteTemplate(templateID); err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, nil)
}
