package controller

import (
	"fmt"
	"net/http"
	"nomad-demo/internal/pkg/nomad"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ServiceInput struct {
	Script bool   `json:"script"`
	Url    string `json:"url" binding:"required"`
}

type ServiceController struct {
	Manager nomad.JobManager
}

func NewDefaultServiceController() (*ServiceController, error) {
	jobManager, err := nomad.NewDefaultJobManager()
	if err != nil {
		return nil, err
	}

	return NewServiceController(jobManager), nil
}

func NewServiceController(manager nomad.JobManager) *ServiceController {
	return &ServiceController{Manager: manager}
}

func (s *ServiceController) CreateService(c *gin.Context) {
	var serviceInput ServiceInput
	err := c.BindJSON(&serviceInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("input parameters are invalid %v", err),
		})
	}

	serviceName := c.Param("name")
	url, err := s.Manager.CreateJob(&nomad.JobParams{
		ServiceName: serviceName,
		Url:         serviceInput.Url,
		Script:      serviceInput.Script,
	})
	if err != nil {
		zap.S().Errorf("failed to run service %v: %v", serviceName, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to run service: %v, double check if your url is valid", serviceName),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}
