package router

import (
	"github.com/gin-gonic/gin"
	appAgent "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/agent"
	appAccount "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/account"
	appApiKey "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/apikey"
	appAuth "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/auth"
	appContainer "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/container"
	appConv "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/conversation"
	appLLM "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/llm"
	appMemory "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/memory"
	appScheduledTask "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/scheduledtask"
	appOrder "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/order"
	appPayment "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/payment"
	appProduct "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/product"
	appChat "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/chat"
	appRag "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/rag"
	appRule "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/rule"
	appTask "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/task"
	appTool "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/tool"
	appTrace "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/trace"
	appUsage "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/usage"
	appUser "github.com/lucky-aeon/agentx/agentx-backend-go/internal/application/user"
	infraPayment "github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/payment"
	infraDocker "github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/docker"
	infraScheduler "github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/scheduler"
	infraEmail "github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/email"
	domainAgent "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/agent"
	domainApiKey "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/apikey"
	domainAuth "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/auth"
	domainContainer "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/container"
	domainConv "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/conversation"
	domainLLM "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/llm"
	domainMemory "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/memory"
	domainOrder "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/order"
	domainProduct "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/product"
	domainRag "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/rag"
	domainRule "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/rule"
	domainScheduledTask "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/scheduledtask"
	domainTask "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/task"
	domainTool "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/tool"
	domainTrace "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/trace"
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
	setupPublicRoutes(api, db, jwtUtils, cfg)

	// ========== 需要认证的接口 ==========
	authenticated := api.Group("")
	authenticated.Use(middleware.AuthMiddleware(jwtUtils))
	setupAuthenticatedRoutes(authenticated, db, cfg)

	// ========== 管理员接口 ==========
	adminGroup := authenticated.Group("/admin")
	adminGroup.Use(middleware.AdminAuthMiddleware())
	setupAdminRoutes(adminGroup, db, cfg)

	// ========== 外部API接口（使用API Key认证）==========
	external := api.Group("/v1")
	setupExternalRoutes(external, db)

	return r
}

