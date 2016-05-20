package main

import (
	"crypto/tls"
	"net/http"
)

type Fetcher interface {
	// Returns a slice of URLs found on the page and
	// a slice of resource URLs, which is not secure, found on the page.
	Fetch(url string) (pageUrls []string, nonHTTPSResourceUrls []string, err error)
}

// UrlFetcher contains a map of fetched urls and resources grouped by .
type InsecureResourceFetcher map[string]*FetchResult

type FetchResult struct {
	pageUrls     []string
	resourceUrls []string
}

// Fetch insecure resources and urls to crawl next.
func (f InsecureResourceFetcher) Fetch(url string) (resourceUrls []string, linkUrls []string, err error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := http.Client{Transport: transport}

	response, err := client.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	insecureResources, links, err := f.Parse(url, response.Body)
	if err != nil {
		return nil, nil, err
	}

	return insecureResources, links, nil
}
