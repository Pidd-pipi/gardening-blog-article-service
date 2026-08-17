# InkLog（个人博客系统）

> 项目类型：全栈 Web 应用（博客）

面向个人博主的内容发布与阅读平台，支持文章发布、分类管理、标签体系、评论互动、全文搜索和归档浏览，博主可登录后台管理文章、分类、标签与评论审核。

## 快速启动（Docker Compose 一键部署，首选）

```bash
# 1. 首次启动前复制环境变量
cp .env.example .env

# 2. 启动全部服务（前端 + 后端 + PostgreSQL + Redis）
docker compose up -d

# 3. 查看健康状态
docker compose ps
```

访问地址：

- 前端：http://localhost:8101
- 后端 API：http://localhost:3101
- 健康检查：http://localhost:3101/healthz

演示账号：

| 角色 | 账号 | 密码 |
| --- | --- | --- |
| 博主（管理员） | admin@gbblog.com | admin123 |
| 游客 | visitor@gbblog.com | user123 |

## 本地开发

```bash
# 前端（React 18 + TS + Vite + Ant Design 5）
cd frontend
npm install
npm run dev        # http://localhost:8101

# 后端（Go）
cd backend
go mod tidy
go run ./cmd/server
```

## 技术栈

| 层次 | 技术 |
| --- | --- |
| 前端 | React 18 + TypeScript + Vite + Ant Design 5 |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | PostgreSQL 15 |
| 缓存 | Redis（go-redis/v9） |
| Markdown | github.com/russross/blackfriday/v2 |
| 认证 | JWT (github.com/golang-jwt/jwt/v5) + RBAC（博主/游客） |
| 校验 | github.com/go-playground/validator/v10 |

## 项目目录结构

```
gb-50-1/
├── docker-compose.yml        # 前端/后端/PostgreSQL/Redis 编排
├── .env.example              # 环境变量示例
├── database/init.sql         # 建表 + 种子数据（首次启动自动执行）
├── frontend/
│   ├── Dockerfile            # Node 构建 + Nginx 托管
│   ├── nginx.conf            # SPA 路由 + /api 反代
│   └── src/
│       ├── api/              # user/article/category/tag/comment/stats
│       ├── stores/           # authStore/userStore/articleStore/categoryStore/commentStore
│       ├── components/common/# ArticleCard/ArticleStatusBadge/CategoryTree/CommentList/MarkdownEditor/StatsCards/TimelineList
│       ├── hooks/            # useAuth/useArticleStats/useArchive
│       ├── pages/            # Home/ArticleDetail/CategoryArticles/TagArticles/Archives/Search/Admin/Login
│       ├── router/           # index.tsx + guards.ts
│       ├── utils/            # dateFormat/markdown/slug/request
│       └── constants/        # article/comment/user/errorCodes
└── backend/
    ├── cmd/server/main.go    # 入口：装配依赖、启动 Gin + Redis
    └── internal/
        ├── config/           # 环境变量解析
        ├── model/            # user/category/tag/article/comment/article_tag
        ├── repository/       # 按实体分文件
        ├── service/          # 按实体分文件 + stats
        ├── handler/          # 按实体分文件
        ├── router/           # router.go + public/admin/users
        ├── middleware/       # auth/rbac/rate_limiter/error_handler/logger/cors
        ├── dto/              # 请求/响应结构体
        ├── constants/        # article/comment/user/error_codes/log_templates/messages
        └── util/             # jwt/logger/formatters/app_error/slug/markdown/search/response
```

## 环境变量

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| COMPOSE_PROJECT_NAME | gbblog | Docker Compose 项目短名 |
| DB_NAME | gbblog_db | 数据库名 |
| DB_USER | gbblog_user | 数据库用户 |
| DB_PASSWORD | gbblog_pwd | 数据库密码 |
| DB_ROOT_PASSWORD | gbblog_root | root 密码（占位） |
| REDIS_PORT | 6501 | Redis 对外端口 |
| JWT_SECRET | change_me_to_a_long_random_string | JWT 签名密钥 |
| APP_CORS_ORIGINS | http://localhost:8101 | CORS 来源白名单（逗号分隔，生产请显式配置，不允许 `*`） |
| FRONTEND_PORT | 8101 | 前端对外端口 |
| BACKEND_PORT | 3101 | 后端对外端口 |
| DB_PORT | 5601 | 数据库对外端口 |

## API 接口清单

