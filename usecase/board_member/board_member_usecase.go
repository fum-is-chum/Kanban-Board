package usecase

import (
	"kanban-board/dto"
	boardMemberRepo "kanban-board/repository/board_member"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type BoardMemberUseCase interface {
	AddMember(data *dto.BoardMemberRequest) error
	DeleteMember(data *dto.BoardMemberRequest) error
}

type boardMemberUseCase struct {
	repo boardMemberRepo.BoardMemberRepository
}

func NewBoardMemberUseCase(repo boardMemberRepo.BoardMemberRepository) *boardMemberUseCase {
	return &boardMemberUseCase{repo}
}

func (b *boardMemberUseCase) AddMember(data *dto.BoardMemberRequest) error {
	if err := validate.Struct(*data); err != nil {
		return err
	}

	if err := b.repo.AddMember(data.BoardID, data.UserID); err != nil {
		return err
	}

	return nil
}

func (b *boardMemberUseCase) DeleteMember(data *dto.BoardMemberRequest) error {
	if err := validate.Struct(*data); err != nil {
		return err
	}
	

	if err := b.repo.DeleteMember(data.BoardID, data.UserID); err != nil {
		return err
	}

	return nil
}
