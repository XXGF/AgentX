package usage

import (
	"time"

	domainUser "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/user"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// UsageRecordDTO 用量记录DTO
type UsageRecordDTO struct {
	ID                 string                 `json:"id"`
	UserID             string                 `json:"userId"`
	ProductID          string                 `json:"productId"`
	QuantityData       map[string]interface{} `json:"quantityData"`
	Cost               float64                `json:"cost"`
	RequestID          string                 `json:"requestId"`
	BilledAt           *time.Time             `json:"billedAt"`
	ServiceName        string                 `json:"serviceName"`
	ServiceType        string                 `json:"serviceType"`
	ServiceDescription string                 `json:"serviceDescription"`
	PricingRule        string                 `json:"pricingRule"`
	RelatedEntityName  string                 `json:"relatedEntityName"`
	CreatedAt          time.Time              `json:"createdAt"`
	UpdatedAt          time.Time              `json:"updatedAt"`
}

// QueryUsageRecordRequest 查询用量记录请求
type QueryUsageRecordRequest struct {
	UserID    string     `form:"userId" json:"userId"`
	ProductID string     `form:"productId" json:"productId"`
	RequestID string     `form:"requestId" json:"requestId"`
	StartTime *time.Time `form:"startTime" json:"startTime"`
	EndTime   *time.Time `form:"endTime" json:"endTime"`
	Page      int        `form:"page" json:"page"`
	PageSize  int        `form:"pageSize" json:"pageSize"`
}

func entityToDTO(e *domainUser.UsageRecordEntity) *UsageRecordDTO {
	if e == nil {
		return nil
	}
	return &UsageRecordDTO{
		ID:                 e.ID,
		UserID:             e.UserID,
		ProductID:          e.ProductID,
		QuantityData:       e.QuantityData,
		Cost:               e.Cost,
		RequestID:          e.RequestID,
		BilledAt:           e.BilledAt,
		ServiceName:        e.ServiceName,
		ServiceType:        e.ServiceType,
		ServiceDescription: e.ServiceDescription,
		PricingRule:        e.PricingRule,
		RelatedEntityName:  e.RelatedEntityName,
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
	}
}

func entitiesToDTOs(entities []domainUser.UsageRecordEntity) []*UsageRecordDTO {
	dtos := make([]*UsageRecordDTO, len(entities))
	for i, e := range entities {
		dtos[i] = entityToDTO(&e)
	}
	return dtos
}

// AppService 用量记录应用服务
type AppService struct {
	usageRecordDomainService *domainUser.UsageRecordDomainService
}

func NewAppService(usageRecordDomainService *domainUser.UsageRecordDomainService) *AppService {
	return &AppService{usageRecordDomainService: usageRecordDomainService}
}

// GetUsageRecordByID 根据ID获取用量记录
func (s *AppService) GetUsageRecordByID(recordID string) (*UsageRecordDTO, error) {
	record, err := s.usageRecordDomainService.GetUsageRecordByID(recordID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, exception.NewBusinessException("使用记录不存在")
	}
	return entityToDTO(record), nil
}

// QueryUsageRecords 按条件查询用量记录
func (s *AppService) QueryUsageRecords(req *QueryUsageRecordRequest) ([]*UsageRecordDTO, int64, error) {
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 15
	}
	entities, total, err := s.usageRecordDomainService.QueryUsageRecords(
		req.UserID, req.ProductID, req.RequestID,
		req.StartTime, req.EndTime, page, pageSize,
	)
	if err != nil {
		return nil, 0, err
	}
	return entitiesToDTOs(entities), total, nil
}

// GetUserTotalCost 获取用户总消费
func (s *AppService) GetUserTotalCost(userID string) (float64, error) {
	return s.usageRecordDomainService.GetUserTotalCost(userID)
}
