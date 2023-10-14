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

func (m *mockBoardMemberRepository) AddMember(board *model.Board, user *model.User) error {
	ret := m.Called(board, user)
	return ret.Error(0)
}

func (m *mockBoardMemberRepository) DeleteMember(board *model.Board, user *model.User) error {
	ret := m.Called(board, user)
	return ret.Error(0)
}
