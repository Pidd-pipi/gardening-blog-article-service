# 执行记录：gb-50-1 功能完善的个人博客系统（gbblog）

- 项目编号/名称：gb-50-1 功能完善的个人博客系统（博客分类，全栈 Web 应用）
- 执行日期：2026-08-16
- 输出目录：/Users/gaobo/repositories/gitlab/评审项目/0-1代码生成提示词/golang-改编提示词/博客主题项目提示词/gb-50-1
- 短名/端口：gbblog，前端 8101 / 后端 3101（容器内 8080）/ PostgreSQL 5601 / Redis 6501
- 技术栈：前端 React 18 + TypeScript + Ant Design 5 + Vite；后端 Go 1.22 + Gin + GORM；PostgreSQL 15 + Redis（go-redis/v9）

## Docker Compose 结果

- `docker compose config --quiet`：通过（中文目录名下）
- `docker compose up -d --build`：成功
- 容器状态：

| 容器 | 状态 | 端口 |
| --- | --- | --- |
| gbblog-db | healthy | 5601->5432 |
| gbblog-redis | healthy | 6501->6379 |
| gbblog-backend | healthy | 3101->8080 |
| gbblog-frontend | healthy | 8101->80 |

- 说明：构建走宿主代理；init.sql 文章种子 JOIN 别名与 slug 不匹配已修复（cat_slug 应为 tech/life）；后端同时提供 `/healthz` 与 `/api/healthz`；视图计数经 Redis 适配器 + DB 自增。

## API 冒烟测试结果（24 项）

| # | 方法 | 路径 | 状态 | 结果摘要 |
| --- | --- | --- | --- | --- |
| 1 | GET | /healthz | 200 | 后端健康检查 |
| 2 | GET | /api/healthz | 200 | Nginx 反代健康检查 |
| 3 | POST | /api/v1/auth/login | 200 | 博主登录获取 JWT |
| 4 | GET | /api/v1/articles | 200 | 已发布文章 2 篇 |
| 5 | GET | /api/v1/articles/top | 200 | 热门文章排行 |
| 6 | GET | /api/v1/articles/:slug | 200 | 详情 views 128→129（计数自增） |
| 7 | GET | /api/v1/categories/tree | 200 | 分类树（技术→Go 语言） |
| 8 | GET | /api/v1/tags | 200 | 标签 4 个 |
| 9 | GET | /api/v1/stats | 200 | 文章 3/分类 3/标签 4/总字数 570 |
| 10 | GET | /api/v1/archives | 200 | 归档月份 2026-08 |
| 11 | GET | /api/v1/archives/articles?month= | 200 | 按月文章 2 篇 |
| 12 | GET | /api/v1/search?q=Go | 200 | 全文搜索命中 1 篇 |
| 13 | POST | /api/v1/comments（游客） | 201 | 评论状态 pending |
| 14 | GET | /api/v1/comments?article_id= | 200 | 已通过评论列表 |
| 15 | POST | /api/v1/articles（admin） | 201 | 创建草稿（Markdown→HTML+字数） |
| 16 | POST | /api/v1/articles/:id/publish | 200 | 草稿→已发布 |
| 17 | POST | /api/v1/categories | 201 | 创建分类 |
| 18 | POST | /api/v1/tags | 201 | 创建标签 |
| 19 | GET | /api/v1/admin/articles | 200 | 后台全量文章 4 篇 |
| 20 | GET | /api/v1/admin/comments | 200 | 评论 3 条（含待审核） |
| 21 | PUT | /api/v1/comments/:id/status | 200 | 评论审核通过 |
| 22 | POST | /api/v1/articles（slug 冲突） | 409 | 别名冲突 |
| 23 | GET | /api/v1/admin/articles（无 token） | 401 | 未授权拦截 |
| 24 | GET | /api/v1/admin/articles（游客） | 403 | RBAC 角色拦截 |

## 浏览器验证（playwright-cli 打包脚本，无外部浏览器）

- http://localhost:8101/：首页渲染——站点统计卡片（文章/分类/标签/评论/总字数）、最新文章（Go 1.22 项目结构、周末小记）、分类导航、热门文章，数据均来自后端 API
- http://localhost:8101/articles/go-122-project-structure：文章详情渲染标题、Markdown 正文、评论（路人甲）与评论提交表单
- http://localhost:8101/admin：后台管理——文章管理表格（redis-cache-practice 已发布、草稿、Go 1.22 已发布）、写文章入口、分类/标签/评论审核 Tab
- http://localhost:8101/search?q=Go：全文搜索命中「Go 1.22 企业级项目结构最佳实践」
- 截图：output/gbblog_home.png、output/gbblog_admin.png
- 结论：页面打开正常、主要功能交互正常、关键业务数据均来自后端 API（Nginx /api 反代）

## README 检查

- 存在 README.md：Docker 一键启动（首选）✅、本地开发 ✅、访问地址与演示账号 ✅、技术栈表格（后端 Go 1.22 + Gin + GORM）✅、目录结构 ✅、环境变量 ✅、API 一览 ✅、Docker 部署说明 ✅、枚举出现位置清单（ArticleStatus/CommentStatus/UserRole）✅、License ✅

## 其他质量项

- `cd backend && go mod tidy && go build ./...`：通过
- `go test ./...`：通过（service/article_service_test、util/slug_test、util/markdown_test 表驱动单测）
- `go vet ./...`：通过
- `cd frontend && npm run build`：通过（tsc -b && vite build 零错误）
- 结构强制清单、严禁合并职责到单一文件、屎山代码设计要求（log_templates ≥25 条、formatters/messages 多耦合、状态机多处定义）均已落实

## 关闭确认

- `docker compose down -v --remove-orphans` 已执行，容器与命名卷清理，无本项目残留。

## Git 提交

- commit 哈希：见仓库 `git log`
