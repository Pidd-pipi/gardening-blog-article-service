package service

import (
	"context"
	"testing"
	"time"

	"github.com/blueship581/gbblog/internal/constants"
	"github.com/blueship581/gbblog/internal/model"
	"github.com/blueship581/gbblog/internal/repository"
)

func TestArticleListByTagCount(t *testing.T) {
	db := newTestDB(t)
	tag := model.Tag{Name: "Go", Slug: "go"}
	if err := db.Create(&tag).Error; err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	withTag := model.Article{Title: "带标签", Slug: "with-tag", Status: constants.ArticlePublished, PublishedAt: &now}
	if err := db.Create(&withTag).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&model.ArticleTag{ArticleID: withTag.ID, TagID: tag.ID}).Error; err != nil {
		t.Fatal(err)
	}
	other := model.Article{Title: "无标签", Slug: "without-tag", Status: constants.ArticlePublished, PublishedAt: &now}
	if err := db.Create(&other).Error; err != nil {
		t.Fatal(err)
	}

	svc := NewArticleService(repository.NewArticleRepository(db), repository.NewTagRepository(db), repository.NewCategoryRepository(db), noopRedis{}, testLogger())
	tagID := tag.ID
	items, total, err := svc.ListPublished(context.Background(), nil, &tagID, "", 1, 10)
	if err != nil {
		t.Fatalf("ListPublished() error = %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("total=%d len=%d, want total=1 len=1", total, len(items))
	}
}
