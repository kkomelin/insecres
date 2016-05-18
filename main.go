package main

import (
	"fmt"
)

func fetch(url string, ch chan []string, quit chan int, depth int, registry *Registry) {
	if depth <= 0 {
		quit <- 0
		return
	}

	// Try to add the url to the registry if it has not yet been added.
	if !registry.IsNew(url) {
		fmt.Printf("depth: %d: skip %s: duplicate\n", depth, url)
		return
	}

	fetcher := UrlFetcher{}

	resourceUrls, pageUrls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Errorf("Error occured: %v", err)
		return
	}
	fmt.Printf("depth: %d: found: %s %q\n", depth, url, pageUrls)

	registry.MarkAsProcessed(url)

	ch <- resourceUrls
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {

	registry := &Registry{processed: make(map[string]int)}

	ch := make(chan []string)
	quit := make(chan int)

	go fetch(url, ch, quit, depth, registry)

	for {
		select {
		case urls := <-ch:
			depth --;
			for _, url := range urls {
				go fetch(url, ch, quit, depth, registry)
			}
		case <-quit:
			fmt.Printf("depth: %d: quit\n", depth)
			fmt.Println("-----")
			fmt.Println(registry)
			return
		}
	}

	return
}

func main() {
	uri := "https://komelin.com"

	fetcher := UrlFetcher{}

	resourceUrls, linkUrls, err := fetcher.Fetch(uri)
	if err != nil {
		fmt.Errorf("Error occured: %v", err)
		return
	}

	fmt.Println(resourceUrls)

	fmt.Println(linkUrls)

	//Crawl("http://golang.org/", 4, fetcher)
}
