package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/viniosilva/swordhealth-api/internal/dto"
	"github.com/viniosilva/swordhealth-api/internal/service"
)

type MiddlewareController interface {
	AccessToken(ctx *gin.Context)
}

type middlewareController struct {
	cryptoService service.CryptoService
}

func NewMiddlewareController(cryptoService service.CryptoService) MiddlewareController {
	impl := &middlewareController{
		cryptoService: cryptoService,
	}

	return impl
}

func (impl *middlewareController) AccessToken(ctx *gin.Context) {
	splitedBearerMiddleware := strings.Split(ctx.Request.Header.Get("authorization"), " ")
	if strings.ToLower(splitedBearerMiddleware[0]) != "bearer" || len(splitedBearerMiddleware) != 2 {
		ctx.JSON(http.StatusBadRequest, dto.ApiError{Error: "invalid authorization"})
		ctx.Abort()
		return
	}

	accessToken := splitedBearerMiddleware[1]
	claims, err := impl.cryptoService.DecryptJwt(ctx, accessToken)
	if err != nil {
		ctx.JSON(http.StatusForbidden, dto.ApiError{Error: err.Error()})
		ctx.Abort()
		return
	}

	for k, v := range claims {
		ctx.Params = append(ctx.Params, gin.Param{Key: k, Value: fmt.Sprint(v)})
	}

	ctx.Next()
}
