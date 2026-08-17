package util

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		title string
		want  string
	}{
		{title: "Hello World", want: "hello-world"},
		{title: "Go 1.22 项目结构", want: "go-1-22"},
		{title: "  Spaces  Trim  ", want: "spaces-trim"},
		{title: "A--B", want: "a-b"},
	}
	for _, tt := range tests {
		if got := Slugify(tt.title); got != tt.want {
			t.Fatalf("Slugify(%q) = %q, want %q", tt.title, got, tt.want)
		}
	}
}

func TestIsValidSlug(t *testing.T) {
	tests := []struct {
		slug string
		want bool
	}{
		{slug: "hello-world", want: true},
		{slug: "go-122", want: true},
		{slug: "Hello World", want: false},
		{slug: "", want: false},
		{slug: "bad/slug", want: false},
	}
	for _, tt := range tests {
		if got := IsValidSlug(tt.slug); got != tt.want {
			t.Fatalf("IsValidSlug(%q) = %v, want %v", tt.slug, got, tt.want)
		}
	}
}
