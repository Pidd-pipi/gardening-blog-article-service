package router

import "github.com/gin-gonic/gin"

// registerAdminRoutes 管理路由（仅 admin）。
func registerAdminRoutes(g *gin.RouterGroup, h Handlers) {
	g.GET("/admin/articles", h.Article.ListAll)
	g.GET("/admin/articles/:id", h.Article.GetAdmin)
	g.POST("/articles", h.Article.Create)
	g.PUT("/articles/:id", h.Article.Update)
	g.DELETE("/articles/:id", h.Article.Delete)
	g.POST("/articles/:id/publish", h.Article.Publish)

	g.POST("/categories", h.Category.Create)
	g.PUT("/categories/:id", h.Category.Update)
	g.DELETE("/categories/:id", h.Category.Delete)

	g.POST("/tags", h.Tag.Create)
	g.PUT("/tags/:id", h.Tag.Update)
	g.DELETE("/tags/:id", h.Tag.Delete)

	g.GET("/admin/comments", h.Comment.ListAll)
	g.PUT("/comments/:id/status", h.Comment.UpdateStatus)
}
