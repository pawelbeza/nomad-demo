package router

import (
	"nomad-demo/internal/api/controller"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	app := gin.Default()

	app.PUT("/services/:name", controller.CreateService)

	return app
}
