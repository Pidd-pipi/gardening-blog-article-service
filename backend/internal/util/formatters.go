package util

import (
	"fmt"
	"time"
)

// formatters.go 日期格式化、字数格式化、状态文本、角色文本、评论状态文本（多处耦合）。

// FormatTime 时间格式化。
func FormatTime(t time.Time) string { return t.Format("2006-01-02 15:04:05") }

// FormatDate 日期格式化。
func FormatDate(t time.Time) string { return t.Format("2006-01-02") }

// FormatWordCount 字数格式化（万/千）。
func FormatWordCount(n int) string {
	switch {
	case n >= 10000:
		return fmt.Sprintf("%.1f 万字", float64(n)/10000)
	case n >= 1000:
		return fmt.Sprintf("%.1f 千字", float64(n)/1000)
	default:
		return fmt.Sprintf("%d 字", n)
	}
}

// ArticleStatusText 文章状态中文文案。
func ArticleStatusText(status string) string {
	switch status {
	case "draft":
		return "草稿"
	case "published":
		return "已发布"
	case "scheduled":
		return "定时发布"
	default:
		return "未知"
	}
}

// CommentStatusText 评论状态中文文案。
func CommentStatusText(status string) string {
	switch status {
	case "pending":
		return "待审核"
	case "approved":
		return "已通过"
	case "deleted":
		return "已删除"
	default:
		return "未知"
	}
}

// RoleText 角色中文文案。
func RoleText(role string) string {
	if role == "admin" {
		return "博主"
	}
	return "游客"
}
