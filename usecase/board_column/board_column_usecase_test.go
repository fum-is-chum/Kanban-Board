package usecase

import (
	"errors"
	"kanban-board/dto"
	"kanban-board/model"
	repo "kanban-board/repository/board_column"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var mockColumnsData = []model.BoardColumn{
	{
		Model:   gorm.Model{ID: 1},
		Label:   "Column 1",
		Desc:    "Column 1 Desc",
		BoardID: uint(1),
	},
	{
		Model:   gorm.Model{ID: 2},
		Label:   "Column 2",
		Desc:    "Column 2 Desc",
		BoardID: uint(1),
	},
}

func TestGetColumns(t *testing.T) {
	t.Run("Success Get Board Columns", func(t *testing.T) {
		mockRepo := repo.NewMockBoardColumnRepo()
		mockRepo.On("Get", uint(1)).Return(mockColumnsData, nil).Once()

		service := NewBoardColumnUseCase(mockRepo)
		columns, err := service.GetColumns(uint(1))

		assert.NoError(t, err)
		assert.NotNil(t, columns)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Get Board Columns (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")
		mockRepo := repo.NewMockBoardColumnRepo()
		mockRepo.On("Get", uint(1)).Return(nil, expectedErr).Once()

		service := NewBoardColumnUseCase(mockRepo)
		columns, err := service.GetColumns(uint(1))

		assert.Error(t, err)
		assert.Nil(t, columns)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetColumnById(t *testing.T) {
	t.Run("Success Get Board Column by Id", func(t *testing.T) {
		mockRepo := repo.NewMockBoardColumnRepo()
		mockRepo.On("GetById", uint(1)).Return(&mockColumnsData[0], nil).Once()

		service := NewBoardColumnUseCase(mockRepo)
		column, err := service.GetColumnById(uint(1))

		assert.NoError(t, err)
		assert.NotNil(t, column)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Get Board Column By Id (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")
		mockRepo := repo.NewMockBoardColumnRepo()
		mockRepo.On("GetById", uint(1)).Return(nil, expectedErr).Once()

		service := NewBoardColumnUseCase(mockRepo)
		column, err := service.GetColumnById(uint(1))

		assert.Error(t, err)
		assert.Nil(t, column)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateColumn(t *testing.T) {
	t.Run("Success Create New Column", func(t *testing.T) {
		mockRequest := &dto.BoardColumnRequest{
			Label:   "Column 1",
			Desc:    "Column 1 description",
			BoardID: uint(1),
		}

		mockModel := &model.BoardColumn{
			Label:   mockRequest.Label,
			Desc:    mockRequest.Desc,
			BoardID: mockRequest.BoardID,
		}

		mockRepo := repo.NewMockBoardColumnRepo()
		mockRepo.On("Create", mockModel).Return(nil).Once()

		service := NewBoardColumnUseCase(mockRepo)
		err := service.CreateColumn(mockRequest)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Create New Column (Invalid struct)", func(t *testing.T) {
		mockRequest := &dto.BoardColumnRequest{
			Label: "Column 1",
			Desc:  "Column 1 description",
		}

		mockRepo := repo.NewMockBoardColumnRepo()
		service := NewBoardColumnUseCase(mockRepo)
		err := service.CreateColumn(mockRequest)

		assert.Error(t, err)
	})

	t.Run("Failed Create New Column (Internal Server Error)", func(t *testing.T) {
		mockRequest := &dto.BoardColumnRequest{
			Label:   "Column 1",
			Desc:    "Column 1 description",
			BoardID: uint(1),
		}

		mockModel := &model.BoardColumn{
			Label:   mockRequest.Label,
			Desc:    mockRequest.Desc,
			BoardID: mockRequest.BoardID,
		}
		expectedErr := errors.New("Internal Server Error")

		mockRepo := repo.NewMockBoardColumnRepo()
		mockRepo.On("Create", mockModel).Return(expectedErr).Once()

		service := NewBoardColumnUseCase(mockRepo)
		err := service.CreateColumn(mockRequest)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateColumn(t *testing.T) {
	t.Run("Success Update Column", func(t *testing.T) {
		mockRequest := &dto.BoardColumnRequest{
			Label: "Column 1",
			Desc:  "Column 1 description",
		}

		mockRepo := repo.NewMockBoardColumnRepo()
		mockRepo.On("Update", uint(1), mockRequest).Return(nil).Once()

		service := NewBoardColumnUseCase(mockRepo)
		err := service.UpdateColumn(uint(1), mockRequest)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Update Column (Internal Server Error)", func(t *testing.T) {
		mockRequest := &dto.BoardColumnRequest{
			Label: "Column 1",
			Desc:  "Column 1 description",
		}
		expectedErr := errors.New("Internal Server Error")

		mockRepo := repo.NewMockBoardColumnRepo()
		mockRepo.On("Update", uint(1), mockRequest).Return(expectedErr).Once()

		service := NewBoardColumnUseCase(mockRepo)
		err := service.UpdateColumn(uint(1), mockRequest)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteColumn(t *testing.T) {
	t.Run("Success Delete Board Column", func(t *testing.T) {
		columnId := uint(1)
		issuerId := uint(1)

		mockRepo := repo.NewMockBoardColumnRepo()
		mockRepo.On("Delete", columnId, issuerId).Return(nil).Once()

		service := NewBoardColumnUseCase(mockRepo)
		err := service.DeleteColumn(columnId, issuerId)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Delete Board Column (Internal Server Error)", func(t *testing.T) {
		columnId := uint(1)
		issuerId := uint(1)
		expectedErr := errors.New("Internal Server Error")

		mockRepo := repo.NewMockBoardColumnRepo()
		mockRepo.On("Delete", columnId, issuerId).Return(expectedErr).Once()

		service := NewBoardColumnUseCase(mockRepo)
		err := service.DeleteColumn(columnId, issuerId)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
