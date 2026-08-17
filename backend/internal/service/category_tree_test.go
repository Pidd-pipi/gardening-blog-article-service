package service

import (
	"context"
	"testing"

	"github.com/blueship581/gbblog/internal/model"
	"github.com/blueship581/gbblog/internal/repository"
)

func TestCategoryTree_AttachesChildren(t *testing.T) {
	db := newTestDB(t)
	root := model.Category{Name: "技术", Slug: "tech", SortOrder: 1}
	if err := db.Create(&root).Error; err != nil {
		t.Fatal(err)
	}
	parentID := root.ID
	child := model.Category{Name: "Go", Slug: "go", ParentID: &parentID, SortOrder: 1}
	if err := db.Create(&child).Error; err != nil {
		t.Fatal(err)
	}

	svc := NewCategoryService(repository.NewCategoryRepository(db), testLogger())
	tree, err := svc.Tree(context.Background())
	if err != nil {
		t.Fatalf("Tree() error = %v", err)
	}
	if len(tree) != 1 {
		t.Fatalf("root categories = %d, want 1", len(tree))
	}
	if len(tree[0].Children) != 1 {
		t.Fatalf("root children = %d, want 1", len(tree[0].Children))
	}
}
