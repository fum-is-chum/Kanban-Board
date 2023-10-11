package usecase

import (
	"errors"
	"kanban-board/dto"
	bcrypt "kanban-board/helpers/bcrypt"
	"kanban-board/model"
	userRepo "kanban-board/repository/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestLogin(t *testing.T) {
	hashedPass, err := bcrypt.HashPassword("123")
	assert.NoError(t, err)

	expectedUser := &model.User{
		Model:    gorm.Model{ID: 3},
		Name:     "Alvin Christ Ardiansyah",
		Email:    "alvinardiansyah2002@gmail.com",
		Password: hashedPass,
	}
	t.Run("Sucess Login", func(t *testing.T) {
		mockRepo := userRepo.NewMockUserRepo()
		mockRepo.On("GetByEmail", expectedUser.Email).Return(expectedUser, nil).Once()

		userUseCase := NewAuthUseCase(mockRepo)
		token, user, err := userUseCase.Login(&dto.LoginRequest{Email: expectedUser.Email, Password: "123"})

		assert.NoError(t, err)
		assert.NotEqual(t, token, "")
		assert.NotNil(t, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Login (email not found)", func(t *testing.T) {
		mockRepo := userRepo.NewMockUserRepo()
		expectedErr := errors.New("Email not found")

		mockRepo.On("GetByEmail", "doe@gmail.com").Return(nil, expectedErr).Once()

		userUseCase := NewAuthUseCase(mockRepo)
		token, user, err := userUseCase.Login(&dto.LoginRequest{Email: "doe@gmail.com", Password: "123"})

		assert.Error(t, err)
		assert.Equal(t, token, "")
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Login (password invalid)", func(t *testing.T) {
		mockRepo := userRepo.NewMockUserRepo()

		mockRepo.On("GetByEmail", expectedUser.Email).Return(expectedUser, nil).Once()

		userUseCase := NewAuthUseCase(mockRepo)
		token, user, err := userUseCase.Login(&dto.LoginRequest{Email: expectedUser.Email, Password: "12345"})

		assert.Error(t, err)
		assert.Equal(t, token, "")
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}
