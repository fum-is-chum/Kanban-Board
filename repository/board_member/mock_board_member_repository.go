package repository

import "github.com/stretchr/testify/mock"

type mockBoardMemberRepository struct {
	mock.Mock
}

func NewMockBoardMemberRepo() *mockBoardMemberRepository {
	return &mockBoardMemberRepository{}
}

func (m *mockBoardMemberRepository) AddMember(boardId uint, userId uint) error {
	ret := m.Called(boardId, userId)
	return ret.Error(0)
}

func (m *mockBoardMemberRepository) DeleteMember(boardId uint, userId uint) error {
	ret := m.Called(boardId, userId)
	return ret.Error(0)
}
