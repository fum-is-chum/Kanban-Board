package usecase

import (
	"errors"
	"fmt"
	"kanban-board/dto"
	bcrypt "kanban-board/helpers/bcrypt"
	fieldHelper "kanban-board/helpers/field"
	"kanban-board/model"
	repository "kanban-board/repository/user"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type UserUseCase interface {
	GetUsers() ([]model.User, error)
	GetUserById(id uint) (*model.User, error)
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
func createUserRequestToUserModel(data *dto.UserRequest) (*model.User, error) {
	hash, err := bcrypt.HashPassword(data.Password)
	if err != nil || data.Password == "" {
		return nil, errors.New("Hash password failed")
	}

	return &model.User{
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

func (u *userUseCase) GetUserById(id uint) (*model.User, error) {
	user, err := u.userRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) CreateUser(data *dto.UserRequest) error {
	// Validate user request first
	err := validate.Struct(data)
	if err != nil {
		return err
	}

	user, err := createUserRequestToUserModel(data)
	if err != nil {
		return err
	}

	if err := u.userRepo.Create(user); err != nil {
		return err
	}
	return nil
}


func (u *userUseCase) UpdateUser(id uint, data *dto.UserRequest) error {
	val := reflect.ValueOf(*data)
	fmt.Print(fieldHelper.IsFieldSet(&val, "Email"))
	fmt.Print(fieldHelper.IsFieldSet(&val, "Password"))

	if fieldHelper.IsFieldSet(&val, "Email") {
		if err := validate.Var(data.Email, "email"); err != nil {
			return errors.New("Email invalid")
		}
	}

	if fieldHelper.IsFieldSet(&val, "Password") {
		hashedPass, err := bcrypt.HashPassword(data.Password)
		if err != nil {
			return err
		}

		data.Password = hashedPass
	}
	
	if err := u.userRepo.Update(id, data); err != nil {
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