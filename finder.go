package main

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ResourceAndLinkFinder encapsulates logic that is used for finding page link urls and resource urls..
type ResourceAndLinkFinder struct{}

// Fetch page by url and returns response body.
func (f ResourceAndLinkFinder) Fetch(url string) (responseBody io.ReadCloser, err error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := http.Client{Transport: transport}

	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	//defer response.Body.Close()

	return response.Body, nil
}

// Parse takes a reader object and returns a slice of insecure resource urls
// found in the HTML.
// It does not close the reader. The reader should be closed from the outside.
func (f ResourceAndLinkFinder) Parse(baseUrl string, httpBody io.Reader) (resourceUrls []string, linkUrls []string, err error) {

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
		case f.isResourceToken(token):
			uri, err := f.processResourceToken(token)
			if err == nil {
				resourceMap[uri] = true
			}
		case f.isLinkToken(token):
			uri, err := f.processLinkToken(token, baseUrl)
			if err == nil {
				linkMap[uri] = true
			}
		}
	}

	resourceUrls = make([]string, 0, len(resourceMap))

	for k := range resourceMap {
		resourceUrls = append(resourceUrls, k)
	}

	linkUrls = make([]string, 0, len(linkMap))

	for k := range linkMap {
		linkUrls = append(linkUrls, k)
	}

	return resourceUrls, linkUrls, nil
}

// Determine whether the token passed is a resource token.
func (f ResourceAndLinkFinder) isResourceToken(token html.Token) bool {
	switch {
	case token.Type == html.SelfClosingTagToken && token.DataAtom.String() == "img":
		return true
	case token.Type == html.StartTagToken && token.DataAtom.String() == "iframe":
		return true
	case token.Type == html.StartTagToken && token.DataAtom.String() == "object":
		return true
	case token.Type == html.StartTagToken && token.DataAtom.String() == "video":
		return true
	case token.Type == html.StartTagToken && token.DataAtom.String() == "audio":
		return true
	case token.Type == html.StartTagToken && token.DataAtom.String() == "source":
		return true
	case token.Type == html.StartTagToken && token.DataAtom.String() == "track":
		return true
	default:
		return false
	}
}

// Process resource token in order to get a url of the resource.
func (f ResourceAndLinkFinder) processResourceToken(token html.Token) (string, error) {

	tag := token.DataAtom.String()

	// Loop for tag attributes.
	for _, attr := range token.Attr {

		if tag == "object" {
			if attr.Key != "data" {
				continue
			}

		} else {
			if attr.Key != "src" && attr.Key != "poster" {
				continue
			}
		}

		uri, err := url.Parse(attr.Val)
		if err != nil {
			return "", err
		}

		// Ignore relative and secure urls.
		if !uri.IsAbs() || uri.Scheme == "https" || (uri.Host != "" && strings.HasPrefix(uri.String(), "//")) {
			return "", fmt.Errorf("Uri is relative or secure. Skipped.")
		}

		return uri.String(), nil
	}

	return "", fmt.Errorf("Src has not been found. Skipped.")
}

// Determine whether the token passed is a link token.
func (f ResourceAndLinkFinder) isLinkToken(token html.Token) bool {
	switch {
	case token.Type == html.StartTagToken && token.DataAtom.String() == "a":
		return true
	default:
		return false
	}
}

// Process <A> token in order to get an absolute url of the link.
func (f ResourceAndLinkFinder) processLinkToken(token html.Token, base string) (string, error) {

	// Loop for tag attributes.
	for _, attr := range token.Attr {
		if attr.Key != "href" {
			continue
		}

		// Ignore anchors.
		if strings.HasPrefix(attr.Val, "#") {
			return "", fmt.Errorf("Url is an anchor. Skipped.")
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
		if uri.IsAbs() || (uri.Host != "" && strings.HasPrefix(uri.String(), "//")) {

			// Ignore external urls considering urls w/ WWW and w/o WWW as the same.
			if strings.TrimPrefix(uri.Host, "www.") != strings.TrimPrefix(baseUrl.Host, "www.") {
				return "", fmt.Errorf("Url is expernal. Skipped.")
			}

			return strings.TrimSuffix(uri.String(), "/"), nil
		}

		// Make it absolute if it's relative.
		absoluteUrl := f.convertToAbsolute(uri, baseUrl)

		return strings.TrimSuffix(absoluteUrl.String(), "/"), nil
	}

	return "", fmt.Errorf("Src has not been found. Skipped.")
}

// Convert a relative url to absolute.
func (f ResourceAndLinkFinder) convertToAbsolute(source, base *url.URL) *url.URL {
	return base.ResolveReference(source)
}
