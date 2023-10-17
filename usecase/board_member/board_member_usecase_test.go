package usecase

import (
	"errors"
	"kanban-board/dto"
	"kanban-board/model"
	repo "kanban-board/repository/board_member"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)


var mockUserData = &model.User{
	Model: gorm.Model{ID: 1},
	Name: "Alvin",
	Email: "alvin@gmail.com",
	Password: "alvin123",
}

var mockBoardData = &model.Board{
	Model: gorm.Model{ID: 1},
	Name: "Golang Project",
	Desc: "Project Description",
	OwnerID: 1,
	Members: nil,
}

func TestAddMember(t *testing.T) {
	t.Run("Success Add Member", func (t* testing.T) {
		mockRequest := &dto.BoardMemberRequest{
			BoardID: mockBoardData.ID,
			UserID: mockUserData.ID,
		}

		mockRepo := repo.NewMockBoardMemberRepo()
		mockRepo.On("AddMember", mockRequest.BoardID, mockRequest.UserID).Return(nil).Once()

		service := NewBoardMemberUseCase(mockRepo)
		err := service.AddMember(mockRequest)

		assert.NoError(t, err, nil)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Add Member (struct error)", func(t *testing.T) {
		mockRequest := &dto.BoardMemberRequest{}

		mockRepo := repo.NewMockBoardMemberRepo()

		service := NewBoardMemberUseCase(mockRepo)
		err := service.AddMember(mockRequest)

		assert.Error(t, err)
	})

	t.Run("Failed Add Member (Internal Server Error)", func(t *testing.T) {
		mockRequest := &dto.BoardMemberRequest{
			BoardID: mockBoardData.ID,
			UserID: mockUserData.ID,
		}
		expectedErr := errors.New("Internal Server Error")
		
		mockRepo := repo.NewMockBoardMemberRepo()
		mockRepo.On("AddMember", mockRequest.BoardID, mockRequest.UserID).Return(expectedErr).Once()

		service := NewBoardMemberUseCase(mockRepo)
		err := service.AddMember(mockRequest)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteMember(t *testing.T) {
	t.Run("Success Delete Member", func (t* testing.T) {
		mockRequest := &dto.BoardMemberRequest{
			BoardID: mockBoardData.ID,
			UserID: mockUserData.ID,
		}

		mockRepo := repo.NewMockBoardMemberRepo()
		mockRepo.On("DeleteMember", mockRequest.BoardID, mockRequest.UserID).Return(nil).Once()

		service := NewBoardMemberUseCase(mockRepo)
		err := service.DeleteMember(mockRequest)

		assert.NoError(t, err, nil)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Delete Member (struct error)", func(t *testing.T) {
		mockRequest := &dto.BoardMemberRequest{}

		mockRepo := repo.NewMockBoardMemberRepo()

		service := NewBoardMemberUseCase(mockRepo)
		err := service.DeleteMember(mockRequest)

		assert.Error(t, err)
	})

	t.Run("Failed Delete Member (Internal Server Error)", func(t *testing.T) {
		mockRequest := &dto.BoardMemberRequest{
			BoardID: mockBoardData.ID,
			UserID: mockUserData.ID,
		}
		expectedErr := errors.New("Internal Server Error")
		
		mockRepo := repo.NewMockBoardMemberRepo()
		mockRepo.On("DeleteMember", mockRequest.BoardID, mockRequest.UserID).Return(expectedErr).Once()

		service := NewBoardMemberUseCase(mockRepo)
		err := service.DeleteMember(mockRequest)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}