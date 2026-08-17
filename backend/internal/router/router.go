package router

import (
	"log/slog"
	"net/http"

	"github.com/blueship581/gbblog/internal/config"
	"github.com/blueship581/gbblog/internal/handler"
	"github.com/blueship581/gbblog/internal/middleware"
	"github.com/gin-gonic/gin"
)

// Handlers 全部接口处理器集合。
type Handlers struct {
	User     *handler.UserHandler
	Category *handler.CategoryHandler
	Tag      *handler.TagHandler
	Article  *handler.ArticleHandler
	Comment  *handler.CommentHandler
	Stats    *handler.StatsHandler
}

// New 装配 Gin 路由。
func New(cfg config.Config, log *slog.Logger, h Handlers, limiter *middleware.RateLimiter) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandler(log))
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger(log))
	r.Use(middleware.Cors(cfg))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": gin.H{"status": "up"}})
	})
	r.GET("/api/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": gin.H{"status": "up"}})
	})

	auth := middleware.AuthRequired(cfg.JWTSecret)
	optionalAuth := middleware.OptionalAuth(cfg.JWTSecret)
	rate := limiter.Limit()

	api := r.Group("/api/v1")
	api.POST("/auth/register", rate, h.User.Register)
	api.POST("/auth/login", rate, h.User.Login)

	// 公开只读路由
	pub := api.Group("")
	pub.Use(rate)
	registerPublicRoutes(pub, h)

	// 可选登录（游客评论）
	pubOpt := api.Group("")
	pubOpt.Use(optionalAuth, rate)
	pubOpt.POST("/comments", h.Comment.Create)

	// 用户路由（任意登录用户）
	user := api.Group("")
	user.Use(auth, rate)
	user.GET("/users/me", h.User.Me)
	user.PUT("/users/me", h.User.UpdateProfile)

	// 管理路由（仅 admin）
	admin := api.Group("")
	admin.Use(auth, middleware.RequireRole("admin"), rate)
	registerAdminRoutes(admin, h)
	return r
}
