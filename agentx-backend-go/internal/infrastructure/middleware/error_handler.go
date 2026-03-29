package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
	"go.uber.org/zap"
)

// ErrorHandler 全局错误处理中间件（对应 Java 的 GlobalExceptionHandler）
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 处理 Gin 中间件/handler 中通过 c.Error() 添加的错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleError(c, err)
			return
		}
	}
}

// handleError 根据错误类型返回对应的响应
func handleError(c *gin.Context, err error) {
	requestURL := c.Request.URL.Path

	// 处理业务异常（对应 Java 的 BusinessException）
	var bizErr *exception.BusinessException
	if errors.As(err, &bizErr) {
		zap.L().Error("业务异常",
			zap.String("message", bizErr.Message),
			zap.String("url", requestURL),
		)
		c.JSON(http.StatusBadRequest, common.BadRequest(bizErr.Message))
		c.Abort()
		return
	}

	// 处理实体未找到异常（对应 Java 的 EntityNotFoundException）
	var notFoundErr *exception.EntityNotFoundException
	if errors.As(err, &notFoundErr) {
		zap.L().Error("实体未找到异常",
			zap.String("message", notFoundErr.Message),
			zap.String("url", requestURL),
		)
		c.JSON(http.StatusNotFound, common.NotFound(notFoundErr.Message))
		c.Abort()
		return
	}

	// 处理参数校验异常（对应 Java 的 ParamValidationException）
	var paramErr *exception.ParamValidationException
	if errors.As(err, &paramErr) {
		zap.L().Error("参数校验异常",
			zap.String("message", paramErr.Message),
			zap.String("url", requestURL),
		)
		c.JSON(http.StatusBadRequest, common.BadRequest(paramErr.Message))
		c.Abort()
		return
	}

	// 处理余额不足异常
	var balanceErr *exception.InsufficientBalanceException
	if errors.As(err, &balanceErr) {
		zap.L().Error("余额不足异常",
			zap.String("message", balanceErr.Message),
			zap.String("url", requestURL),
		)
		c.JSON(http.StatusBadRequest, common.BadRequest(balanceErr.Message))
		c.Abort()
		return
	}

	// 处理限流异常
	var rateLimitErr *exception.RateLimitException
	if errors.As(err, &rateLimitErr) {
		zap.L().Error("限流异常",
			zap.String("message", rateLimitErr.Message),
			zap.String("url", requestURL),
		)
		c.JSON(http.StatusTooManyRequests, common.Error(429, rateLimitErr.Message))
		c.Abort()
		return
	}

	// 处理未预期的异常
	zap.L().Error("未预期的异常",
		zap.Error(err),
		zap.String("url", requestURL),
	)
	c.JSON(http.StatusInternalServerError, common.ServerError("服务器内部错误: "+err.Error()))
	c.Abort()
}

// Recovery Panic恢复中间件（补充Gin默认的Recovery）
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		zap.L().Error("Panic恢复",
			zap.Any("error", recovered),
			zap.String("url", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, common.ServerError("服务器内部错误"))
		c.Abort()
	})
}
