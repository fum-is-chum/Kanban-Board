package usecase

import (
	"errors"
	"kanban-board/dto"
	bcrypt "kanban-board/helpers/bcrypt"
	"kanban-board/model"
	userRepo "kanban-board/repository/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestGetUsers(t *testing.T) {
	returnData := []model.User{
		{Name: "User1", Email: "user1@example.com", Password: "password1"},
		{Name: "User2", Email: "user2@example.com", Password: "password2"},
	}

	t.Run("Sucess Get Users", func(t *testing.T) {
		mockRepo := userRepo.NewMockUserRepo()
		mockRepo.On("Get", mock.Anything).Return(returnData, nil).Once()
		service := NewUserUseCase(mockRepo)

		res, err := service.GetUsers()
		assert.NoError(t, err)
		assert.Equal(t, len(returnData), len(res))
		assert.Equal(t, returnData[0].Name, res[0].Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Get Users", func(t *testing.T) {
		mockRepo := userRepo.NewMockUserRepo()
		expectedErr := errors.New("internal server error") // Define the expected error
		mockRepo.On("Get", mock.Anything).Return(nil, expectedErr).Once()
		service := NewUserUseCase(mockRepo)

		res, err := service.GetUsers()
		assert.Error(t, err)
		assert.Nil(t, res) // No users should be returned in case of an error
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("Success create user", func(t *testing.T) {
		data := dto.UserRequest{
			Name:     "Alvin",
			Email:    "alvinardiansyah2002@gmail.com",
			Password: "123",
		}

		user, err := createUserRequestToUser(&data)
		assert.Equal(t, err, nil)

		mockRepo := userRepo.NewMockUserRepo()
		mockRepo.On("Create", mock.MatchedBy(func(u *model.User) bool {
			return u.Name == user.Name &&
				u.Email == user.Email
		})).Return(nil)

		service := NewUserUseCase(mockRepo)
		err = service.CreateUser(&data)

		assert.Equal(t, err, nil)

		// Use userData when asserting the Create method is called
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Create User missing name", func(t *testing.T) {
		data := dto.UserRequest{
			Name:     "",
			Email:    "alvinardiansyah2002@gmail.com",
			Password: "123",
		}

		mockRepo := userRepo.NewMockUserRepo()
		service := NewUserUseCase(mockRepo)
		err := service.CreateUser(&data)
		assert.NotEqual(t, err, nil)
	})

	t.Run("Failed Create User missing email", func(t *testing.T) {
		data := dto.UserRequest{
			Name:     "Alvin",
			Email:    "",
			Password: "123",
		}

		mockRepo := userRepo.NewMockUserRepo()
		service := NewUserUseCase(mockRepo)
		err := service.CreateUser(&data)
		assert.NotEqual(t, err, nil)
	})

	t.Run("Failed Create User missing password", func(t *testing.T) {
		data := dto.UserRequest{
			Name:     "Alvin",
			Email:    "alvinardiansyah2002@gmail.com",
			Password: "",
		}

		mockRepo := userRepo.NewMockUserRepo()
		service := NewUserUseCase(mockRepo)
		err := service.CreateUser(&data)
		assert.NotEqual(t, err, nil)
	})
}

func TestUpdateUser(t *testing.T) {
	userToUpdate := &model.User{
		Name:     "Alvin Christ Ardiansyah",
		Email:    "alvinardiansyah2002@gmail.com",
		Password: "123",
	}

	t.Run("Success Update User Name", func(t *testing.T) {
		userRequest := &dto.UserRequest{
			Name: "Alvin (Updated)",
		}

		updateMap := &map[string]interface{}{
			"Name": "Alvin (Updated)",
		}

		mockRepo := userRepo.NewMockUserRepo()
		mockRepo.On("Update", userToUpdate.ID, updateMap).Return(nil).Once()

		userUseCase := NewUserUseCase(mockRepo)
		err := userUseCase.UpdateUser(userToUpdate.ID, userRequest)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Update User Name (Name Empty)", func(t *testing.T) {
		userRequest := &dto.UserRequest{
			Name: "",
		}

		updateMap := &map[string]interface{}{}

		mockRepo := userRepo.NewMockUserRepo()
		mockRepo.On("Update", userToUpdate.ID, updateMap).Return(nil).Once()

		userUseCase := NewUserUseCase(mockRepo)
		err := userUseCase.UpdateUser(userToUpdate.ID, userRequest)

		assert.Error(t, err)
	})

	t.Run("Success Update User Email", func(t *testing.T) {
		userRequest := &dto.UserRequest{
			Email: "johndoe@gmail.com",
		}

		updateMap := &map[string]interface{}{
			"Email": "johndoe@gmail.com",
		}

		mockRepo := userRepo.NewMockUserRepo()
		mockRepo.On("Update", userToUpdate.ID, updateMap).Return(nil).Once()

		userUseCase := NewUserUseCase(mockRepo)
		err := userUseCase.UpdateUser(userToUpdate.ID, userRequest)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Update User Email (Email invalid)", func(t *testing.T) {
		userRequest := &dto.UserRequest{
			Email: "johndoegmail.com",
		}

		mockRepo := userRepo.NewMockUserRepo()
		userUseCase := NewUserUseCase(mockRepo)
		err := userUseCase.UpdateUser(userToUpdate.ID, userRequest)

		assert.Error(t, err)
	})

	t.Run("Failed Update User Email (Email Empty)", func(t *testing.T) {
		userRequest := &dto.UserRequest{
			Email: "",
		}

		mockRepo := userRepo.NewMockUserRepo()
		userUseCase := NewUserUseCase(mockRepo)
		err := userUseCase.UpdateUser(userToUpdate.ID, userRequest)

		assert.Error(t, err)
	})

	t.Run("Success Update User Password", func(t *testing.T) {
		userRequest := &dto.UserRequest{
			Password: "12345",
		}

		updateMap := &map[string]interface{}{
			"Password": "12345",
		}

		mockRepo := userRepo.NewMockUserRepo()
		mockRepo.On("Update", userToUpdate.ID, updateMap).Return(nil).Once()

		userUseCase := NewUserUseCase(mockRepo)
		err := userUseCase.UpdateUser(userToUpdate.ID, userRequest)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Update User Password (Password Empty)", func(t *testing.T) {
		userRequest := &dto.UserRequest{
			Password: "",
		}

		updateMap := &map[string]interface{}{}

		mockRepo := userRepo.NewMockUserRepo()
		mockRepo.On("Update", userToUpdate.ID, updateMap).Return(nil).Once()

		userUseCase := NewUserUseCase(mockRepo)
		err := userUseCase.UpdateUser(userToUpdate.ID, userRequest)

		assert.Error(t, err)
	})
}

func TestGetUserById(t *testing.T) {
	expectedUser := &model.User{
		Model:    gorm.Model{ID: 3},
		Name:     "Alvin Christ Ardiansyah",
		Email:    "alvinardiansyah2002@gmail.com",
		Password: "123",
	}

	t.Run("Success Get User By Id", func(t *testing.T) {
		mockRepo := userRepo.NewMockUserRepo()
		mockRepo.On("GetById", expectedUser.ID).Return(expectedUser, nil).Once()

		userUseCase := NewUserUseCase(mockRepo)
		user, err := userUseCase.GetUserById(expectedUser.ID)

		assert.NoError(t, err)
		assert.Equal(t, user.Email, expectedUser.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Get User By Id (ID not found)", func(t *testing.T) {
		mockRepo := userRepo.NewMockUserRepo()
		expectedErr := errors.New("internal server error") // Define the expected error
		mockRepo.On("GetById", uint(9)).Return(nil, expectedErr).Once()

		userUseCase := NewUserUseCase(mockRepo)
		_, err := userUseCase.GetUserById(uint(9))

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteUser(t *testing.T) {
	toDeleteUser := &model.User{
		Model:    gorm.Model{ID: 3},
		Name:     "Alvin Christ Ardiansyah",
		Email:    "alvinardiansyah2002@gmail.com",
		Password: "123",
	}
	t.Run("Sucess Delete User", func(t *testing.T) {
		mockRepo := userRepo.NewMockUserRepo()
		mockRepo.On("Delete", toDeleteUser.ID).Return(nil).Once()

		userUseCase := NewUserUseCase(mockRepo)
		err := userUseCase.DeleteUser(toDeleteUser.ID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Delete User (ID not found)", func(t *testing.T) {
		mockRepo := userRepo.NewMockUserRepo()
		expectedErr := errors.New("internal server error") // Define the expected error

		mockRepo.On("Delete", toDeleteUser.ID).Return(expectedErr).Once()

		userUseCase := NewUserUseCase(mockRepo)
		err := userUseCase.DeleteUser(toDeleteUser.ID)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

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

		userUseCase := NewUserUseCase(mockRepo)
		err := userUseCase.Login(expectedUser.Email, "123")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Login (email not found)", func(t *testing.T) {
		mockRepo := userRepo.NewMockUserRepo()
		expectedErr := errors.New("Email not found")

		mockRepo.On("GetByEmail", "doe@gmail.com").Return(nil, expectedErr).Once()

		userUseCase := NewUserUseCase(mockRepo)
		err := userUseCase.Login("doe@gmail.com", "123")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed LOgin (password invalid)", func(t *testing.T) {
		mockRepo := userRepo.NewMockUserRepo()

		mockRepo.On("GetByEmail", expectedUser.Email).Return(expectedUser, nil).Once()

		userUseCase := NewUserUseCase(mockRepo)
		err := userUseCase.Login(expectedUser.Email, "12345")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
