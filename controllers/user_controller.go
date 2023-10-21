package controller

import (
	"errors"
	"fmt"
	"kanban-board/dto"
	responseHelper "kanban-board/helpers/response"
	userUsecase "kanban-board/usecase/user"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type userController struct {
	useCase userUsecase.UserUseCase
}

func NewUserController(userUsecase userUsecase.UserUseCase) *userController {
	return &userController{userUsecase}
}

func (u *userController) GetUsers(c echo.Context) error {
	users, err := u.useCase.GetUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	var response []dto.UserResponse
	for _, value := range users {
		response = append(response, dto.UserResponse{
			ID:        value.ID,
			Name:      value.Name,
			Email:     value.Email,
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessWithDataResponse("Sucess get users", response))
}

func (u *userController) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse("Bad Request: ID invalid"))
	}

	user, err := u.useCase.GetUserById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, responseHelper.FailedResponse("User not found"))
		}
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error on fetching user with id %d", id)))
	}

	var memberOf []*dto.BoardResponse
	for _, board := range user.MemberOf {
		memberOf = append(memberOf, &dto.BoardResponse{
			ID:   board.ID,
			Name: board.Name,
			Desc: board.Desc,
		})
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessWithDataResponse(fmt.Sprintf("Sucess fetch user with id %d", id), dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		MemberOf:  memberOf,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}))
}

func (u *userController) CreateUser(c echo.Context) error {
	var payload dto.UserRequest

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad Request: %s", err.Error())))
	}

	err := u.useCase.CreateUser(&payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusCreated, responseHelper.SuccessResponse("Success create user"))
}

func (u *userController) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse("Bad request: ID invalid"))
	}

	var payload dto.UserRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse(fmt.Sprintf("Bad request: %s", err.Error())))
	}

	if err := u.useCase.UpdateUser(uint(id), &payload); err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessResponse("Success update user"))
}

func (u *userController) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responseHelper.FailedResponse("Bad request: ID invalid"))
	}

	if err := u.useCase.DeleteUser(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, responseHelper.FailedResponse(fmt.Sprintf("Error: %s", err.Error())))
	}

	return c.JSON(http.StatusOK, responseHelper.SuccessResponse("Success delete user"))
}
