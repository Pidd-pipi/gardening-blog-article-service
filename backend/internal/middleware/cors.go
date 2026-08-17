package middleware

import (
	"github.com/blueship581/gbblog/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors 跨域配置：来源白名单从 APP_CORS_ORIGINS 环境变量读取，生产默认不允许通配符。
func Cors(cfg config.Config) gin.HandlerFunc {
	allowOrigins := cfg.CORSOriginsList()
	if len(allowOrigins) == 0 {
		allowOrigins = []string{"http://localhost:8101"}
	}
	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		AllowCredentials: true,
	})
}