| 方法 | 路径 | 认证 | 说明 |
| --- | --- | --- | --- |
| GET | `/healthz` | 无 | 健康检查 |
| GET | `/api/healthz` | 无 | 健康检查 |
| POST | `/api/v1/auth/register` | 无 | 注册（JWT） |
| POST | `/api/v1/auth/login` | 无 | 登录（JWT） |
| GET | `/api/v1/articles` | 无 | 前台已发布文章列表 |
| GET | `/api/v1/articles/top` | 无 | 热门文章 Top N |
| GET | `/api/v1/articles/:slug` | 无 | 按 slug 查看文章 |
| GET | `/api/v1/categories/tree` | 无 | 分类树 |
| GET | `/api/v1/tags` | 无 | 标签列表 |
| GET | `/api/v1/tags/:slug` | 无 | 标签详情 |
| GET | `/api/v1/comments` | 无 | 按文章查询评论 |
| POST | `/api/v1/comments` | 可选登录 | 发表评论 |
| GET | `/api/v1/search` | 无 | 全文搜索 |
| GET | `/api/v1/stats` | 无 | 站点统计 |
| GET | `/api/v1/archives` | 无 | 文章按月归档 |
| GET | `/api/v1/archives/articles` | 无 | 按月查询已发布文章 |
| GET | `/api/v1/users/me` | JWT | 当前用户资料 |
| PUT | `/api/v1/users/me` | JWT | 更新当前用户资料 |
| GET | `/api/v1/admin/articles` | JWT（admin） | 后台文章列表 |
| GET | `/api/v1/admin/articles/:id` | JWT（admin） | 后台文章详情 |
| POST | `/api/v1/articles` | JWT（admin） | 创建文章 |
| PUT | `/api/v1/articles/:id` | JWT（admin） | 更新文章 |
| DELETE | `/api/v1/articles/:id` | JWT（admin） | 删除文章 |
| POST | `/api/v1/articles/:id/publish` | JWT（admin） | 发布文章 |
| POST | `/api/v1/categories` | JWT（admin） | 创建分类 |
| PUT | `/api/v1/categories/:id` | JWT（admin） | 更新分类 |
| DELETE | `/api/v1/categories/:id` | JWT（admin） | 删除分类 |
| POST | `/api/v1/tags` | JWT（admin） | 创建标签 |
| PUT | `/api/v1/tags/:id` | JWT（admin） | 更新标签 |
| DELETE | `/api/v1/tags/:id` | JWT（admin） | 删除标签 |
| GET | `/api/v1/admin/comments` | JWT（admin） | 全部评论列表 |
| PUT | `/api/v1/comments/:id/status` | JWT（admin） | 评论审核 |

> 所有响应头均携带 `X-Request-ID`，响应体统一为 `{code, message, data}`；除 `/healthz`、`/api/healthz`、注册/登录与公开只读接口外，其余接口需携带 `Authorization: Bearer <JWT>`。

## Docker 部署说明

- 端口映射：前端 `8101:80`、后端 `3101:8080`、数据库 `5601:5432`、Redis `6501:6379`
- 数据卷：`db_data` 命名卷持久化
- 依赖顺序：backend 等待 db/redis 健康；frontend 等待 backend 健康
- 常见问题：端口冲突修改 `.env`；数据重置 `docker compose down -v` 后重建

## 枚举出现位置清单

### ArticleStatus（文章状态：draft/published/scheduled）

- 后端：`backend/internal/constants/article.go`（定义）、`backend/internal/model/article.go`（模型）、`backend/internal/service/article_service.go`（状态机：创建/发布/更新）、`backend/internal/repository/article_repository.go`（按状态查询）、`backend/internal/util/formatters.go`（中文文案）、`backend/internal/constants/log_templates.go`（日志）、`backend/internal/constants/error_codes.go`（错误码）
- 前端：`frontend/src/constants/article.ts`（定义）、`frontend/src/components/common/ArticleStatusBadge.tsx`（徽标）、`frontend/src/pages/Admin.tsx`（文章管理/按钮显隐）、`frontend/src/types/index.ts`（类型）

### CommentStatus（评论状态：pending/approved/deleted）

- 后端：`backend/internal/constants/comment.go`（定义）、`backend/internal/model/comment.go`（模型）、`backend/internal/service/comment_service.go`（审核状态机）、`backend/internal/repository/comment_repository.go`（按状态查询）、`backend/internal/util/formatters.go`（中文文案）、`backend/internal/constants/log_templates.go`（日志）
- 前端：`frontend/src/constants/comment.ts`（定义）、`frontend/src/pages/Admin.tsx`（评论审核）、`frontend/src/components/common/CommentList.tsx`（展示）、`frontend/src/types/index.ts`（类型）

### UserRole（用户角色：admin/visitor）

- 后端：`backend/internal/constants/user.go`、`backend/internal/model/user.go`、`backend/internal/middleware/rbac.go`、`backend/internal/router/admin.go`（RequireRole("admin")）、`backend/internal/util/formatters.go`
- 前端：`frontend/src/constants/user.ts`、`frontend/src/hooks/useAuth.ts`、`frontend/src/router/guards.ts`、`frontend/src/components/Shell.tsx`（后台入口显隐）、`frontend/src/pages/Admin.tsx`

## License

MIT
