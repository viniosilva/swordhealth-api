package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/viniosilva/swordhealth-api/internal/dto"
	"github.com/viniosilva/swordhealth-api/internal/exception"
	"github.com/viniosilva/swordhealth-api/internal/model"
	"github.com/viniosilva/swordhealth-api/internal/service"
)

type TaskController interface {
	CreateTask(ctx *gin.Context)
	ListTasks(ctx *gin.Context)
}

type taskController struct {
	taskService         service.TaskService
	notificationService service.NotificationService
}

func NewTaskController(router *gin.RouterGroup, taskService service.TaskService, notificationService service.NotificationService) TaskController {
	impl := &taskController{
		taskService:         taskService,
		notificationService: notificationService,
	}

	router.POST("/tasks", impl.CreateTask)
	router.GET("/tasks", impl.ListTasks)

	return impl
}

// @Summary create task
// @Schemes
// @Tags task
// @Accept json
// @Produce json
// @Param request body dto.CreateTaskDto true "task"
// @Success 201 {object} dto.TaskResponse
// @Router /tasks [post]
func (impl *taskController) CreateTask(ctx *gin.Context) {
	var data dto.CreateTaskDto
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ApiError{Error: exception.FormatBindingErrors(err)})
		return
	}

	task, err := impl.taskService.CreateTask(ctx, data)
	if err != nil {
		if _, ok := err.(*exception.ForeignKeyConstraintException); ok {
			ctx.JSON(http.StatusBadRequest, dto.ApiError{Error: err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ApiError{Error: "internal server error"})
		return
	}

	go impl.notificationService.NotifyAdminUserOnSaveTask(ctx, task, task.UserID)

	ctx.JSON(http.StatusCreated, dto.TaskResponse{Data: impl.ParseTaskDto(task)})
}

// @Summary list tasks
// @Schemes
// @Tags task
// @Accept json
// @Produce json
// @Param request body dto.CreateTaskDto true "task"
// @Success 200 {array} []dto.TasksResponse
// @Router /tasks [get]
func (impl *taskController) ListTasks(ctx *gin.Context) {
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		offset = 0
	}

	tasks, total, err := impl.taskService.ListTasks(ctx, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiError{Error: "internal server error"})
		return
	}

	data := []dto.TaskDto{}
	for _, t := range tasks {
		data = append(data, impl.ParseTaskDto(&t))
	}

	ctx.JSON(http.StatusOK, dto.TasksResponse{
		Pagination: dto.Pagination{
			Count: len(data),
			Total: total,
		},
		Data: data,
	})
}

func (impl *taskController) ParseTaskDto(task *model.Task) dto.TaskDto {
	dto := dto.TaskDto{
		ID:        task.ID,
		CreatedAt: task.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: task.UpdatedAt.Format("2006-01-02 15:04:05"),
		User:      dto.UserDto{ID: task.UserID},
		Summary:   task.Summary,
		Status:    task.Status,
	}

	return dto
}
