package main

import (
	"net/http"
	"crypto/tls"
)

type Fetcher interface {
	// Fetch returns a slice of Page URLs and
	// a slice of resource URLs, which is not secure, found on that page.
	Fetch(url string) (pageUrls []string, nonHTTPSResourceUrls []string, err error)
}

// UrlFetcher is Fetcher that returns canned results.
type UrlFetcher map[string]*FetchResult

type FetchResult struct {
	pageUrls     []string
	resourceUrls []string
}

func (f UrlFetcher) Fetch(url string) (resourceUrls []string, linkUrls []string, err error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := http.Client{Transport: transport}

	resp, err := client.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	insecureResources, links, err := f.Parse(url, resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return insecureResources, links, nil
}
