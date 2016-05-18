package main

import (
	"fmt"
	"time"
)

func fetch(url string, ch chan []string, quit chan int, registry *Registry) {
	//if depth <= 0 {
	//	quit <- 0
	//	return
	//}

	// Try to add the url to the registry if it has not yet been added.
	if !registry.IsNew(url) {
		//fmt.Printf("skip %s: duplicate\n", url)
		return
	}

	fetcher := InsecureResourceFetcher{}

	insecureResourceUrls, pageUrls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Errorf("Error occured: %v", err)
		return
	}

	for _, insecureResourceUrl := range insecureResourceUrls {
		fmt.Printf("%s: %s\n", url, insecureResourceUrl)
	}

	registry.MarkAsProcessed(url)

	ch <- pageUrls
}

// Uses fetcher to recursively crawl
// pages starting with url.
func Crawl(url string, fetcher Fetcher) {

	registry := &Registry{processed: make(map[string]int)}

	ch := make(chan []string)
	quit := make(chan int)

	go fetch(url, ch, quit, registry)

	tick := time.Tick(1000 * time.Millisecond)

	i := 0
	for {
		select {
		case urls := <-ch:
			i = 0
			for _, url := range urls {
				go fetch(url, ch, quit, registry)
			}
		case <-tick:
			if i > 1 {
				fmt.Println("-----")
				fmt.Printf("log:\n")
				fmt.Println(registry)
				return
			} else {
				i++
				fmt.Println(i)
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
