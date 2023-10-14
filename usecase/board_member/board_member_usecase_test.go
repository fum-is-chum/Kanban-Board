package usecase

import (
	"errors"
	"kanban-board/dto"
	"kanban-board/model"
	boardRepo "kanban-board/repository/board"
	boardMemberRepo "kanban-board/repository/board_member"
	userRepo "kanban-board/repository/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var mockBoardData = &model.Board{
	Name: "Board 1",
	Desc: "Board 1 Desc",
	OwnerID: uint(3),
}

var mockUserData = &model.User{
	Model: gorm.Model{ID: uint(2)},
	Name: "Alvin",
	Email: "alvin@gmail.com",
	Password: "alvin123",
	Boards: []model.Board{},
}

func TestAddMember(t *testing.T) {
	t.Run("Success Add Member", func(t *testing.T) {
		

		mockData := &dto.BoardMemberRequest{
			BoardId: uint(1),
			UserId:  uint(2),
		}

		issuerUserId := mockBoardData.OwnerID

		// mock repos
		mockMemberRepo := boardMemberRepo.NewMockBoardMemberRepo()
		mockUserRepo := userRepo.NewMockUserRepo()
		mockBoardRepo := boardRepo.NewMockBoardRepo()

		mockBoardRepo.On("GetById", mockData.BoardId).Return(mockBoardData, nil).Once()
		mockMemberRepo.On("AddMember", mockData.BoardId, mockData.UserId).Return(nil).Once()
		mockUserRepo.On("GetById", mockData.UserId).Return(mockUserData, nil).Once()

		multiRepos := &BoardMemberMultiRepos{
			MemberRepo: mockMemberRepo,
			BoardRepo: mockBoardRepo,
			UserRepo: mockUserRepo,
		}

		service := NewBoardMemberUseCase(multiRepos)
		err := service.AddNewMember(mockData, issuerUserId)

		assert.NoError(t, err)
		mockMemberRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Failed Add Member (userId missing)", func(t *testing.T) {
		mockData := &dto.BoardMemberRequest{
			BoardId: uint(1),
		}
		issuerUserId := mockBoardData.OwnerID

		// mock repos
		mockMemberRepo := boardMemberRepo.NewMockBoardMemberRepo()
		mockUserRepo := userRepo.NewMockUserRepo()
		mockBoardRepo := boardRepo.NewMockBoardRepo()

		multiRepos := &BoardMemberMultiRepos{
			MemberRepo: mockMemberRepo,
			BoardRepo: mockBoardRepo,
			UserRepo: mockUserRepo,
		}

		service := NewBoardMemberUseCase(multiRepos)
		err := service.AddNewMember(mockData, issuerUserId)

		assert.Error(t, err)
	})

	t.Run("Failed Add Member (boardId missing)", func(t *testing.T) {
		mockData := &dto.BoardMemberRequest{
			UserId: uint(1),
		}
		issuerUserId := mockBoardData.OwnerID

		// mock repos
		mockMemberRepo := boardMemberRepo.NewMockBoardMemberRepo()
		mockUserRepo := userRepo.NewMockUserRepo()
		mockBoardRepo := boardRepo.NewMockBoardRepo()

		multiRepos := &BoardMemberMultiRepos{
			MemberRepo: mockMemberRepo,
			BoardRepo: mockBoardRepo,
			UserRepo: mockUserRepo,
		}

		service := NewBoardMemberUseCase(multiRepos)
		err := service.AddNewMember(mockData, issuerUserId)

		assert.Error(t, err)
	})

	t.Run("Failed Add Member (Internal Server Error)", func(t *testing.T) {
		mockData := &dto.BoardMemberRequest{
			BoardId: uint(1),
			UserId:  uint(2),
		}
		issuerUserId := mockBoardData.OwnerID
		expectedErr := errors.New("Database Error")

		// mock repos
		mockMemberRepo := boardMemberRepo.NewMockBoardMemberRepo()
		mockUserRepo := userRepo.NewMockUserRepo()
		mockBoardRepo := boardRepo.NewMockBoardRepo()

		mockMemberRepo.On("AddMember", mockData.BoardId, mockData.UserId).Return(expectedErr).Once()
		mockBoardRepo.On("GetById", mockData.BoardId).Return(mockBoardData, nil).Once()
		mockUserRepo.On("GetById", mockData.UserId).Return(mockUserData, nil).Once()

		multiRepos := &BoardMemberMultiRepos{
			MemberRepo: mockMemberRepo,
			BoardRepo: mockBoardRepo,
			UserRepo: mockUserRepo,
		}

		service := NewBoardMemberUseCase(multiRepos)
		err := service.AddNewMember(mockData, issuerUserId)

		assert.Error(t, err)
		mockMemberRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Failed Add Member (not board owner)", func(t *testing.T) {
		mockData := &dto.BoardMemberRequest{
			BoardId: uint(1),
			UserId:  uint(2),
		}
		issuerUserId := uint(10)

		// mock repos
		mockMemberRepo := boardMemberRepo.NewMockBoardMemberRepo()
		mockUserRepo := userRepo.NewMockUserRepo()
		mockBoardRepo := boardRepo.NewMockBoardRepo()

		mockBoardRepo.On("GetById", mockData.BoardId).Return(mockBoardData, nil).Once()

		multiRepos := &BoardMemberMultiRepos{
			MemberRepo: mockMemberRepo,
			BoardRepo: mockBoardRepo,
			UserRepo: mockUserRepo,
		}

		service := NewBoardMemberUseCase(multiRepos)
		err := service.AddNewMember(mockData, issuerUserId)

		assert.Error(t, err)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Add Member (board not found)", func(t *testing.T) {
		mockData := &dto.BoardMemberRequest{
			BoardId: uint(1),
			UserId:  uint(2),
		}
		issuerUserId := mockBoardData.OwnerID
		expectedErr := errors.New("Board Not Found!")

		// mock repos
		mockMemberRepo := boardMemberRepo.NewMockBoardMemberRepo()
		mockUserRepo := userRepo.NewMockUserRepo()
		mockBoardRepo := boardRepo.NewMockBoardRepo()

		mockBoardRepo.On("GetById", mockData.BoardId).Return(nil, expectedErr).Once()

		multiRepos := &BoardMemberMultiRepos{
			MemberRepo: mockMemberRepo,
			BoardRepo: mockBoardRepo,
			UserRepo: mockUserRepo,
		}

		service := NewBoardMemberUseCase(multiRepos)
		err := service.AddNewMember(mockData, issuerUserId)

		assert.Error(t, err)
		mockBoardRepo.AssertExpectations(t)
	})
}

