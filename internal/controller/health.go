package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniosilva/swordhealth-api/internal/dto"
	"github.com/viniosilva/swordhealth-api/internal/service"
)

type HealthController interface {
	Health(ctx *gin.Context)
}

type healthController struct {
	healthService service.HealthService
}

func NewHealthController(router *gin.RouterGroup, healthService service.HealthService) HealthController {
	impl := &healthController{
		healthService: healthService,
	}

	router.GET("/healthcheck", impl.Health)

	return impl
}

func (impl *healthController) Health(ctx *gin.Context) {
	err := impl.healthService.Health(ctx)
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, dto.HealthResponse{Status: dto.HealshStatusDown})
		return
	}

	ctx.JSON(http.StatusOK, dto.HealthResponse{Status: dto.HealshStatusUp})
}
