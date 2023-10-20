package controller

import (
	"fmt"
	"kanban-board/dto"
	responseHelper "kanban-board/helpers/response"
	authUsecase "kanban-board/usecase/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

type authController struct {
	useCase authUsecase.AuthUseCase
}

func NewAuthController(authUsecase authUsecase.AuthUseCase) *authController {
	return &authController{authUsecase}
}

func (u *authController) Login(c echo.Context) error {
	var loginRequest dto.LoginRequest

	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	token, user, err := u.useCase.Login(&loginRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	loginResponse := dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessWithDataResponse("Sucess Login", loginResponse))
}
