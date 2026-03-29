package main

import (
	"fmt"
	"os"

	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/config"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/database"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/logger"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/router"
	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 2. 初始化日志
	logger.Init(cfg.Logging.Level)
	defer logger.Sync()

	// 3. 初始化数据库
	db, err := database.Init(cfg)
	if err != nil {
		zap.L().Fatal("初始化数据库失败", zap.Error(err))
	}

	// 4. 初始化路由
	r := router.Setup(db, cfg)

	// 5. 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	zap.L().Info("AgentX Go 后端服务启动", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		zap.L().Fatal("服务启动失败", zap.Error(err))
	}
}
