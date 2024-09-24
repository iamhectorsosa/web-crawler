package main

import (
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		description string
		htmlBody    string
		inputURL    string
		want        []string
	}{
		{
			description: "returns unique absolute, relative, valid hashes, and external URLs",

			htmlBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Path one</span>
		</a>
		<a href="https://hectorsosa.me/path/two">
			<span>Path two</span>
		</a>
		<a href="https://hectorsosa.me/path/two#hash">
			<span>Path two with hash</span>
		</a>
		<a href="../up-a-level">
			<span>Up a Level</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Another path</span>
		</a>
	</body>
</html>
`,
			inputURL: "https://hectorsosa.me",
			want:     []string{"https://hectorsosa.me/path/one", "https://hectorsosa.me/path/two", "https://hectorsosa.me/up-a-level", "https://other.com/path/one"},
		},
		{
			description: "ignores malformed, missing hrefs and hashes",
			htmlBody: `
<html>
		<body>
			<a>
				<span>No Href</span>
			</a>
			<a href="ftp://malformed.url/path">
				Malformed link
			</a>
			<a href="">
				Empty href
			</a>
			<a href="#hash">
				Just hash
			</a>
			<a href="https://hectorsosa.me/about#hash">
				Also a hash
			</a>
		</body>
</html>
`,
			inputURL: "https://hectorsosa.me/about",
			want:     []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got, err := getURLsFromHTML(tt.htmlBody, tt.inputURL)
			if err != nil {
				t.Errorf("error running function: %v", err)
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("unexpected URLs: want %v, got %v", tt.want, got)
			}
		})
	}
}
