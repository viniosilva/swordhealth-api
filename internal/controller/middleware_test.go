package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/swordhealth-api/internal/controller"
	"github.com/viniosilva/swordhealth-api/internal/dto"
	"github.com/viniosilva/swordhealth-api/internal/model"
	"github.com/viniosilva/swordhealth-api/mock"
)

func TestMiddlewareControllerAccessToken(t *testing.T) {
	var cases = map[string]struct {
		inputAuthorization  string
		mocking             func(cryptoService *mock.MockCryptoService)
		expectedUserIDParam string
		expectedStatusCode  int
		expectedErrorBody   dto.ApiError
	}{
		"should next": {
			inputAuthorization: "bearer 123.abc.x0z",
			mocking: func(cryptoService *mock.MockCryptoService) {
				cryptoService.EXPECT().DecryptJwt(gomock.Any(), gomock.Any()).
					Return(map[string]interface{}{
						"sub": 1,
					}, nil)
			},
			expectedUserIDParam: "1",
			expectedStatusCode:  http.StatusOK,
		},
		"should throw invalid authorization error": {
			inputAuthorization: "error",
			mocking:            func(cryptoService *mock.MockCryptoService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorBody:  dto.ApiError{Error: "invalid authorization"},
		},
		"should throw error on decrypt jwt": {
			inputAuthorization: "bearer 123.abc.x0z",
			mocking: func(cryptoService *mock.MockCryptoService) {
				cryptoService.EXPECT().DecryptJwt(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("error"))
			},
			expectedStatusCode: http.StatusForbidden,
			expectedErrorBody:  dto.ApiError{Error: "error"},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gin.SetMode(gin.TestMode)
			res := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(res)
			ctx.Request = httptest.NewRequest("GET", "/api/healthcheck", nil)
			ctx.Request.Header.Add("authorization", cs.inputAuthorization)

			cryptoServiceMock := mock.NewMockCryptoService(ctrl)
			middlewareController := controller.NewMiddlewareController(cryptoServiceMock)

			cs.mocking(cryptoServiceMock)

			// when
			middlewareController.AccessToken(ctx)
			userIDParam, _ := ctx.Params.Get("sub")

			var errorBody dto.ApiError
			json.Unmarshal(res.Body.Bytes(), &errorBody)

			// then
			assert.Equal(t, cs.expectedUserIDParam, userIDParam)
			assert.Equal(t, cs.expectedStatusCode, res.Result().StatusCode)
			assert.Equal(t, cs.expectedErrorBody, errorBody)
		})
	}
}

func TestMiddlewareControllerHealth(t *testing.T) {
	var cases = map[string]struct {
		inputRole          model.UserRole
		expectedStatusCode int
		expectedErrorBody  dto.ApiError
	}{
		"should next": {
			inputRole:          model.UserRoleManager,
			expectedStatusCode: http.StatusOK,
		},
		"should throw unauthorized exception when user role is not manager": {
			inputRole:          model.UserRoleTechnician,
			expectedStatusCode: http.StatusUnauthorized,
			expectedErrorBody:  dto.ApiError{Error: "unauthorized user role"},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gin.SetMode(gin.TestMode)
			res := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(res)
			ctx.Request = httptest.NewRequest("GET", "/api/healthcheck", nil)
			ctx.Params = append(ctx.Params, gin.Param{Key: "role", Value: string(cs.inputRole)})

			middlewareController := controller.NewMiddlewareController(nil)

			// when
			middlewareController.UserManager(ctx)

			var errorBody dto.ApiError
			json.Unmarshal(res.Body.Bytes(), &errorBody)

			// then
			assert.Equal(t, cs.expectedStatusCode, res.Result().StatusCode)
			assert.Equal(t, cs.expectedErrorBody, errorBody)
		})
	}
}
