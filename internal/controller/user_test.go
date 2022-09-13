package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/swordhealth-api/internal/controller"
	"github.com/viniosilva/swordhealth-api/internal/dto"
	"github.com/viniosilva/swordhealth-api/internal/model"
	"github.com/viniosilva/swordhealth-api/mock"
)

func TestUserControllerCreateUser(t *testing.T) {
	now := time.Now()
	user := &model.User{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		Username:  "username",
		Email:     "email@email.com",
		Password:  "1122334455",
		Role:      model.UserRoleTechnician,
	}

	var cases = map[string]struct {
		inputPayload       string
		mocking            func(userService *mock.MockUserService)
		expectedStatusCode int
		expectedBody       dto.UserResponse
		expectedErrorBody  dto.ApiError
	}{
		"should create user": {
			inputPayload: `{
				"username": "username",
				"email": "email@email.com",
				"password": "1122334455"
			}`,
			mocking: func(userService *mock.MockUserService) {
				userService.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(user, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody: dto.UserResponse{Data: dto.UserDto{
				ID:        user.ID,
				CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
				Username:  user.Username,
				Email:     user.Email,
				Role:      user.Role,
			}},
		},
		"should throw bad request when payload is invalid": {
			inputPayload: `{
				"username": 123,
				"email": "email",
				"password": "a",
				"role": "unknown"
			}`,
			mocking:            func(userService *mock.MockUserService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorBody:  dto.ApiError{Error: "invalid payload"},
		},
		"should throw bad request when payload data is invalid": {
			inputPayload: `{
				"username": "1",
				"email": "email",
				"password": "a",
				"role": "unknown"
			}`,
			mocking:            func(userService *mock.MockUserService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorBody: dto.ApiError{Error: strings.Join([]string{
				"Key: 'CreateUserDto.Username' Error:Field validation for 'Username' failed on the 'min' tag",
				"Key: 'CreateUserDto.Email' Error:Field validation for 'Email' failed on the 'email' tag",
				"Key: 'CreateUserDto.Password' Error:Field validation for 'Password' failed on the 'min' tag",
				"Key: 'CreateUserDto.Role' Error:Field validation for 'Role' failed on the 'enum' tag",
			}, "; ")},
		},
		"should throw internal server error": {
			inputPayload: `{
				"username": "username",
				"email": "email@email.com",
				"password": "1122334455",
				"role": "technician"
			}`,
			mocking: func(userService *mock.MockUserService) {
				userService.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorBody:  dto.ApiError{Error: "internal server error"},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gin.SetMode(gin.TestMode)
			res := httptest.NewRecorder()
			ctx, r := gin.CreateTestContext(res)
			ctx.Request = httptest.NewRequest("POST", "/api/users", strings.NewReader(cs.inputPayload))

			userServiceMock := mock.NewMockUserService(ctrl)
			userController := controller.NewUserController(r.Group("/api"), userServiceMock)

			cs.mocking(userServiceMock)

			// when
			userController.CreateUser(ctx)

			var body dto.UserResponse
			json.Unmarshal(res.Body.Bytes(), &body)

			var errorBody dto.ApiError
			json.Unmarshal(res.Body.Bytes(), &errorBody)

			// then
			assert.Equal(t, res.Result().StatusCode, cs.expectedStatusCode)
			assert.Equal(t, body, cs.expectedBody)
			assert.Equal(t, errorBody, cs.expectedErrorBody)
		})
	}
}
