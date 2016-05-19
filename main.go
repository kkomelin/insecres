package main

import (
	"fmt"
	"time"
)

// Goroutine function fetches and parses the passed url in order to find insecure resources and next urls to fetch from.
func fetchUrl(url string, queue chan []string, registry *Registry) {

	// Lock url so that no one other goroutine can process it.
	registry.MarkAsProcessed(url)

	fetcher := InsecureResourceFetcher{}

	insecureResourceUrls, pageUrls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Errorf("Error occured: %v", err)
		return
	}

	for _, insecureResourceUrl := range insecureResourceUrls {
		fmt.Printf("%s: %s\n", url, insecureResourceUrl)
	}

	// TODO: calculate speed of processing by batches (as it is now) or by single url and adjust if necessary.
	queue <- pageUrls
}

// Crawl pages starting with url and find insecure resources.
func crawl(url string, fetcher Fetcher) {

	registry := &Registry{processed: make(map[string]int)}

	queue := make(chan []string)

	go fetchUrl(url, queue, registry)

	tick := time.Tick(1000 * time.Millisecond)

	flag := false
	for {
		select {
		case urls := <-queue:
			flag = false
			for _, url := range urls {

				// Ignore processed urls.
				if !registry.IsNew(url) {
					continue
				}
				go fetchUrl(url, queue, registry)
			}
		case <-tick:
			if flag {
				fmt.Println("-----")
				fmt.Printf("log:\n")
				fmt.Println(registry)
				return
			} else {
				flag = true
			}
		}
	}

	return
}

func main() {
	// TODO: Pass site url as an argument.
	uri := "http://drupal7"

	fetcher := InsecureResourceFetcher{}

	crawl(uri, fetcher)
}
