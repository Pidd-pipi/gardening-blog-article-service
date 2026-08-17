package util

import (
	"regexp"
	"strings"
)

var (
	slugRe      = regexp.MustCompile(`[^a-z0-9-]`)
	hyphenRe    = regexp.MustCompile(`-{2,}`)
)

// Slugify 标题转 slug（保留 ascii，中文转拼音占位+序号由调用方补充）。
func Slugify(title string) string {
	s := strings.ToLower(strings.TrimSpace(title))
	s = strings.ReplaceAll(s, " ", "-")
	s = slugRe.ReplaceAllString(s, "-")
	s = hyphenRe.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// IsValidSlug 校验 slug 格式。
func IsValidSlug(slug string) bool {
	if slug == "" || len(slug) > 100 {
		return false
	}
	return slugRe.ReplaceAllString(slug, "") == slug
}
