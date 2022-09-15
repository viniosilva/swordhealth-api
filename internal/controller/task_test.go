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
		inputUserID        int
		inputPayload       string
		mocking            func(taskService *mock.MockTaskService, notificationService *mock.MockNotificationService, async chan bool)
		expectedStatusCode int
		expectedBody       dto.TaskResponse
		expectedErrorBody  dto.ApiError
	}{
		"should create task": {
			inputUserID: 1,
			inputPayload: `{
				"summary": "summary"
			}`,
			mocking: func(taskService *mock.MockTaskService, notificationService *mock.MockNotificationService, async chan bool) {
				taskService.EXPECT().CreateTask(gomock.Any(), gomock.Any(), gomock.Any()).Return(task, nil)
				notificationService.EXPECT().NotifyAdminUserOnSaveTask(gomock.Any(), gomock.Any(), gomock.Any()).
					Do(func(arg interface{}, arg2 interface{}, arg3 interface{}) {
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
			inputUserID: 1,
			inputPayload: `{
				"summary": 123
			}`,
			mocking: func(taskService *mock.MockTaskService, notificationService *mock.MockNotificationService, async chan bool) {
				async <- true
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorBody:  dto.ApiError{Error: "invalid payload"},
		},
		"should throw bad request when payload data is invalid": {
			inputUserID: 1,
			inputPayload: `{
				"summary": ""
			}`,
			mocking: func(taskService *mock.MockTaskService, notificationService *mock.MockNotificationService, async chan bool) {
				async <- true
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorBody:  dto.ApiError{Error: "Key: 'CreateTaskDto.Summary' Error:Field validation for 'Summary' failed on the 'required' tag"},
		},
		"should throw bad request when user not found": {
			inputUserID: 1,
			inputPayload: `{
				"summary": "summary"
			}`,
			mocking: func(taskService *mock.MockTaskService, notificationService *mock.MockNotificationService, async chan bool) {
				taskService.EXPECT().CreateTask(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, &exception.ForeignKeyConstraintException{Message: "user not found"})
				async <- true
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorBody:  dto.ApiError{Error: "user not found"},
		},
		"should throw internal server error": {
			inputUserID: 1,
			inputPayload: `{
				"summary": "summary"
			}`,
			mocking: func(taskService *mock.MockTaskService, notificationService *mock.MockNotificationService, async chan bool) {
				taskService.EXPECT().CreateTask(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
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
			ctx.Params = append(ctx.Params, gin.Param{Key: "sub", Value: fmt.Sprint(cs.inputUserID)})
			ctx.Request = httptest.NewRequest("POST", "/api/tasks", strings.NewReader(cs.inputPayload))

			taskServiceMock := mock.NewMockTaskService(ctrl)
			notificationServiceMock := mock.NewMockNotificationService(ctrl)
			taskController := controller.NewTaskController(r.Group("/api"), taskServiceMock, nil, notificationServiceMock, nil)

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
		inputUserID        string
		inputLimit         int
		inputOffset        int
		mocking            func(taskService *mock.MockTaskService, userService *mock.MockUserService)
		expectedStatusCode int
		expectedBody       dto.TasksResponse
		expectedErrorBody  dto.ApiError
	}{
		"should list tasks": {
			inputUserID: "1",
			mocking: func(taskService *mock.MockTaskService, userService *mock.MockUserService) {
				userService.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(&model.User{
						ID:   1,
						Role: model.UserRoleManager,
					}, nil)
				taskService.EXPECT().ListTasks(gomock.Any(), 10, 0, gomock.Any()).Return([]model.Task{task}, 1, nil)
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
			inputUserID: "1",
			inputLimit:  1,
			inputOffset: 2,
			mocking: func(taskService *mock.MockTaskService, userService *mock.MockUserService) {
				userService.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(&model.User{
						ID:   1,
						Role: model.UserRoleTechnician,
					}, nil)
				taskService.EXPECT().ListTasks(gomock.Any(), 1, 2, gomock.Any()).Return([]model.Task{task}, 10, nil)
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
		"should throw internal server error on get user by id": {
			inputUserID: "1",
			inputLimit:  10,
			inputOffset: 0,
			mocking: func(taskService *mock.MockTaskService, userService *mock.MockUserService) {
				userService.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorBody:  dto.ApiError{Error: "internal server error"},
		},
		"should throw internal server error on user not found": {
			inputUserID: "1",
			inputLimit:  10,
			inputOffset: 0,
			mocking: func(taskService *mock.MockTaskService, userService *mock.MockUserService) {
				userService.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(nil, &exception.NotFoundException{Message: "user not found"})
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorBody:  dto.ApiError{Error: "internal server error"},
		},
		"should throw internal server error on list tasks": {
			inputUserID: "1",
			inputLimit:  10,
			inputOffset: 0,
			mocking: func(taskService *mock.MockTaskService, userService *mock.MockUserService) {
				userService.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(&model.User{
						ID:   1,
						Role: model.UserRoleManager,
					}, nil)
				taskService.EXPECT().ListTasks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, 0, fmt.Errorf("error"))
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
			ctx.Params = append(ctx.Params, gin.Param{Key: "sub", Value: cs.inputUserID})

			taskServiceMock := mock.NewMockTaskService(ctrl)
			userServiceMock := mock.NewMockUserService(ctrl)
			taskController := controller.NewTaskController(r.Group("/api"), taskServiceMock, userServiceMock, nil, nil)

			cs.mocking(taskServiceMock, userServiceMock)

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
