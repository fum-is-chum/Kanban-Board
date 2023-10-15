package repository

import (
	"kanban-board/dto"
	"kanban-board/model"

	"gorm.io/gorm"
)

type BoardRepository interface {
	Get() ([]model.Board, error)
	GetById(id uint) (*model.Board, error)
	Create(data *model.Board) error
	Update(id uint, data *dto.BoardRequest) error
	Delete(id uint) error
}

type boardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) *boardRepository {
	return &boardRepository{db}
}

func (b *boardRepository) Get() ([]model.Board, error) {
	var boards []model.Board

	tx := b.db.Preload("Owner").Order("created_at desc").Find(&boards)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return boards, nil
}

func (b *boardRepository) GetById(id uint) (*model.Board, error) {
	var board model.Board

	tx := b.db.Preload("Owner").Preload("Members").Where("id = ?", id).First(&board)
	if tx.Error != nil {
		return nil, tx.Error
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

func (b *boardRepository) Update(id uint, data *dto.BoardRequest) error {
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
