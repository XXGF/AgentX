package router

import (
	"github.com/gin-gonic/gin"
	appAgent "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/agent"
	appApiKey "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/apikey"
	appAuth "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/auth"
	appConv "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/conversation"
	appLLM "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/llm"
	appMemory "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/memory"
	appScheduledTask "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/scheduledtask"
	appTool "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/tool"
	appUser "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/user"
	domainAgent "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/agent"
	domainApiKey "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/apikey"
	domainAuth "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/auth"
	domainConv "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/conversation"
	domainLLM "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/llm"
	domainMemory "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/memory"
	domainScheduledTask "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/scheduledtask"
	domainTool "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/tool"
	domainUser "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/user"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/auth"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/config"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/middleware"
	infraSso "github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/sso"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/admin"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/common"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/interfaces/api/portal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Setup 初始化路由（对应 Java 的 WebMvcConfig + 各 Controller 的 @RequestMapping）
func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(middleware.Recovery())          // Panic 恢复
	r.Use(middleware.CORSMiddleware())    // 跨域
	r.Use(middleware.ErrorHandler())      // 全局错误处理
	r.Use(gin.Logger())                   // 请求日志

	// JWT 工具
	jwtUtils := auth.NewJWTUtils(cfg.JWT.Secret, cfg.JWT.Expiration)

	// API 路由组（对应 Java 的 server.servlet.context-path: /api）
	api := r.Group(cfg.Server.ContextPath)

	// 健康检查（不需要认证）
	api.GET("/health", func(c *gin.Context) {
		common.SuccessJSON(c, gin.H{"status": "UP"})
	})

	// ========== 公开接口（不需要认证）==========
	setupPublicRoutes(api, db, jwtUtils)

	// ========== 需要认证的接口 ==========
	authenticated := api.Group("")
	authenticated.Use(middleware.AuthMiddleware(jwtUtils))
	setupAuthenticatedRoutes(authenticated, db)

	// ========== 管理员接口 ==========
	adminGroup := authenticated.Group("/admin")
	adminGroup.Use(middleware.AdminAuthMiddleware())
	setupAdminRoutes(adminGroup, db)

	// ========== 外部API接口（使用API Key认证）==========
	external := api.Group("/v1")
	setupExternalRoutes(external, db)

	return r
}

// setupPublicRoutes 注册公开路由（登录、注册等）
func setupPublicRoutes(rg *gin.RouterGroup, db *gorm.DB, jwtUtils *auth.JWTUtils) {
	// === 初始化 User 模块依赖 ===
	userRepo := domainUser.NewUserRepository(db)
	userSettingsRepo := domainUser.NewUserSettingsRepository(db)
	userDomainService := domainUser.NewDomainService(userRepo, userSettingsRepo)
	loginAppService := appUser.NewLoginAppService(userDomainService, jwtUtils)

	// === 登录注册路由（对应 Java 的 LoginController）===
	loginController := portal.NewLoginController(loginAppService)
	rg.POST("/login", loginController.Login)
	rg.POST("/register", loginController.Register)
	rg.POST("/get-captcha", loginController.GetCaptcha)
	rg.POST("/send-email-code", loginController.SendEmailCode)
	rg.POST("/send-reset-password-code", loginController.SendResetPasswordCode)
	rg.POST("/verify-email-code", loginController.VerifyEmailCode)
	rg.POST("/reset-password", loginController.ResetPassword)

	// === 初始化 Auth 模块依赖 ===
	authSettingRepo := domainAuth.NewAuthSettingRepository(db)
	authDomainService := domainAuth.NewDomainService(authSettingRepo)
	authAppService := appAuth.NewAppService(authDomainService)

	// === 认证配置路由（对应 Java 的 AuthConfigController）===
	authConfigController := portal.NewAuthConfigController(authAppService)
	authGroup := rg.Group("/auth")
	{
		authGroup.GET("/config", authConfigController.GetAuthConfig)
	}

	// === SSO 路由（对应 Java 的 SsoController）===
	logger, _ := zap.NewProduction()
	ssoConfigProvider := infraSso.NewConfigProvider(authDomainService)
	githubSsoService := infraSso.NewGitHubSsoService(ssoConfigProvider, logger)
	ssoFactory := infraSso.NewServiceFactory(authDomainService, githubSsoService)
	ssoAppService := appAuth.NewSsoAppService(ssoFactory, userDomainService, jwtUtils)
	ssoController := portal.NewSsoController(ssoAppService)
	ssoGroup := rg.Group("/sso")
	{
		ssoGroup.GET("/:provider/login", ssoController.GetSsoLoginUrl)
		ssoGroup.GET("/:provider/callback", ssoController.HandleSsoCallback)
	}
}

