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

// TagService 标签服务。
type TagService struct {
	repo *repository.TagRepository
	log  *slog.Logger
}

// NewTagService 构造标签服务。
func NewTagService(repo *repository.TagRepository, log *slog.Logger) *TagService {
	return &TagService{repo: repo, log: log}
}

// Create 创建标签。
func (s *TagService) Create(ctx context.Context, name, slug string) (*model.Tag, error) {
	if existing, err := s.repo.FindBySlug(slug); err == nil && existing != nil {
		return nil, util.ConflictError(constants.MsgSlugConflict, errors.New("slug exists"))
	}
	t := &model.Tag{Name: name, Slug: slug}
	if err := s.repo.Create(t); err != nil {
		return nil, util.LogError(s.log, constants.LOG_TAG_CREATED, fmt.Errorf("create tag: %w", err))
	}
	s.log.InfoContext(ctx, constants.LOG_TAG_CREATED, "tag_id", t.ID)
	return t, nil
}

// List 标签列表。
func (s *TagService) List(ctx context.Context) ([]model.Tag, error) {
	return s.repo.List()
}

// GetBySlug 按 slug 查询标签。
func (s *TagService) GetBySlug(ctx context.Context, slug string) (*model.Tag, error) {
	t, err := s.repo.FindBySlug(slug)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			return nil, util.NotFoundError(constants.MsgTagNotFound, err)
		}
		return nil, err
	}
	return t, nil
}

// Update 更新标签。
func (s *TagService) Update(ctx context.Context, id uint, name, slug string) (*model.Tag, error) {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.NotFoundError(constants.MsgTagNotFound, err)
	}
	t.Name = name
	t.Slug = slug
	if err := s.repo.Update(t); err != nil {
		return nil, util.LogError(s.log, constants.LOG_TAG_UPDATED, fmt.Errorf("update tag: %w", err))
	}
	s.log.InfoContext(ctx, constants.LOG_TAG_UPDATED, "tag_id", id)
	return t, nil
}

// Delete 删除标签。
func (s *TagService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return util.LogError(s.log, constants.LOG_TAG_DELETED, fmt.Errorf("delete tag: %w", err))
	}
	s.log.InfoContext(ctx, constants.LOG_TAG_DELETED, "tag_id", id)
	return nil
}
