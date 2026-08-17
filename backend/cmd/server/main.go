package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/blueship581/gbblog/internal/config"
	"github.com/blueship581/gbblog/internal/constants"
	"github.com/blueship581/gbblog/internal/handler"
	"github.com/blueship581/gbblog/internal/middleware"
	"github.com/blueship581/gbblog/internal/model"
	"github.com/blueship581/gbblog/internal/repository"
	"github.com/blueship581/gbblog/internal/router"
	"github.com/blueship581/gbblog/internal/service"
	"github.com/blueship581/gbblog/internal/util"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("load config: %w", err))
	}
	log := util.NewLogger()

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
	if err != nil {
		panic(fmt.Errorf("open database: %w", err))
	}
	if err := migrateAndSeed(db, log); err != nil {
		panic(fmt.Errorf("migrate database: %w", err))
	}
	log.Info(constants.LOG_DB_INITIALIZED)

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr(), Password: cfg.RedisPassword})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Warn(constants.LOG_REDIS_ERROR, "error", err)
	} else {
		log.Info(constants.LOG_REDIS_CONNECTED, "addr", cfg.RedisAddr())
	}

	userRepo := repository.NewUserRepository(db)
	catRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)
	articleRepo := repository.NewArticleRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	userSvc := service.NewUserService(userRepo, cfg.JWTSecret, cfg.JWTExpireHours, log)
	catSvc := service.NewCategoryService(catRepo, log)
	tagSvc := service.NewTagService(tagRepo, log)
	articleSvc := service.NewArticleService(articleRepo, tagRepo, catRepo, &redisAdapter{client: rdb, log: log}, log)
	commentSvc := service.NewCommentService(commentRepo, log)
	statsSvc := service.NewStatsService(articleRepo, catRepo, tagRepo, commentRepo, log)

	h := router.Handlers{
		User:     handler.NewUserHandler(userSvc, log),
		Category: handler.NewCategoryHandler(catSvc, log),
		Tag:      handler.NewTagHandler(tagSvc, log),
		Article:  handler.NewArticleHandler(articleSvc, log),
		Comment:  handler.NewCommentHandler(commentSvc, log),
		Stats:    handler.NewStatsHandler(statsSvc, log),
	}
	r := router.New(cfg, log, h, middleware.NewRateLimiter(cfg.RateLimitPerMin))

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		log.Info(constants.LOG_SERVER_STARTED, "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(constants.LOG_SERVER_STARTED, "error", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error(constants.LOG_SERVER_SHUTDOWN, "error", err)
	}
}

// redisAdapter 适配 go-redis 客户端到 service.redisClient。
type redisAdapter struct {
	client *redis.Client
	log    *slog.Logger
}

func (a *redisAdapter) Incr(ctx context.Context, key string) (int64, error) {
	n, err := a.client.Incr(ctx, key).Result()
	if err != nil {
		a.log.Warn(constants.LOG_REDIS_ERROR, "key", key, "error", err)
	}
	return n, err
}

// migrateAndSeed 自动迁移并注入种子数据；init.sql 已建表则跳过。
func migrateAndSeed(db *gorm.DB, log *slog.Logger) error {
	var tableCount int64
	if err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public' AND table_name='articles'").Scan(&tableCount).Error; err != nil {
		return err
	}
	if tableCount > 0 {
		return nil
	}
	if err := db.AutoMigrate(
		&model.User{}, &model.Category{}, &model.Tag{}, &model.Article{}, &model.Comment{}, &model.ArticleTag{},
	); err != nil {
		return err
	}
	adminHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	userHash, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	users := []model.User{
		{Username: "admin", Email: "admin@gbblog.com", PasswordHash: string(adminHash), Nickname: "博主", Role: constants.RoleAdmin, Bio: "热爱写作与分享"},
		{Username: "visitor", Email: "visitor@gbblog.com", PasswordHash: string(userHash), Nickname: "访客", Role: constants.RoleVisitor},
	}
	if err := db.Create(&users).Error; err != nil {
		return err
	}
	cats := []model.Category{
		{Name: "技术", Slug: "tech", Description: "技术分享", SortOrder: 1},
		{Name: "生活", Slug: "life", Description: "生活随笔", SortOrder: 2},
		{Name: "Go 语言", Slug: "golang", ParentID: ptr(1), Description: "Go 相关", SortOrder: 1},
	}
	if err := db.Create(&cats).Error; err != nil {
		return err
	}
	tags := []model.Tag{
		{Name: "Go", Slug: "go"}, {Name: "Gin", Slug: "gin"}, {Name: "随笔", Slug: "essay"}, {Name: "教程", Slug: "tutorial"},
	}
	if err := db.Create(&tags).Error; err != nil {
		return err
	}
	now := time.Now()
	articles := []model.Article{
		{
			UserID: 1, CategoryID: ptr(1), Title: "Go 1.22 企业级项目结构最佳实践", Slug: "go-122-project-structure",
			Summary: "介绍 cmd/internal 分层、依赖注入与错误链规范。",
			ContentMarkdown: "# 项目结构\n\n采用 `cmd/` + `internal/` 布局，handler → service → repository 单向依赖。",
			ContentHTML: util.RenderMarkdown("# 项目结构\n\n采用 `cmd/` + `internal/` 布局，handler → service → repository 单向依赖。"),
			Status: "published", IsTop: true, ViewCount: 128, WordCount: 300, PublishedAt: &now,
		},
		{
			UserID: 1, CategoryID: ptr(2), Title: "周末小记：城市的烟火气", Slug: "weekend-notes",
			Summary: "记录周末散步的所见所闻。",
			ContentMarkdown: "# 周末小记\n\n巷口的早餐店、街角的猫，都是生活的注脚。",
			ContentHTML: util.RenderMarkdown("# 周末小记\n\n巷口的早餐店、街角的猫，都是生活的注脚。"),
			Status: "published", IsTop: false, ViewCount: 56, WordCount: 180, PublishedAt: &now,
		},
		{
			UserID: 1, CategoryID: ptr(1), Title: "GORM 迁移与种子数据实践", Slug: "gorm-migration-seed",
			Summary: "草稿：待补充。",
			ContentMarkdown: "# GORM 迁移\n\nAutoMigrate 与 init.sql 的配合。",
			ContentHTML: util.RenderMarkdown("# GORM 迁移\n\nAutoMigrate 与 init.sql 的配合。"),
			Status: "draft", IsTop: false, WordCount: 90,
		},
	}
	if err := db.Create(&articles).Error; err != nil {
		return err
	}
	if err := db.Create(&[]model.ArticleTag{{ArticleID: 1, TagID: 1}, {ArticleID: 1, TagID: 2}, {ArticleID: 2, TagID: 3}}).Error; err != nil {
		return err
	}
	comments := []model.Comment{
		{ArticleID: 1, Nickname: "路人甲", Email: "a@example.com", Content: "写得很好，学习了！", Status: constants.CommentApproved},
		{ArticleID: 1, ParentID: nil, Nickname: "新人", Email: "b@example.com", Content: "请问可以转载吗？", Status: constants.CommentPending},
	}
	if err := db.Create(&comments).Error; err != nil {
		return err
	}
	log.Info(constants.LOG_DB_INITIALIZED, "seed", "ok")
	return nil
}

func ptr(v int) *uint {
	u := uint(v)
	return &u
}
