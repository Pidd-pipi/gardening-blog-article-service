package model

import "time"

// Comment 评论（最多三层嵌套回复）。
type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ArticleID uint      `gorm:"index;not null" json:"article_id"`
	ParentID  *uint     `gorm:"index" json:"parent_id"`
	UserID    *uint     `gorm:"index" json:"user_id"`
	Nickname  string    `gorm:"size:50" json:"nickname"`
	Email     string    `gorm:"size:100" json:"email"`
	Content   string    `gorm:"size:1000;not null" json:"content"`
	Status    string    `gorm:"size:20;default:pending;index" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Article   *Article  `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
	Replies   []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}
