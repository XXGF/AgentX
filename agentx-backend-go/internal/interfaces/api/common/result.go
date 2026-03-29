package common

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Result 统一API响应结果（对应 Java 的 org.xhy.interfaces.api.common.Result<T>）
type Result struct {
	Code      int         `json:"code"`      // 状态码
	Message   string      `json:"message"`   // 响应消息
	Data      interface{} `json:"data"`      // 响应数据
	Timestamp int64       `json:"timestamp"` // 时间戳
}

// Success 成功响应（无数据）
func Success() *Result {
	return &Result{
		Code:      http.StatusOK,
		Message:   "操作成功",
		Timestamp: time.Now().UnixMilli(),
	}
}

// SuccessWithData 成功响应（有数据）
func SuccessWithData(data interface{}) *Result {
	return &Result{
		Code:      http.StatusOK,
		Message:   "操作成功",
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	}
}

// SuccessWithMessage 成功响应（自定义消息和数据）
func SuccessWithMessage(message string, data interface{}) *Result {
	return &Result{
		Code:      http.StatusOK,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	}
}

// Error 失败响应
func Error(code int, message string) *Result {
	return &Result{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().UnixMilli(),
	}
}

// BadRequest 参数错误（400）
func BadRequest(message string) *Result {
	return Error(http.StatusBadRequest, message)
}

// Unauthorized 未授权（401）
func Unauthorized(message string) *Result {
	return Error(http.StatusUnauthorized, message)
}

// Forbidden 禁止访问（403）
func Forbidden(message string) *Result {
	return Error(http.StatusForbidden, message)
}

// NotFound 资源不存在（404）
func NotFound(message string) *Result {
	return Error(http.StatusNotFound, message)
}

// ServerError 服务器内部错误（500）
func ServerError(message string) *Result {
	return Error(http.StatusInternalServerError, message)
}

// SetMessage 链式设置消息（对应 Java 的 Result.message()）
func (r *Result) SetMessage(message string) *Result {
	r.Message = message
	return r
}

// JSON 将 Result 写入 Gin 响应
func (r *Result) JSON(c *gin.Context) {
	c.JSON(r.Code, r)
}

// SuccessJSON 快捷方法：成功响应并写入 Gin
func SuccessJSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessWithData(data))
}

// ErrorJSON 快捷方法：错误响应并写入 Gin
func ErrorJSON(c *gin.Context, code int, message string) {
	c.JSON(code, Error(code, message))
}
