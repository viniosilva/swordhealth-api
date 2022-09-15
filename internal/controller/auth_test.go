package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/swordhealth-api/internal/controller"
	"github.com/viniosilva/swordhealth-api/internal/dto"
	"github.com/viniosilva/swordhealth-api/internal/exception"
	"github.com/viniosilva/swordhealth-api/internal/model"
	"github.com/viniosilva/swordhealth-api/mock"
)

func TestAuthControllerLogin(t *testing.T) {
	accessTokenMock := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
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
		inputBasicAuth     string
		mocking            func(authService *mock.MockAuthService, userService *mock.MockUserService, cryptoService *mock.MockCryptoService)
		expectedStatusCode int
		expectedBody       dto.AuthLoginResponse
		expectedErrorBody  dto.ApiError
	}{
		"should return access controll": {
			inputBasicAuth: "Basic dXNlcm5hbWU6MTEyMjMzNDQ1NQ==",
			mocking: func(authService *mock.MockAuthService, userService *mock.MockUserService, cryptoService *mock.MockCryptoService) {
				authService.EXPECT().DecodeBasicAuth(gomock.Any(), gomock.Any()).Return("username", "1122334455", nil)
				cryptoService.EXPECT().Hash(gomock.Any()).Return("aabbccddee")
				userService.EXPECT().GetUserByUsernameAndPassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(user, nil)
				cryptoService.EXPECT().EncryptJwt(gomock.Any(), gomock.Any(), gomock.Any()).Return(accessTokenMock, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       dto.AuthLoginResponse{AccessToken: accessTokenMock},
		},
		"should throw error on decode basic auth": {
			inputBasicAuth: "Basic ",
			mocking: func(authService *mock.MockAuthService, userService *mock.MockUserService, cryptoService *mock.MockCryptoService) {
				authService.EXPECT().DecodeBasicAuth(gomock.Any(), gomock.Any()).Return("", "", fmt.Errorf("invalid basic auth"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorBody: dto.ApiError{
				Error: "invalid basic auth",
			},
		},
		"should throw forbidden exception on get user by username and password": {
			inputBasicAuth: "Basic dXNlcm5hbWU6MTEyMjMzNDQ1NQ==",
			mocking: func(authService *mock.MockAuthService, userService *mock.MockUserService, cryptoService *mock.MockCryptoService) {
				authService.EXPECT().DecodeBasicAuth(gomock.Any(), gomock.Any()).Return("username", "1122334455", nil)
				cryptoService.EXPECT().Hash(gomock.Any()).Return("aabbccddee")
				userService.EXPECT().GetUserByUsernameAndPassword(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, &exception.NotFoundException{Message: "user not found"})
			},
			expectedStatusCode: http.StatusForbidden,
			expectedErrorBody:  dto.ApiError{Error: "user not found"},
		},
		"should throw error on get user by username and password": {
			inputBasicAuth: "Basic dXNlcm5hbWU6MTEyMjMzNDQ1NQ==",
			mocking: func(authService *mock.MockAuthService, userService *mock.MockUserService, cryptoService *mock.MockCryptoService) {
				authService.EXPECT().DecodeBasicAuth(gomock.Any(), gomock.Any()).Return("username", "1122334455", nil)
				cryptoService.EXPECT().Hash(gomock.Any()).Return("aabbccddee")
				userService.EXPECT().GetUserByUsernameAndPassword(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorBody:  dto.ApiError{Error: "internal server error"},
		},
		"should throw error on encrypt jwt": {
			inputBasicAuth: "Basic dXNlcm5hbWU6MTEyMjMzNDQ1NQ==",
			mocking: func(authService *mock.MockAuthService, userService *mock.MockUserService, cryptoService *mock.MockCryptoService) {
				authService.EXPECT().DecodeBasicAuth(gomock.Any(), gomock.Any()).Return("username", "1122334455", nil)
				cryptoService.EXPECT().Hash(gomock.Any()).Return("aabbccddee")
				userService.EXPECT().GetUserByUsernameAndPassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(user, nil)
				cryptoService.EXPECT().EncryptJwt(gomock.Any(), gomock.Any(), gomock.Any()).Return(accessTokenMock, fmt.Errorf("error"))
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

			ctx.Request = httptest.NewRequest("POST", "/api/auth/login", nil)
			ctx.Request.Header.Add("Authorization", cs.inputBasicAuth)

			authServiceMock := mock.NewMockAuthService(ctrl)
			userServiceMock := mock.NewMockUserService(ctrl)
			cryptoServiceMock := mock.NewMockCryptoService(ctrl)

			authController := controller.NewAuthController(r.Group("/api"),
				authServiceMock, userServiceMock, cryptoServiceMock)

			cs.mocking(authServiceMock, userServiceMock, cryptoServiceMock)

			// when
			authController.Login(ctx)

			var body dto.AuthLoginResponse
			json.Unmarshal(res.Body.Bytes(), &body)

			var errorBody dto.ApiError
			json.Unmarshal(res.Body.Bytes(), &errorBody)

			// then
			assert.Equal(t, cs.expectedStatusCode, res.Result().StatusCode)
			assert.Equal(t, cs.expectedBody, body)
			assert.Equal(t, cs.expectedErrorBody, errorBody)
		})
	}
}
