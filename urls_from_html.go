package main

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, inputURL string) ([]string, error) {
	htmlNode, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		return nil, err
	}

	hrefsMap := make(map[string]bool)

	var extractHrefs func(*html.Node)
	extractHrefs = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				// Only hrefs that aren't empty
				if attr.Key == "href" && strings.TrimSpace(attr.Val) != "" {
					hrefValue := attr.Val
					parsedURL, err := url.Parse(hrefValue)
					if err != nil {
						break
					}
					// Only valid schemas
					hasValidScheme := (parsedURL.Scheme == "" || parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
					// Only hashes belonging to other pages
					isHash := parsedURL.Fragment != "" && (parsedURL.Path == "" || baseURL.Path == parsedURL.Path)
					if hasValidScheme && !isHash {
						// If ref is absolute URL, then ignores base and return copy of ref
						resolvedURL := baseURL.ResolveReference(parsedURL)
						// Strip hashes from the resolved URL
						resolvedURL.Fragment = ""
						// Check for duplication
						if _, exists := hrefsMap[resolvedURL.String()]; !exists {
							hrefsMap[resolvedURL.String()] = true
						}
					}
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			extractHrefs(child)
		}
	}

	extractHrefs(htmlNode)

	hrefs := make([]string, 0, len(hrefsMap))
	for href := range hrefsMap {
		hrefs = append(hrefs, href)
	}

	return hrefs, nil
}
