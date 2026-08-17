package model

import "time"

// Article 文章。
type Article struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	UserID         uint       `gorm:"index;not null" json:"user_id"`
	CategoryID     *uint      `gorm:"index" json:"category_id"`
	Title          string     `gorm:"size:200;not null" json:"title"`
	Slug           string     `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	Summary        string     `gorm:"size:500" json:"summary"`
	ContentMarkdown string    `gorm:"type:text" json:"content_markdown"`
	ContentHTML    string     `gorm:"type:text" json:"content_html"`
	Cover          string     `gorm:"size:255" json:"cover"`
	Status         string     `gorm:"size:20;default:draft;index" json:"status"`
	IsTop          bool       `gorm:"default:false" json:"is_top"`
	PublishedAt    *time.Time `gorm:"index" json:"published_at"`
	ViewCount      int        `gorm:"default:0" json:"view_count"`
	WordCount      int        `gorm:"default:0" json:"word_count"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	User           User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Category       *Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags           []Tag      `gorm:"many2many:article_tags;" json:"tags,omitempty"`
}
