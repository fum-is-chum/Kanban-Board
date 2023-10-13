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
	boards, err := b.useCase.GetBoards()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	var response []dto.BoardResponse
	for _, value := range boards {
		response = append(response, dto.BoardResponse{
			Id:      value.ID,
			Name:    value.Name,
			Desc:    value.Desc,
			OwnerID: value.OwnerID,
		})
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessWithDataResponse("Success get boards", response))
}

func (b *boardController) GetBoardById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse("Bad Request: Invalid Id"))
	}

	board, err := b.useCase.GetBoardById(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessWithDataResponse(fmt.Sprintf("Success get board with id %d", id), dto.BoardResponse{
		Id: board.ID,
		Name: board.Name,
		Desc: board.Desc,
		OwnerID: board.OwnerID,
	}))
}

func (b *boardController) CreateBoard(c echo.Context) error {
	var payload dto.BoardRequest

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	if err := b.useCase.CreateBoard(&payload); err != nil {
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

	if err := b.useCase.UpdateBoard(uint(id), &payload); err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessResponse("Success update board"))
}

func (b *boardController) UpdateBoardOwner(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse("Bad Request: Invalid Id"))
	}

	var payload dto.BoardRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	if err := b.useCase.UpdateBoardOwnership(uint(id), uint(payload.OwnerID)); err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}
	
	return c.JSON(http.StatusOK, responseHelper.SuccessResponse("Success update board owner"))
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