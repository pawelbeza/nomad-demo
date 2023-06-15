package api

import (
	"nomad-demo/internal/api/router"
	"nomad-demo/internal/pkg/config"
	"nomad-demo/internal/pkg/logger"

	"go.uber.org/zap"
)

func Run() {
	config.Setup()
	conf := config.GetConfig()

	logger.Setup(conf.Server.LogLevel)

	web := router.Setup()
	zap.L().Info("Running nomad-demo on port " + conf.Server.Port)
	if err := web.Run(":" + conf.Server.Port); err != nil {
		zap.S().Fatalf("webserver crashed with error: %v", err)
	}
}
