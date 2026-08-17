package util

import (
	"bytes"
	"strings"

	"github.com/russross/blackfriday/v2"
)

// RenderMarkdown Markdown 转 HTML（安全转义 + 常用扩展）。
func RenderMarkdown(md string) string {
	if strings.TrimSpace(md) == "" {
		return ""
	}
	html := blackfriday.Run([]byte(md),
		blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.AutoHeadingIDs|blackfriday.HardLineBreak))
	return string(html)
}

// PlainText 提取纯文本（用于摘要/搜索）。
func PlainText(html string) string {
	var buf bytes.Buffer
	inTag := false
	for _, r := range html {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			buf.WriteRune(r)
		}
	}
	return strings.TrimSpace(buf.String())
}
