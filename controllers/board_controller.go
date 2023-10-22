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
	for _, value := range boards {
		response = append(response, dto.BoardResponse{
			ID:   value.ID,
			Name: value.Name,
			Desc: value.Desc,
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

	var members []*dto.MemberResponse
	for _, value := range board.Members {
		members = append(members, &dto.MemberResponse{
			ID:        value.ID,
			Name:      value.Name,
			Email:     value.Email,
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		})
	}

	var columns []*dto.BoardColumnResponse
	for _, value := range board.Columns {
		columns = append(columns, &dto.BoardColumnResponse{
			ID:    value.ID,
			Label: value.Label,
			Desc:  value.Desc,
		})
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessWithDataResponse(fmt.Sprintf("Success get board with id %d", id), dto.BoardResponse{
		ID:   board.ID,
		Name: board.Name,
		Desc: board.Desc,
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
