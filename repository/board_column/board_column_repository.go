package repository

import (
	"errors"
	"kanban-board/dto"
	"kanban-board/model"

	"gorm.io/gorm"
)

type BoardColumnRepository interface {
	Get(boardId uint) ([]model.BoardColumn, error)
	GetById(id uint) (*model.BoardColumn, error)
	Create(data *model.BoardColumn) error
	Update(id uint, data *dto.BoardColumnRequest) error
	Delete(id uint, issuerId uint) error
}

type boardColumnRepository struct {
	db *gorm.DB
}

func NewBoardColumnRepository(db *gorm.DB) *boardColumnRepository {
	return &boardColumnRepository{db}
}

func (b *boardColumnRepository) Get(boardId uint) ([]model.BoardColumn, error) {
	var columns []model.BoardColumn
	
	if err := b.db.Where("board_id = ?", boardId).Find(&columns).Error; err != nil {
		return nil, err
	}
	
	return columns, nil
}

func (b *boardColumnRepository) GetById(id uint) (*model.BoardColumn, error) {
	var column model.BoardColumn

	if err := b.db.Where("id = ?", id).First(&column).Error; err != nil {
		return nil, err
	}

	return &column, nil
}

func (b *boardColumnRepository) Create(data *model.BoardColumn) error {
	if err := b.db.Create(data).Error; err != nil {
		return err
	}

	return nil
}

func (b *boardColumnRepository) Update(id uint, data *dto.BoardColumnRequest) error {
	if err := b.db.Model(&model.BoardColumn{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func (b *boardColumnRepository) Delete(id uint, issuerId uint) error {
	var column model.BoardColumn
	var board model.Board

	// get column
	if err := b.db.Preload("Board").Where("id = ?", id).First(&column).Error; err != nil {
		return err
	}

	board = *column.Board

	// get board members
	if err := b.db.Preload("Members").Where("id = ?", board.ID).First(&board).Error; err != nil {
		return err
	}

	var issuerExist bool

	for _, member := range board.Members {
		if member.ID == issuerId {
			issuerExist = true
			break
		}
	}

	if !issuerExist {
		return errors.New("Issuer is not member of this board!")
	}

	if err := b.db.Unscoped().Delete(&model.BoardColumn{}, id).Error; err != nil {
		return err
	}

	return nil
}