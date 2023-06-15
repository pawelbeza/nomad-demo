package logger

import (
	"log"

	"go.uber.org/zap"
)

func Setup(logLevel string) {
	loggerConf := zap.NewProductionConfig()
	err := loggerConf.Level.UnmarshalText([]byte(logLevel))

	if err != nil {
		log.Fatalf("invalid logger level %v: %v", logLevel, err)
	}

	logger, err := loggerConf.Build()
	if err != nil {
		log.Fatalf("couldn't initialize zap logger: %v", err)
	}

	zap.ReplaceGlobals(logger)
}
