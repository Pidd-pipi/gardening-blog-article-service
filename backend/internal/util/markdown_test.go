package util

import (
	"strings"
	"testing"
)

func TestRenderMarkdown(t *testing.T) {
	html := RenderMarkdown("# 标题\n\n**加粗** 文本")
	if !strings.Contains(html, "<h1") || !strings.Contains(html, "<strong>") {
		t.Fatalf("markdown not rendered: %s", html)
	}
	if RenderMarkdown("") != "" {
		t.Fatal("empty markdown should render empty")
	}
}

func TestPlainText(t *testing.T) {
	got := PlainText("<p>Hello <b>World</b></p>")
	if got != "Hello World" {
		t.Fatalf("PlainText() = %q", got)
	}
}
