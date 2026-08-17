package model

import "time"

// ArticleTag 文章-标签关联。
type ArticleTag struct {
	ArticleID uint      `gorm:"primaryKey" json:"article_id"`
	TagID     uint      `gorm:"primaryKey" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}
