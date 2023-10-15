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

func TestGetBoards(t *testing.T) {
	returnData := []model.Board{
		{Name: "board1", Desc: "Board number 1", OwnerID: 0},
		{Name: "board2", Desc: "Board number 2", OwnerID: 0},
	}

	t.Run("Success Get Boards", func(t *testing.T) {
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("Get", mock.Anything).Return(returnData, nil).Once()

		service := NewBoardUseCase(mockRepo)
		res, err := service.GetBoards()

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
		res, err := service.GetBoards()

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

	t.Run("Success Get Board By Id", func(t *testing.T) {
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("GetById", expectedBoard.ID).Return(expectedBoard, nil).Once()

		service := NewBoardUseCase(mockRepo)
		board, err := service.GetBoardById(expectedBoard.ID)

		assert.NoError(t, err)
		assert.NotNil(t, board)
		assert.Equal(t, board.ID, expectedBoard.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Get Board By Id (Id not found)", func(t *testing.T) {
		mockRepo := boardRepo.NewMockBoardRepo()
		expectedErr := errors.New("ID not found")
		mockRepo.On("GetById", mock.Anything).Return(nil, expectedErr).Once()

		service := NewBoardUseCase(mockRepo)
		board, err := service.GetBoardById(uint(10))

		assert.Error(t, err)
		assert.Nil(t, board)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Get Board By Id (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Database Error")
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("GetById", expectedBoard.ID).Return(nil, expectedErr).Once()

		service := NewBoardUseCase(mockRepo)
		board, err := service.GetBoardById(expectedBoard.ID)

		assert.Error(t, err)
		assert.Nil(t, board)
		mockRepo.AssertExpectations(t)
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

		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("Create", dataModel).Return(nil).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.CreateBoard(data)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Create Board (Missing board name)", func(t *testing.T) {
		data := &dto.BoardRequest{
			Name:    "",
			Desc:    "Board 1 Description",
			OwnerID: 3,
		}

		mockRepo := boardRepo.NewMockBoardRepo()

		service := NewBoardUseCase(mockRepo)
		err := service.CreateBoard(data)

		assert.Error(t, err)
	})

	t.Run("Failed Create Board (Missing board desc)", func(t *testing.T) {
		data := &dto.BoardRequest{
			Name:    "Board 1",
			Desc:    "",
			OwnerID: 3,
		}

		mockRepo := boardRepo.NewMockBoardRepo()

		service := NewBoardUseCase(mockRepo)
		err := service.CreateBoard(data)

		assert.Error(t, err)
	})

	t.Run("Failed Create Board (Missing OwnerId)", func(t *testing.T) {
		data := &dto.BoardRequest{
			Name: "Board 1",
			Desc: "Board 1 Desc",
		}

		mockRepo := boardRepo.NewMockBoardRepo()

		service := NewBoardUseCase(mockRepo)
		err := service.CreateBoard(data)

		assert.Error(t, err)
	})

	t.Run("Failed Create Board (Internal Server Error)", func(t *testing.T) {
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

		expectedErr := errors.New("Database Error")
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("Create", dataModel).Return(expectedErr).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.CreateBoard(data)

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

		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("Update", boardToUpdate.ID, boardRequest).Return(nil).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.UpdateBoard(boardToUpdate.ID, boardRequest)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success Update Board Desc", func(t *testing.T) {
		boardRequest := &dto.BoardRequest{
			Desc: "New Description",
		}

		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("Update", boardToUpdate.ID, boardRequest).Return(nil).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.UpdateBoard(boardToUpdate.ID, boardRequest)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Update Board (Internal Server Error)", func(t *testing.T) {
		boardRequest := &dto.BoardRequest{
			Name: "Board 2",
		}

		expectedErr := errors.New("Database Error")
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("Update", boardToUpdate.ID, boardRequest).Return(expectedErr).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.UpdateBoard(boardToUpdate.ID, boardRequest)

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

	t.Run("Success delete board", func(t *testing.T) {
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("GetById", boardToDelete.ID).Return(boardToDelete, nil).Once()
		mockRepo.On("Delete", boardToDelete.ID).Return(nil).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.DeleteBoard(boardToDelete.ID, boardToDelete.OwnerID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed delete board (boardId not found)", func(t *testing.T) {
		expectedErr := errors.New("BoardID not found")
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("GetById", boardToDelete.ID).Return(nil, expectedErr).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.DeleteBoard(boardToDelete.ID, boardToDelete.OwnerID)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed delete board (OwnerID != userID)", func(t *testing.T) {
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("GetById", boardToDelete.ID).Return(boardToDelete, nil).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.DeleteBoard(boardToDelete.ID, uint(10))

		assert.Error(t, err)
		assert.NotEqual(t, uint(10), boardToDelete.OwnerID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed delete board (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Database Error")
		mockRepo := boardRepo.NewMockBoardRepo()
		mockRepo.On("GetById", boardToDelete.ID).Return(boardToDelete, nil).Once()
		mockRepo.On("Delete", boardToDelete.ID).Return(expectedErr).Once()

		service := NewBoardUseCase(mockRepo)
		err := service.DeleteBoard(boardToDelete.ID, boardToDelete.OwnerID)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
