package repository

import (
	"errors"

	"github.com/blueship581/gbblog/internal/model"
	"github.com/blueship581/gbblog/internal/util"
	"gorm.io/gorm"
)

// CategoryRepository 分类仓储。
type CategoryRepository struct{ db *gorm.DB }

// NewCategoryRepository 构造分类仓储。
func NewCategoryRepository(db *gorm.DB) *CategoryRepository { return &CategoryRepository{db: db} }

func (r *CategoryRepository) Create(c *model.Category) error { return r.db.Create(c).Error }

func (r *CategoryRepository) FindByID(id uint) (*model.Category, error) {
	var c model.Category
	if err := r.db.First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) FindBySlug(slug string) (*model.Category, error) {
	var c model.Category
	if err := r.db.Where("slug = ?", slug).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) ExistsBySlug(slug string, excludeID uint) (bool, error) {
	var count int64
	q := r.db.Model(&model.Category{}).Where("slug = ?", slug)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

// Tree 查询全部分类并组装为树。
func (r *CategoryRepository) Tree() ([]model.Category, error) {
	var cats []model.Category
	if err := r.db.Order("sort_order asc").Find(&cats).Error; err != nil {
		return nil, err
	}
	return buildTree(cats, nil), nil
}

func buildTree(cats []model.Category, parentID *uint) []model.Category {
	tree := make([]model.Category, 0)
	for _, c := range cats {
		if (c.ParentID == nil && parentID == nil) || (c.ParentID != nil && parentID != nil && *c.ParentID == *parentID) {
			c.Children = buildTree(cats, &c.ID)
			tree = append(tree, c)
		}
	}
	return tree
}

func (r *CategoryRepository) Update(c *model.Category) error { return r.db.Save(c).Error }

func (r *CategoryRepository) Delete(id uint) error { return r.db.Delete(&model.Category{}, id).Error }

func (r *CategoryRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Category{}).Count(&count).Error
	return count, err
}
