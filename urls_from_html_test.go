package main

import (
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		description string
		htmlBody    string
		rawBaseURL  string
		want        []string
	}{
		{
			description: "returns absolute and relative URLs",

			htmlBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			rawBaseURL: "https://blog.boot.dev",
			want:       []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			description: "returns only relative URLs",
			htmlBody: `
<html>
		<body>
			<a href="/internal/page1">
				<span>Page 1</span>
			</a>
			<a href="about">
				<span>About</span>
			</a>
		</body>
</html>
`,
			rawBaseURL: "https://example.com",
			want:       []string{"https://example.com/internal/page1", "https://example.com/about"},
		},
		{
			description: "handles malformed and missing hrefs",
			htmlBody: `
<html>
		<body>
			<a href="https://valid.com/path">
				<span>Valid Link</span>
			</a>
			<a>
				<span>No Href</span>
			</a>
			<a href="ftp://malformed.url/path">
				<span>Malformed Link</span>
			</a>
			<a href="">
				<span>Empty Href</span>
			</a>
		</body>
</html>
`,
			rawBaseURL: "https://example.com",
			want:       []string{"https://valid.com/path"},
		},
		{
			description: "returns base URL with a path and relative URLs",
			htmlBody: `
<html>
		<body>
			<a href="relative/path1">
				<span>Relative Link 1</span>
			</a>
			<a href="/absolute/path2">
				<span>Absolute Link 2</span>
			</a>
			<a href="../up-a-level">
				<span>Up a Level</span>
			</a>
		</body>
</html>
`,
			rawBaseURL: "https://example.com/blog/",
			want:       []string{"https://example.com/blog/relative/path1", "https://example.com/absolute/path2", "https://example.com/up-a-level"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got, err := getURLsFromHTML(tt.htmlBody, tt.rawBaseURL)
			if err != nil {
				t.Errorf("error getting URLs: %v", err)
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("unexpected URLs: want %v, got %v", tt.want, got)
			}
		})
	}
}
