package repository

import (
	"kanban-board/model"

	"gorm.io/gorm"
)

type BoardMemberRepository interface {
	GetBoardOwner(boardId uint) (*uint, error)
	GetBoardMembers(boardId uint) ([]model.BoardMember, error)
	AddMember(boardId uint, userId uint) error
	DeleteMember(boardId uint, userId uint) error
}

type boardMemberRepository struct {
	db *gorm.DB
}

func NewBoardMemberRepository(db *gorm.DB) *boardMemberRepository {
	return &boardMemberRepository{db}
}

func (b *boardMemberRepository) GetBoardOwner(boardId uint) (*uint, error) {
	var board model.Board
	if err := b.db.First(&board, boardId).Error; err != nil {
		return nil, err
	}

	return &board.OwnerID, nil
}

func (b *boardMemberRepository) GetBoardMembers(boardId uint) ([]model.BoardMember, error) {
	var boardMembers []model.BoardMember
	if err := b.db.Model(&model.BoardMember{}).Where("board_id = ?", boardId).Find(&boardMembers).Error; err != nil {
		return nil, err
	}

	return boardMembers, nil
}

func (b *boardMemberRepository) AddMember(boardId uint, userId uint) error {
	var user model.User
	var board model.Board

	if err := b.db.First(&user, userId).Error; err != nil {
		return err
	}

	if err := b.db.First(&board, boardId).Error; err != nil {
		return err
	}

	if err := b.db.Model(&board).Association("Members").Append(&user); err != nil {
		return err
	}

	return nil
}

func (b *boardMemberRepository) DeleteMember(boardId uint, userId uint) error {
	var user model.User
	var board model.Board

	if err := b.db.First(&user, userId).Error; err != nil {
		return err
	}

	if err := b.db.First(&board, boardId).Error; err != nil {
		return err
	}

	if err := b.db.Model(&board).Association("Members").Delete(&user); err != nil {
		return err
	}

	return nil
}
