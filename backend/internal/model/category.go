package model

import "time"

// Category 文章分类（多级自关联）。
type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ParentID    *uint     `gorm:"index" json:"parent_id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Slug        string    `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	Description string    `gorm:"size:300" json:"description"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	Parent      *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}
