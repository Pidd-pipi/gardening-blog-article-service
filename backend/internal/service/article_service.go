package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/blueship581/gbblog/internal/constants"
	"github.com/blueship581/gbblog/internal/model"
	"github.com/blueship581/gbblog/internal/repository"
	"github.com/blueship581/gbblog/internal/util"
	"gorm.io/gorm"
)

// ArticleService 文章服务：撰写、草稿、发布、置顶、搜索、归档。
type ArticleService struct {
	repo     *repository.ArticleRepository
	tagRepo  *repository.TagRepository
	catRepo  *repository.CategoryRepository
	redis    redisClient
	log      *slog.Logger
}

// redisClient Redis 操作接口（供视图计数缓存）。
type redisClient interface {
	Incr(ctx context.Context, key string) (int64, error)
}

// NewArticleService 构造文章服务。
func NewArticleService(repo *repository.ArticleRepository, tagRepo *repository.TagRepository, catRepo *repository.CategoryRepository, r redisClient, log *slog.Logger) *ArticleService {
	return &ArticleService{repo: repo, tagRepo: tagRepo, catRepo: catRepo, redis: r, log: log}
}

// ArticleInput 文章入参。
type ArticleInput struct {
	Title      string `json:"title" binding:"required,max=200"`
	Slug       string `json:"slug" binding:"required,max=100"`
	CategoryID *uint  `json:"category_id"`
	Summary    string `json:"summary" binding:"max=500"`
	ContentMD  string `json:"content_markdown"`
	Cover      string `json:"cover" binding:"max=255"`
	Status     string `json:"status" binding:"oneof=draft published scheduled"`
	IsTop      bool   `json:"is_top"`
	TagIDs     []uint `json:"tag_ids"`
}

// Create 创建文章（草稿或发布）。
func (s *ArticleService) Create(ctx context.Context, userID uint, input ArticleInput) (*model.Article, error) {
	if err := s.validateSlug(input.Slug, 0); err != nil {
		return nil, err
	}
	if input.CategoryID != nil {
		if _, err := s.catRepo.FindByID(*input.CategoryID); err != nil {
			return nil, util.NotFoundError(constants.MsgCategoryNotFound, err)
		}
	}
	html := util.RenderMarkdown(input.ContentMD)
	now := time.Now()
	a := &model.Article{
		UserID: userID, CategoryID: input.CategoryID, Title: input.Title, Slug: input.Slug,
		Summary: input.Summary, ContentMarkdown: input.ContentMD, ContentHTML: html,
		Cover: input.Cover, Status: input.Status, IsTop: input.IsTop,
		WordCount: wordCount(input.ContentMD),
	}
	if input.Status == constants.ArticlePublished {
		a.PublishedAt = &now
	}
	err := s.repo.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.WithTx(tx).Create(a); err != nil {
			return fmt.Errorf("create article: %w", err)
		}
		if err := s.repo.WithTx(tx).ReplaceTags(a.ID, input.TagIDs); err != nil {
			return fmt.Errorf("replace article tags: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, util.LogError(s.log, constants.LOG_ARTICLE_CREATED, err)
	}
	s.log.InfoContext(ctx, constants.LOG_ARTICLE_CREATED, "article_id", a.ID, "status", a.Status)
	return s.repo.FindByID(a.ID)
}

// Update 更新文章。
func (s *ArticleService) Update(ctx context.Context, id uint, input ArticleInput) (*model.Article, error) {
	if err := s.validateSlug(input.Slug, id); err != nil {
		return nil, err
	}
	a, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.NotFoundError(constants.MsgArticleNotFound, err)
	}
	a.Title = input.Title
	a.Slug = input.Slug
	a.CategoryID = input.CategoryID
	a.Summary = input.Summary
	a.ContentMarkdown = input.ContentMD
	a.ContentHTML = util.RenderMarkdown(input.ContentMD)
	a.Cover = input.Cover
	a.IsTop = input.IsTop
	a.WordCount = wordCount(input.ContentMD)
	if a.Status != constants.ArticlePublished && input.Status == constants.ArticlePublished {
		now := time.Now()
		a.PublishedAt = &now
	}
	a.Status = input.Status
	err = s.repo.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.WithTx(tx).Update(a); err != nil {
			return fmt.Errorf("update article: %w", err)
		}
		if err := s.repo.WithTx(tx).ReplaceTags(a.ID, input.TagIDs); err != nil {
			return fmt.Errorf("replace article tags: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, util.LogError(s.log, constants.LOG_ARTICLE_UPDATED, err)
	}
	s.log.InfoContext(ctx, constants.LOG_ARTICLE_UPDATED, "article_id", id)
	return s.repo.FindByID(id)
}

// Delete 删除文章（文章与标签关联一并删除）。
func (s *ArticleService) Delete(ctx context.Context, id uint) error {
	err := s.repo.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.WithTx(tx).Delete(id); err != nil {
			return fmt.Errorf("delete article: %w", err)
		}
		return nil
	})
	if err != nil {
		return util.LogError(s.log, constants.LOG_ARTICLE_DELETED, err)
	}
	s.log.InfoContext(ctx, constants.LOG_ARTICLE_DELETED, "article_id", id)
	return nil
}

