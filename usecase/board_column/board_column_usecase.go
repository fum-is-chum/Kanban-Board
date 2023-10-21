package usecase

import (
	"kanban-board/dto"
	"kanban-board/model"
	boardColumnRepo "kanban-board/repository/board_column"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type BoardColumnUseCase interface {
	GetColumns(boardId uint) ([]model.BoardColumn, error)
	GetColumnById(id uint) (*model.BoardColumn, error)
	CreateColumn(issuerId uint, data *dto.BoardColumnRequest) error
	UpdateColumn(id uint, issuerId uint, data *dto.BoardColumnRequest) error
	DeleteColumn(id uint, issuerId uint) error
}

type boardColumnUseCase struct {
	repo boardColumnRepo.BoardColumnRepository
}

func NewBoardColumnUseCase(repo boardColumnRepo.BoardColumnRepository) *boardColumnUseCase {
	return &boardColumnUseCase{repo}
}

func (b *boardColumnUseCase) GetColumns(boardId uint) ([]model.BoardColumn, error) {
	columns, err := b.repo.Get(boardId)
	if err != nil {
		return nil, err
	}

	return columns, nil
}

func (b *boardColumnUseCase) GetColumnById(id uint) (*model.BoardColumn, error) {
	column, err := b.repo.GetById(id)
	if err != nil {
		return nil, err
	}

	return column, nil
}

func (b *boardColumnUseCase) CreateColumn(issuerId uint, data *dto.BoardColumnRequest) error {
	if err := validate.Struct(*data); err != nil {
		return err
	}

	columnModel := &model.BoardColumn{
		Label:   data.Label,
		Desc:    data.Desc,
		BoardID: data.BoardID,
	}

	if err := b.repo.Create(issuerId, columnModel); err != nil {
		return err
	}

	return nil
}

func (b *boardColumnUseCase) UpdateColumn(id uint, issuerId uint, data *dto.BoardColumnRequest) error {
	updatedData := &dto.BoardColumnRequest{
		Label: data.Label,
		Desc: data.Desc,
	}

	if err := b.repo.Update(id, issuerId, updatedData); err != nil {
		return err
	}

	return nil
}

func (b *boardColumnUseCase) DeleteColumn(id uint, issuerId uint) error {
	if err := b.repo.Delete(id, issuerId); err != nil {
		return err
	}

	return nil
}
