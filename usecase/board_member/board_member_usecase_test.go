package usecase

import (
	"errors"
	"kanban-board/dto"
	"kanban-board/model"
	boardRepo "kanban-board/repository/board"
	memberRepo "kanban-board/repository/board_member"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var mockUserData = &model.User{
	Model:    gorm.Model{ID: 1},
	Name:     "Alvin",
	Email:    "alvin@gmail.com",
	Password: "alvin123",
}

var mockBoardData = &model.Board{
	Model:   gorm.Model{ID: 1},
	Name:    "Golang Project",
	Desc:    "Project Description",
	OwnerID: 1,
	Members: nil,
}

var boardMembers = []model.BoardMember{
	{
		UserID:  1,
		BoardID: 1,
	},
	{
		UserID:  2,
		BoardID: 1,
	},
}

var issuerId = uint(1)

func TestIsOwner(t *testing.T) {
	t.Run("Success isOwner", func(t *testing.T) {
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := memberRepo.NewMockBoardMemberRepo()
		mockBoardRepo.On("GetBoardOwner", mockBoardData.ID).Return(&issuerId, nil).Once()

		service := NewBoardMemberUseCase(mockBoardRepo, mockMemberRepo)
		err := service.isOwner(issuerId, mockBoardData.ID)

		assert.NoError(t, err)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed isOwner (Not owner)", func(t *testing.T) {
		expectedErr := errors.New("User is not owner of this board!")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := memberRepo.NewMockBoardMemberRepo()
		mockBoardRepo.On("GetBoardOwner", mockBoardData.ID).Return(&issuerId, nil).Once()

		service := NewBoardMemberUseCase(mockBoardRepo, mockMemberRepo)
		err := service.isOwner(uint(10), mockBoardData.ID)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), expectedErr.Error())
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed isOwner (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := memberRepo.NewMockBoardMemberRepo()
		mockBoardRepo.On("GetBoardOwner", mockBoardData.ID).Return(nil, expectedErr).Once()

		service := NewBoardMemberUseCase(mockBoardRepo, mockMemberRepo)
		err := service.isOwner(issuerId, mockBoardData.ID)

		assert.Error(t, err)
		mockBoardRepo.AssertExpectations(t)
	})
}

func TestIsMember(t *testing.T) {
	t.Run("Success isMember", func(t *testing.T) {
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := memberRepo.NewMockBoardMemberRepo()
		mockBoardRepo.On("GetBoardMembers", mockBoardData.ID).Return(boardMembers, nil).Once()

		service := NewBoardMemberUseCase(mockBoardRepo, mockMemberRepo)
		err := service.isMember(issuerId, mockBoardData.ID)

		assert.NoError(t, err)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed isMember (Not a member)", func(t *testing.T) {
		expectedErr := errors.New("User is not member of this board!")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := memberRepo.NewMockBoardMemberRepo()
		mockBoardRepo.On("GetBoardMembers", mockBoardData.ID).Return(boardMembers, nil).Once()

		service := NewBoardMemberUseCase(mockBoardRepo, mockMemberRepo)
		err := service.isMember(uint(10), mockBoardData.ID)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), expectedErr.Error())
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed isMember (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := memberRepo.NewMockBoardMemberRepo()
		mockBoardRepo.On("GetBoardMembers", mockBoardData.ID).Return(nil, expectedErr).Once()

		service := NewBoardMemberUseCase(mockBoardRepo, mockMemberRepo)
		err := service.isMember(issuerId, mockBoardData.ID)

		assert.Error(t, err)
		mockBoardRepo.AssertExpectations(t)
	})
}

func TestAddMember(t *testing.T) {
	t.Run("Success Add Member", func(t *testing.T) {
		mockRequest := &dto.BoardMemberRequest{
			BoardID: mockBoardData.ID,
			UserID:  mockUserData.ID,
		}

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := memberRepo.NewMockBoardMemberRepo()
		mockBoardRepo.On("GetBoardOwner", mockRequest.BoardID).Return(&issuerId, nil).Once()
		mockMemberRepo.On("AddMember", mockRequest.BoardID, mockRequest.UserID).Return(nil).Once()

		service := NewBoardMemberUseCase(mockBoardRepo, mockMemberRepo)
		err := service.AddMember(issuerId, mockRequest)

		assert.NoError(t, err, nil)
		mockMemberRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Add Member (struct error)", func(t *testing.T) {
		mockRequest := &dto.BoardMemberRequest{}

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := memberRepo.NewMockBoardMemberRepo()

		service := NewBoardMemberUseCase(mockBoardRepo, mockMemberRepo)
		err := service.AddMember(issuerId, mockRequest)

		assert.Error(t, err)
	})

	t.Run("Failed Add Member (User not owner of this board)", func(t *testing.T) {
		mockRequest := &dto.BoardMemberRequest{
			BoardID: mockBoardData.ID,
			UserID:  mockUserData.ID,
		}
		expectedErr := errors.New("User is not owner of this board!")

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := memberRepo.NewMockBoardMemberRepo()
		mockBoardRepo.On("GetBoardOwner", mockRequest.BoardID).Return(&issuerId, nil).Once()

		service := NewBoardMemberUseCase(mockBoardRepo, mockMemberRepo)
		err := service.AddMember(uint(10), mockRequest)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), expectedErr.Error())
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Add Member (Internal Server Error)", func(t *testing.T) {
		mockRequest := &dto.BoardMemberRequest{
			BoardID: mockBoardData.ID,
			UserID:  mockUserData.ID,
		}
		expectedErr := errors.New("Internal Server Error")

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := memberRepo.NewMockBoardMemberRepo()
		mockBoardRepo.On("GetBoardOwner", mockRequest.BoardID).Return(&issuerId, nil).Once()
		mockMemberRepo.On("AddMember", mockRequest.BoardID, mockRequest.UserID).Return(expectedErr).Once()

		service := NewBoardMemberUseCase(mockBoardRepo, mockMemberRepo)
		err := service.AddMember(issuerId, mockRequest)

		assert.Error(t, err)
		mockMemberRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})
}

