package controllers

import (
	"net/http"

	"github.com/Mangrover007/banking-backend/cmd/entities"
	"github.com/Mangrover007/banking-backend/cmd/services"
	"github.com/gin-gonic/gin"
)

type VideoController interface {
	FindAll(ctx *gin.Context)
	Save(ctx *gin.Context)
}

type controller struct {
	videoService services.VideoService
}

func New(service services.VideoService) *controller {
	return &controller{
		videoService: service,
	}
}

func (c *controller) FindAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.videoService.FindAll())
}

func (c *controller) Save(ctx *gin.Context) {
	var video entities.Video
	err := ctx.ShouldBindJSON(&video)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	c.videoService.Save(video)
	ctx.JSON(http.StatusOK, video)
}
