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
	userService service.UserService
}

func NewUserController(router *gin.RouterGroup, userService service.UserService) UserController {
	impl := &userController{
		userService: userService,
	}

	router.POST("/users", impl.CreateUser)

	return impl
}

func (impl *userController) CreateUser(ctx *gin.Context) {
	var data dto.CreateUserDto

	impl.RegisterValidationUserEnum()

	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiError{Error: exception.FormatBindingErrors(err)})
		return
	}

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
