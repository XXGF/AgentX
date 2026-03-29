package exception

import "fmt"

// BusinessException 业务异常（对应 Java 的 org.xhy.infrastructure.exception.BusinessException）
type BusinessException struct {
	ErrorCode string
	Message   string
	Cause     error
}

func (e *BusinessException) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// NewBusinessException 创建业务异常
func NewBusinessException(message string) *BusinessException {
	return &BusinessException{Message: message}
}

// NewBusinessExceptionWithCode 创建带错误码的业务异常
func NewBusinessExceptionWithCode(errorCode, message string) *BusinessException {
	return &BusinessException{ErrorCode: errorCode, Message: message}
}

// NewBusinessExceptionWithCause 创建带原因的业务异常
func NewBusinessExceptionWithCause(message string, cause error) *BusinessException {
	return &BusinessException{Message: message, Cause: cause}
}

// EntityNotFoundException 实体未找到异常（对应 Java 的 EntityNotFoundException）
type EntityNotFoundException struct {
	Message string
}

func (e *EntityNotFoundException) Error() string {
	return e.Message
}

// NewEntityNotFoundException 创建实体未找到异常
func NewEntityNotFoundException(message string) *EntityNotFoundException {
	return &EntityNotFoundException{Message: message}
}

// ParamValidationException 参数校验异常（对应 Java 的 ParamValidationException）
type ParamValidationException struct {
	BusinessException
}

// NewParamValidationException 创建参数校验异常
func NewParamValidationException(message string) *ParamValidationException {
	return &ParamValidationException{
		BusinessException: BusinessException{
			ErrorCode: "PARAM_VALIDATION_ERROR",
			Message:   message,
		},
	}
}

// NewParamValidationExceptionWithParam 创建带参数名的校验异常
func NewParamValidationExceptionWithParam(paramName, message string) *ParamValidationException {
	return &ParamValidationException{
		BusinessException: BusinessException{
			ErrorCode: "PARAM_VALIDATION_ERROR",
			Message:   fmt.Sprintf("参数[%s]无效: %s", paramName, message),
		},
	}
}

// InsufficientBalanceException 余额不足异常（对应 Java 的 InsufficientBalanceException）
type InsufficientBalanceException struct {
	BusinessException
}

// NewInsufficientBalanceException 创建余额不足异常
func NewInsufficientBalanceException(message string) *InsufficientBalanceException {
	return &InsufficientBalanceException{
		BusinessException: BusinessException{Message: message},
	}
}

// RateLimitException 限流异常（对应 Java 的 RateLimitException）
type RateLimitException struct {
	BusinessException
}

// NewRateLimitException 创建限流异常
func NewRateLimitException(message string) *RateLimitException {
	return &RateLimitException{
		BusinessException: BusinessException{Message: message},
	}
}
