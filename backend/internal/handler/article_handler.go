package handler

import (
	"log/slog"
	"strconv"

	"github.com/blueship581/gbblog/internal/service"
	"github.com/blueship581/gbblog/internal/util"
	"github.com/gin-gonic/gin"
)

// ArticleHandler 文章接口。
type ArticleHandler struct {
	svc *service.ArticleService
	log *slog.Logger
}

// NewArticleHandler 构造文章接口。
func NewArticleHandler(svc *service.ArticleService, log *slog.Logger) *ArticleHandler {
	return &ArticleHandler{svc: svc, log: log}
}

func ptrUint(v uint) *uint { return &v }

// Create 创建文章。
func (h *ArticleHandler) Create(c *gin.Context) {
	var req service.ArticleInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.BadRequest("文章（Article）参数不合法", err))
		return
	}
	a, err := h.svc.Create(c.Request.Context(), userID(c), req)
	if err != nil {
		c.Error(err)
		return
	}
	util.Created(c, a)
}

// Update 更新文章。
func (h *ArticleHandler) Update(c *gin.Context) {
	id := parseUint(c.Param("id"))
	var req service.ArticleInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(util.BadRequest("文章（Article）参数不合法", err))
		return
	}
	a, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, a)
}

// Delete 删除文章。
func (h *ArticleHandler) Delete(c *gin.Context) {
	id := parseUint(c.Param("id"))
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	util.OK(c, gin.H{"deleted": id})
}

// Publish 发布文章。
func (h *ArticleHandler) Publish(c *gin.Context) {
	id := parseUint(c.Param("id"))
	a, err := h.svc.Publish(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, a)
}

// GetPublished 前台查看文章。
func (h *ArticleHandler) GetPublished(c *gin.Context) {
	a, err := h.svc.GetPublishedBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, a)
}

// GetAdmin 后台查看文章。
func (h *ArticleHandler) GetAdmin(c *gin.Context) {
	id := parseUint(c.Param("id"))
	a, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, a)
}

// ListPublished 前台文章列表。
func (h *ArticleHandler) ListPublished(c *gin.Context) {
	page := parseQueryInt(c.Query("page"), 1)
	pageSize := parseQueryInt(c.Query("page_size"), 10)
	var categoryID, tagID *uint
	if v := c.Query("category_id"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			categoryID = ptrUint(uint(n))
		}
	}
	if v := c.Query("tag_id"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			tagID = ptrUint(uint(n))
		}
	}
	items, total, err := h.svc.ListPublished(c.Request.Context(), categoryID, tagID, c.Query("keyword"), page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, util.PageData{List: items, Total: total, Page: page, Size: pageSize})
}

// ListAll 后台文章列表。
func (h *ArticleHandler) ListAll(c *gin.Context) {
	page := parseQueryInt(c.Query("page"), 1)
	pageSize := parseQueryInt(c.Query("page_size"), 10)
	items, total, err := h.svc.ListAll(c.Request.Context(), c.Query("status"), page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, util.PageData{List: items, Total: total, Page: page, Size: pageSize})
}

// Search 全文搜索。
func (h *ArticleHandler) Search(c *gin.Context) {
	page := parseQueryInt(c.Query("page"), 1)
	pageSize := parseQueryInt(c.Query("page_size"), 10)
	items, total, err := h.svc.Search(c.Request.Context(), c.Query("q"), page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, util.PageData{List: items, Total: total, Page: page, Size: pageSize})
}

// TopViewed 热门文章。
func (h *ArticleHandler) TopViewed(c *gin.Context) {
	items, err := h.svc.TopViewed(c.Request.Context(), 10)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, items)
}
