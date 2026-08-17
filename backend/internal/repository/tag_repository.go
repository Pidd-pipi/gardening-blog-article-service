package repository

import (
	"errors"

	"github.com/blueship581/gbblog/internal/model"
	"github.com/blueship581/gbblog/internal/util"
	"gorm.io/gorm"
)

// TagRepository 标签仓储。
type TagRepository struct{ db *gorm.DB }

// NewTagRepository 构造标签仓储。
func NewTagRepository(db *gorm.DB) *TagRepository { return &TagRepository{db: db} }

func (r *TagRepository) Create(t *model.Tag) error { return r.db.Create(t).Error }

func (r *TagRepository) FindByID(id uint) (*model.Tag, error) {
	var t model.Tag
	if err := r.db.First(&t, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *TagRepository) FindBySlug(slug string) (*model.Tag, error) {
	var t model.Tag
	if err := r.db.Where("slug = ?", slug).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *TagRepository) List() ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Order("name asc").Find(&tags).Error
	return tags, err
}

func (r *TagRepository) ListByIDs(ids []uint) ([]model.Tag, error) {
	var tags []model.Tag
	if len(ids) == 0 {
		return tags, nil
	}
	err := r.db.Where("id IN ?", ids).Find(&tags).Error
	return tags, err
}

func (r *TagRepository) Update(t *model.Tag) error { return r.db.Save(t).Error }

func (r *TagRepository) Delete(id uint) error { return r.db.Delete(&model.Tag{}, id).Error }

func (r *TagRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Tag{}).Count(&count).Error
	return count, err
}
