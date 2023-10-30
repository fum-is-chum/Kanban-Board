package controller

import (
	"fmt"
	"kanban-board/dto"
	responseHelper "kanban-board/helpers/response"
	m "kanban-board/middlewares"
	assigneeUsecase "kanban-board/usecase/task_assignee"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type taskAssigneeController struct {
	useCase assigneeUsecase.TaskAssigneeUseCase
}

func NewTaskAssigneeController(useCase assigneeUsecase.TaskAssigneeUseCase) *taskAssigneeController {
	return &taskAssigneeController{useCase}
}

func (t *taskAssigneeController) AddAssignee(c echo.Context) error {
	var request dto.TaskAssigneeRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	userID := m.ExtractTokenUserId(c)
	if err := t.useCase.AddAssignee(uint(userID), &request); err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessResponse("Success add assignee to task"))
}

func (t *taskAssigneeController) RemoveAssignee(c echo.Context) error {
	var request dto.TaskAssigneeRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	userID := m.ExtractTokenUserId(c)
	if err := t.useCase.DeleteAssignee(uint(userID), &request); err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessResponse("Success delete assignee from task"))
}

func (t *taskAssigneeController) ExitTask(c echo.Context) error {
	taskId, err := strconv.Atoi(c.QueryParam("taskId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	userID := m.ExtractTokenUserId(c)
	if err := t.useCase.ExitTask(uint(userID), uint(taskId)); err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessResponse("Success exit from task"))
}
