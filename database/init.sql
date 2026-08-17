-- gb-50-1 gbblog 功能完善的个人博客系统 初始化脚本（容器首次启动自动执行）
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(50) DEFAULT '',
    avatar VARCHAR(255) DEFAULT '',
    bio VARCHAR(500) DEFAULT '',
    social_links TEXT DEFAULT '',
    role VARCHAR(20) DEFAULT 'visitor',
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT,
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description VARCHAR(300) DEFAULT '',
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS articles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    category_id BIGINT,
    title VARCHAR(200) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    summary VARCHAR(500) DEFAULT '',
    content_markdown TEXT DEFAULT '',
    content_html TEXT DEFAULT '',
    cover VARCHAR(255) DEFAULT '',
    status VARCHAR(20) DEFAULT 'draft',
    is_top BOOLEAN DEFAULT FALSE,
    published_at TIMESTAMPTZ,
    view_count INT DEFAULT 0,
    word_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_article_status ON articles(status, published_at);

CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    article_id BIGINT NOT NULL,
    parent_id BIGINT,
    user_id BIGINT,
    nickname VARCHAR(50) DEFAULT '',
    email VARCHAR(100) DEFAULT '',
    content VARCHAR(1000) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_comment_article ON comments(article_id, status);

CREATE TABLE IF NOT EXISTS article_tags (
    article_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (article_id, tag_id)
);

-- 种子数据（密码 bcrypt：admin123 / user123）
INSERT INTO users (username, email, password_hash, nickname, bio, role) VALUES
('admin', 'admin@gbblog.com', '$2a$10$lZmZzCvRaYR.pAKTsPaXnuqcYHGKshT.v9FnKf2a7YCSMd6PmOPq.', '博主', '热爱写作与分享', 'admin'),
('visitor', 'visitor@gbblog.com', '$2a$10$USF27Faq3QyN/MoRKahwlu.q4ZKxNEufvWfFQX4NOfOGr8ChWvoES', '访客', '', 'visitor')
ON CONFLICT (email) DO NOTHING;

INSERT INTO categories (name, slug, description, sort_order) VALUES
('技术', 'tech', '技术分享', 1),
('生活', 'life', '生活随笔', 2)
ON CONFLICT (slug) DO NOTHING;

INSERT INTO categories (name, slug, description, sort_order, parent_id)
SELECT 'Go 语言', 'golang', 'Go 相关', 1, id FROM categories WHERE slug='tech'
ON CONFLICT (slug) DO NOTHING;

INSERT INTO tags (name, slug) VALUES
('Go', 'go'), ('Gin', 'gin'), ('随笔', 'essay'), ('教程', 'tutorial')
ON CONFLICT (slug) DO NOTHING;

INSERT INTO articles (user_id, category_id, title, slug, summary, content_markdown, content_html, status, is_top, view_count, word_count, published_at)
SELECT u.id, c.id, a.title, a.slug, a.summary, a.md, a.html, a.status, a.is_top, a.views, a.words, now()
FROM (VALUES
    ('tech','Go 1.22 企业级项目结构最佳实践','go-122-project-structure','介绍 cmd/internal 分层、依赖注入与错误链规范。','# 项目结构','采用 cmd/ + internal/ 布局。','published',TRUE,128,300),
    ('life','周末小记：城市的烟火气','weekend-notes','记录周末散步的所见所闻。','# 周末小记','巷口的早餐店、街角的猫。','published',FALSE,56,180),
    ('tech','GORM 迁移与种子数据实践','gorm-migration-seed','草稿：待补充。','# GORM 迁移','AutoMigrate 与 init.sql 的配合。','draft',FALSE,0,90)
) AS a(cat_slug, title, slug, summary, md, html, status, is_top, views, words)
JOIN users u ON u.username='admin'
JOIN categories c ON c.slug = a.cat_slug
ON CONFLICT (slug) DO NOTHING;

INSERT INTO article_tags (article_id, tag_id)
SELECT a.id, t.id FROM articles a JOIN tags t ON t.slug='go' WHERE a.slug='go-122-project-structure'
ON CONFLICT DO NOTHING;
INSERT INTO article_tags (article_id, tag_id)
SELECT a.id, t.id FROM articles a JOIN tags t ON t.slug='gin' WHERE a.slug='go-122-project-structure'
ON CONFLICT DO NOTHING;
INSERT INTO article_tags (article_id, tag_id)
SELECT a.id, t.id FROM articles a JOIN tags t ON t.slug='essay' WHERE a.slug='weekend-notes'
ON CONFLICT DO NOTHING;

INSERT INTO comments (article_id, nickname, email, content, status)
SELECT a.id, '路人甲', 'a@example.com', '写得很好，学习了！', 'approved' FROM articles a WHERE a.slug='go-122-project-structure'
ON CONFLICT DO NOTHING;
INSERT INTO comments (article_id, nickname, email, content, status)
SELECT a.id, '新人', 'b@example.com', '请问可以转载吗？', 'pending' FROM articles a WHERE a.slug='go-122-project-structure'
ON CONFLICT DO NOTHING;
