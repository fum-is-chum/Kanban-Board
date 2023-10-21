package repository

import (
	"kanban-board/dto"
	"kanban-board/model"

	"github.com/stretchr/testify/mock"
)

type mockBoardRepository struct {
	mock.Mock
}

func NewMockBoardRepo() *mockBoardRepository {
	return &mockBoardRepository{}
}

func (m *mockBoardRepository) Get() ([]model.Board, error) {
	ret := m.Called()
	if boards, ok := ret.Get(0).([]model.Board); ok {
		return boards, ret.Error(1)
	}

	return nil, ret.Error(1)
}

func (m *mockBoardRepository) Create(data *model.Board) error {
	ret := m.Called(data)
	return ret.Error(0)
}

func (m *mockBoardRepository) Update(id uint, issuerId uint, data *dto.BoardRequest) error {
	ret := m.Called(id, issuerId, data)
	return ret.Error(0)
}

func (m *mockBoardRepository) GetById(id uint) (*model.Board, error) {
	ret := m.Called(id)
	if board, ok := ret.Get(0).(*model.Board); ok {
		return board, ret.Error(1)
	}

	return nil, ret.Error(1)
}

func (m *mockBoardRepository) Delete(id uint) error {
	ret := m.Called(id)
	return ret.Error(0)
}
