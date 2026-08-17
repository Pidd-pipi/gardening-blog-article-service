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

// CommentService 评论服务：提交、嵌套回复、审核。
type CommentService struct {
	repo   *repository.CommentRepository
	log    *slog.Logger
}

// NewCommentService 构造评论服务。
func NewCommentService(repo *repository.CommentRepository, log *slog.Logger) *CommentService {
	return &CommentService{repo: repo, log: log}
}

// CommentInput 评论入参。
type CommentInput struct {
	ArticleID uint   `json:"article_id" binding:"required"`
	ParentID  *uint  `json:"parent_id"`
	Nickname  string `json:"nickname" binding:"max=50"`
	Email     string `json:"email" binding:"max=100"`
	Content   string `json:"content" binding:"required,max=1000"`
}

// Create 提交评论（游客默认待审核；登录用户直接通过）。
func (s *CommentService) Create(ctx context.Context, userID *uint, input CommentInput) (*model.Comment, error) {
	status := constants.CommentPending
	if userID != nil {
		status = constants.CommentApproved
	}
	c := &model.Comment{
		ArticleID: input.ArticleID, ParentID: input.ParentID, UserID: userID,
		Nickname: input.Nickname, Email: input.Email, Content: input.Content, Status: status,
	}
	if err := s.repo.Create(c); err != nil {
		return nil, util.LogError(s.log, constants.LOG_COMMENT_CREATED, fmt.Errorf("create comment: %w", err))
	}
	s.log.InfoContext(ctx, constants.LOG_COMMENT_CREATED, "comment_id", c.ID, "status", status)
	return c, nil
}

// ListByArticle 文章评论。
func (s *CommentService) ListByArticle(ctx context.Context, articleID uint) ([]model.Comment, error) {
	return s.repo.ListByArticle(articleID)
}

// ListAll 后台评论管理。
func (s *CommentService) ListAll(ctx context.Context, status string, page, pageSize int) ([]model.Comment, int64, error) {
	return s.repo.ListAll(status, page, pageSize)
}

// UpdateStatus 审核/删除评论。
func (s *CommentService) UpdateStatus(ctx context.Context, id uint, status string) error {
	valid := map[string]bool{constants.CommentPending: true, constants.CommentApproved: true, constants.CommentDeleted: true}
	if !valid[status] {
		return util.BadRequest("评论状态（Comment.status）不合法", errors.New("invalid status"))
	}
	if _, err := s.repo.FindByID(id); err != nil {
		return util.NotFoundError(constants.MsgCommentNotFound, err)
	}
	if err := s.repo.UpdateStatus(id, status); err != nil {
		return util.LogError(s.log, constants.LOG_COMMENT_APPROVED, fmt.Errorf("update comment status: %w", err))
	}
	s.log.InfoContext(ctx, statusLogTemplate(status), "comment_id", id, "status", status)
	return nil
}

func statusLogTemplate(status string) string {
	switch status {
	case constants.CommentApproved:
		return constants.LOG_COMMENT_APPROVED
	case constants.CommentDeleted:
		return constants.LOG_COMMENT_DELETED
	default:
		return constants.LOG_COMMENT_REJECTED
	}
}
