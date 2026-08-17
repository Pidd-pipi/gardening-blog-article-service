package handler

import (
	"log/slog"

	"github.com/blueship581/gbblog/internal/dto"
	"github.com/blueship581/gbblog/internal/model"
	"github.com/blueship581/gbblog/internal/service"
	"github.com/blueship581/gbblog/internal/util"
	"github.com/gin-gonic/gin"
)

// CategoryHandler 分类接口。
type CategoryHandler struct {
	svc *service.CategoryService
	log *slog.Logger
}

// NewCategoryHandler 构造分类接口。
func NewCategoryHandler(svc *service.CategoryService, log *slog.Logger) *CategoryHandler {
	return &CategoryHandler{svc: svc, log: log}
}

// Create 创建分类。
func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.BadRequest("分类（Category）参数不合法", err))
		return
	}
	cat, err := h.svc.Create(c.Request.Context(), &model.Category{ParentID: req.ParentID, Name: req.Name, Slug: req.Slug, Description: req.Description, SortOrder: req.SortOrder})
	if err != nil {
		c.Error(err)
		return
	}
	util.Created(c, cat)
}

// Tree 分类树。
func (h *CategoryHandler) Tree(c *gin.Context) {
	tree, err := h.svc.Tree(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, tree)
}

// Detail 分类详情。
func (h *CategoryHandler) Detail(c *gin.Context) {
	cat, err := h.svc.Get(c.Request.Context(), parseUint(c.Param("id")))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, cat)
}

// Update 更新分类。
func (h *CategoryHandler) Update(c *gin.Context) {
	id := parseUint(c.Param("id"))
	var req dto.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.BadRequest("分类（Category）参数不合法", err))
		return
	}
	cat, err := h.svc.Update(c.Request.Context(), id, &model.Category{Name: req.Name, Slug: req.Slug, Description: req.Description, SortOrder: req.SortOrder})
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, cat)
}

// Delete 删除分类。
func (h *CategoryHandler) Delete(c *gin.Context) {
	id := parseUint(c.Param("id"))
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	util.OK(c, gin.H{"deleted": id})
}
