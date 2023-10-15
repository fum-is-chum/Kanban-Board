package controller

import (
	"fmt"
	"kanban-board/dto"
	responseHelper "kanban-board/helpers/response"
	boardMemberUsecase "kanban-board/usecase/board_member"
	"net/http"

	"github.com/labstack/echo/v4"
)

type boardMemberController struct {
	useCase boardMemberUsecase.BoardMemberUseCase
}

func NewBoardMemberController(useCase boardMemberUsecase.BoardMemberUseCase) *boardMemberController {
	return &boardMemberController{useCase}
}

func (b *boardMemberController) AddMember(c echo.Context) error {
	var request dto.BoardMemberRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	if err := b.useCase.AddMember(&request); err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessResponse("Success add member to board"))
}

func (b *boardMemberController) RemoveMember(c echo.Context) error {
	var request dto.BoardMemberRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	if err := b.useCase.DeleteMember(&request); err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessResponse("Success delete member from board"))
}
