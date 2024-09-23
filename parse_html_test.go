package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseHTML(t *testing.T) {
	t.Run("returns parsed HTML", func(t *testing.T) {
		want := "<html><body>Hello World</body></html>"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(want))
		}))
		defer server.Close()

		got, err := parseHTML(server.URL)

		if err != nil {
			t.Errorf("error parsing HTML: %v", err)
		}

		if want != got {
			t.Errorf("return error: want %s, got %s", want, got)
		}
	})

	t.Run("errors with status code other than 200", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		_, err := parseHTML(server.URL)

		if err == nil {
			t.Error("expected status code error")
		}
	})

	t.Run("errors with non text/html content", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello World"))
		}))
		defer server.Close()

		_, err := parseHTML(server.URL)

		if err == nil {
			t.Error("expected content-type error")
		}
	})
}
