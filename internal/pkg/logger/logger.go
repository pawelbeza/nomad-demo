package logger

import (
	"log"

	"go.uber.org/zap"
)

func Setup(logLevel string) {
	loggerConf := zap.NewProductionConfig()
	loggerConf.Level.UnmarshalText([]byte(logLevel))

	logger, err := loggerConf.Build()
	if err != nil {
		log.Fatalf("couldn't initialize zap logger: %v", err)
	}

	zap.ReplaceGlobals(logger)
}
