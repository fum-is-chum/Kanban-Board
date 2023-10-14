package repository

import (
	"kanban-board/model"

	"gorm.io/gorm"
)

type BoardMemberRepository interface {
	AddMember(boardId uint, userId uint) error
	DeleteMember(boardId uint, userId uint) error
}

type boardMemberRepository struct {
	db *gorm.DB
}

func NewBoardMemberRepostory(db *gorm.DB) *boardMemberRepository {
	return &boardMemberRepository{db}
}

func (b *boardMemberRepository) AddMember(boardId uint, userId uint) error {
	member := &model.BoardMember{
		BoardID: boardId,
		UserID:  userId,
	}

	if err := b.db.Create(member).Error; err != nil {
		return err
	}

	return nil
}

func (b *boardMemberRepository) DeleteMember(boardId uint, userId uint) error {
	if err := b.db.Unscoped().Where("board_id = ? AND user_id = ?", boardId, userId).Delete(&model.BoardMember{}).Error; err != nil {
		return err
	}

	return nil
}
