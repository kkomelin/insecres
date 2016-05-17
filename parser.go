package main

import (
	"golang.org/x/net/html"
	"io"
	"net/url"
)

type Parser interface {
	// Parse HTML in order to find non-HTTPS resources.
	Parse(httpBody io.Reader) (resourceUrls []string, err error)
}

type Page struct {
	insecureResources []string
}

// Takes a reader object and returns a slice of strings of insecure resource urls
// found in the HTML.
// It does not close the reader passed to it.
func (p *Page) Parse(httpBody io.Reader) (resourceUrls []string, err error) {

	linkMap := make(map[string]bool)

	page := html.NewTokenizer(httpBody)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			break
		}
		token := page.Token()

		// Ignore all non IMG tags.
		if tokenType != html.SelfClosingTagToken || token.DataAtom.String() != "img" {
			continue
		}

		// Loop for tag attributes.
		for _, attr := range token.Attr {
			if attr.Key != "src" {
				continue
			}

			// TODO: extract method.

			uri, err := url.Parse(attr.Val)
			if err != nil {
				break
			}

			// Check absolute urls only.
			if !uri.IsAbs() {
				break
			}

			// Ignore secure links.
			if uri.Scheme == "https" {
				break
			}

			linkMap[uri.String()] = true
			break
		}
	}

	links := make([]string, 0, len(linkMap))

	for k, _ := range linkMap {
		links = append(links, k)
	}

	return links, nil
}
