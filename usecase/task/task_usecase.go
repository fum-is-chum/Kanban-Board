package usecase

import (
	"kanban-board/dto"
	"kanban-board/model"
	taskRepo "kanban-board/repository/task"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type TaskUseCase interface {
	GetTasks(boardId uint) ([]model.Task, error)
	GetTaskById(id uint, issuerId uint) (*model.Task, error)
	CreateTask(issuerId uint, data *dto.TaskCreateRequest) error
	UpdateTask(id uint, issuerId uint, data *dto.TaskUpdateRequest) error
	DeleteTask(id uint, issuerId uint) error
}

type taskUseCase struct {
	repo taskRepo.TaskRepository
}

func NewTaskUseCase(repo taskRepo.TaskRepository) *taskUseCase {
	return &taskUseCase{repo}
}

func (t *taskUseCase) GetTasks(boardId uint) ([]model.Task, error) {
	tasks, err := t.repo.Get(boardId)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *taskUseCase) GetTaskById(id uint, issuerId uint) (*model.Task, error) {
	task, err := t.repo.GetById(id, issuerId)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (t *taskUseCase) CreateTask(issuerId uint, data *dto.TaskCreateRequest) error {
	if err := validate.Struct(*data); err != nil {
		return err
	}

	taskModel := &model.Task{
		Title:         data.Title,
		Desc:          data.Desc,
		BoardID:       data.BoardID,
		BoardColumnID: data.BoardColumnID,
	}

	if err := t.repo.Create(issuerId, taskModel); err != nil {
		return err
	}

	return nil
}

func (t *taskUseCase) UpdateTask(id uint, issuerId uint, data *dto.TaskUpdateRequest) error {
	if err := t.repo.Update(id, issuerId, data); err != nil {
		return err
	}

	return nil
}

func (t *taskUseCase) DeleteTask(id uint, issuerId uint) error {
	if err := t.repo.Delete(id, issuerId); err != nil {
		return err
	}

	return nil
}
