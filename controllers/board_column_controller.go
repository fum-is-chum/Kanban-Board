package controller

import (
	"fmt"
	"kanban-board/dto"
	responseHelper "kanban-board/helpers/response"
	m "kanban-board/middlewares"
	useCase "kanban-board/usecase/board_column"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type boardColumnController struct {
	useCase useCase.BoardColumnUseCase
}

func NewBoardColumnController(useCase useCase.BoardColumnUseCase) *boardColumnController {
	return &boardColumnController{useCase}
}

func (b *boardColumnController) GetColumns(c echo.Context) error {
	id, err := strconv.Atoi(c.QueryParam("boardId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	columns, err := b.useCase.GetColumns(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	var response []dto.BoardColumnResponse

	for _, column := range columns {
		response = append(response, dto.BoardColumnResponse{
			ID:      column.ID,
			Label:   column.Label,
			Desc:    column.Desc,
			BoardID: column.BoardID,
		})
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessWithDataResponse("Success get columns", response))
}

func (b *boardColumnController) GetColumn(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	column, err := b.useCase.GetColumnById(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessWithDataResponse(fmt.Sprintf("Success get column with id %d", id), dto.BoardColumnResponse{
		ID:      column.ID,
		Label:   column.Label,
		Desc:    column.Desc,
		BoardID: column.BoardID,
	}))
}

func (b *boardColumnController) CreateColumn(c echo.Context) error {
	var payload dto.BoardColumnRequest

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	if err := b.useCase.CreateColumn(&payload); err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessResponse("Success create column"))
}

func (b *boardColumnController) UpdateColumn(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	var payload dto.BoardColumnRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	if err := b.useCase.UpdateColumn(uint(id), &payload); err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessResponse("Success update column"))
}

func (b *boardColumnController) DeleteColumn(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	userId := m.ExtractTokenUserId(c)

	if err := b.useCase.DeleteColumn(uint(id), uint(userId)); err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessResponse("Success delete column"))
}
