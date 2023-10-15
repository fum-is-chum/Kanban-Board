package usecase

import (
	"errors"
	"kanban-board/dto"
	fieldHelper "kanban-board/helpers/field"
	"kanban-board/model"
	boardRepo "kanban-board/repository/board"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type BoardUseCase interface {
	GetBoards() ([]model.Board, error)
	GetBoardById(id uint) (*model.Board, error)
	CreateBoard(data *dto.BoardRequest) error
	UpdateBoard(id uint, data *dto.BoardRequest) error
	DeleteBoard(id uint, issuerUserId uint) error
}

type boardUseCase struct {
	repo boardRepo.BoardRepository
}

func NewBoardUseCase(repo boardRepo.BoardRepository) *boardUseCase {
	return &boardUseCase{repo}
}

func (b *boardUseCase) GetBoards() ([]model.Board, error) {
	boards, err := b.repo.Get()
	if err != nil {
		return nil, err
	}

	return boards, err
}

func (b *boardUseCase) GetBoardById(id uint) (*model.Board, error) {
	board, err := b.repo.GetById(id)
	if err != nil {
		return nil, err
	}

	return board, nil
}

func (b *boardUseCase) CreateBoard(data *dto.BoardRequest) error {
	err := validate.Struct(*data)
	if err != nil {
		return err
	}

	boardModel := &model.Board{
		Name:    data.Name,
		Desc:    data.Desc,
		OwnerID: data.OwnerID,
	}

	if err := b.repo.Create(boardModel); err != nil {
		return err
	}

	return nil
}

func (b *boardUseCase) UpdateBoard(id uint, data *dto.BoardRequest) error {
	val := reflect.ValueOf(*data)

	if fieldHelper.IsFieldSet(&val, "OwnerID") {
		return errors.New("Cannot update OwnerID from this endpoint")
	}

	if err := b.repo.Update(id, data); err != nil {
		return err
	}

	return nil
}

func (b *boardUseCase) DeleteBoard(id uint, issuerUserId uint) error {
	// ensure user cannot delete other user's board
	board, err := b.repo.GetById(id)
	if err != nil {
		return err
	}

	if board.OwnerID != issuerUserId {
		return errors.New("User only can delete board they owned!")
	}

	if err := b.repo.Delete(id); err != nil {
		return err
	}

	return nil
}
