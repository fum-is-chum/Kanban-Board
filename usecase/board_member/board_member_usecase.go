package usecase

import (
	"errors"
	"kanban-board/dto"
	boardMemberRepo "kanban-board/repository/board_member"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type BoardMemberUseCase interface {
	AddMember(issuerId uint, data *dto.BoardMemberRequest) error
	DeleteMember(issuerId uint, data *dto.BoardMemberRequest) error
	ExitBoard(issuerId uint, boardId uint) error
}

type boardMemberUseCase struct {
	repo boardMemberRepo.BoardMemberRepository
}

// --------------------------- helper function ----------------------------
func (b *boardMemberUseCase) isOwner(userId uint, boardId uint) error {
	// check if user is owner of the board 
	ownerId, err := b.repo.GetBoardOwner(boardId)
	if err != nil {
		return err
	}

	if *ownerId != userId {
		return errors.New("User is not owner of this board!")
	}

	return nil
}

func (b *boardMemberUseCase) isMember(userId uint, boardId uint) error {
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

func NewBoardMemberUseCase(repo boardMemberRepo.BoardMemberRepository) *boardMemberUseCase {
	return &boardMemberUseCase{repo}
}

func (b *boardMemberUseCase) AddMember(issuerId uint, data *dto.BoardMemberRequest) error {
	if err := validate.Struct(*data); err != nil {
		return err
	}

	if err := b.isOwner(issuerId, data.BoardID); err != nil {
		return err
	}

	if err := b.repo.AddMember(data.BoardID, data.UserID); err != nil {
		return err
	}

	return nil
}

func (b *boardMemberUseCase) DeleteMember(issuerId uint, data *dto.BoardMemberRequest) error {
	if err := validate.Struct(*data); err != nil {
		return err
	}

	if err := b.isOwner(issuerId, data.BoardID); err != nil {
		return err
	}

	if err := b.repo.DeleteMember(data.BoardID, data.UserID); err != nil {
		return err
	}

	return nil
}

func (b *boardMemberUseCase) ExitBoard(issuerId uint, boardId uint) error {
	if err := b.isMember(issuerId, boardId); err != nil {
		return err
	}

	if err := b.repo.DeleteMember(boardId, issuerId); err != nil {
		return err
	}

	return nil
}