package usecase

import (
	"errors"
	"kanban-board/dto"
	"kanban-board/model"
	boardRepo "kanban-board/repository/board"
	userRepo "kanban-board/repository/user"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type BoardUseCase interface {
	GetBoards() ([]model.Board, error)
	GetBoardById(id uint) (*model.Board, error)
	CreateBoard(data *dto.BoardRequest) error
	UpdateBoard(id uint, data *dto.BoardRequest) error
	UpdateBoardOwnership(boardId uint, ownerId uint) error
	DeleteBoard(id uint) error
}

type boardUseCase struct {
	userRepo  userRepo.UserRepository
	boardRepo boardRepo.BoardRepository
}

func NewBoardUseCase(userRepo userRepo.UserRepository, boardRepo boardRepo.BoardRepository) *boardUseCase {
	return &boardUseCase{userRepo, boardRepo}
}

func (b *boardUseCase) GetBoards() ([]model.Board, error) {
	boards, err := b.boardRepo.Get()
	if err != nil {
		return nil, err
	}

	return boards, err
}

func (b *boardUseCase) GetBoardById(id uint) (*model.Board, error) {
	board, err := b.boardRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	return board, nil
}

func (b *boardUseCase) CreateBoard(data *dto.BoardRequest) error {
	err := validate.Struct(data)
	if err != nil {
		return err
	}

	boardModel := &model.Board{
		Name:    data.Name,
		Desc:    data.Desc,
		OwnerID: data.OwnerID,
	}

	if err := b.boardRepo.Create(boardModel); err != nil {
		return err
	}

	return nil
}

func (b *boardUseCase) UpdateBoard(id uint, data *dto.BoardRequest) error {
	updates := make(map[string]interface{})

	structValue := reflect.ValueOf(*data)
	for i := 0; i < structValue.NumField(); i++ {
		key := structValue.Type().Field(i).Name
		value := structValue.Field(i).Interface()

		if value != reflect.Zero(structValue.Type().Field(i).Type).Interface() {
			if key == "OwnerID" {
				return errors.New("Cannot update board OwnerID from this endpoint!")
			}
			updates[key] = value
		}
	}

	// check if there is no fields to update
	if len(updates) == 0 {
		return errors.New("No fields to update or fields value is empty!")
	}

	if err := b.boardRepo.Update(id, &updates); err != nil {
		return err
	}

	return nil
}

func (b *boardUseCase) UpdateBoardOwnership(boardId uint, ownerId uint) error {
	// check board exist or not
	if _, err := b.boardRepo.GetById(boardId); err != nil {
		return err
	}

	// check user exist or not
	if _, err := b.userRepo.GetById(ownerId); err != nil {
		return err
	}

	updateMap := &map[string]interface{}{
		"OwnerID": ownerId,
	}

	if err := b.boardRepo.Update(boardId, updateMap); err != nil {
		return err
	}

	return nil
}

func (b *boardUseCase) DeleteBoard(id uint) error {
	if err := b.boardRepo.Delete(id); err != nil {
		return err
	}

	return nil
}
