package router

import (
	"nomad-demo/internal/api/controller"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Setup() *gin.Engine {
	app := gin.Default()

	serviceController, err := controller.NewDefaultServiceController()
	if err != nil {
		zap.S().Fatalf("failed to create service controller: %v", err)
	}

	app.PUT("/services/:name", serviceController.CreateService)

	return app
}
