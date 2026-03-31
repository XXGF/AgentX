package account

import (
	"time"

	domainUser "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/user"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
)

// AccountDTO 账户DTO
type AccountDTO struct {
	ID                string     `json:"id"`
	UserID            string     `json:"userId"`
	Balance           float64    `json:"balance"`
	Credit            float64    `json:"credit"`
	TotalConsumed     float64    `json:"totalConsumed"`
	AvailableBalance  float64    `json:"availableBalance"`
	LastTransactionAt *time.Time `json:"lastTransactionAt"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

func entityToDTO(e *domainUser.AccountEntity) *AccountDTO {
	if e == nil {
		return nil
	}
	return &AccountDTO{
		ID:                e.ID,
		UserID:            e.UserID,
		Balance:           e.Balance,
		Credit:            e.Credit,
		TotalConsumed:     e.TotalConsumed,
		AvailableBalance:  e.GetAvailableBalance(),
		LastTransactionAt: e.LastTransactionAt,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}

// RechargeRequest 充值请求
type RechargeRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// AddCreditRequest 增加信用额度请求
type AddCreditRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// AppService 账户应用服务
type AppService struct {
	accountDomainService *domainUser.AccountDomainService
}

func NewAppService(accountDomainService *domainUser.AccountDomainService) *AppService {
	return &AppService{accountDomainService: accountDomainService}
}

// GetUserAccount 获取用户账户信息
func (s *AppService) GetUserAccount(userID string) (*AccountDTO, error) {
	account, err := s.accountDomainService.GetOrCreateAccount(userID)
	if err != nil {
		return nil, err
	}
	return entityToDTO(account), nil
}

// GetAccountByID 根据账户ID获取账户信息
func (s *AppService) GetAccountByID(accountID string) (*AccountDTO, error) {
	account, err := s.accountDomainService.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, exception.NewBusinessException("账户不存在")
	}
	return entityToDTO(account), nil
}

// Recharge 账户充值
func (s *AppService) Recharge(userID string, amount float64) (*AccountDTO, error) {
	if amount <= 0 {
		return nil, exception.NewBusinessException("充值金额必须大于0")
	}
	if err := s.accountDomainService.RechargeBalance(userID, amount); err != nil {
		return nil, err
	}
	account, err := s.accountDomainService.GetAccountByUserID(userID)
	if err != nil {
		return nil, err
	}
	return entityToDTO(account), nil
}

// AddCredit 增加信用额度
func (s *AppService) AddCredit(userID string, amount float64) (*AccountDTO, error) {
	if amount <= 0 {
		return nil, exception.NewBusinessException("信用额度必须大于0")
	}
	if err := s.accountDomainService.AddCredit(userID, amount); err != nil {
		return nil, err
	}
	account, err := s.accountDomainService.GetAccountByUserID(userID)
	if err != nil {
		return nil, err
	}
	return entityToDTO(account), nil
}

// CheckSufficientBalance 检查余额是否充足
func (s *AppService) CheckSufficientBalance(userID string, amount float64) (bool, error) {
	return s.accountDomainService.CheckSufficientBalance(userID, amount)
}

// GetAvailableBalance 获取可用余额
func (s *AppService) GetAvailableBalance(userID string) (float64, error) {
	return s.accountDomainService.GetAvailableBalance(userID)
}
