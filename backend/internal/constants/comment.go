package constants

// CommentStatus 评论状态枚举（README 枚举出现位置清单必列）。
const (
	CommentPending  = "pending"  // 待审核
	CommentApproved = "approved" // 已通过
	CommentDeleted  = "deleted"  // 已删除
)

// MaxReplyDepth 评论最多支持两层嵌套回复。
const MaxReplyDepth = 2

// CommentStatuses 全部评论状态。
var CommentStatuses = []string{CommentPending, CommentApproved, CommentDeleted}
