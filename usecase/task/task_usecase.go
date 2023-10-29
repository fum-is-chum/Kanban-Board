package usecase

import (
	"errors"
	"kanban-board/dto"
	"kanban-board/model"
	boardRepo "kanban-board/repository/board"
	taskRepo "kanban-board/repository/task"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type TaskUseCase interface {
	GetTasks(boardId uint, issuerId uint) ([]model.Task, error)
	GetTaskById(id uint, issuerId uint) (*model.Task, error)
	CreateTask(issuerId uint, data *dto.TaskCreateRequest) error
	UpdateTask(id uint, issuerId uint, data *dto.TaskUpdateRequest) error
	DeleteTask(id uint, issuerId uint) error
}

type taskUseCase struct {
	boardRepo boardRepo.BoardRepository
	taskRepo  taskRepo.TaskRepository
}

func NewTaskUseCase(boardRepo boardRepo.BoardRepository, taskRepo taskRepo.TaskRepository) *taskUseCase {
	return &taskUseCase{boardRepo, taskRepo}
}

// ------------------------------------------------------------------
func (t *taskUseCase) isOwner(userId uint, boardId uint) error {
	// check if user is owner of the board
	ownerId, err := t.boardRepo.GetBoardOwner(boardId)
	if err != nil {
		return err
	}

	if *ownerId != userId {
		return errors.New("User is not owner of this board!")
	}

	return nil
}

func (t *taskUseCase) isMember(userId uint, boardId uint) error {
	// check if user is member of the board
	members, err := t.boardRepo.GetBoardMembers(boardId)
	if err != nil {
		return err
	}

	for _, member := range members {
		if member.UserID == userId {
			return nil
		}
	}

	return errors.New("User is not member of this board!")
}

//------------------------------------------------------------------

func (t *taskUseCase) GetTasks(boardId uint, issuerId uint) ([]model.Task, error) {
	if err := t.isMember(issuerId, boardId); err != nil {
		return nil, err
	}

	tasks, err := t.taskRepo.Get(boardId)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *taskUseCase) GetTaskById(id uint, issuerId uint) (*model.Task, error) {
	boardId, err := t.taskRepo.GetBoardIdByTaskId(id)
	if err != nil {
		return nil, err
	}

	if err := t.isMember(issuerId, *boardId); err != nil {
		return nil, err
	}

	task, err := t.taskRepo.GetById(id, issuerId)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (t *taskUseCase) CreateTask(issuerId uint, data *dto.TaskCreateRequest) error {
	if err := t.isMember(issuerId, data.BoardID); err != nil {
		return err
	}

	if err := validate.Struct(*data); err != nil {
		return err
	}

	taskModel := &model.Task{
		Title:         data.Title,
		Desc:          data.Desc,
		BoardID:       data.BoardID,
		BoardColumnID: data.BoardColumnID,
	}

	if err := t.taskRepo.Create(issuerId, taskModel); err != nil {
		return err
	}

	return nil
}

func (t *taskUseCase) UpdateTask(id uint, issuerId uint, data *dto.TaskUpdateRequest) error {
	boardID, err := t.taskRepo.GetBoardIdByColumnId(data.BoardColumnID)
	if err != nil {
		return err
	}

	if err := t.isMember(issuerId, *boardID); err != nil {
		return err
	}

	if err := t.taskRepo.Update(id, issuerId, data); err != nil {
		return err
	}

	return nil
}

func (t *taskUseCase) DeleteTask(id uint, issuerId uint) error {
	boardId, err := t.taskRepo.GetBoardIdByTaskId(id)
	if err != nil {
		return err
	}

	if err := t.isMember(issuerId, *boardId); err != nil {
		return err
	}

	if err := t.taskRepo.Delete(id, issuerId); err != nil {
		return err
	}

	return nil
}
