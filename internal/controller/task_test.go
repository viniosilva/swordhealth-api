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
		mocking            func(taskService *mock.MockTaskService, notificationService *mock.MockNotificationService, async chan bool)
		expectedStatusCode int
		expectedBody       dto.TaskResponse
		expectedErrorBody  dto.ApiError
	}{
		"should create task": {
			inputPayload: `{
				"user_id": 1,
				"summary": "summary"
			}`,
			mocking: func(taskService *mock.MockTaskService, notificationService *mock.MockNotificationService, async chan bool) {
				taskService.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(task, nil)
				notificationService.EXPECT().NotifyAdminUserOnSaveTask(gomock.Any(), gomock.Any()).
					Do(func(arg interface{}, arg2 interface{}) {
						async <- true
					})
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
			mocking: func(taskService *mock.MockTaskService, notificationService *mock.MockNotificationService, async chan bool) {
				async <- true
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorBody:  dto.ApiError{Error: "invalid payload"},
		},
		"should throw bad request when payload data is invalid": {
			inputPayload: `{
				"user_id": -1,
				"summary": ""
			}`,
			mocking: func(taskService *mock.MockTaskService, notificationService *mock.MockNotificationService, async chan bool) {
				async <- true
			},
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
			mocking: func(taskService *mock.MockTaskService, notificationService *mock.MockNotificationService, async chan bool) {
				taskService.EXPECT().CreateTask(gomock.Any(), gomock.Any()).
					Return(nil, &exception.ForeignKeyConstraintException{Message: "user not found"})
				async <- true
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorBody:  dto.ApiError{Error: "user not found"},
		},
		"should throw internal server error": {
			inputPayload: `{
				"user_id": 1,
				"summary": "summary"
			}`,
			mocking: func(taskService *mock.MockTaskService, notificationService *mock.MockNotificationService, async chan bool) {
				taskService.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
				async <- true
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
			notificationServiceMock := mock.NewMockNotificationService(ctrl)
			taskController := controller.NewTaskController(r.Group("/api"), taskServiceMock, notificationServiceMock)

			async := make(chan bool, 1)
			cs.mocking(taskServiceMock, notificationServiceMock, async)

			// when
			taskController.CreateTask(ctx)
			<-async

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

func TestTaskControllerListTasks(t *testing.T) {
	now := time.Now()
	task := model.Task{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    1,
		Summary:   "summary",
		Status:    model.TaskStatusOpened,
	}

	var cases = map[string]struct {
		inputLimit         int
		inputOffset        int
		mocking            func(taskService *mock.MockTaskService)
		expectedStatusCode int
		expectedBody       dto.TasksResponse
		expectedErrorBody  dto.ApiError
	}{
		"should list tasks": {
			mocking: func(taskService *mock.MockTaskService) {
				taskService.EXPECT().ListTasks(gomock.Any(), 10, 0).Return([]model.Task{task}, 1, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: dto.TasksResponse{
				Pagination: dto.Pagination{
					Count: 1,
					Total: 1,
				},
				Data: []dto.TaskDto{
					{
						ID:        task.ID,
						CreatedAt: task.CreatedAt.Format("2006-01-02 15:04:05"),
						UpdatedAt: task.UpdatedAt.Format("2006-01-02 15:04:05"),
						User:      dto.UserDto{ID: task.UserID},
						Summary:   task.Summary,
						Status:    task.Status,
					},
				}},
		},
		"should list tasks when page is 3": {
			inputLimit:  1,
			inputOffset: 2,
			mocking: func(taskService *mock.MockTaskService) {
				taskService.EXPECT().ListTasks(gomock.Any(), 1, 2).Return([]model.Task{task}, 10, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: dto.TasksResponse{
				Pagination: dto.Pagination{
					Count: 1,
					Total: 10,
				},
				Data: []dto.TaskDto{
					{
						ID:        task.ID,
						CreatedAt: task.CreatedAt.Format("2006-01-02 15:04:05"),
						UpdatedAt: task.UpdatedAt.Format("2006-01-02 15:04:05"),
						User:      dto.UserDto{ID: task.UserID},
						Summary:   task.Summary,
						Status:    task.Status,
					},
				}},
		},
		"should throw internal server error": {
			inputLimit:  10,
			inputOffset: 0,
			mocking: func(taskService *mock.MockTaskService) {
				taskService.EXPECT().ListTasks(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, 0, fmt.Errorf("error"))
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

			var query = []string{}
			if cs.inputLimit > 0 {
				query = append(query, fmt.Sprintf("limit=%d", cs.inputLimit))
			}
			if cs.inputOffset > 0 {
				query = append(query, fmt.Sprintf("offset=%d", cs.inputOffset))
			}

			url := strings.Join([]string{
				"/api/tasks",
				strings.Join(query, "&"),
			}, "?")

			ctx.Request = httptest.NewRequest("GET", url, nil)

			taskServiceMock := mock.NewMockTaskService(ctrl)
			taskController := controller.NewTaskController(r.Group("/api"), taskServiceMock, nil)

			cs.mocking(taskServiceMock)

			// when
			taskController.ListTasks(ctx)

			var body dto.TasksResponse
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
