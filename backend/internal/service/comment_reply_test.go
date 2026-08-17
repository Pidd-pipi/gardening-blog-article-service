package service

import (
	"context"
	"testing"

	"github.com/blueship581/gbblog/internal/constants"
	"github.com/blueship581/gbblog/internal/model"
	"github.com/blueship581/gbblog/internal/repository"
)

func TestCommentService_ListByArticleLoadsReplies(t *testing.T) {
	db := newTestDB(t)
	article := model.Article{Title: "文章", Slug: "article-reply", Status: constants.ArticlePublished}
	if err := db.Create(&article).Error; err != nil {
		t.Fatal(err)
	}
	root := model.Comment{ArticleID: article.ID, Content: "根评论", Status: constants.CommentApproved}
	if err := db.Create(&root).Error; err != nil {
		t.Fatal(err)
	}
	parentID := root.ID
	reply := model.Comment{ArticleID: article.ID, ParentID: &parentID, Content: "回复", Status: constants.CommentApproved}
	if err := db.Create(&reply).Error; err != nil {
		t.Fatal(err)
	}

	svc := NewCommentService(repository.NewCommentRepository(db), testLogger())
	roots, err := svc.ListByArticle(context.Background(), article.ID)
	if err != nil {
		t.Fatalf("ListByArticle() error = %v", err)
	}
	if len(roots) != 1 || len(roots[0].Replies) != 1 {
		t.Fatalf("roots=%d replies=%d, want 1/1", len(roots), len(roots[0].Replies))
	}
}
