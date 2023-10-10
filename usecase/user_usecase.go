package usecase

import (
	"errors"
	"kanban-board/dto"
	"kanban-board/helpers"
	"kanban-board/model"
	"kanban-board/repository"
	"reflect"
)

type UserUseCase interface {
	GetUsers() ([]model.User, error)
	GetUserById(id uint) (model.User, error)
	CreateUser(data *dto.UserRequest) error
	UpdateUser(id uint, data *dto.UserRequest) error
	DeleteUser(id uint) error
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *userUseCase {
	return &userUseCase{userRepo}
}

// ---------------- utility -------------------
func createUserRequestToUser(data *dto.UserRequest) (model.User, error) {
	hash, err := helpers.HashPassword(data.Password)
	if err != nil {
		return model.User{}, errors.New("Hash password failed")
	}

    return model.User{
        Name:     data.Name,
        Email:    data.Email,
        Password: hash,
    }, nil
}
// ---------------------------------------------

func (u *userUseCase) GetUsers() ([]model.User, error) {
	users, err := u.userRepo.Get()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userUseCase) GetUserById(id uint) (model.User, error) {
	user, err := u.userRepo.GetById(id)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *userUseCase) CreateUser(data *dto.UserRequest) error {
	user, err := createUserRequestToUser(data)
	if err != nil {
		return err
	}

	if err := u.userRepo.Create(&user); err != nil {
		return err
	}
	return nil
}

func (u *userUseCase) UpdateUser(id uint, data *dto.UserRequest) error {
	var updates map[string]interface{}

	structValue := reflect.ValueOf(data)

	for i := 0; i < structValue.NumField(); i++ {
		key := structValue.Type().Field(i).Name
		value := structValue.Field(i).Interface()

		if value != nil && value != "" {
			updates[key] = value
		}
	}

	if err := u.userRepo.Update(id, &updates); err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) DeleteUser(id uint) error {
	if err := u.userRepo.Delete(id); err != nil {
		return err
	}

	return nil
}
