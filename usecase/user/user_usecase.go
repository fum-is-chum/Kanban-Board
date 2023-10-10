package usecase

import (
	"errors"
	"kanban-board/dto"
	bcrypt "kanban-board/helpers/bcrypt"
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
	Login(email string, password string) error
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *userUseCase {
	return &userUseCase{userRepo}
}

// ---------------- utility -------------------
func createUserRequestToUser(data *dto.UserRequest) (model.User, error) {
	hash, err := bcrypt.HashPassword(data.Password)
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

func (u *userUseCase) GetUserById(id uint) (*model.User, error) {
	user, err := u.userRepo.GetById(id)
	if err != nil {
		return &model.User{}, err
	}

	return user, nil
}

func (u *userUseCase) CreateUser(data *dto.UserRequest) error {
	// Validate user request first
	err := validate.Struct(data)
	if err != nil {
		return err
	}

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
	updates := make(map[string]interface{})
	
	structValue := reflect.ValueOf(*data)
	for i := 0; i < structValue.NumField(); i++ {
		key := structValue.Type().Field(i).Name
		value := structValue.Field(i).Interface()

		// Check if the field is updated (non-empty) and not an empty string
		if value != nil && value != "" {
			if key == "Email" {
				if err := validate.Var(value, "email"); err != nil {
					return errors.New("email invalid!")
				}
			}
			updates[key] = value
		}
	}

	// Check if there are no fields to update
	if len(updates) == 0 {
		return errors.New("No fields to update or fields value is empty!")
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

func (u *userUseCase) Login(email string, password string) error {
	// find user by email
	user, err := u.userRepo.GetByEmail(email);
	if err != nil {
		return err
	}

	// validate password
	if err := bcrypt.VerifyPassword(user.Password, password); err != nil {
		return err
	}

	return nil
}