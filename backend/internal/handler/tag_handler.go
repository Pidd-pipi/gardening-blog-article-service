package handler

import (
	"log/slog"

	"github.com/blueship581/gbblog/internal/dto"
	"github.com/blueship581/gbblog/internal/service"
	"github.com/blueship581/gbblog/internal/util"
	"github.com/gin-gonic/gin"
)

// TagHandler 标签接口。
type TagHandler struct {
	svc *service.TagService
	log *slog.Logger
}

// NewTagHandler 构造标签接口。
func NewTagHandler(svc *service.TagService, log *slog.Logger) *TagHandler {
	return &TagHandler{svc: svc, log: log}
}

// Create 创建标签。
func (h *TagHandler) Create(c *gin.Context) {
	var req dto.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.BadRequest("标签（Tag）参数不合法", err))
		return
	}
	tag, err := h.svc.Create(c.Request.Context(), req.Name, req.Slug)
	if err != nil {
		c.Error(err)
		return
	}
	util.Created(c, tag)
}

// List 标签列表。
func (h *TagHandler) List(c *gin.Context) {
	tags, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, tags)
}

// Detail 标签详情。
func (h *TagHandler) Detail(c *gin.Context) {
	tag, err := h.svc.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, tag)
}

// Update 更新标签。
func (h *TagHandler) Update(c *gin.Context) {
	id := parseUint(c.Param("id"))
	var req dto.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.BadRequest("标签（Tag）参数不合法", err))
		return
	}
	tag, err := h.svc.Update(c.Request.Context(), id, req.Name, req.Slug)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, tag)
}

// Delete 删除标签。
func (h *TagHandler) Delete(c *gin.Context) {
	id := parseUint(c.Param("id"))
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	util.OK(c, gin.H{"deleted": id})
}
