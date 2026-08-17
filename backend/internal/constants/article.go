package constants

// ArticleStatus 文章状态枚举（README 枚举出现位置清单必列）。
const (
	ArticleDraft     = "draft"     // 草稿
	ArticlePublished = "published" // 已发布
	ArticleScheduled = "scheduled" // 定时发布
)

// ArticleStatuses 全部文章状态。
var ArticleStatuses = []string{ArticleDraft, ArticlePublished, ArticleScheduled}
