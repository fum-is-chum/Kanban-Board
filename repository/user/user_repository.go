package repository

import (
	"kanban-board/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Get() ([]model.User, error)
	GetById(id uint) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Create(data *model.User) error
	Update(id uint, data *map[string]interface{}) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (u *userRepository) Get() ([]model.User, error) {
	var users []model.User

	tx := u.db.Order("created_at desc").Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return users, nil
}

func (u *userRepository) GetById(id uint) (*model.User, error) {
	var user model.User

	tx := u.db.Where("id = ?", id).First(&user)
	if tx.Error != nil {
		return &model.User{}, tx.Error
	}

	return &user, nil
}

func (u *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User

	tx := u.db.Where("email = ?", email).First(&user)
	if tx.Error != nil {
		return &model.User{}, nil
	}

	return &user, nil
}

func (u *userRepository) Create(data *model.User) error {
	if err := u.db.Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (u *userRepository) Update(id uint, data *map[string]interface{}) error {
	tx := u.db.Model(&model.User{}).Where("id = ?", id).Updates(&data)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (u *userRepository) Delete(id uint) error {
	if err := u.db.Delete(&model.User{}, id).Error; err != nil {
		return err
	}
	return nil
}