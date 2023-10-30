package controller

import (
	"fmt"
	"kanban-board/dto"
	responseHelper "kanban-board/helpers/response"
	m "kanban-board/middlewares"
	useCase "kanban-board/usecase/task"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type taskController struct {
	useCase useCase.TaskUseCase
}

func NewTaskController(useCase useCase.TaskUseCase) *taskController {
	return &taskController{useCase}
}

func (t *taskController) GetTasks(c echo.Context) error {
	boardId, err := strconv.Atoi(c.QueryParam("boardId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse("Bad Request: ID invalid"))
	}

	userID := m.ExtractTokenUserId(c)
	tasks, err := t.useCase.GetTasks(uint(boardId), uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	var response []*dto.TaskResponse
	for _, task := range tasks {
		response = append(response, &dto.TaskResponse{
			ID:            task.ID,
			Title:         task.Title,
			Desc:          task.Desc,
			BoardColumnID: task.BoardColumnID,
		})
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessWithDataResponse(fmt.Sprintf("Success Get Tasks with board id %d", boardId), response))
}

func (t *taskController) GetTaskById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	userID := m.ExtractTokenUserId(c)
	task, err := t.useCase.GetTaskById(uint(id), uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	var assignees []*dto.TaskAssigneeResponse
	for _, assignee := range task.Assignees {
		assignees = append(assignees, &dto.TaskAssigneeResponse{
			ID: assignee.ID,
			Name: assignee.Name,
			Email: assignee.Email,
		})
	}
	

	return c.JSON(http.StatusOK, responseHelper.SuccessWithDataResponse(fmt.Sprintf("Success Get Task with id %d", id), &dto.TaskResponse{
		ID:            task.ID,
		Title:         task.Title,
		Desc:          task.Desc,
		BoardColumnID: task.BoardColumnID,
		Assignees: assignees,
	}))
}

func (t *taskController) CreateTask(c echo.Context) error {
	var createTaskRequest dto.TaskCreateRequest
	if err := c.Bind(&createTaskRequest); err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	userID := m.ExtractTokenUserId(c)
	if err := t.useCase.CreateTask(uint(userID), &createTaskRequest); err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessResponse("Success create task"))
}

func (t *taskController) UpdateTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	var updateTaskRequest dto.TaskUpdateRequest
	if err := c.Bind(&updateTaskRequest); err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	userID := m.ExtractTokenUserId(c)
	if err := t.useCase.UpdateTask(uint(id), uint(userID), &updateTaskRequest); err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessResponse("Success update task"))
}

func (t *taskController) DeleteTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	userID := m.ExtractTokenUserId(c)
	if err := t.useCase.DeleteTask(uint(id), uint(userID)); err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessResponse(fmt.Sprintf("Success delete task with id %d", id)))
}
