package main

import (
	"fmt"
	"time"
)

func process(url string, queue chan []string, registry *Registry) {

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

// Uses fetcher to recursively crawl
// pages starting with url.
func Crawl(url string, fetcher Fetcher) {

	registry := &Registry{processed: make(map[string]int)}

	queue := make(chan []string)

	go process(url, queue, registry)

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
				go process(url, queue, registry)
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
	uri := "http://drupal7"

	fetcher := InsecureResourceFetcher{}

	Crawl(uri, fetcher)
}
