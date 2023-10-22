package repository

import (
	"errors"
	"kanban-board/dto"
	"kanban-board/model"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Get(board_id uint) ([]model.Task, error)
	GetById(id uint, issuerId uint) (*model.Task, error)
	Create(issuerId uint, data *model.Task) error
	Update(id uint, issuerId uint, data *dto.TaskUpdateRequest) error
	Delete(id uint, issuerId uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *taskRepository {
	return &taskRepository{db}
}

func (t *taskRepository) Get(board_id uint) ([]model.Task, error) {
	var tasks []model.Task
	if err := t.db.Where("board_id = ?", board_id).Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *taskRepository) GetById(id uint, issuerId uint) (*model.Task, error) {
	var task model.Task
	if err := t.db.Where("id = ?", id).First(&task).Error; err != nil {
		return nil, err
	}

	// check if issuer is member of the board
	var members []model.BoardMember
	if err := t.db.Table("board_members").Where("board_id = ?", task.BoardID).Find(&members).Error; err != nil {
		return nil, err
	}

	var issuerIsMember bool
	for _, member := range members {
		if member.UserID == issuerId {
			issuerIsMember = true
			break
		}
	}

	if !issuerIsMember {
		return nil, errors.New("Issuer is not member of this board!")
	}

	return &task, nil
}

func (t *taskRepository) Create(issuerId uint, data *model.Task) error {
	var board model.Board
	if err := t.db.Preload("Members").Where("id = ?", data.BoardID).Find(&board).Error; err != nil {
		return nil
	}

	var issuerIsMember bool
	for _, member := range board.Members {
		if member.ID == issuerId {
			issuerIsMember = true
			break
		}
	}

	if !issuerIsMember {
		return errors.New("Issuer is not member of this board!")
	}

	if err := t.db.Create(data).Error; err != nil {
		return err
	}

	return nil
}

func (t *taskRepository) Update(id uint, issuerId uint, data *dto.TaskUpdateRequest) error {
	// check member status
	var column model.BoardColumn
	var board model.Board
	if err := t.db.First(&column, data.BoardColumnID).Error; err != nil {
		return err
	}

	if err := t.db.Preload("Members").Where("id = ?", column.BoardID).Find(&board).Error; err != nil {
		return nil
	}

	var issuerIsMember bool
	for _, member := range board.Members {
		if member.ID == issuerId {
			issuerIsMember = true
			break
		}
	}

	if !issuerIsMember {
		return errors.New("Issuer is not member of this board!")
	}

	// update
	if err := t.db.Model(&model.Task{}).Where("id =  ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func (t *taskRepository) Delete(id uint, issuerId uint) error {
	var column model.BoardColumn
	var board model.Board

	if err := t.db.Where(&model.BoardColumn{}, id).First(&column).Error; err != nil {
		return err
	}

	if err := t.db.Preload("Members").Where("id = ?", column.BoardID).Find(&board).Error; err != nil {
		return nil
	}

	var issuerIsMember bool
	for _, member := range board.Members {
		if member.ID == issuerId {
			issuerIsMember = true
			break
		}
	}

	if !issuerIsMember {
		return errors.New("Issuer is not member of this board!")
	}

	if err := t.db.Unscoped().Delete(&model.Task{}, id).Error; err != nil {
		return err
	}

	return nil
}
