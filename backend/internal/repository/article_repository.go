package repository

import (
	"errors"

	"github.com/blueship581/gbblog/internal/model"
	"github.com/blueship581/gbblog/internal/util"
	"gorm.io/gorm"
)

// ArticleRepository 文章仓储。
type ArticleRepository struct{ db *gorm.DB }

// NewArticleRepository 构造文章仓储。
func NewArticleRepository(db *gorm.DB) *ArticleRepository { return &ArticleRepository{db: db} }

// WithTx 使用事务连接构造仓储。
func (r *ArticleRepository) WithTx(tx *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: tx}
}

// Transaction 在事务内执行 fn，任一步返回 error 则整体回滚。
func (r *ArticleRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *ArticleRepository) Create(a *model.Article) error { return r.db.Create(a).Error }

func (r *ArticleRepository) FindByID(id uint) (*model.Article, error) {
	var a model.Article
	if err := r.db.Preload("User").Preload("Category").Preload("Tags").First(&a, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.ErrNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *ArticleRepository) FindBySlug(slug string) (*model.Article, error) {
	var a model.Article
	if err := r.db.Preload("User").Preload("Category").Preload("Tags").Where("slug = ?", slug).First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.ErrNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *ArticleRepository) ExistsBySlug(slug string, excludeID uint) (bool, error) {
	var count int64
	q := r.db.Model(&model.Article{}).Where("slug = ?", slug)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

// ListPublished 前台文章列表（仅 published，可按时序/置顶）。
func (r *ArticleRepository) ListPublished(categoryID, tagID *uint, keyword string, page, pageSize int) ([]model.Article, int64, error) {
	query := func(db *gorm.DB) *gorm.DB {
		q := db.Model(&model.Article{}).Where("articles.status = ?", "published")
		if categoryID != nil {
			q = q.Where("articles.category_id = ?", *categoryID)
		}
		if tagID != nil {
			q = q.Joins("JOIN article_tags ON article_tags.article_id = articles.id").Where("article_tags.tag_id = ?", *tagID)
		}
		if keyword != "" {
			like := "%" + keyword + "%"
			q = q.Where("(LOWER(articles.title) LIKE LOWER(?) OR LOWER(articles.summary) LIKE LOWER(?))", like, like)
		}
		return q
	}
	var total int64
	if err := query(r.db).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.Article
	err := query(r.db).Preload("User").Preload("Category").Preload("Tags").
		Order("is_top desc, published_at desc").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

// ListAll 后台全量文章（含草稿/定时）。
func (r *ArticleRepository) ListAll(status string, page, pageSize int) ([]model.Article, int64, error) {
	q := r.db.Model(&model.Article{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.Article
	q2 := r.db.Preload("User").Preload("Category").Preload("Tags").Order("updated_at desc")
	if status != "" {
		q2 = q2.Where("status = ?", status)
	}
	err := q2.Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *ArticleRepository) Update(a *model.Article) error { return r.db.Save(a).Error }

// ReplaceTags 替换文章标签关联。
func (r *ArticleRepository) ReplaceTags(articleID uint, tagIDs []uint) error {
	if err := r.db.Where("article_id = ?", articleID).Delete(&model.ArticleTag{}).Error; err != nil {
		return err
	}
	for _, tid := range tagIDs {
		if err := r.db.Create(&model.ArticleTag{ArticleID: articleID, TagID: tid}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *ArticleRepository) Delete(id uint) error {
	if err := r.db.Where("article_id = ?", id).Delete(&model.ArticleTag{}).Error; err != nil {
		return err
	}
	return r.db.Delete(&model.Article{}, id).Error
}

func (r *ArticleRepository) IncrementView(id uint) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// TopViewed 热门文章 Top N。
func (r *ArticleRepository) TopViewed(limit int) ([]model.Article, error) {
	var items []model.Article
	err := r.db.Where("status = ?", "published").Order("view_count desc").Limit(limit).Find(&items).Error
	return items, err
}

// Count 统计。
func (r *ArticleRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Article{}).Count(&count).Error
	return count, err
}

// SumWordCount 总字数。
func (r *ArticleRepository) SumWordCount() (int64, error) {
	var sum int64
	err := r.db.Model(&model.Article{}).Select("COALESCE(sum(word_count),0)").Scan(&sum).Error
	return sum, err
}

// Search 全文搜索（标题/摘要/正文），按相关度排序。
func (r *ArticleRepository) Search(keyword string, page, pageSize int) ([]model.Article, int64, error) {
	like := "%" + keyword + "%"
	q := r.db.Model(&model.Article{}).Where("status = ? AND LOWER(title) LIKE LOWER(?)", "published", like)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.Article
	err := r.db.Preload("User").Preload("Category").Preload("Tags").
		Where("articles.status = ?", "published").
		Where("LOWER(articles.title) LIKE LOWER(?)", like).
		Order("published_at desc").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

// Archives 按月归档（published）。
func (r *ArticleRepository) Archives() ([]model.ArchiveMonth, error) {
	var months []model.ArchiveMonth
	err := r.db.Model(&model.Article{}).
		Select("to_char(published_at, 'YYYY-MM') as month, count(*) as count").
		Where("status = ?", "published").
		Group("to_char(published_at, 'YYYY-MM')").
		Order("month desc").Scan(&months).Error
	return months, err
}

// ListByMonth 按月查询已发布文章。
func (r *ArticleRepository) ListByMonth(month string) ([]model.Article, error) {
	var items []model.Article
	err := r.db.Preload("User").Preload("Category").Preload("Tags").
		Where("status = ? AND to_char(published_at, 'YYYY-MM') = ?", "published", month).
		Order("published_at desc").Find(&items).Error
	return items, err
}
