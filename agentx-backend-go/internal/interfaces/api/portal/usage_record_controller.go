package portal

import (
	"github.com/gin-gonic/gin"
	appUsage "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/usage"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
)

// UsageRecordController 用量记录控制器（对应 Java 的 PortalUsageRecordController）
type UsageRecordController struct {
	appService *appUsage.AppService
}

func NewUsageRecordController(appService *appUsage.AppService) *UsageRecordController {
	return &UsageRecordController{appService: appService}
}

// GetUsageRecordByID 根据ID获取使用记录
// GET /usage-records/:recordId
func (ctrl *UsageRecordController) GetUsageRecordByID(c *gin.Context) {
	recordID := c.Param("recordId")
	record, err := ctrl.appService.GetUsageRecordByID(recordID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, record)
}

// QueryUsageRecords 按条件查询当前用户使用记录
// GET /usage-records
func (ctrl *UsageRecordController) QueryUsageRecords(c *gin.Context) {
	var req appUsage.QueryUsageRecordRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ErrorJSON(c, 400, "参数错误: "+err.Error())
		return
	}
	// 前台API只能查询当前用户的记录
	req.UserID = auth.GetCurrentUserID(c)
	records, total, err := ctrl.appService.QueryUsageRecords(&req)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, gin.H{
		"records": records,
		"total":   total,
	})
}

// GetCurrentUserTotalCost 获取当前用户的总消费金额
// GET /usage-records/current/total-cost
func (ctrl *UsageRecordController) GetCurrentUserTotalCost(c *gin.Context) {
	userID := auth.GetCurrentUserID(c)
	totalCost, err := ctrl.appService.GetUserTotalCost(userID)
	if err != nil {
		common.ErrorJSON(c, 500, err.Error())
		return
	}
	common.SuccessJSON(c, totalCost)
}
