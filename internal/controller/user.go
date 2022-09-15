package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/viniosilva/swordhealth-api/internal/dto"
	"github.com/viniosilva/swordhealth-api/internal/exception"
	"github.com/viniosilva/swordhealth-api/internal/model"
	"github.com/viniosilva/swordhealth-api/internal/service"
	"golang.org/x/exp/slices"
)

type UserController interface {
	CreateUser(ctx *gin.Context)
}

type userController struct {
	userService   service.UserService
	cryptoService service.CryptoService
}

func NewUserController(router *gin.RouterGroup, userService service.UserService, cryptoService service.CryptoService,
	middlewareAccessToken, middlewareUserManager func(ctx *gin.Context)) UserController {
	impl := &userController{
		userService:   userService,
		cryptoService: cryptoService,
	}

	router.POST("/users", middlewareAccessToken, middlewareUserManager, impl.CreateUser)

	return impl
}

// @Summary create user
// @Schemes
// @Tags user
// @Accept json
// @Produce json
// @Security JwtAuth
// @Param request body dto.CreateUserDto true "user"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} dto.ApiError
// @Failure 401 {object} dto.ApiError
// @Failure 403 {object} dto.ApiError
// @Failure 500 {object} dto.ApiError
// @Router /users [post]
func (impl *userController) CreateUser(ctx *gin.Context) {
	impl.RegisterValidationUserEnum()

	var data dto.CreateUserDto
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiError{Error: exception.FormatBindingErrors(err)})
		return
	}

	data.Password = impl.cryptoService.Hash(data.Password)
	user, err := impl.userService.CreateUser(ctx, data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiError{Error: "internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, dto.UserResponse{Data: impl.ParseUserDto(user)})
}

func (impl *userController) ParseUserDto(user *model.User) dto.UserDto {
	dto := dto.UserDto{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),

		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	return dto
}

func (impl *userController) RegisterValidationUserEnum() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("enum", func(fl validator.FieldLevel) bool {
			if v := fl.Field().Interface().(model.UserRole); v != "" {
				return slices.Contains(
					[]model.UserRole{model.UserRoleManager, model.UserRoleTechnician}, v,
				)
			}

			return true
		})
	}
}
