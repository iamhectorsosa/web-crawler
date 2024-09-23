package main

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		description string
		inputURL    string
		want        string
	}{
		{
			description: "remove https scheme",
			inputURL:    "https://blog.boot.dev/path/",
			want:        "blog.boot.dev/path",
		},
		{
			description: "remove http scheme",
			inputURL:    "http://blog.boot.dev/path",
			want:        "blog.boot.dev/path",
		},
		{
			description: "remove query params",
			inputURL:    "blog.boot.dev/path?sort=asc",
			want:        "blog.boot.dev/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got, err := normalizeURL(tt.inputURL)
			if err != nil {
				t.Errorf("error normalizing URL: %v", err)
			}

			if tt.want != got {
				t.Errorf("unexpected URL: want %s, got %s", tt.want, got)
			}
		})
	}
}
