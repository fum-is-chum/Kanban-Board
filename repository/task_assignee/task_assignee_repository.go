package repository

import (
	"kanban-board/dto"
	"kanban-board/model"

	"gorm.io/gorm"
)

type TaskAssigneeRepository interface {
	GetBoardIdByTaskId(taskId uint) (*uint, error)
	AddAssignee(data *dto.TaskAssigneeRequest) error
	DeleteAssignee(taskId uint, userId uint) error
}

type taskAssigneeRepository struct {
	db *gorm.DB
}

func NewTaskAssigneeRepository(db *gorm.DB) *taskAssigneeRepository {
	return &taskAssigneeRepository{db}
}

func (t *taskAssigneeRepository) GetBoardIdByTaskId(taskId uint) (*uint, error) {
	var task model.Task

	if err := t.db.First(&task, taskId).Error; err != nil {
		return nil, err
	}

	return &task.BoardID, nil
}

func (t *taskAssigneeRepository) AddAssignee(data *dto.TaskAssigneeRequest) error {
	var task model.Task
	var user model.User

	if err := t.db.First(&task, data.TaskID).Error; err != nil {
		return err
	}

	if err := t.db.First(&user, data.UserID).Error; err != nil {
		return err
	}

	if err := t.db.Model(&task).Association("Assignees").Append(&user); err != nil {
		return err
	}

	return nil
}

func (t *taskAssigneeRepository) DeleteAssignee(taskId uint, userId uint) error {
	var user model.User
	var task model.Task

	if err := t.db.First(&user, userId).Error; err != nil {
		return err
	}

	if err := t.db.First(&task, taskId).Error; err != nil {
		return err
	}

	if err := t.db.Model(&task).Association("Assignees").Delete(&user); err != nil {
		return err
	}

	return nil
}