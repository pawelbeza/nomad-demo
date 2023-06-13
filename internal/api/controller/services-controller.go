package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceInput struct {
	Script bool   `json:"script"`
	Url    string `json:"url" binding:"required"`
}

func CreateService(c *gin.Context) {
	var serviceInput *ServiceInput
	err := c.BindJSON(serviceInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Input parameters are invalid %v", err),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"url": "http://IP:777",
	})
}
