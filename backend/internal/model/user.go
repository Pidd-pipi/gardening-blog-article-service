package model

import "time"

// User 用户：博主管理员或游客。
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email        string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Nickname     string    `gorm:"size:50" json:"nickname"`
	Avatar       string    `gorm:"size:255" json:"avatar"`
	Bio          string    `gorm:"size:500" json:"bio"`
	SocialLinks  string    `gorm:"type:text" json:"social_links"`
	Role         string    `gorm:"size:20;default:visitor" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
