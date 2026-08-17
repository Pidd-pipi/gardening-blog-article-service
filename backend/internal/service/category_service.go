package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/blueship581/gbblog/internal/constants"
	"github.com/blueship581/gbblog/internal/model"
	"github.com/blueship581/gbblog/internal/repository"
	"github.com/blueship581/gbblog/internal/util"
)

// CategoryService 分类服务：多级分类树。
type CategoryService struct {
	repo *repository.CategoryRepository
	log  *slog.Logger
}

// NewCategoryService 构造分类服务。
func NewCategoryService(repo *repository.CategoryRepository, log *slog.Logger) *CategoryService {
	return &CategoryService{repo: repo, log: log}
}

// Create 创建分类。
func (s *CategoryService) Create(ctx context.Context, c *model.Category) (*model.Category, error) {
	exists, err := s.repo.ExistsBySlug(c.Slug, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.ConflictError(constants.MsgSlugConflict, errors.New("slug exists"))
	}
	if err := s.repo.Create(c); err != nil {
		return nil, util.LogError(s.log, constants.LOG_CATEGORY_CREATED, fmt.Errorf("create category: %w", err))
	}
	s.log.InfoContext(ctx, constants.LOG_CATEGORY_CREATED, "category_id", c.ID)
	return c, nil
}

// Tree 分类树。
func (s *CategoryService) Tree(ctx context.Context) ([]model.Category, error) {
	return s.repo.Tree()
}

// Get 查询分类。
func (s *CategoryService) Get(ctx context.Context, id uint) (*model.Category, error) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			return nil, util.NotFoundError(constants.MsgCategoryNotFound, err)
		}
		return nil, err
	}
	return c, nil
}

// GetBySlug 按 slug 查询分类。
func (s *CategoryService) GetBySlug(ctx context.Context, slug string) (*model.Category, error) {
	c, err := s.repo.FindBySlug(slug)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			return nil, util.NotFoundError(constants.MsgCategoryNotFound, err)
		}
		return nil, err
	}
	return c, nil
}

// Update 更新分类。
func (s *CategoryService) Update(ctx context.Context, id uint, c *model.Category) (*model.Category, error) {
	exists, err := s.repo.ExistsBySlug(c.Slug, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.ConflictError(constants.MsgSlugConflict, errors.New("slug exists"))
	}
	cur, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.NotFoundError(constants.MsgCategoryNotFound, err)
	}
	cur.Name = c.Name
	cur.Slug = c.Slug
	cur.Description = c.Description
	cur.SortOrder = c.SortOrder
	if err := s.repo.Update(cur); err != nil {
		return nil, util.LogError(s.log, constants.LOG_CATEGORY_UPDATED, fmt.Errorf("update category: %w", err))
	}
	s.log.InfoContext(ctx, constants.LOG_CATEGORY_UPDATED, "category_id", id)
	return cur, nil
}

// Delete 删除分类。
func (s *CategoryService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return util.LogError(s.log, constants.LOG_CATEGORY_DELETED, fmt.Errorf("delete category: %w", err))
	}
	s.log.InfoContext(ctx, constants.LOG_CATEGORY_DELETED, "category_id", id)
	return nil
}
