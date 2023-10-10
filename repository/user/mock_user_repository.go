package repository

import (
	"kanban-board/model"

	mock "github.com/stretchr/testify/mock"
)

type mockUserRepository struct {
	mock.Mock
}

func NewMockUserRepo() *mockUserRepository {
	return &mockUserRepository{}
}

func (m *mockUserRepository) Get() ([]model.User, error) {
	ret := m.Called()
	if users, ok := ret.Get(0).([]model.User); ok {
		return users, ret.Error(1)
	}

	return nil, ret.Error(1)
}

func (m *mockUserRepository) Create(data *model.User) error {
	ret := m.Called(data)
	return ret.Error(0)
}

func (m *mockUserRepository) Update(id uint, data *map[string]interface{}) error {
	ret := m.Called(id, data)
	return ret.Error(0)
}

func (m *mockUserRepository) GetByEmail(email string) (*model.User, error) {
	ret := m.Called(email)
	if user, ok := ret.Get(0).(*model.User); ok {
		return user, ret.Error(1)
	}

	return nil, ret.Error(1)
}

func (m *mockUserRepository) GetById(id uint) (*model.User, error) {
	ret := m.Called(id)
	if user, ok := ret.Get(0).(*model.User); ok {
		return user, ret.Error(1)
	}
	return nil, ret.Error(1)
}

func (m *mockUserRepository) Delete(id uint) error {
	ret := m.Called(id)
	return ret.Error(0)
}