func TestDeleteMember(t *testing.T) {
	t.Run("Success Delete Member", func(t *testing.T) {
		mockRequest := &dto.BoardMemberRequest{
			BoardID: mockBoardData.ID,
			UserID:  mockUserData.ID,
		}

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := memberRepo.NewMockBoardMemberRepo()
		mockBoardRepo.On("GetBoardOwner", mockRequest.BoardID).Return(&issuerId, nil).Once()
		mockMemberRepo.On("DeleteMember", mockRequest.BoardID, mockRequest.UserID).Return(nil).Once()

		service := NewBoardMemberUseCase(mockBoardRepo, mockMemberRepo)
		err := service.DeleteMember(issuerId, mockRequest)

		assert.NoError(t, err, nil)
		mockMemberRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Delete Member (struct error)", func(t *testing.T) {
		mockRequest := &dto.BoardMemberRequest{}

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := memberRepo.NewMockBoardMemberRepo()

		service := NewBoardMemberUseCase(mockBoardRepo, mockMemberRepo)
		err := service.DeleteMember(issuerId, mockRequest)

		assert.Error(t, err)
	})

	t.Run("Failed Add Member (User not member of this board)", func(t *testing.T) {
		mockRequest := &dto.BoardMemberRequest{
			BoardID: mockBoardData.ID,
			UserID:  mockUserData.ID,
		}
		expectedErr := errors.New("User is not owner of this board!")

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := memberRepo.NewMockBoardMemberRepo()
		mockBoardRepo.On("GetBoardOwner", mockRequest.BoardID).Return(&issuerId, nil).Once()

		service := NewBoardMemberUseCase(mockBoardRepo, mockMemberRepo)
		err := service.DeleteMember(uint(10), mockRequest)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), expectedErr.Error())
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Delete Member (Internal Server Error)", func(t *testing.T) {
		mockRequest := &dto.BoardMemberRequest{
			BoardID: mockBoardData.ID,
			UserID:  mockUserData.ID,
		}
		expectedErr := errors.New("Internal Server Error")

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := memberRepo.NewMockBoardMemberRepo()
		mockBoardRepo.On("GetBoardOwner", mockRequest.BoardID).Return(&issuerId, nil).Once()
		mockMemberRepo.On("DeleteMember", mockRequest.BoardID, mockRequest.UserID).Return(expectedErr).Once()

		service := NewBoardMemberUseCase(mockBoardRepo, mockMemberRepo)
		err := service.DeleteMember(issuerId, mockRequest)

		assert.Error(t, err)
		mockMemberRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})
}

func TestExitBoard(t *testing.T) {
	t.Run("Success exit board", func(t *testing.T) {
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := memberRepo.NewMockBoardMemberRepo()
		mockBoardRepo.On("GetBoardMembers", mockBoardData.ID).Return(boardMembers, nil).Once()
		mockMemberRepo.On("DeleteMember", mockBoardData.ID, issuerId).Return(nil).Once()

		service := NewBoardMemberUseCase(mockBoardRepo, mockMemberRepo)
		err := service.ExitBoard(issuerId, mockBoardData.ID)

		assert.NoError(t, err)
		mockMemberRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed exit board (user is not member)", func(t *testing.T) {
		expectedErr := errors.New("User is not member of this board!")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := memberRepo.NewMockBoardMemberRepo()
		mockBoardRepo.On("GetBoardMembers", mockBoardData.ID).Return(boardMembers, nil).Once()

		service := NewBoardMemberUseCase(mockBoardRepo, mockMemberRepo)
		err := service.ExitBoard(uint(10), mockBoardData.ID)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), expectedErr.Error())
		mockBoardRepo.AssertExpectations(t)
	})

  	t.Run("Failed exit board (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := memberRepo.NewMockBoardMemberRepo()
		mockBoardRepo.On("GetBoardMembers", mockBoardData.ID).Return(boardMembers, nil).Once()
		mockMemberRepo.On("DeleteMember", mockBoardData.ID, issuerId).Return(expectedErr).Once()

		service := NewBoardMemberUseCase(mockBoardRepo, mockMemberRepo)
		err := service.ExitBoard(issuerId, mockBoardData.ID)

		assert.Error(t, err)
		mockMemberRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})
}