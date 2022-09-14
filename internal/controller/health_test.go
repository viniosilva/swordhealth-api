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
	"github.com/viniosilva/swordhealth-api/mock"
)

func TestHealthControllerHealth(t *testing.T) {
	var cases = map[string]struct {
		mocking            func(healthService *mock.MockHealthService)
		expectedStatusCode int
		expectedBody       dto.HealthResponse
	}{
		"should return status up": {
			mocking: func(healthService *mock.MockHealthService) {
				healthService.EXPECT().Health(gomock.Any()).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       dto.HealthResponse{Status: dto.HealshStatusUp},
		},
		"should return status down": {
			mocking: func(healthService *mock.MockHealthService) {
				healthService.EXPECT().Health(gomock.Any()).Return(fmt.Errorf("error"))
			},
			expectedStatusCode: http.StatusServiceUnavailable,
			expectedBody:       dto.HealthResponse{Status: dto.HealshStatusDown},
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
			ctx.Request = httptest.NewRequest("GET", "/api/healthcheck", nil)

			healthServiceMock := mock.NewMockHealthService(ctrl)
			healthController := controller.NewHealthController(r.Group("/api"), healthServiceMock)

			cs.mocking(healthServiceMock)

			// when
			healthController.Health(ctx)

			var body dto.HealthResponse
			json.Unmarshal(res.Body.Bytes(), &body)

			// then
			assert.Equal(t, cs.expectedStatusCode, res.Result().StatusCode)
			assert.Equal(t, cs.expectedBody, body)
		})
	}
}
