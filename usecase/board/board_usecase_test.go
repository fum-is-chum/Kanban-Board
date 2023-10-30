package usecase

import (
	"errors"
	"kanban-board/dto"
	"kanban-board/model"
	boardRepo "kanban-board/repository/board"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var mockMembersData = []model.BoardMember{
	{
		BoardID: uint(1),
		UserID:  uint(1),
	},
	{
		BoardID: uint(1),
		UserID:  uint(2),
	},
	{
		BoardID: uint(1),
		UserID:  uint(3),
	},
}

func TestIsMember(t *testing.T) {
	boardId := uint(1)
	issuerId := uint(1)
	t.Run("Success isMember", func(t *testing.T) {
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockBoardRepo.On("GetBoardMembers", boardId).Return(mockMembersData, nil).Once()

		service := NewBoardUseCase(mockBoardRepo)
		err := service.isMember(issuerId, boardId)

		assert.NoError(t, err)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed isMember (Not a member)", func(t *testing.T) {
		expectedErr := errors.New("User is not member of this board!")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockBoardRepo.On("GetBoardMembers", boardId).Return(mockMembersData, nil).Once()

		service := NewBoardUseCase(mockBoardRepo)
		err := service.isMember(uint(10), boardId)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), expectedErr.Error())
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed isMember (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockBoardRepo.On("GetBoardMembers", boardId).Return(nil, expectedErr).Once()

		service := NewBoardUseCase(mockBoardRepo)
		err := service.isMember(issuerId, boardId)

		assert.Error(t, err)
		mockBoardRepo.AssertExpectations(t)
	})
}


