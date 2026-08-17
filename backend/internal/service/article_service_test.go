package service

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/blueship581/gbblog/internal/constants"
	"github.com/blueship581/gbblog/internal/model"
	"github.com/blueship581/gbblog/internal/repository"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type noopRedis struct{}

func (noopRedis) Incr(ctx context.Context, key string) (int64, error) { return 1, nil }

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Category{}, &model.Tag{}, &model.Article{}, &model.Comment{}, &model.ArticleTag{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestArticleService_CreatePublishSearch(t *testing.T) {
	db := newTestDB(t)
	cat := model.Category{Name: "技术", Slug: "tech"}
	if err := db.Create(&cat).Error; err != nil {
		t.Fatal(err)
	}
	tag := model.Tag{Name: "Go", Slug: "go"}
	if err := db.Create(&tag).Error; err != nil {
		t.Fatal(err)
	}
	svc := NewArticleService(
		repository.NewArticleRepository(db), repository.NewTagRepository(db),
		repository.NewCategoryRepository(db), noopRedis{}, testLogger(),
	)
	ctx := context.Background()
	catID := cat.ID
	created, err := svc.Create(ctx, 1, ArticleInput{
		Title: "测试文章", Slug: "test-article", CategoryID: &catID,
		ContentMD: "# 标题", Status: constants.ArticlePublished, TagIDs: []uint{tag.ID},
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.Status != constants.ArticlePublished || len(created.Tags) != 1 {
		t.Fatalf("created article mismatch: %+v", created)
	}
	// slug 冲突
	_, err = svc.Create(ctx, 1, ArticleInput{Title: "重复", Slug: "test-article", Status: constants.ArticleDraft})
	if err == nil {
		t.Fatal("expected slug conflict")
	}
	// 前台列表
	items, total, err := svc.ListPublished(ctx, nil, nil, "", 1, 10)
	if err != nil {
		t.Fatalf("ListPublished() error = %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("ListPublished total=%d len=%d", total, len(items))
	}
	// 搜索
	results, _, err := svc.Search(ctx, "测试", 1, 10)
	if err != nil || len(results) != 1 {
		t.Fatalf("Search() error=%v results=%d", err, len(results))
	}
}
