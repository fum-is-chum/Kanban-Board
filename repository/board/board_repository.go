package repository

import (
	"kanban-board/dto"
	"kanban-board/model"

	"gorm.io/gorm"
)

type BoardRepository interface {
	Get(issuerId uint) ([]model.Board, error)
	GetById(id uint, issuerId uint) (*model.Board, error)
	Create(data *model.Board) error
	Update(id uint, issuerId uint, data *dto.BoardRequest) error
	Delete(id uint) error
	GetBoardOwner(id uint) (*uint, error)
	GetBoardMembers(id uint) ([]model.BoardMember, error)
}

type boardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) *boardRepository {
	return &boardRepository{db}
}

func (b *boardRepository) Get(issuerId uint) ([]model.Board, error) {
	var members []model.BoardMember
	var boardIDs []uint
	var boards []model.Board

	if err := b.db.Where("user_id = ?", issuerId).Find(&members).Error; err != nil {
		return nil, err
	}

	for _, member := range members {
		boardIDs = append(boardIDs, member.BoardID)
	}

	if err := b.db.Where("id IN (?)", boardIDs).Preload("Owner").Order("created_at desc").Find(&boards).Error; err != nil {
		return nil, err
	}

	return boards, nil
}

func (b *boardRepository) GetById(id uint, issuerId uint) (*model.Board, error) {
	// check if issuer is member of the board
	var member model.BoardMember
	if err := b.db.Where("board_id = ? AND user_id = ?", id, issuerId).First(&member).Error; err != nil {
		return nil, err
	}

	var board model.Board
	if err := b.db.Where("id = ?", id).Preload("Owner").Preload("Members").Preload("Columns.Tasks.Assignees").First(&board).Error; err != nil {
		return nil, err
	}

	return &board, nil
}

func (b *boardRepository) Create(data *model.Board) error {
	if err := b.db.Create(&data).Error; err != nil {
		return err
	}

	// add user as member
	var user model.User
	var board model.Board

	if err := b.db.First(&user, data.OwnerID).Error; err != nil {
		return err
	}

	if err := b.db.First(&board, data.ID).Error; err != nil {
		return err
	}

	if err := b.db.Model(&board).Association("Members").Append(&user); err != nil {
		return err
	}

	return nil
}

func (b *boardRepository) Update(id uint, issuerId uint, data *dto.BoardRequest) error {
	tx := b.db.Model(&model.Board{}).Where("id = ?", id).Updates(&data)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (b *boardRepository) Delete(id uint) error {
	// Delete all board member first
	if err := b.db.Model(&model.Board{Model: gorm.Model{ID: id}}).Association("Members").Clear(); err != nil {
		return err
	}

	if err := b.db.Unscoped().Delete(&model.Board{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (b *boardRepository) GetBoardOwner(id uint) (*uint, error) {
	var board model.Board
	if err := b.db.First(&board, id).Error; err != nil {
		return nil, err
	}

	return &board.OwnerID, nil
}

func (b *boardRepository) GetBoardMembers(id uint) ([]model.BoardMember, error) {
	var boardMembers []model.BoardMember
	if err := b.db.Model(&model.BoardMember{}).Where("board_id = ?", id).Find(&boardMembers).Error; err != nil {
		return nil, err
	}

	return boardMembers, nil
}