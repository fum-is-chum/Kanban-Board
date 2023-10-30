package usecase

import (
	"errors"
	"kanban-board/dto"
	boardRepo "kanban-board/repository/board"
	assigneeRepo "kanban-board/repository/task_assignee"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type TaskAssigneeUseCase interface {
	AddAssignee(issuerId uint, data *dto.TaskAssigneeRequest) error
	DeleteAssignee(issuerId uint, data *dto.TaskAssigneeRequest) error
	ExitTask(issuerId uint, taskId uint) error
}

type taskAssigneeUseCase struct {
	assigneeRepo assigneeRepo.TaskAssigneeRepository
	boardRepo    boardRepo.BoardRepository
}

// ------------------------ helper function ---------------------------------
func (t *taskAssigneeUseCase) isMember(userId uint, boardId uint) error {
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

// -------------------------------------------------------------------------

func NewTaskAssigneeUseCase(boardRepo boardRepo.BoardRepository, assigneeRepo assigneeRepo.TaskAssigneeRepository) *taskAssigneeUseCase {
	return &taskAssigneeUseCase{assigneeRepo, boardRepo}
}

func (t *taskAssigneeUseCase) AddAssignee(issuerId uint, data *dto.TaskAssigneeRequest) error {
	if err := validate.Struct(*data); err != nil {
		return err
	}

	boardId, err := t.assigneeRepo.GetBoardIdByTaskId(data.TaskID)
	if err != nil {
		return err
	}

	if err := t.isMember(issuerId, *boardId); err != nil {
		return err
	}

	if err := t.assigneeRepo.AddAssignee(data); err != nil {
		return err
	}

	return nil
}

func (t *taskAssigneeUseCase) DeleteAssignee(issuerId uint, data *dto.TaskAssigneeRequest) error {
	if err := validate.Struct(*data); err != nil {
		return err
	}

	boardId, err := t.assigneeRepo.GetBoardIdByTaskId(data.TaskID)
	if err != nil {
		return err
	}

	if err := t.isMember(issuerId, *boardId); err != nil {
		return err
	}

	if err := t.assigneeRepo.DeleteAssignee(data.TaskID, data.UserID); err != nil {
		return err
	}

	return nil
}

func (t *taskAssigneeUseCase) ExitTask(issuerId uint, taskId uint) error {
	boardId, err := t.assigneeRepo.GetBoardIdByTaskId(taskId)
	if err != nil {
		return err
	}

	if err := t.isMember(issuerId, *boardId); err != nil {
		return err
	}

	if err := t.assigneeRepo.DeleteAssignee(taskId, issuerId); err != nil {
		return err
	}

	return nil
}
