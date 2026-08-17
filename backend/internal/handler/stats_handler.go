package handler

import (
	"log/slog"

	"github.com/blueship581/gbblog/internal/service"
	"github.com/blueship581/gbblog/internal/util"
	"github.com/gin-gonic/gin"
)

// StatsHandler 统计接口。
type StatsHandler struct {
	svc *service.StatsService
	log *slog.Logger
}

// NewStatsHandler 构造统计接口。
func NewStatsHandler(svc *service.StatsService, log *slog.Logger) *StatsHandler {
	return &StatsHandler{svc: svc, log: log}
}

// Stats 站点统计。
func (h *StatsHandler) Stats(c *gin.Context) {
	stats, err := h.svc.Stats(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, stats)
}

// Archives 归档。
func (h *StatsHandler) Archives(c *gin.Context) {
	items, err := h.svc.Archives(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, items)
}

// ArticlesByMonth 按月文章。
func (h *StatsHandler) ArticlesByMonth(c *gin.Context) {
	items, err := h.svc.ArticlesByMonth(c.Request.Context(), c.Query("month"))
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, items)
}
