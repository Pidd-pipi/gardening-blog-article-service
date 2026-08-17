package service

import (
	"context"
	"testing"

	"github.com/blueship581/gbblog/internal/model"
	"github.com/blueship581/gbblog/internal/repository"
)

func TestArticleService_UpdateDeduplicatesTagIDs(t *testing.T) {
	db := newTestDB(t)
	tag := model.Tag{Name: "Go", Slug: "go"}
	if err := db.Create(&tag).Error; err != nil {
		t.Fatal(err)
	}
	svc := NewArticleService(
		repository.NewArticleRepository(db), repository.NewTagRepository(db),
		repository.NewCategoryRepository(db), noopRedis{}, testLogger(),
	)
	ctx := context.Background()

	created, err := svc.Create(ctx, 1, ArticleInput{Title: "测试文章", Slug: "test-article", ContentMD: "# 标题", TagIDs: []uint{tag.ID}})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	updated, err := svc.Update(ctx, created.ID, ArticleInput{Title: "测试文章", Slug: "test-article", ContentMD: "# 标题", TagIDs: []uint{tag.ID, tag.ID}})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if len(updated.Tags) != 1 {
		t.Fatalf("updated tags = %d, want 1", len(updated.Tags))
	}
}