// setupAuthenticatedRoutes 注册需要认证的路由
func setupAuthenticatedRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// === User 模块（对应 Java 的 PortalUserController）===
	userRepo := domainUser.NewUserRepository(db)
	userSettingsRepo := domainUser.NewUserSettingsRepository(db)
	userDomainService := domainUser.NewDomainService(userRepo, userSettingsRepo)
	userAppService := appUser.NewAppService(userDomainService)
	userController := portal.NewUserController(userAppService)

	users := rg.Group("/users")
	{
		users.GET("", userController.GetUserInfo)
		users.POST("", userController.UpdateUserInfo)
		users.PUT("/password", userController.ChangePassword)
		// TODO: 后续迁移 UserSettings 相关接口
		// users.GET("/settings", userController.GetUserSettings)
		// users.PUT("/settings", userController.UpdateUserSettings)
		// users.GET("/settings/default-model", userController.GetUserDefaultModelId)
		// users.GET("/settings/ocr-models", userController.GetOcrModels)
		// users.GET("/settings/embedding-models", userController.GetEmbeddingModels)
	}

	// === LLM 模块（对应 Java 的 PortalLLMController）===
	providerRepo := domainLLM.NewProviderRepository(db)
	modelRepo := domainLLM.NewModelRepository(db)
	llmDomainService := domainLLM.NewDomainService(providerRepo, modelRepo)
	llmAppService := appLLM.NewAppService(llmDomainService, userDomainService)
	llmController := portal.NewLLMController(llmAppService)

	llmProviders := rg.Group("/llms/providers")
	{
		llmProviders.GET("", llmController.GetProviders)
		llmProviders.GET("/protocols", llmController.GetProviderProtocols)
		llmProviders.GET("/:providerId", llmController.GetProviderDetail)
		llmProviders.POST("", llmController.CreateProvider)
		llmProviders.PUT("", llmController.UpdateProvider)
		llmProviders.POST("/:providerId/status", llmController.UpdateProviderStatus)
		llmProviders.DELETE("/:providerId", llmController.DeleteProvider)
	}

	llmModels := rg.Group("/llms/models")
	{
		llmModels.GET("", llmController.GetModels)
		llmModels.GET("/types", llmController.GetModelTypes)
		llmModels.GET("/default", llmController.GetDefaultModel)
		llmModels.POST("", llmController.CreateModel)
		llmModels.PUT("", llmController.UpdateModel)
		llmModels.PUT("/:modelId/status", llmController.UpdateModelStatus)
		llmModels.DELETE("/:modelId", llmController.DeleteModel)
	}

	// === Agent 模块（对应 Java 的 PortalAgentController）===
	agentRepo := domainAgent.NewAgentRepository(db)
	agentVersionRepo := domainAgent.NewAgentVersionRepository(db)
	agentWorkspaceRepo := domainAgent.NewAgentWorkspaceRepository(db)
	agentDomainService := domainAgent.NewDomainService(agentRepo, agentVersionRepo, agentWorkspaceRepo)
	agentWorkspaceDomainService := domainAgent.NewWorkspaceDomainService(agentWorkspaceRepo, agentRepo)
	agentAppService := appAgent.NewAppService(agentDomainService, agentWorkspaceDomainService)
	agentController := portal.NewAgentController(agentAppService)

	agents := rg.Group("/agents")
	{
		agents.POST("", agentController.CreateAgent)
		agents.GET("/user", agentController.GetUserAgents)
		agents.GET("/published", agentController.GetPublishedAgents)
		agents.GET("/:agentId", agentController.GetAgent)
		agents.PUT("/:agentId", agentController.UpdateAgent)
		agents.PUT("/:agentId/toggle-status", agentController.ToggleAgentStatus)
		agents.DELETE("/:agentId", agentController.DeleteAgent)
		agents.POST("/:agentId/publish", agentController.PublishAgentVersion)
		agents.GET("/:agentId/versions", agentController.GetAgentVersions)
		agents.GET("/:agentId/versions/latest", agentController.GetLatestAgentVersion)
		agents.GET("/:agentId/versions/:versionNumber", agentController.GetAgentVersion)
	}

	// === Conversation 模块（对应 Java 的 PortalAgentSessionController）===
	sessionRepo := domainConv.NewSessionRepository(db)
	messageRepo := domainConv.NewMessageRepository(db)
	contextRepo := domainConv.NewContextRepository(db)
	sessionDomainService := domainConv.NewSessionDomainService(sessionRepo)
	conversationDomainService := domainConv.NewConversationDomainService(messageRepo)
	_ = domainConv.NewMessageDomainService(messageRepo, contextRepo) // 后续Chat功能使用
	_ = domainConv.NewContextDomainService(contextRepo)               // 后续Chat功能使用

	agentSessionAppService := appConv.NewAgentSessionAppService(
		agentWorkspaceDomainService, agentDomainService,
		sessionDomainService, conversationDomainService,
	)
	conversationAppService := appConv.NewConversationAppService(conversationDomainService, sessionDomainService)
	sessionController := portal.NewSessionController(agentSessionAppService, conversationAppService)

	sessions := rg.Group("/agents/sessions")
	{
		sessions.GET("/:sessionId/messages", sessionController.GetConversationMessages)
		sessions.GET("/:agentId", sessionController.GetAgentSessionList)
		sessions.POST("/:agentId", sessionController.CreateSession)
		sessions.PUT("/:id", sessionController.UpdateSession)
		sessions.DELETE("/:id", sessionController.DeleteSession)
		sessions.POST("/chat", sessionController.Chat)
		sessions.POST("/:sessionId/interrupt", sessionController.InterruptSession)
	}

	// === Tool 模块（对应 Java 的 PortalToolController）===
	toolRepo := domainTool.NewToolRepository(db)
	toolVersionRepo := domainTool.NewToolVersionRepository(db)
	userToolRepo := domainTool.NewUserToolRepository(db)
	toolDomainService := domainTool.NewToolDomainService(toolRepo, toolVersionRepo, userToolRepo)
	userToolDomainService := domainTool.NewUserToolDomainService(userToolRepo)
	toolVersionDomainService := domainTool.NewToolVersionDomainService(toolVersionRepo)
	toolAppService := appTool.NewAppService(toolDomainService, userToolDomainService, toolVersionDomainService, userDomainService)
	toolController := portal.NewToolController(toolAppService)

	tools := rg.Group("/tools")
	{
		tools.POST("", toolController.CreateTool)
		tools.GET("/user", toolController.GetUserTools)
		tools.GET("/recommend", toolController.GetRecommendTools)
		tools.GET("/installed", toolController.GetInstalledTools)
		tools.POST("/market", toolController.MarketTool)
		tools.GET("/market", toolController.Market)
		tools.GET("/market/:toolId/versions", toolController.GetToolVersions)
		tools.GET("/market/:toolId/:version", toolController.GetToolVersionDetail)
		tools.POST("/install/:toolId/:version", toolController.InstallTool)
		tools.POST("/uninstall/:toolId", toolController.UninstallTool)
		tools.POST("/user/:toolId/:version/status", toolController.UpdateToolVersionStatus)
		tools.GET("/:toolId", toolController.GetToolDetail)
		tools.GET("/:toolId/latest", toolController.GetLatestToolVersion)
		tools.PUT("/:toolId", toolController.UpdateTool)
		tools.DELETE("/:toolId", toolController.DeleteTool)
	}

	// === ScheduledTask 模块（对应 Java 的 PortalScheduledTaskController）===
	scheduledTaskRepo := domainScheduledTask.NewScheduledTaskRepository(db)
	scheduledTaskDomainService := domainScheduledTask.NewDomainService(scheduledTaskRepo)
	scheduledTaskAppService := appScheduledTask.NewAppService(scheduledTaskDomainService)
	scheduledTaskController := portal.NewScheduledTaskController(scheduledTaskAppService)

	scheduledTasks := rg.Group("/scheduled-tasks")
	{
		scheduledTasks.POST("", scheduledTaskController.CreateScheduledTask)
		scheduledTasks.GET("", scheduledTaskController.GetScheduledTasks)
		scheduledTasks.GET("/agent/:agentId", scheduledTaskController.GetScheduledTasksByAgent)
		scheduledTasks.GET("/:taskId", scheduledTaskController.GetScheduledTask)
		scheduledTasks.PUT("/:taskId", scheduledTaskController.UpdateScheduledTask)
		scheduledTasks.DELETE("/:taskId", scheduledTaskController.DeleteScheduledTask)
		scheduledTasks.POST("/:taskId/pause", scheduledTaskController.PauseTask)
		scheduledTasks.POST("/:taskId/resume", scheduledTaskController.ResumeTask)
	}

	// === ApiKey 模块（对应 Java 的 PortalApiKeyController）===
	apiKeyRepo := domainApiKey.NewApiKeyRepository(db)
	apiKeyDomainService := domainApiKey.NewDomainService(apiKeyRepo)
	apiKeyAppService := appApiKey.NewAppService(apiKeyDomainService, agentDomainService)
	apiKeyController := portal.NewApiKeyController(apiKeyAppService)

	apiKeys := rg.Group("/api-keys")
	{
		apiKeys.POST("", apiKeyController.CreateApiKey)
		apiKeys.GET("", apiKeyController.GetUserApiKeys)
		apiKeys.GET("/agent/:agentId", apiKeyController.GetAgentApiKeys)
		apiKeys.GET("/:apiKeyId", apiKeyController.GetApiKey)
		apiKeys.PUT("/:apiKeyId/status", apiKeyController.UpdateApiKeyStatus)
		apiKeys.DELETE("/:apiKeyId", apiKeyController.DeleteApiKey)
		apiKeys.POST("/:apiKeyId/reset", apiKeyController.ResetApiKey)
	}

	// === Memory 模块（对应 Java 的 PortalMemoryController）===
	memoryRepo := domainMemory.NewMemoryItemRepository(db)
	memoryDomainService := domainMemory.NewDomainService(memoryRepo)
	memoryAppService := appMemory.NewAppService(memoryDomainService)
	memoryController := portal.NewMemoryController(memoryAppService)

	memoryGroup := rg.Group("/portal/memory")
	{
		memoryGroup.GET("/items", memoryController.ListMemories)
		memoryGroup.POST("/items", memoryController.CreateMemory)
		memoryGroup.DELETE("/items/:itemId", memoryController.DeleteMemory)
	}

	// === 其他模块 ===
	// TODO: 后续迁移
}

