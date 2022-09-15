package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniosilva/swordhealth-api/internal/dto"
	"github.com/viniosilva/swordhealth-api/internal/exception"
	"github.com/viniosilva/swordhealth-api/internal/service"
)

type AuthController interface {
	Login(ctx *gin.Context)
}

type authController struct {
	authService   service.AuthService
	userService   service.UserService
	cryptoService service.CryptoService
}

func NewAuthController(router *gin.RouterGroup, authService service.AuthService,
	userService service.UserService, cryptoService service.CryptoService) AuthController {
	impl := &authController{
		authService:   authService,
		userService:   userService,
		cryptoService: cryptoService,
	}

	router.POST("/auth/login", impl.Login)

	return impl
}

// @Summary login
// @Schemes
// @Tags auth
// @Accept json
// @Produce json
// @Security BasicAuth
// @Success 200 {object} dto.AuthLoginResponse
// @Failure 400 {object} dto.ApiError
// @Failure 403 {object} dto.ApiError
// @Router /auth/login [post]
func (impl *authController) Login(ctx *gin.Context) {
	username, password, err := impl.authService.DecodeBasicAuth(ctx, ctx.Request.Header.Get("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiError{Error: err.Error()})
		return
	}

	encryptedPassword := impl.cryptoService.Hash(password)
	user, err := impl.userService.GetUserByUsernameAndPassword(ctx, username, encryptedPassword)
	if err != nil {
		if _, ok := err.(*exception.NotFoundException); ok {
			ctx.JSON(http.StatusForbidden, dto.ApiError{Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiError{Error: "internal server error"})
		return
	}

	accessToken, err := impl.cryptoService.EncryptJwt(ctx, user.ID, map[string]interface{}{
		"username": user.Username,
		"role":     user.Role,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiError{Error: "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, dto.AuthLoginResponse{AccessToken: accessToken})
}
