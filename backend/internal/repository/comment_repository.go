package repository

import (
	"errors"

	"github.com/blueship581/gbblog/internal/model"
	"github.com/blueship581/gbblog/internal/util"
	"gorm.io/gorm"
)

// CommentRepository 评论仓储。
type CommentRepository struct{ db *gorm.DB }

// NewCommentRepository 构造评论仓储。
func NewCommentRepository(db *gorm.DB) *CommentRepository { return &CommentRepository{db: db} }

func (r *CommentRepository) Create(c *model.Comment) error { return r.db.Create(c).Error }

func (r *CommentRepository) FindByID(id uint) (*model.Comment, error) {
	var c model.Comment
	if err := r.db.First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

// ListByArticle 按文章查询已通过评论（嵌套两层回复）。
func (r *CommentRepository) ListByArticle(articleID uint) ([]model.Comment, error) {
	var roots []model.Comment
	if err := r.db.Preload("User").Where("article_id = ? AND parent_id IS NULL AND status = ?", articleID, "approved").Order("created_at asc").Find(&roots).Error; err != nil {
		return nil, err
	}
	for i := range roots {
		var replies []model.Comment
		if err := r.db.Preload("User").Where("article_id = ? AND parent_id = ? AND status = ?", articleID, roots[i].ID, "approved").Order("created_at asc").Find(&replies).Error; err != nil {
			return nil, err
		}
		roots[i].Replies = replies
	}
	return roots, nil
}

// ListAll 后台全量评论（含待审核）。
func (r *CommentRepository) ListAll(status string, page, pageSize int) ([]model.Comment, int64, error) {
	q := r.db.Model(&model.Comment{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.Comment
	q2 := r.db.Preload("User").Preload("Article")
	if status != "" {
		q2 = q2.Where("status = ?", status)
	}
	err := q2.Order("created_at desc").Offset((page-1)*pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *CommentRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", id).Update("status", status).Error
}

func (r *CommentRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Comment{}).Count(&count).Error
	return count, err
}
