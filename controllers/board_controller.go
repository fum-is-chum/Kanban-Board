package controller

import (
	"fmt"
	"kanban-board/dto"
	responseHelper "kanban-board/helpers/response"
	m "kanban-board/middlewares"
	boardUsecase "kanban-board/usecase/board"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type boardController struct {
	useCase boardUsecase.BoardUseCase
}

func NewBoardController(boardUsecase boardUsecase.BoardUseCase) *boardController {
	return &boardController{boardUsecase}
}

func (b *boardController) GetBoards(c echo.Context) error {
	userID := m.ExtractTokenUserId(c)
	boards, err := b.useCase.GetBoards(uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	var response []dto.BoardResponse
	for _, board := range boards {
		response = append(response, dto.BoardResponse{
			ID:   board.ID,
			Name: board.Name,
			Desc: board.Desc,
			Owner: &dto.BoardMemberResponse{
				ID:    board.Owner.ID,
				Name:  board.Owner.Name,
				Email: board.Owner.Email,
			},
		})
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessWithDataResponse("Success get boards", response))
}

func (b *boardController) GetBoardById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse("Bad Request: Invalid Id"))
	}

	userID := m.ExtractTokenUserId(c)
	board, err := b.useCase.GetBoardById(uint(id), uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	var members []*dto.BoardMemberResponse
	for _, member := range board.Members {
		members = append(members, &dto.BoardMemberResponse{
			ID:        member.ID,
			Name:      member.Name,
			Email:     member.Email,
			CreatedAt: member.CreatedAt,
			UpdatedAt: member.UpdatedAt,
		})
	}

	var columns []*dto.BoardColumnResponse
	for _, column := range board.Columns {
		var tasks []*dto.TaskResponse // Create a slice to store tasks for this column
		for _, task := range column.Tasks {
			var assignees []*dto.TaskAssigneeResponse
			for _, assignee := range task.Assignees {
				assignees = append(assignees, &dto.TaskAssigneeResponse{
					ID: assignee.ID,
					Name: assignee.Name,
					Email: assignee.Email,
				})
			}

			tasks = append(tasks, &dto.TaskResponse{
				ID:    task.ID,
				Title: task.Title,
				Desc:  task.Desc,
				Assignees: assignees,
			})
		}

		columns = append(columns, &dto.BoardColumnResponse{
			ID:    column.ID,
			Label: column.Label,
			Desc:  column.Desc,
			Tasks: tasks, // Include the tasks in the column response
		})
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessWithDataResponse(fmt.Sprintf("Success get board with id %d", id), dto.BoardResponse{
		ID:   board.ID,
		Name: board.Name,
		Desc: board.Desc,
		Owner: &dto.BoardMemberResponse{
			ID:    board.Owner.ID,
			Name:  board.Owner.Name,
			Email: board.Owner.Email,
		},
		Members: members,
		Columns: columns,
	}))
}

func (b *boardController) CreateBoard(c echo.Context) error {
	var payload dto.BoardRequest

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	userId := m.ExtractTokenUserId(c)

	if err := b.useCase.CreateBoard(uint(userId), &payload); err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessResponse("Success add new board"))
}

func (b *boardController) UpdateBoard(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse("Bad Request: Invalid Id"))
	}

	var payload dto.BoardRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	userId := m.ExtractTokenUserId(c)

	if err := b.useCase.UpdateBoard(uint(id), uint(userId), &payload); err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessResponse("Success update board"))
}

func (b *boardController) DeleteBoard(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse("Bad Request: Invalid Id"))
	}

	userId := m.ExtractTokenUserId(c)

	if err := b.useCase.DeleteBoard(uint(id), uint(userId)); err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessResponse("Succcess delete board"))
}