// setupPublicRoutes 注册公开路由（登录、注册等）
func setupPublicRoutes(rg *gin.RouterGroup, db *gorm.DB, jwtUtils *auth.JWTUtils, cfg *config.Config) {
	// === 初始化 User 模块依赖 ===
	userRepo := domainUser.NewUserRepository(db)
	userSettingsRepo := domainUser.NewUserSettingsRepository(db)
	userDomainService := domainUser.NewDomainService(userRepo, userSettingsRepo)

	// 初始化邮件和验证码服务
	pubLogger, _ := zap.NewProduction()
	emailService := infraEmail.NewSMTPEmailService(&cfg.Mail, pubLogger)
	verificationService := infraEmail.NewVerificationCodeService(pubLogger)
	captchaService := infraEmail.NewCaptchaService()

	loginAppService := appUser.NewLoginAppService(userDomainService, jwtUtils, emailService, verificationService, captchaService)

	// === 登录注册路由（对应 Java 的 LoginController）===
	loginController := portal.NewLoginController(loginAppService, captchaService)
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
func setupAuthenticatedRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
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
	messageDomainService := domainConv.NewMessageDomainService(messageRepo, contextRepo)
	contextDomainService := domainConv.NewContextDomainService(contextRepo)

	// 创建 Chat 应用服务（供 SessionController 和 ChatController 使用）
	logger, _ := zap.NewProduction()
	chatAppService := appChat.NewChatAppService(
		conversationDomainService,
		sessionDomainService,
		agentDomainService,
		agentWorkspaceDomainService,
		llmDomainService,
		contextDomainService,
		messageDomainService,
		logger,
	)

	agentSessionAppService := appConv.NewAgentSessionAppService(
		agentWorkspaceDomainService, agentDomainService,
		sessionDomainService, conversationDomainService,
	)
	conversationAppService := appConv.NewConversationAppService(conversationDomainService, sessionDomainService)
	sessionController := portal.NewSessionController(agentSessionAppService, conversationAppService, chatAppService, logger)

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
	taskScheduleService := infraScheduler.NewTaskScheduleService()
	scheduledTaskAppService := appScheduledTask.NewAppService(scheduledTaskDomainService, taskScheduleService)
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

	// === Task 模块（对应 Java 的 TaskController）===
	taskRepo := domainTask.NewTaskRepository(db)
	taskDomainService := domainTask.NewDomainService(taskRepo)
	taskAppService := appTask.NewAppService(taskDomainService)
	taskController := portal.NewTaskController(taskAppService)

	tasks := rg.Group("/tasks")
	{
		tasks.GET("/session/:sessionId/latest", taskController.GetSessionTasks)
	}

	// === Trace 模块（对应 Java 的 PortalTraceController）===
	summaryRepo := domainTrace.NewExecutionSummaryRepository(db)
	detailRepo := domainTrace.NewExecutionDetailRepository(db)
	traceDomainService := domainTrace.NewDomainService(summaryRepo, detailRepo)
	traceAppService := appTrace.NewAppService(traceDomainService)
	traceController := portal.NewTraceController(traceAppService)

	traces := rg.Group("/traces")
	{
		traces.GET("/history", traceController.GetExecutionHistory)
		traces.GET("/statistics", traceController.GetUserExecutionStatistics)
		traces.GET("/sessions/:sessionId", traceController.GetSessionExecutionHistory)
		traces.GET("/:traceId", traceController.GetTraceDetail)
		traces.GET("/:traceId/details", traceController.GetExecutionDetails)
	}

	// === Order 模块（对应 Java 的 OrderController）===
	orderRepo := domainOrder.NewOrderRepository(db)
	orderDomainService := domainOrder.NewDomainService(orderRepo)
	orderAppService := appOrder.NewAppService(orderDomainService)
	orderController := portal.NewOrderController(orderAppService)

	orders := rg.Group("/orders")
	{
		orders.GET("", orderController.GetUserOrders)
		orders.GET("/:orderId", orderController.GetOrderDetail)
	}

	// === Product 模块（对应 Java 的 PortalProductController）===
	productRepo := domainProduct.NewProductRepository(db)
	productDomainService := domainProduct.NewDomainService(productRepo)
	productAppService := appProduct.NewAppService(productDomainService)
	productController := portal.NewProductController(productAppService)

	products := rg.Group("/products")
	{
		products.GET("/:productId", productController.GetProductByID)
		products.GET("/business", productController.GetProductByBusinessKey)
		products.GET("/active", productController.GetActiveProducts)
		products.GET("/business/active", productController.IsProductActive)
	}

	// === Account 模块（对应 Java 的 AccountController）===
	accountRepo := domainUser.NewAccountRepository(db)
	accountDomainService := domainUser.NewAccountDomainService(accountRepo)
	accountAppService := appAccount.NewAppService(accountDomainService)
	accountController := portal.NewAccountController(accountAppService)

	accounts := rg.Group("/accounts")
	{
		accounts.GET("/current", accountController.GetCurrentUserAccount)
	}

	// === Payment 模块（对应 Java 的 PaymentController）===
	paymentLogger, _ := zap.NewProduction()
	paymentFactory := infraPayment.NewPaymentProviderFactory(&cfg.Payment, paymentLogger)
	paymentAppService := appPayment.NewAppService(orderDomainService, accountDomainService, paymentFactory)
	paymentController := portal.NewPaymentController(paymentAppService)

	payments := rg.Group("/payments")
	{
		payments.POST("/recharge", paymentController.CreateRechargePayment)
		payments.GET("/orders/:orderNo/status", paymentController.QueryOrderStatus)
		payments.GET("/methods", paymentController.GetAvailablePaymentMethods)
	}

	// === Usage 模块（对应 Java 的 PortalUsageRecordController）===
	usageRecordRepo := domainUser.NewUsageRecordRepository(db)
	usageRecordDomainService := domainUser.NewUsageRecordDomainService(usageRecordRepo)
	usageAppService := appUsage.NewAppService(usageRecordDomainService)
	usageRecordController := portal.NewUsageRecordController(usageAppService)

	usageRecords := rg.Group("/usage-records")
	{
		usageRecords.GET("/:recordId", usageRecordController.GetUsageRecordByID)
		usageRecords.GET("", usageRecordController.QueryUsageRecords)
		usageRecords.GET("/current/total-cost", usageRecordController.GetCurrentUserTotalCost)
	}

	// === Upload 模块（对应 Java 的 UploadController）===
	uploadController := portal.NewUploadController(&cfg.Upload)

	upload := rg.Group("/upload")
	{
		upload.GET("/credential", uploadController.GetUploadCredential)
	}

	// === 其他模块 ===

	// === RAG 模块（对应 Java 的 RAG 相关 Controller）===
	userRagRepo := domainRag.NewUserRagRepository(db)
	ragVersionRepo := domainRag.NewRagVersionRepository(db)
	fileDetailRepo := domainRag.NewFileDetailRepository(db)
	docUnitRepo := domainRag.NewDocumentUnitRepository(db)
	qaRepo := domainRag.NewRagQaDatasetRepository(db)

	fileOpAppService := appRag.NewFileOperationAppService(userRagRepo, fileDetailRepo, docUnitRepo)
	fileOpController := portal.NewFileOperationController(fileOpAppService)

	ragMarketAppService := appRag.NewRagMarketAppService(ragVersionRepo, userRagRepo)
	ragMarketController := portal.NewRagMarketController(ragMarketAppService)

	ragPublishAppService := appRag.NewRagPublishAppService(ragVersionRepo, userRagRepo)
	ragPublishController := portal.NewRagPublishController(ragPublishAppService)

	ragSearchAppService := appRag.NewRAGSearchAppService(userRagRepo)
	ragSearchController := portal.NewRagSearchController(ragSearchAppService)

	ragQaAppService := appRag.NewRagQaDatasetAppService(qaRepo, userRagRepo)
	ragQaController := portal.NewRagQaDatasetController(ragQaAppService)

	ragFiles := rg.Group("/rag/files")
	{
		ragFiles.GET("/dataset/:dataSetId", fileOpController.GetFilesByDataSetID)
		ragFiles.GET("/:fileId/document-units", fileOpController.GetDocumentUnits)
		ragFiles.DELETE("/:fileId", fileOpController.DeleteFile)
	}

	ragMarket := rg.Group("/rag/market")
	{
		ragMarket.GET("", ragMarketController.GetMarketList)
		ragMarket.POST("/install", ragMarketController.InstallRag)
		ragMarket.DELETE("/:userRagId", ragMarketController.UninstallRag)
		ragMarket.GET("/my", ragMarketController.GetUserRags)
	}

	ragPublish := rg.Group("/rag/publish")
	{
		ragPublish.POST("", ragPublishController.PublishRag)
		ragPublish.GET("/:originalRagId/versions", ragPublishController.GetVersions)
	}

	ragSearch := rg.Group("/rag/search")
	{
		ragSearch.POST("", ragSearchController.RagSearch)
		ragSearch.POST("/user-rag/:userRagId", ragSearchController.RagSearchByUserRag)
	}

	ragQa := rg.Group("/rag/qa")
	{
		ragQa.GET("/:userRagId", ragQaController.GetQaDatasets)
		ragQa.POST("/:userRagId", ragQaController.CreateQaDataset)
	}

	// === Chat 模块（对应 Java 的 ConversationAppService Chat 功能）===
	chatController := portal.NewChatController(chatAppService, logger)

	chat := rg.Group("/chat")
	{
		chat.POST("/stream", chatController.StreamChat)
		chat.POST("", chatController.Chat)
		chat.POST("/stop", chatController.StopChat)
	}
}

