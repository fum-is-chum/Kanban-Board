package usecase

import (
	"kanban-board/dto"
	bcrypt "kanban-board/helpers/bcrypt"
	"kanban-board/middlewares"
	"kanban-board/model"
	repository "kanban-board/repository/user"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type AuthUseCase interface {
	Login(data *dto.LoginRequest) (string, *model.User, error)
}

type authUseCase struct {
	userRepo repository.UserRepository
}

func NewAuthUseCase(userRepo repository.UserRepository) *authUseCase {
	return &authUseCase{userRepo}
}

func (u *authUseCase) Login(data *dto.LoginRequest) (string, *model.User, error) {
	// find user by email
	user, err := u.userRepo.GetByEmail(data.Email)
	if err != nil {
		return "", nil, err
	}

	// validate password
	if err := bcrypt.VerifyPassword(user.Password, data.Password); err != nil {
		return "", nil, err
	}

	token, err := middlewares.CreateToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