func TestDeleteMember(t *testing.T) {
	t.Run("Success Delete Member", func(t *testing.T) {
		mockData := &dto.BoardMemberRequest{
			BoardId: uint(1),
			UserId:  uint(2),
		}
		issuerUserId := mockBoardData.OwnerID

		// mock repos
		mockMemberRepo := boardMemberRepo.NewMockBoardMemberRepo()
		mockUserRepo := userRepo.NewMockUserRepo()
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		
		mockMemberRepo.On("DeleteMember", mockData.BoardId, mockData.UserId).Return(nil).Once()
		mockBoardRepo.On("GetById", mockData.BoardId).Return(mockBoardData, nil).Once()

		multiRepos := &BoardMemberMultiRepos{
			MemberRepo: mockMemberRepo,
			BoardRepo: mockBoardRepo,
			UserRepo: mockUserRepo,
		}

		service := NewBoardMemberUseCase(multiRepos)
		err := service.DeleteMember(mockData, issuerUserId)

		assert.NoError(t, err)
		mockMemberRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Delete Member (userId missing)", func(t *testing.T) {
		mockData := &dto.BoardMemberRequest{
			BoardId: uint(1),
		}
		issuerUserId := mockBoardData.OwnerID

		// mock repos
		mockMemberRepo := boardMemberRepo.NewMockBoardMemberRepo()
		mockUserRepo := userRepo.NewMockUserRepo()
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		
		multiRepos := &BoardMemberMultiRepos{
			MemberRepo: mockMemberRepo,
			BoardRepo: mockBoardRepo,
			UserRepo: mockUserRepo,
		}


		service := NewBoardMemberUseCase(multiRepos)
		err := service.DeleteMember(mockData, issuerUserId)

		assert.Error(t, err)
	})

	t.Run("Failed Delete Member (boardId missing)", func(t *testing.T) {
		mockData := &dto.BoardMemberRequest{
			UserId: uint(1),
		}
		issuerUserId := mockBoardData.OwnerID

		// mock repos
		mockMemberRepo := boardMemberRepo.NewMockBoardMemberRepo()
		mockUserRepo := userRepo.NewMockUserRepo()
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		
		multiRepos := &BoardMemberMultiRepos{
			MemberRepo: mockMemberRepo,
			BoardRepo: mockBoardRepo,
			UserRepo: mockUserRepo,
		}

		service := NewBoardMemberUseCase(multiRepos)
		err := service.DeleteMember(mockData, issuerUserId)

		assert.Error(t, err)
	})

	t.Run("Failed Delete Member (Internal Server Error)", func(t *testing.T) {
		mockData := &dto.BoardMemberRequest{
			BoardId: uint(1),
			UserId:  uint(2),
		}
		issuerUserId := mockBoardData.OwnerID
		expectedErr := errors.New("Database Error")

		// mock repos
		mockMemberRepo := boardMemberRepo.NewMockBoardMemberRepo()
		mockUserRepo := userRepo.NewMockUserRepo()
		mockBoardRepo := boardRepo.NewMockBoardRepo()

		mockMemberRepo.On("DeleteMember", mockData.BoardId, mockData.UserId).Return(expectedErr).Once()
		mockBoardRepo.On("GetById", mockData.BoardId).Return(mockBoardData, nil).Once()

		multiRepos := &BoardMemberMultiRepos{
			MemberRepo: mockMemberRepo,
			BoardRepo: mockBoardRepo,
			UserRepo: mockUserRepo,
		}

		service := NewBoardMemberUseCase(multiRepos)
		err := service.DeleteMember(mockData, issuerUserId)

		assert.Error(t, err)
		mockMemberRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Delete Member (not board owner)", func(t *testing.T) {
		mockData := &dto.BoardMemberRequest{
			BoardId: uint(1),
			UserId:  uint(2),
		}
		issuerUserId := uint(10)

		// mock repos
		mockMemberRepo := boardMemberRepo.NewMockBoardMemberRepo()
		mockUserRepo := userRepo.NewMockUserRepo()
		mockBoardRepo := boardRepo.NewMockBoardRepo()

		mockBoardRepo.On("GetById", mockData.BoardId).Return(mockBoardData, nil).Once()

		multiRepos := &BoardMemberMultiRepos{
			MemberRepo: mockMemberRepo,
			BoardRepo: mockBoardRepo,
			UserRepo: mockUserRepo,
		}

		service := NewBoardMemberUseCase(multiRepos)
		err := service.DeleteMember(mockData, issuerUserId)

		assert.Error(t, err)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Delete Member (board not found)", func(t *testing.T) {
		mockData := &dto.BoardMemberRequest{
			BoardId: uint(1),
			UserId:  uint(2),
		}
		issuerUserId := mockBoardData.OwnerID
		expectedErr := errors.New("Board Not Found!")

		// mock repos
		mockMemberRepo := boardMemberRepo.NewMockBoardMemberRepo()
		mockUserRepo := userRepo.NewMockUserRepo()
		mockBoardRepo := boardRepo.NewMockBoardRepo()

		mockBoardRepo.On("GetById", mockData.BoardId).Return(nil, expectedErr).Once()

		multiRepos := &BoardMemberMultiRepos{
			MemberRepo: mockMemberRepo,
			BoardRepo: mockBoardRepo,
			UserRepo: mockUserRepo,
		}

		service := NewBoardMemberUseCase(multiRepos)
		err := service.DeleteMember(mockData, issuerUserId)

		assert.Error(t, err)
		mockBoardRepo.AssertExpectations(t)
	})
}