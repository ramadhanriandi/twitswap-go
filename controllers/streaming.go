package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"twitswap-go/request"
)

type StreamingController struct{}

var (
	errInvalidRequest       = errors.New("invalid request body")
	errFailedStartStreaming = errors.New("failed to start streaming")
	errFailedStopStreaming  = errors.New("failed to stop streaming")
)

func (s *StreamingController) StartStreaming(c *gin.Context) {
	var data request.StartStreaming

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": errInvalidRequest.Error(),
			"error":   err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Successfully start streaming"})
}

func (s *StreamingController) StopStreaming(c *gin.Context) {
	var data request.StopStreaming

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": errInvalidRequest.Error(),
			"error":   err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Successfully stop streaming"})
}
