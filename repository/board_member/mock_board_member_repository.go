package repository

import (
	"kanban-board/model"

	"github.com/stretchr/testify/mock"
)

type mockBoardMemberRepository struct {
	mock.Mock
}

func NewMockBoardMemberRepo() *mockBoardMemberRepository {
	return &mockBoardMemberRepository{}
}

func (m *mockBoardMemberRepository) GetBoardOwner(boardId uint) (*uint, error) {
	ret := m.Called(boardId)
	if ownerId, ok := ret.Get(0).(*uint); ok {
		return ownerId, ret.Error(1)
	}
	return nil, ret.Error(1)
}

func (m *mockBoardMemberRepository) GetBoardMembers(boardId uint) ([]model.BoardMember, error) {
	ret := m.Called(boardId)
	if members, ok := ret.Get(0).([]model.BoardMember); ok {
		return members, ret.Error(1)
	}
	return nil, ret.Error(1)
}

func (m *mockBoardMemberRepository) AddMember(boardId uint, userId uint) error {
	ret := m.Called(boardId, userId)
	return ret.Error(0)
}

func (m *mockBoardMemberRepository) DeleteMember(boardId uint, userId uint) error {
	ret := m.Called(boardId, userId)
	return ret.Error(0)
}
