package auth

import (
	"github.com/gin-gonic/gin"
)

// UserContext 相关的 Gin Context Key
const (
	// UserIDKey 用户ID在Gin Context中的Key
	UserIDKey = "currentUserId"
)

// SetCurrentUserID 设置当前用户ID到 Gin Context（对应 Java 的 UserContext.setCurrentUserId）
func SetCurrentUserID(c *gin.Context, userID string) {
	c.Set(UserIDKey, userID)
}

// GetCurrentUserID 从 Gin Context 获取当前用户ID（对应 Java 的 UserContext.getCurrentUserId）
// 如果未设置则返回空字符串
func GetCurrentUserID(c *gin.Context) string {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return ""
	}
	return userID.(string)
}

// ExternalAPIContext 外部API上下文相关的 Key
const (
	ExternalUserIDKey  = "externalUserId"
	ExternalAgentIDKey = "externalAgentId"
)

// SetExternalUserID 设置外部API的用户ID
func SetExternalUserID(c *gin.Context, userID string) {
	c.Set(ExternalUserIDKey, userID)
}

// GetExternalUserID 获取外部API的用户ID
func GetExternalUserID(c *gin.Context) string {
	userID, exists := c.Get(ExternalUserIDKey)
	if !exists {
		return ""
	}
	return userID.(string)
}

// SetExternalAgentID 设置外部API的AgentID
func SetExternalAgentID(c *gin.Context, agentID string) {
	c.Set(ExternalAgentIDKey, agentID)
}

// GetExternalAgentID 获取外部API的AgentID
func GetExternalAgentID(c *gin.Context) string {
	agentID, exists := c.Get(ExternalAgentIDKey)
	if !exists {
		return ""
	}
	return agentID.(string)
}
