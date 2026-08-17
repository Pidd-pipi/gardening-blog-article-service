package dto

// RegisterRequest 注册请求。
type RegisterRequest struct {
	Username string `json:"username" binding:"required,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname" binding:"max=50"`
}

// LoginRequest 登录请求。
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest 更新资料请求。
type UpdateProfileRequest struct {
	Nickname    string `json:"nickname" binding:"max=50"`
	Avatar      string `json:"avatar" binding:"max=255"`
	Bio         string `json:"bio" binding:"max=500"`
	SocialLinks string `json:"social_links"`
}

// TokenResponse 令牌响应。
type TokenResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// CategoryRequest 分类请求。
type CategoryRequest struct {
	ParentID    *uint  `json:"parent_id"`
	Name        string `json:"name" binding:"required,max=50"`
	Slug        string `json:"slug" binding:"required,max=100"`
	Description string `json:"description" binding:"max=300"`
	SortOrder   int    `json:"sort_order"`
}

// TagRequest 标签请求。
type TagRequest struct {
	Name string `json:"name" binding:"required,max=50"`
	Slug string `json:"slug" binding:"required,max=100"`
}

// CommentRequest 评论请求。
type CommentRequest struct {
	ArticleID uint   `json:"article_id" binding:"required"`
	ParentID  *uint  `json:"parent_id"`
	Nickname  string `json:"nickname" binding:"max=50"`
	Email     string `json:"email" binding:"max=100"`
	Content   string `json:"content" binding:"required,max=1000"`
}

// UpdateCommentStatusRequest 评论审核请求。
type UpdateCommentStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending approved deleted"`
}
