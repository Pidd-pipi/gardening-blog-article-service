package router

import "github.com/gin-gonic/gin"

// registerPublicRoutes 公开只读路由。
func registerPublicRoutes(g *gin.RouterGroup, h Handlers) {
	g.GET("/articles", h.Article.ListPublished)
	g.GET("/articles/top", h.Article.TopViewed)
	g.GET("/articles/:slug", h.Article.GetPublished)
	g.GET("/categories/tree", h.Category.Tree)
	g.GET("/tags", h.Tag.List)
	g.GET("/tags/:slug", h.Tag.Detail)
	g.GET("/comments", h.Comment.ListByArticle)
	g.GET("/stats", h.Stats.Stats)
	g.GET("/archives", h.Stats.Archives)
	g.GET("/archives/articles", h.Stats.ArticlesByMonth)
	g.GET("/search", h.Article.Search)
}
