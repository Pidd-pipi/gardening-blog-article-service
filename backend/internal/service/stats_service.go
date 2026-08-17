package service

import (
	"context"
	"log/slog"

	"github.com/blueship581/gbblog/internal/constants"
	"github.com/blueship581/gbblog/internal/model"
	"github.com/blueship581/gbblog/internal/repository"
)

// StatsService 站点统计服务。
type StatsService struct {
	articleRepo *repository.ArticleRepository
	catRepo     *repository.CategoryRepository
	tagRepo     *repository.TagRepository
	commentRepo *repository.CommentRepository
	log         *slog.Logger
}

// NewStatsService 构造统计服务。
func NewStatsService(articleRepo *repository.ArticleRepository, catRepo *repository.CategoryRepository, tagRepo *repository.TagRepository, commentRepo *repository.CommentRepository, log *slog.Logger) *StatsService {
	return &StatsService{articleRepo: articleRepo, catRepo: catRepo, tagRepo: tagRepo, commentRepo: commentRepo, log: log}
}

// SiteStats 站点统计。
type SiteStats struct {
	ArticleCount int64 `json:"article_count"`
	CategoryCount int64 `json:"category_count"`
	TagCount     int64 `json:"tag_count"`
	CommentCount int64 `json:"comment_count"`
	WordCount    int64 `json:"word_count"`
}

// Stats 返回站点统计。
func (s *StatsService) Stats(ctx context.Context) (*SiteStats, error) {
	stats := &SiteStats{}
	stats.ArticleCount, _ = s.articleRepo.Count()
	stats.CategoryCount, _ = s.catRepo.Count()
	stats.TagCount, _ = s.tagRepo.Count()
	stats.CommentCount, _ = s.commentRepo.Count()
	stats.WordCount, _ = s.articleRepo.SumWordCount()
	s.log.InfoContext(ctx, constants.LOG_STATS_GENERATED, "articles", stats.ArticleCount)
	return stats, nil
}

// Archives 文章按月归档。
func (s *StatsService) Archives(ctx context.Context) ([]model.ArchiveMonth, error) {
	s.log.InfoContext(ctx, constants.LOG_ARCHIVE_GENERATED)
	return s.articleRepo.Archives()
}

// ArticlesByMonth 按月文章。
func (s *StatsService) ArticlesByMonth(ctx context.Context, month string) ([]model.Article, error) {
	return s.articleRepo.ListByMonth(month)
}
