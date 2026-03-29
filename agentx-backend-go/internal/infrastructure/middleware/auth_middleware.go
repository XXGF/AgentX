package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
	"go.uber.org/zap"
)

// 不需要认证的路径列表（对应 Java 的 WebMvcConfig 中的 excludePathPatterns）
var excludedPaths = []string{
	"/login",
	"/health",
	"/register",
	"/auth/config",
	"/send-email-code",
	"/verify-email-code",
	"/get-captcha",
	"/reset-password",
	"/send-reset-password-code",
	"/oauth/github/authorize",
	"/oauth/github/callback",
	"/sso/",
	"/widget/",
	"/v1/",
	"/payments/callback/",
}

// AuthMiddleware JWT 认证中间件（对应 Java 的 UserAuthInterceptor）
func AuthMiddleware(jwtUtils *auth.JWTUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestURI := c.Request.URL.Path
		method := c.Request.Method

		// 检查是否为排除路径
		for _, path := range excludedPaths {
			if strings.HasPrefix(requestURI, path) || requestURI == path {
				c.Next()
				return
			}
		}

		// OPTIONS 请求跳过（CORS 预检请求）
		if method == http.MethodOptions {
			c.Next()
			return
		}

		// 从请求头中获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			zap.L().Warn("认证失败 - 缺少Authorization头",
				zap.String("method", method),
				zap.String("uri", requestURI),
			)
			c.JSON(http.StatusUnauthorized, common.Unauthorized("缺少认证头"))
			c.Abort()
			return
		}

		// 检查 Bearer 前缀
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			zap.L().Warn("认证失败 - Authorization头格式错误",
				zap.String("method", method),
				zap.String("uri", requestURI),
			)
			c.JSON(http.StatusUnauthorized, common.Unauthorized("认证头格式错误"))
			c.Abort()
			return
		}

		// 提取 Token
		token := strings.TrimPrefix(authHeader, bearerPrefix)
		if token == "" {
			zap.L().Warn("认证失败 - Token为空",
				zap.String("method", method),
				zap.String("uri", requestURI),
			)
			c.JSON(http.StatusUnauthorized, common.Unauthorized("Token为空"))
			c.Abort()
			return
		}

		// 验证 Token
		if !jwtUtils.ValidateToken(token) {
			zap.L().Warn("认证失败 - Token验证失败",
				zap.String("method", method),
				zap.String("uri", requestURI),
			)
			c.JSON(http.StatusUnauthorized, common.Unauthorized("Token无效或已过期"))
			c.Abort()
			return
		}

		// 从 Token 中获取用户ID并设置到上下文
		userID, err := jwtUtils.GetUserIDFromToken(token)
		if err != nil || userID == "" {
			zap.L().Warn("认证失败 - 无法从Token中获取用户ID",
				zap.String("method", method),
				zap.String("uri", requestURI),
			)
			c.JSON(http.StatusUnauthorized, common.Unauthorized("无效的用户信息"))
			c.Abort()
			return
		}

		// 设置用户ID到上下文（对应 Java 的 UserContext.setCurrentUserId）
		auth.SetCurrentUserID(c, userID)

		zap.L().Debug("认证成功",
			zap.String("method", method),
			zap.String("uri", requestURI),
			zap.String("userId", userID),
		)

		c.Next()
	}
}

// AdminAuthMiddleware 管理员权限中间件（对应 Java 的 AdminAuthInterceptor）
// 注意：需要在 AuthMiddleware 之后使用，因为依赖 UserContext 中的 userId
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 从数据库查询用户是否为管理员
		// 这里需要注入 UserDomainService，暂时先预留接口
		// userID := auth.GetCurrentUserID(c)
		// user := userDomainService.GetUserInfo(userID)
		// if user == nil || !user.IsAdmin {
		//     c.JSON(http.StatusForbidden, common.Forbidden("无权限访问管理功能"))
		//     c.Abort()
		//     return
		// }
		c.Next()
	}
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
