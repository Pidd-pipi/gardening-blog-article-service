package service

import (
	"context"
	"testing"
	"time"

	"github.com/blueship581/gbblog/internal/constants"
	"github.com/blueship581/gbblog/internal/model"
	"github.com/blueship581/gbblog/internal/repository"
)

func TestSearchBySummary(t *testing.T) {
	db := newTestDB(t)
	now := time.Now()
	a := model.Article{Title: "标题", Summary: "这里是 Go 实战笔记", Slug: "search-summary", Status: constants.ArticlePublished, PublishedAt: &now}
	if err := db.Create(&a).Error; err != nil {
		t.Fatal(err)
	}

	svc := NewArticleService(repository.NewArticleRepository(db), repository.NewTagRepository(db), repository.NewCategoryRepository(db), noopRedis{}, testLogger())
	items, total, err := svc.Search(context.Background(), "go", 1, 10)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("total=%d len=%d, want 1/1", total, len(items))
	}
}
