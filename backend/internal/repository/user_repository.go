package repository

import (
	"errors"

	"github.com/blueship581/gbblog/internal/model"
	"github.com/blueship581/gbblog/internal/util"
	"gorm.io/gorm"
)

// UserRepository 用户仓储。
type UserRepository struct{ db *gorm.DB }

// NewUserRepository 构造用户仓储。
func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) Create(user *model.User) error { return r.db.Create(user).Error }

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error { return r.db.Save(user).Error }

func (r *UserRepository) CountByEmail(email string) (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count, err
}
