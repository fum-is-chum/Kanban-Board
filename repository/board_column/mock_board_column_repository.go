package repository

import (
	"kanban-board/dto"
	"kanban-board/model"

	"github.com/stretchr/testify/mock"
)

type mockBoardColumnRepository struct {
	mock.Mock
}

func NewMockBoardColumnRepo() *mockBoardColumnRepository {
	return &mockBoardColumnRepository{}
}

func (m *mockBoardColumnRepository) Get(boardId uint) ([]model.BoardColumn, error) {
	ret := m.Called(boardId)
	if columns, ok := ret.Get(0).([]model.BoardColumn); ok {
		return columns, ret.Error(1)
	}

	return nil, ret.Error(1)
}

func (m *mockBoardColumnRepository) GetById(id uint) (*model.BoardColumn, error) {
	ret := m.Called(id)
	if column, ok := ret.Get(0).(*model.BoardColumn); ok {
		return column, nil
	}

	return nil, ret.Error(1)
}

func (m *mockBoardColumnRepository) Create(data *model.BoardColumn) error {
	ret := m.Called(data)
	return ret.Error(0)
}

func (m *mockBoardColumnRepository) Update(id uint, data *dto.BoardColumnRequest) error {
	ret := m.Called(id, data)
	return ret.Error(0)
}

func (m *mockBoardColumnRepository) Delete(id uint, issuerId uint) error {
	ret := m.Called(id, issuerId)
	return ret.Error(0)
}
