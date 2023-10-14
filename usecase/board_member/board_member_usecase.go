package usecase

import (
	"errors"
	"kanban-board/dto"
	boardRepo "kanban-board/repository/board"
	boardMemberRepo "kanban-board/repository/board_member"
	userRepo "kanban-board/repository/user"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type BoardMemberUseCase interface {
	AddNewMember(data *dto.BoardMemberRequest, issuerUserId uint) error
	DeleteMember(data *dto.BoardMemberRequest, issuerUserId uint) error
}

type BoardMemberMultiRepos struct {
	MemberRepo boardMemberRepo.BoardMemberRepository
	BoardRepo  boardRepo.BoardRepository
	UserRepo   userRepo.UserRepository
}

type boardMemberUseCase struct {
	Repos BoardMemberMultiRepos
}

func NewBoardMemberUseCase(repos *BoardMemberMultiRepos) *boardMemberUseCase {
	return &boardMemberUseCase{*repos}
}

func (b *boardMemberUseCase) AddNewMember(data *dto.BoardMemberRequest, issuerUserId uint) error {
	if err := validate.Struct(data); err != nil {
		return err
	}

	// get board owner
	board, err := b.Repos.BoardRepo.GetById(data.BoardId)
	if err != nil {
		return err
	}

	if board.OwnerID != issuerUserId {
		return errors.New("Only board owner can add member!")
	}

	// check if user exist
	if _, err := b.Repos.UserRepo.GetById(data.UserId); err != nil {
		return err
	}

	if err := b.Repos.MemberRepo.AddMember(data.BoardId, data.UserId); err != nil {
		return err
	}

	return nil
}

func (b *boardMemberUseCase) DeleteMember(data *dto.BoardMemberRequest, issuerUserId uint) error {
	if err := validate.Struct(data); err != nil {
		return err
	}

	// get board owner
	board, err := b.Repos.BoardRepo.GetById(data.BoardId)
	if err != nil {
		return err
	}

	if board.OwnerID != issuerUserId {
		return errors.New("Only board owner can delete member!")
	}

	if err := b.Repos.MemberRepo.DeleteMember(data.BoardId, data.UserId); err != nil {
		return err
	}

	return nil
}