func TestGetBoards(t *testing.T) {
	returnData := []model.Board{
		{Name: "board1", Desc: "Board number 1", OwnerID: 0},
		{Name: "board2", Desc: "Board number 2", OwnerID: 0},
	}

	issuerId := uint(1)

	t.Run("Success Get Boards", func(t *testing.T) {
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("Get", mock.Anything).Return(returnData, nil).Once()

		service := NewBoardUseCase(mockRepo)
		res, err := service.GetBoards(issuerId)

		assert.NoError(t, err)
		assert.Equal(t, len(returnData), len(res))
		assert.Equal(t, returnData[0].Name, res[0].Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Get Boards (Internal Server Error)", func(t *testing.T) {
		mockRepo := boardRepo.NewMockBoardRepo()
		expectedErr := errors.New("Internal Server Error")
		mockRepo.On("Get", mock.Anything).Return(nil, expectedErr).Once()

		service := NewBoardUseCase(mockRepo)
		res, err := service.GetBoards(issuerId)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})

}

func TestGetBoardById(t *testing.T) {
	expectedBoard := &model.Board{
		Model: gorm.Model{ID: 1},
		Name:  "Board 1",
		Desc:  "Board 1 Desc",
	}
	issuerId := uint(1)

	t.Run("Success Get Board By Id", func(t *testing.T) {
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("GetById", expectedBoard.ID, issuerId).Return(expectedBoard, nil).Once()

		service := NewBoardUseCase(mockRepo)
		board, err := service.GetBoardById(expectedBoard.ID, issuerId)

		assert.NoError(t, err)
		assert.NotNil(t, board)
		assert.Equal(t, board.ID, expectedBoard.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Get Board By Id (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("GetById", expectedBoard.ID, issuerId).Return(nil, expectedErr).Once()

		service := NewBoardUseCase(mockRepo)
		board, err := service.GetBoardById(expectedBoard.ID, issuerId)

		assert.Error(t, err)
		assert.Nil(t, board)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateBoard(t *testing.T) {
	t.Run("Success Create Board", func(t *testing.T) {
		data := &dto.BoardRequest{
			Name: "Board 1",
			Desc: "Board 1 Description",
		}
		issuerId := uint(3)
		dataModel := &model.Board{
			Name:    data.Name,
			Desc:    data.Desc,
			OwnerID: issuerId,
		}

		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("Create", dataModel).Return(nil).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.CreateBoard(issuerId, data)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Create Board (Missing board name)", func(t *testing.T) {
		data := &dto.BoardRequest{
			Name: "",
			Desc: "Board 1 Description",
		}
		issuerId := uint(3)

		mockRepo := boardRepo.NewMockBoardRepo()

		service := NewBoardUseCase(mockRepo)
		err := service.CreateBoard(issuerId, data)

		assert.Error(t, err)
	})

	t.Run("Failed Create Board (Missing board desc)", func(t *testing.T) {
		data := &dto.BoardRequest{
			Name: "Board 1",
			Desc: "",
		}
		issuerId := uint(3)

		mockRepo := boardRepo.NewMockBoardRepo()

		service := NewBoardUseCase(mockRepo)
		err := service.CreateBoard(issuerId, data)

		assert.Error(t, err)
	})

	t.Run("Failed Create Board (Internal Server Error)", func(t *testing.T) {
		data := &dto.BoardRequest{
			Name: "Board 1",
			Desc: "Board 1 Description",
		}
		issuerId := uint(3)
		dataModel := &model.Board{
			Name:    data.Name,
			Desc:    data.Desc,
			OwnerID: issuerId,
		}

		expectedErr := errors.New("Database Error")
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("Create", dataModel).Return(expectedErr).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.CreateBoard(issuerId, data)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateBoard(t *testing.T) {
	boardToUpdate := &model.Board{
		Model:   gorm.Model{ID: 1},
		Name:    "Board 1",
		Desc:    "Board 1 Description",
		OwnerID: 3,
	}

	t.Run("Success Update Board Name", func(t *testing.T) {
		boardRequest := &dto.BoardRequest{
			Name: "Board 2",
		}
		issuerId := uint(3)

		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("GetBoardMembers", boardToUpdate.ID).Return(mockMembersData, nil).Once()
		mockRepo.On("Update", boardToUpdate.ID, issuerId, boardRequest).Return(nil).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.UpdateBoard(boardToUpdate.ID, issuerId, boardRequest)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success Update Board Desc", func(t *testing.T) {
		boardRequest := &dto.BoardRequest{
			Desc: "New Description",
		}
		issuerId := uint(3)

		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("GetBoardMembers", boardToUpdate.ID).Return(mockMembersData, nil).Once()
		mockRepo.On("Update", boardToUpdate.ID, issuerId, boardRequest).Return(nil).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.UpdateBoard(boardToUpdate.ID, issuerId, boardRequest)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Update Board (isMember false)", func(t *testing.T) {
		boardRequest := &dto.BoardRequest{
			Desc: "New Description",
		}
		issuerId := uint(10)
		expectedErr := errors.New("User is not member of this board!")

		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("GetBoardMembers", boardToUpdate.ID).Return(mockMembersData, nil).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.UpdateBoard(boardToUpdate.ID, issuerId, boardRequest)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Update Board (Internal Server Error)", func(t *testing.T) {
		boardRequest := &dto.BoardRequest{
			Name: "Board 2",
		}
		issuerId := uint(3)

		expectedErr := errors.New("Database Error")
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("GetBoardMembers", boardToUpdate.ID).Return(mockMembersData, nil).Once()
		mockRepo.On("Update", boardToUpdate.ID, issuerId, boardRequest).Return(expectedErr).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.UpdateBoard(boardToUpdate.ID, issuerId, boardRequest)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteBoard(t *testing.T) {
	boardToDelete := &model.Board{
		Model:   gorm.Model{ID: 1},
		Name:    "Board 1",
		Desc:    "Board 1 Description",
		OwnerID: 3,
	}
	issuerId := uint(3)

	t.Run("Success delete board", func(t *testing.T) {
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("GetById", boardToDelete.ID, issuerId).Return(boardToDelete, nil).Once()
		mockRepo.On("Delete", boardToDelete.ID).Return(nil).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.DeleteBoard(boardToDelete.ID, boardToDelete.OwnerID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Failed delete board (OwnerID != userID)", func(t *testing.T) {
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("GetById", boardToDelete.ID, uint(10)).Return(boardToDelete, nil).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.DeleteBoard(boardToDelete.ID, uint(10))

		assert.Error(t, err)
		assert.NotEqual(t, uint(10), boardToDelete.OwnerID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed delete board (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Database Error")
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("GetById", boardToDelete.ID, issuerId).Return(nil, expectedErr).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.DeleteBoard(boardToDelete.ID, boardToDelete.OwnerID)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed delete board (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Database Error")
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("GetById", boardToDelete.ID, issuerId).Return(boardToDelete, nil).Once()
		mockRepo.On("Delete", boardToDelete.ID).Return(expectedErr).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.DeleteBoard(boardToDelete.ID, boardToDelete.OwnerID)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
