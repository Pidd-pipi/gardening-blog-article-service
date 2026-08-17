package util

import "strings"

// ScoreMatch 关键词匹配评分（标题命中 > 摘要命中 > 正文命中）。
func ScoreMatch(title, summary, content, keyword string) int {
	kw := strings.ToLower(keyword)
	score := 0
	if strings.Contains(strings.ToLower(title), kw) {
		score += 10
	}
	if strings.Contains(strings.ToLower(summary), kw) {
		score += 5
	}
	if strings.Contains(strings.ToLower(content), kw) {
		score += 1
	}
	return score
}
