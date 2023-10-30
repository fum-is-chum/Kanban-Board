package repository

import (
	"kanban-board/dto"
	"kanban-board/model"

	"github.com/stretchr/testify/mock"
)

type mockTaskRepository struct {
	mock.Mock
}

func NewMockTaskRepo() *mockTaskRepository {
	return &mockTaskRepository{}
}

func (m *mockTaskRepository) Get(boardId uint) ([]model.Task, error) {
	ret := m.Called(boardId)
	if tasks, ok := ret.Get(0).([]model.Task); ok {
		return tasks, ret.Error(1)
	}

	return nil, ret.Error(1)
}

func (m *mockTaskRepository) GetById(id uint, issuerId uint) (*model.Task, error) {
	ret := m.Called(id, issuerId)
	if task, ok := ret.Get(0).(*model.Task); ok {
		return task, ret.Error(1)
	}

	return nil, ret.Error(1)
}

func (m *mockTaskRepository) Create(issuerId uint, data *model.Task) error {
	ret := m.Called(issuerId, data)
	return ret.Error(0)
}

func (m *mockTaskRepository) Update(id uint, issuerId uint, data *dto.TaskUpdateRequest) error {
	ret := m.Called(id, issuerId, data)
	return ret.Error(0)
}

func (m *mockTaskRepository) Delete(id uint, issuerId uint) error {
	ret := m.Called(id, issuerId)
	return ret.Error(0)
}

func (m *mockTaskRepository) GetBoardIdByColumnId(columnId uint) (*uint, error) {
	ret := m.Called(columnId)
	if boardId, ok := ret.Get(0).(*uint); ok {
		return boardId, ret.Error(1)
	}

	return nil, ret.Error(1)
}

func (m *mockTaskRepository) GetBoardIdByTaskId(taskId uint) (*uint, error) {
	ret := m.Called(taskId)
	if boardId, ok := ret.Get(0).(*uint); ok {
		return boardId, ret.Error(1)
	}

	return nil, ret.Error(1)
}