// setupAdminRoutes 注册管理员路由
func setupAdminRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
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

	// === Admin Container 管理（对应 Java 的 AdminContainerController）===
	containerRepo := domainContainer.NewContainerRepository(db)
	containerTemplateRepo := domainContainer.NewContainerTemplateRepository(db)
	containerDomainService := domainContainer.NewDomainService(containerRepo, containerTemplateRepo)

	// 初始化Docker服务
	adminLogger, _ := zap.NewProduction()
	dockerService := infraDocker.NewDockerEngineService(&cfg.Docker, adminLogger)
	containerLifecycleService := infraDocker.NewContainerLifecycleService(containerDomainService, dockerService, adminLogger)

	containerAppService := appContainer.NewAppService(containerDomainService, containerLifecycleService)
	containerTemplateAppService := appContainer.NewTemplateAppService(containerDomainService)
	adminContainerController := admin.NewAdminContainerController(containerAppService)
	adminContainerTemplateController := admin.NewAdminContainerTemplateController(containerTemplateAppService)

	adminContainers := rg.Group("/containers")
	{
		adminContainers.GET("", adminContainerController.GetContainersPage)
		adminContainers.GET("/:containerId", adminContainerController.GetContainerByID)
		adminContainers.POST("/:containerId/start", adminContainerController.StartContainer)
		adminContainers.POST("/:containerId/stop", adminContainerController.StopContainer)
		adminContainers.DELETE("/:containerId", adminContainerController.DeleteContainer)
		adminContainers.GET("/:containerId/logs", adminContainerController.GetContainerLogs)
	}

	adminContainerTemplates := rg.Group("/container-templates")
	{
		adminContainerTemplates.GET("", adminContainerTemplateController.GetEnabledTemplates)
		adminContainerTemplates.GET("/:templateId", adminContainerTemplateController.GetTemplate)
		adminContainerTemplates.DELETE("/:templateId", adminContainerTemplateController.DeleteTemplate)
	}

	// === Admin Product 管理（对应 Java 的 AdminProductController）===
	adminProductRepo := domainProduct.NewProductRepository(db)
	adminProductDomainService := domainProduct.NewDomainService(adminProductRepo)
	adminProductAppService := appProduct.NewAppService(adminProductDomainService)
	adminProductController := admin.NewAdminProductController(adminProductAppService)

	adminProducts := rg.Group("/products")
	{
		adminProducts.GET("", adminProductController.GetProducts)
		adminProducts.GET("/all", adminProductController.GetAllProducts)
		adminProducts.GET("/:productId", adminProductController.GetProductByID)
		adminProducts.DELETE("/:productId", adminProductController.DeleteProduct)
	}

	// === Admin Order 管理（对应 Java 的 AdminOrderController）===
	adminOrderRepo := domainOrder.NewOrderRepository(db)
	adminOrderDomainService := domainOrder.NewDomainService(adminOrderRepo)
	adminOrderAppService := appOrder.NewAppService(adminOrderDomainService)
	adminOrderController := admin.NewAdminOrderController(adminOrderAppService)

	adminOrders := rg.Group("/orders")
	{
		adminOrders.GET("", adminOrderController.GetAllOrders)
		adminOrders.GET("/:orderId", adminOrderController.GetOrderDetail)
	}

	// === Admin Rule 管理（对应 Java 的 AdminRuleController）===
	adminRuleRepo := domainRule.NewRuleRepository(db)
	adminRuleDomainService := domainRule.NewDomainService(adminRuleRepo)
	adminRuleAppService := appRule.NewAppService(adminRuleDomainService)
	adminRuleController := admin.NewAdminRuleController(adminRuleAppService)

	adminRules := rg.Group("/rules")
	{
		adminRules.POST("", adminRuleController.CreateRule)
		adminRules.PUT("/:ruleId", adminRuleController.UpdateRule)
		adminRules.GET("/:ruleId", adminRuleController.GetRuleByID)
		adminRules.GET("", adminRuleController.GetRules)
		adminRules.GET("/all", adminRuleController.GetAllRules)
		adminRules.DELETE("/:ruleId", adminRuleController.DeleteRule)
	}
}

// setupExternalRoutes 注册外部API路由（使用API Key认证）
func setupExternalRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// TODO: 后续迁移 external API 模块时注册
	// rg.POST("/chat/completions", externalController.ChatCompletions)
}
