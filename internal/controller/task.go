package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniosilva/swordhealth-api/internal/dto"
	"github.com/viniosilva/swordhealth-api/internal/exception"
	"github.com/viniosilva/swordhealth-api/internal/model"
	"github.com/viniosilva/swordhealth-api/internal/service"
)

type TaskController interface {
	CreateTask(ctx *gin.Context)
}

type taskController struct {
	taskService service.TaskService
}

func NewTaskController(router *gin.RouterGroup, taskService service.TaskService) TaskController {
	impl := &taskController{
		taskService: taskService,
	}

	router.POST("/tasks", impl.CreateTask)

	return impl
}

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

	ctx.JSON(http.StatusCreated, dto.TaskResponse{Data: impl.ParseTaskDto(task)})
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