// setupAdminRoutes 注册管理员路由
func setupAdminRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// === Admin User 管理（对应 Java 的 AdminUserController）===
	userRepo := domainUser.NewUserRepository(db)
	userSettingsRepo := domainUser.NewUserSettingsRepository(db)
	userDomainService := domainUser.NewDomainService(userRepo, userSettingsRepo)
	userAppService := appUser.NewAppService(userDomainService)
	adminUserController := admin.NewAdminUserController(userAppService)

	adminUsers := rg.Group("/users")
	{
		adminUsers.GET("", adminUserController.GetUsers)
	}

	// === Admin Auth Setting 管理（对应 Java 的 AdminAuthSettingController）===
	authSettingRepo := domainAuth.NewAuthSettingRepository(db)
	authDomainService := domainAuth.NewDomainService(authSettingRepo)
	authAppService := appAuth.NewAppService(authDomainService)
	adminAuthController := admin.NewAdminAuthSettingController(authAppService)

	authSettings := rg.Group("/auth-settings")
	{
		authSettings.GET("", adminAuthController.GetAllAuthSettings)
		authSettings.GET("/:id", adminAuthController.GetAuthSettingByID)
		authSettings.PUT("/:id/toggle", adminAuthController.ToggleAuthSetting)
		authSettings.PUT("/:id", adminAuthController.UpdateAuthSetting)
		authSettings.DELETE("/:id", adminAuthController.DeleteAuthSetting)
	}

	// === Admin LLM 管理（对应 Java 的 AdminLLMController）===
	adminProviderRepo := domainLLM.NewProviderRepository(db)
	adminModelRepo := domainLLM.NewModelRepository(db)
	adminLLMDomainService := domainLLM.NewDomainService(adminProviderRepo, adminModelRepo)
	adminLLMAppService := appLLM.NewAdminAppService(adminLLMDomainService)
	adminLLMController := admin.NewAdminLLMController(adminLLMAppService)

	adminLLMProviders := rg.Group("/llms/providers")
	{
		adminLLMProviders.GET("", adminLLMController.GetProviders)
		adminLLMProviders.GET("/protocols", adminLLMController.GetProviderProtocols)
		adminLLMProviders.GET("/:providerId", adminLLMController.GetProviderDetail)
		adminLLMProviders.POST("", adminLLMController.CreateProvider)
		adminLLMProviders.PUT("/:id", adminLLMController.UpdateProvider)
		adminLLMProviders.POST("/:id/status", adminLLMController.ToggleProviderStatus)
		adminLLMProviders.DELETE("/:id", adminLLMController.DeleteProvider)
	}

	adminLLMModels := rg.Group("/llms/models")
	{
		adminLLMModels.GET("", adminLLMController.GetModels)
		adminLLMModels.GET("/types", adminLLMController.GetModelTypes)
		adminLLMModels.POST("", adminLLMController.CreateModel)
		adminLLMModels.PUT("/:id", adminLLMController.UpdateModel)
		adminLLMModels.POST("/:id/status", adminLLMController.ToggleModelStatus)
		adminLLMModels.DELETE("/:id", adminLLMController.DeleteModel)
	}

	// === Admin Agent 管理（对应 Java 的 AdminAgentController）===
	adminAgentRepo := domainAgent.NewAgentRepository(db)
	adminAgentVersionRepo := domainAgent.NewAgentVersionRepository(db)
	adminAgentWorkspaceRepo := domainAgent.NewAgentWorkspaceRepository(db)
	adminAgentDomainService := domainAgent.NewDomainService(adminAgentRepo, adminAgentVersionRepo, adminAgentWorkspaceRepo)
	adminAgentWorkspaceDomainService := domainAgent.NewWorkspaceDomainService(adminAgentWorkspaceRepo, adminAgentRepo)
	adminAgentAppService := appAgent.NewAppService(adminAgentDomainService, adminAgentWorkspaceDomainService)
	adminAgentController := admin.NewAdminAgentController(adminAgentAppService)

	adminAgents := rg.Group("/agents")
	{
		adminAgents.GET("", adminAgentController.GetAgents)
		adminAgents.GET("/statistics", adminAgentController.GetAgentStatistics)
		adminAgents.GET("/versions", adminAgentController.GetVersions)
		adminAgents.POST("/versions/:versionId/status", adminAgentController.UpdateVersionStatus)
	}

	// === Admin Tool 管理（对应 Java 的 AdminToolController）===
	adminToolRepo := domainTool.NewToolRepository(db)
	adminToolVersionRepo := domainTool.NewToolVersionRepository(db)
	adminUserToolRepo := domainTool.NewUserToolRepository(db)
	adminToolDomainService := domainTool.NewToolDomainService(adminToolRepo, adminToolVersionRepo, adminUserToolRepo)
	adminUserToolDomainService := domainTool.NewUserToolDomainService(adminUserToolRepo)
	adminToolVersionDomainService := domainTool.NewToolVersionDomainService(adminToolVersionRepo)
	adminToolAppService := appTool.NewAppService(adminToolDomainService, adminUserToolDomainService, adminToolVersionDomainService, userDomainService)
	adminToolController := admin.NewAdminToolController(adminToolAppService, adminToolDomainService)

	adminTools := rg.Group("/tools")
	{
		adminTools.GET("", adminToolController.GetTools)
		adminTools.GET("/statistics", adminToolController.GetToolStatistics)
		adminTools.POST("/official", adminToolController.CreateOfficialTool)
		adminTools.POST("/:toolId/status", adminToolController.UpdateStatus)
		adminTools.PUT("/:toolId/global-status", adminToolController.UpdateGlobalStatus)
	}
}

// setupExternalRoutes 注册外部API路由（使用API Key认证）
func setupExternalRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// TODO: 后续迁移 external API 模块时注册
	// rg.POST("/chat/completions", externalController.ChatCompletions)
}
