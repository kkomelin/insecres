package main

import (
	"golang.org/x/net/html"
	"io"
	"net/url"
	"errors"
	"strings"
)

type Parser interface {
	// Parse HTML in order to find non-HTTPS resources.
	Parse(baseUrl string, httpBody io.Reader) (resourceUrls []string, linkUrls []string, err error)
}

// Takes a reader object and returns a slice of insecure resource urls
// found in the HTML.
// It does not close the reader. The reader should be closed from the outside.
func (f UrlFetcher) Parse(baseUrl string, httpBody io.Reader) (resourceUrls []string, linkUrls []string, err error) {

	resourceMap := make(map[string]bool)
	linkMap := make(map[string]bool)

	page := html.NewTokenizer(httpBody)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			break
		}
		token := page.Token()

		switch {
		// Find all insecure urls.
		case tokenType == html.SelfClosingTagToken && token.DataAtom.String() == "img":
			uri, err := f.processImageToken(token)
			if err == nil {
				resourceMap[uri] = true
			}
		case tokenType == html.StartTagToken && token.DataAtom.String() == "a":
			uri, err := f.processLinkToken(token, baseUrl)
			if err == nil {
				linkMap[uri] = true
			}
		}
	}

	resourceUrls = make([]string, 0, len(resourceMap))

	for k, _ := range resourceMap {
		resourceUrls = append(resourceUrls, k)
	}

	linkUrls = make([]string, 0, len(linkMap))

	for k, _ := range linkMap {
		linkUrls = append(linkUrls, k)
	}

	return resourceUrls, linkUrls, nil
}

func (f UrlFetcher) processImageToken(token html.Token) (string, error) {
	// Loop for tag attributes.
	for _, attr := range token.Attr {
		if attr.Key != "src" {
			continue
		}

		uri, err := url.Parse(attr.Val)
		if err != nil {
			return "", err
		}

		// Ignore relative and secure urls.
		if !uri.IsAbs() || uri.Scheme == "https" {
			return "", errors.New("Uri is relative or secure. Skipped.")
		}

		return uri.String(), nil
	}

	return "", errors.New("Src has not been found. Skipped.");
}

func (f UrlFetcher) processLinkToken(token html.Token, base string) (string, error) {
	// Loop for tag attributes.
	for _, attr := range token.Attr {
		if attr.Key != "href" {
			continue
		}

		// Ignore anchors.
		if strings.HasPrefix(attr.Val, "#") {
			return "",  errors.New("Url is an anchor. Skipped.")
		}

		uri, err := url.Parse(attr.Val)
		if err != nil {
			return "", err
		}

		baseUrl, err := url.Parse(base)
		if err != nil {
			return "", err
		}

		// Return result if the uri is absolute.
		if uri.IsAbs() {

			// Ignore external urls.
			// TODO: consider urls with WWW and with WWW as the same.
			if uri.Host != baseUrl.Host {
				return "", errors.New("Url is expernal. Skipped.")
			}

			return uri.String(), nil
		}

		// Make it absolute if it's relative.
		absoluteUrl := f.convertToAbsolute(uri, baseUrl)

		return absoluteUrl.String(), nil
	}

	return "", errors.New("Src has not been found. Skipped.");
}

func (f UrlFetcher) convertToAbsolute(href, base *url.URL) (*url.URL) {
	return base.ResolveReference(href)
}
