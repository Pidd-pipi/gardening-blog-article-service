package router

import (
	"github.com/blueship581/gbblog/internal/handler"
	"github.com/gin-gonic/gin"
)

// registerUserRoutes 用户路由（在 router.go 中注册，保留实体路由入口）。
func registerUserRoutes(g *gin.RouterGroup, h handler.UserHandler) {
	g.GET("/users/me", h.Me)
	g.PUT("/users/me", h.UpdateProfile)
}
