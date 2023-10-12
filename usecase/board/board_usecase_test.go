package usecase

import (
	"errors"
	"kanban-board/dto"
	"kanban-board/model"
	boardRepo "kanban-board/repository/board"
	userRepo "kanban-board/repository/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var mockUserRepo = userRepo.NewMockUserRepo()

func TestGetBoards(t *testing.T) {
	returnData := []model.Board{
		{Name: "board1", Desc: "Board number 1", OwnerID: 0},
		{Name: "board2", Desc: "Board number 2", OwnerID: 0},
	}

	t.Run("Success Get Boards", func(t *testing.T) {
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockBoardRepo.On("Get", mock.Anything).Return(returnData, nil).Once()

		service := NewBoardUseCase(mockUserRepo, mockBoardRepo)
		res, err := service.GetBoards()

		assert.NoError(t, err)
		assert.Equal(t, len(returnData), len(res))
		assert.Equal(t, returnData[0].Name, res[0].Name)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Get Boards", func(t *testing.T) {
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		expectedErr := errors.New("Internal Server Error")
		mockBoardRepo.On("Get", mock.Anything).Return(nil, expectedErr).Once()

		service := NewBoardUseCase(mockUserRepo, mockBoardRepo)
		res, err := service.GetBoards()

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, expectedErr, err)
		mockBoardRepo.AssertExpectations(t)
	})
}

func TestGetBoardById(t *testing.T) {
	expectedBoard := &model.Board{
		Model: gorm.Model{ID: 1},
		Name:  "Board 1",
		Desc:  "Board 1 Desc",
	}

	t.Run("Success Get Board By Id", func(t *testing.T) {
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockBoardRepo.On("GetById", expectedBoard.ID).Return(expectedBoard, nil).Once()

		service := NewBoardUseCase(mockUserRepo, mockBoardRepo)
		board, err := service.GetBoardById(expectedBoard.ID)

		assert.NoError(t, err)
		assert.NotNil(t, board)
		assert.Equal(t, board.ID, expectedBoard.ID)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Get Board By Id (Id not found)", func(t *testing.T) {
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		expectedErr := errors.New("ID not found")
		mockBoardRepo.On("GetById", mock.Anything).Return(nil, expectedErr).Once()

		service := NewBoardUseCase(mockUserRepo, mockBoardRepo)
		board, err := service.GetBoardById(uint(10))

		assert.Error(t, err)
		assert.Nil(t, board)
		mockBoardRepo.AssertExpectations(t)
	})
}

func TestCreateBoard(t *testing.T) {
	t.Run("Success Create Board", func(t *testing.T) {
		data := &dto.BoardRequest{
			Name:    "Board 1",
			Desc:    "Board 1 Description",
			OwnerID: 3,
		}

		dataModel := &model.Board{
			Name:    data.Name,
			Desc:    data.Desc,
			OwnerID: data.OwnerID,
		}

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockBoardRepo.On("Create", dataModel).Return(nil).Once()

		service := NewBoardUseCase(mockUserRepo, mockBoardRepo)
		err := service.CreateBoard(data)

		assert.NoError(t, err)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Create Board (Missing board name)", func(t *testing.T) {
		data := &dto.BoardRequest{
			Name:    "",
			Desc:    "Board 1 Description",
			OwnerID: 3,
		}

		mockBoardRepo := boardRepo.NewMockBoardRepo()

		service := NewBoardUseCase(mockUserRepo, mockBoardRepo)
		err := service.CreateBoard(data)

		assert.Error(t, err)
	})

	t.Run("Failed Create Board (Missing OwnerId)", func(t *testing.T) {
		data := &dto.BoardRequest{
			Name: "Board 1",
			Desc: "Board 1 Desc",
		}

		mockBoardRepo := boardRepo.NewMockBoardRepo()

		service := NewBoardUseCase(mockUserRepo, mockBoardRepo)
		err := service.CreateBoard(data)

		assert.Error(t, err)
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

		updateMap := &map[string]interface{}{
			"Name": "Board 2",
		}

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockBoardRepo.On("Update", boardToUpdate.ID, updateMap).Return(nil).Once()

		service := NewBoardUseCase(mockUserRepo, mockBoardRepo)
		err := service.UpdateBoard(boardToUpdate.ID, boardRequest)

		assert.NoError(t, err)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Update Board Name (board name empty)", func(t *testing.T) {
		boardRequest := &dto.BoardRequest{
			Name: "",
		}

		mockBoardRepo := boardRepo.NewMockBoardRepo()

		service := NewBoardUseCase(mockUserRepo, mockBoardRepo)
		err := service.UpdateBoard(boardToUpdate.ID, boardRequest)

		assert.Error(t, err)
	})

	t.Run("Success Update Board Desc", func(t *testing.T) {
		boardRequest := &dto.BoardRequest{
			Desc: "New Description",
		}

		updateMap := &map[string]interface{}{
			"Desc": "New Description",
		}

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockBoardRepo.On("Update", boardToUpdate.ID, updateMap).Return(nil).Once()

		service := NewBoardUseCase(mockUserRepo, mockBoardRepo)
		err := service.UpdateBoard(boardToUpdate.ID, boardRequest)

		assert.NoError(t, err)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Update Board Name (board Desc empty)", func(t *testing.T) {
		boardRequest := &dto.BoardRequest{
			Desc: "",
		}

		mockBoardRepo := boardRepo.NewMockBoardRepo()

		service := NewBoardUseCase(mockUserRepo, mockBoardRepo)
		err := service.UpdateBoard(boardToUpdate.ID, boardRequest)

		assert.Error(t, err)
	})
}

