package handler

import (
	"log/slog"

	"github.com/blueship581/gbblog/internal/dto"
	"github.com/blueship581/gbblog/internal/middleware"
	"github.com/blueship581/gbblog/internal/service"
	"github.com/blueship581/gbblog/internal/util"
	"github.com/gin-gonic/gin"
)

// CommentHandler 评论接口。
type CommentHandler struct {
	svc *service.CommentService
	log *slog.Logger
}

// NewCommentHandler 构造评论接口。
func NewCommentHandler(svc *service.CommentService, log *slog.Logger) *CommentHandler {
	return &CommentHandler{svc: svc, log: log}
}

// Create 提交评论（可选登录）。
func (h *CommentHandler) Create(c *gin.Context) {
	var req dto.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.BadRequest("评论（Comment）参数不合法", err))
		return
	}
	var uid *uint
	if v, ok := c.Get(middleware.UserIDKey); ok {
		if u, ok := v.(uint); ok {
			uid = &u
		}
	}
	comment, err := h.svc.Create(c.Request.Context(), uid, service.CommentInput{
		ArticleID: req.ArticleID, ParentID: req.ParentID,
		Nickname: req.Nickname, Email: req.Email, Content: req.Content,
	})
	if err != nil {
		c.Error(err)
		return
	}
	util.Created(c, comment)
}

// ListByArticle 文章评论。
func (h *CommentHandler) ListByArticle(c *gin.Context) {
	items, err := h.svc.ListByArticle(c.Request.Context(), parseUint(c.Query("article_id")))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, items)
}

// ListAll 后台评论管理。
func (h *CommentHandler) ListAll(c *gin.Context) {
	page := parseQueryInt(c.Query("page"), 1)
	pageSize := parseQueryInt(c.Query("page_size"), 10)
	items, total, err := h.svc.ListAll(c.Request.Context(), c.Query("status"), page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, util.PageData{List: items, Total: total, Page: page, Size: pageSize})
}

// UpdateStatus 审核/删除评论。
func (h *CommentHandler) UpdateStatus(c *gin.Context) {
	id := parseUint(c.Param("id"))
	var req dto.UpdateCommentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.BadRequest("评论状态（Comment.status）不合法", err))
		return
	}
	if err := h.svc.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		c.Error(err)
		return
	}
	util.OK(c, gin.H{"id": id, "status": req.Status})
}
