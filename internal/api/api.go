package api

import (
	"log"
	"nomad-demo/internal/api/router"
	"nomad-demo/internal/pkg/config"

	"go.uber.org/zap"
)

func initLogger(logLevel string) {
	loggerConf := zap.NewProductionConfig()
	loggerConf.Level.UnmarshalText([]byte(logLevel))

	logger, err := loggerConf.Build()
	if err != nil {
		log.Fatalf("couldn't initialize zap logger: %v", err)
	}
	zap.ReplaceGlobals(logger)
}

func Run() {
	config.Setup()
	conf := config.GetConfig()

	initLogger(conf.Server.LogLevel)

	web := router.Setup()
	zap.L().Info("Running nomad-demo on port " + conf.Server.Port)
	if err := web.Run(":" + conf.Server.Port); err != nil {
		zap.S().Fatalf("webserver crashed with error: %v", err)
	}
}