func TestUpdateBoardOwnership(t *testing.T) {
	boardToUpdate := &model.Board{
		Model:   gorm.Model{ID: 1},
		Name:    "Board 1",
		Desc:    "Board 1 Description",
		OwnerID: 3,
	}

	expectedUser := &model.User{
		Name:     "Alvin Christ Ardiansyah",
		Email:    "alvinardiansyah2002@gmail.com",
		Password: "123",
	}

	t.Run("Success Update Board Ownership", func(t *testing.T) {
		boardRequest := &dto.BoardRequest{
			OwnerID: uint(3),
		}

		updateMap := &map[string]interface{}{
			"OwnerID": uint(3),
		}

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockBoardRepo.On("Update", boardToUpdate.ID, updateMap).Return(nil).Once()
		mockBoardRepo.On("GetById", boardToUpdate.ID).Return(boardToUpdate, nil).Once()
		mockUserRepo.On("GetById", boardRequest.OwnerID).Return(expectedUser, nil).Once()

		service := NewBoardUseCase(mockUserRepo, mockBoardRepo)
		err := service.UpdateBoardOwnership(boardToUpdate.ID, boardRequest.OwnerID)

		assert.NoError(t, err)
		mockBoardRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Failed Update Board Ownership (BoardID not found)", func(t *testing.T) {
		boardRequest := &dto.BoardRequest{
			OwnerID: uint(3),
		}

		expectedError := errors.New("BoardID not found")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockBoardRepo.On("GetById", boardToUpdate.ID).Return(nil, expectedError).Once()

		service := NewBoardUseCase(mockUserRepo, mockBoardRepo)
		err := service.UpdateBoardOwnership(boardToUpdate.ID, boardRequest.OwnerID)

		assert.Error(t, err)
	})

	t.Run("Failed Update Board Ownership (OwnerID not found)", func(t *testing.T) {
		boardRequest := &dto.BoardRequest{
			OwnerID: uint(3),
		}

		expectedError := errors.New("UserID not found")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockBoardRepo.On("GetById", boardToUpdate.ID).Return(boardToUpdate, nil).Once()
		mockUserRepo.On("GetById", boardRequest.OwnerID).Return(nil, expectedError).Once()

		service := NewBoardUseCase(mockUserRepo, mockBoardRepo)
		err := service.UpdateBoardOwnership(boardToUpdate.ID, boardRequest.OwnerID)

		assert.Error(t, err)
	})
}

func TestDeleteBoard(t *testing.T) {
	t.Run("Success delete board", func(t *testing.T) {
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockBoardRepo.On("Delete", uint(3)).Return(nil).Once()

		service := NewBoardUseCase(mockUserRepo, mockBoardRepo)
		err := service.DeleteBoard(uint(3))

		assert.NoError(t, err)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed delete board (boardId not found)", func(t *testing.T) {
		expectedErr := errors.New("BoardID not found")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockBoardRepo.On("Delete", uint(3)).Return(expectedErr).Once()

		service := NewBoardUseCase(mockUserRepo, mockBoardRepo)
		err := service.DeleteBoard(uint(3))

		assert.Error(t, err)
		mockBoardRepo.AssertExpectations(t)

	})
}