// Publish 发布文章。
func (s *ArticleService) Publish(ctx context.Context, id uint) (*model.Article, error) {
	a, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.NotFoundError(constants.MsgArticleNotFound, err)
	}
	now := time.Now()
	a.Status = constants.ArticlePublished
	if a.PublishedAt == nil {
		a.PublishedAt = &now
	}
	if err := s.repo.Update(a); err != nil {
		return nil, util.LogError(s.log, constants.LOG_ARTICLE_PUBLISHED, fmt.Errorf("publish article: %w", err))
	}
	s.log.InfoContext(ctx, constants.LOG_ARTICLE_PUBLISHED, "article_id", id)
	return a, nil
}

// GetPublishedBySlug 前台查看文章（计数 +1）。
func (s *ArticleService) GetPublishedBySlug(ctx context.Context, slug string) (*model.Article, error) {
	a, err := s.repo.FindBySlug(slug)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			return nil, util.NotFoundError(constants.MsgArticleNotFound, err)
		}
		return nil, err
	}
	if a.Status != constants.ArticlePublished {
		return nil, util.NotFoundError(constants.MsgArticleNotFound, errors.New("article not published"))
	}
	if err := s.repo.IncrementView(a.ID); err != nil {
		s.log.WarnContext(ctx, constants.LOG_REDIS_ERROR, "error", err)
	}
	a.ViewCount++
	s.log.InfoContext(ctx, constants.LOG_ARTICLE_VIEWED, "article_id", a.ID, "slug", slug)
	return a, nil
}

// GetByID 后台查看（含草稿）。
func (s *ArticleService) GetByID(ctx context.Context, id uint) (*model.Article, error) {
	a, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			return nil, util.NotFoundError(constants.MsgArticleNotFound, err)
		}
		return nil, err
	}
	return a, nil
}

// ListPublished 前台列表。
func (s *ArticleService) ListPublished(ctx context.Context, categoryID, tagID *uint, keyword string, page, pageSize int) ([]model.Article, int64, error) {
	return s.repo.ListPublished(categoryID, tagID, keyword, page, pageSize)
}

// ListAll 后台列表。
func (s *ArticleService) ListAll(ctx context.Context, status string, page, pageSize int) ([]model.Article, int64, error) {
	return s.repo.ListAll(status, page, pageSize)
}

// Search 全文搜索。
func (s *ArticleService) Search(ctx context.Context, keyword string, page, pageSize int) ([]model.Article, int64, error) {
	keyword = util.NormalizeKeyword(keyword)
	if len(keyword) > constants.MaxSearchKeywordLen {
		return nil, 0, util.BadRequest("搜索关键词（keyword）过长", errors.New("search keyword too long"))
	}
	s.log.InfoContext(ctx, constants.LOG_ARTICLE_SEARCHED, "keyword", keyword)
	return s.repo.Search(keyword, page, pageSize)
}

// TopViewed 热门文章。
func (s *ArticleService) TopViewed(ctx context.Context, limit int) ([]model.Article, error) {
	return s.repo.TopViewed(limit)
}

func (s *ArticleService) validateSlug(slug string, excludeID uint) error {
	if !util.IsValidSlug(slug) {
		return util.BadRequest("文章别名（Article.slug）不合法", errors.New("invalid slug"))
	}
	exists, err := s.repo.ExistsBySlug(slug, excludeID)
	if err != nil {
		return err
	}
	if exists {
		return util.ConflictError(constants.MsgSlugConflict, errors.New("slug exists"))
	}
	return nil
}

func wordCount(md string) int {
	return len([]rune(util.PlainText(util.RenderMarkdown(md))))
}
