package main

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	htmlReader := strings.NewReader(htmlBody)
	htmlNode, err := html.Parse(htmlReader)
	if err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}

	hrefs := make([]string, 0)

	var extractHrefs func(*html.Node)
	extractHrefs = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				// Only hrefs that aren't empty
				if attr.Key == "href" && strings.TrimSpace(attr.Val) != "" {
					hrefUrl := attr.Val
					parsedURL, err := url.Parse(hrefUrl)
					hasValidScheme := (parsedURL.Scheme == "" || parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
					isHash := parsedURL.Fragment != "" && parsedURL.Path == "" && parsedURL.Host == ""
					if err == nil && hasValidScheme && !isHash {
						// If ref is absolute URL, then ignores base and return copy of ref
						resolvedURL := baseURL.ResolveReference(parsedURL)
						hrefs = append(hrefs, resolvedURL.String())
					}
					break
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			extractHrefs(child)
		}
	}

	extractHrefs(htmlNode)

	return hrefs, nil
}
