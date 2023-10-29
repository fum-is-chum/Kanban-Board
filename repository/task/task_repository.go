package repository

import (
	"errors"
	"kanban-board/dto"
	"kanban-board/model"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Get(boardId uint) ([]model.Task, error)
	GetById(id uint, issuerId uint) (*model.Task, error)
	Create(issuerId uint, data *model.Task) error
	Update(id uint, issuerId uint, data *dto.TaskUpdateRequest) error
	Delete(id uint, issuerId uint) error
	GetBoardIdByColumnId(columnId uint) (*uint, error)
	GetBoardIdByTaskId(taskId uint) (*uint, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *taskRepository {
	return &taskRepository{db}
}

func (t *taskRepository) Get(boardId uint) ([]model.Task, error) {
	var tasks []model.Task
	if err := t.db.Where("board_id = ?", boardId).Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *taskRepository) GetById(id uint, issuerId uint) (*model.Task, error) {
	var task model.Task
	if err := t.db.Where("id = ?", id).First(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (t *taskRepository) Create(issuerId uint, data *model.Task) error {
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

func (t *taskRepository) GetBoardIdByColumnId(columnId uint) (*uint, error) {
	var column model.BoardColumn
	if err := t.db.First(&column, columnId).Error; err != nil {
		return nil, err
	}

	return &column.BoardID, nil
}

func (t *taskRepository) GetBoardIdByTaskId(taskId uint) (*uint, error) {
	var task model.Task
	if err := t.db.First(&task, taskId).Error; err != nil {
		return nil, err
	}

	boardId, err := t.GetBoardIdByColumnId(task.BoardColumnID)
	if err != nil {
		return nil, err
	}

	return boardId, nil
}