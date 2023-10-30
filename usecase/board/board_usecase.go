package usecase

import (
	"errors"
	"kanban-board/dto"
	"kanban-board/model"
	boardRepo "kanban-board/repository/board"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type BoardUseCase interface {
	GetBoards(issuerId uint) ([]model.Board, error)
	GetBoardById(id uint, issuerId uint) (*model.Board, error)
	CreateBoard(issuerId uint, data *dto.BoardRequest) error
	UpdateBoard(id uint, issuerId uint, data *dto.BoardRequest) error
	DeleteBoard(id uint, issuerId uint) error
}

type boardUseCase struct {
	repo boardRepo.BoardRepository
}

func (b *boardUseCase) isMember(userId uint, boardId uint) error {
	// check if user is member of the board
	members, err := b.repo.GetBoardMembers(boardId)
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

// ------------------------------------------------------------------------

func NewBoardUseCase(repo boardRepo.BoardRepository) *boardUseCase {
	return &boardUseCase{repo}
}

func (b *boardUseCase) GetBoards(issuerId uint) ([]model.Board, error) {
	boards, err := b.repo.Get(issuerId)
	if err != nil {
		return nil, err
	}

	return boards, err
}

func (b *boardUseCase) GetBoardById(id uint, issuerId uint) (*model.Board, error) {
	board, err := b.repo.GetById(id, issuerId)
	if err != nil {
		return nil, err
	}

	return board, nil
}

func (b *boardUseCase) CreateBoard(issuerId uint, data *dto.BoardRequest) error {
	err := validate.Struct(*data)
	if err != nil {
		return err
	}

	boardModel := &model.Board{
		Name:    data.Name,
		Desc:    data.Desc,
		OwnerID: issuerId,
	}

	if err := b.repo.Create(boardModel); err != nil {
		return err
	}

	return nil
}

func (b *boardUseCase) UpdateBoard(id uint, issuerId uint, data *dto.BoardRequest) error {
	// check board member
	if err := b.isMember(issuerId, id); err != nil {
		return err
	}

	updatedData := &dto.BoardRequest{
		Name: data.Name,
		Desc: data.Desc,
	}

	if err := b.repo.Update(id, issuerId, updatedData); err != nil {
		return err
	}

	return nil
}

func (b *boardUseCase) DeleteBoard(id uint, issuerId uint) error {
	// ensure user cannot delete other user's board
	board, err := b.repo.GetById(id, issuerId)
	if err != nil {
		return err
	}

	if board.OwnerID != issuerId {
		return errors.New("User only can delete board they owned!")
	}

	if err := b.repo.Delete(id); err != nil {
		return err
	}

	return nil
}
