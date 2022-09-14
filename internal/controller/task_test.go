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
	"github.com/viniosilva/swordhealth-api/internal/exception"
	"github.com/viniosilva/swordhealth-api/internal/model"
	"github.com/viniosilva/swordhealth-api/mock"
)

func TestTaskControllerCreateTask(t *testing.T) {
	now := time.Now()
	task := &model.Task{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    1,
		Summary:   "summary",
		Status:    model.TaskStatusOpened,
	}

	var cases = map[string]struct {
		inputPayload       string
		mocking            func(taskService *mock.MockTaskService)
		expectedStatusCode int
		expectedBody       dto.TaskResponse
		expectedErrorBody  dto.ApiError
	}{
		"should create task": {
			inputPayload: `{
				"user_id": 1,
				"summary": "summary"
			}`,
			mocking: func(taskService *mock.MockTaskService) {
				taskService.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(task, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody: dto.TaskResponse{Data: dto.TaskDto{
				ID:        task.ID,
				CreatedAt: task.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: task.UpdatedAt.Format("2006-01-02 15:04:05"),
				User:      dto.UserDto{ID: task.UserID},
				Summary:   task.Summary,
				Status:    task.Status,
			}},
		},
		"should throw bad request when payload is invalid": {
			inputPayload: `{
				"user_id": "user_id",
				"summary": 123
			}`,
			mocking:            func(taskService *mock.MockTaskService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorBody:  dto.ApiError{Error: "invalid payload"},
		},
		"should throw bad request when payload data is invalid": {
			inputPayload: `{
				"user_id": -1,
				"summary": ""
			}`,
			mocking:            func(taskService *mock.MockTaskService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorBody: dto.ApiError{Error: strings.Join([]string{
				"Key: 'CreateTaskDto.UserID' Error:Field validation for 'UserID' failed on the 'min' tag",
				"Key: 'CreateTaskDto.Summary' Error:Field validation for 'Summary' failed on the 'required' tag",
			}, "; ")},
		},
		"should throw bad request when user not found": {
			inputPayload: `{
				"user_id": 1,
				"summary": "summary"
			}`,
			mocking: func(taskService *mock.MockTaskService) {
				taskService.EXPECT().CreateTask(gomock.Any(), gomock.Any()).
					Return(nil, &exception.ForeignKeyConstraintException{Message: "user not found"})
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorBody:  dto.ApiError{Error: "user not found"},
		},
		"should throw internal server error": {
			inputPayload: `{
				"user_id": 1,
				"summary": "summary"
			}`,
			mocking: func(taskService *mock.MockTaskService) {
				taskService.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
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
			ctx.Request = httptest.NewRequest("POST", "/api/tasks", strings.NewReader(cs.inputPayload))

			taskServiceMock := mock.NewMockTaskService(ctrl)
			taskController := controller.NewTaskController(r.Group("/api"), taskServiceMock)

			cs.mocking(taskServiceMock)

			// when
			taskController.CreateTask(ctx)

			var body dto.TaskResponse
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
