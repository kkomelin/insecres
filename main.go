package main

import (
	"fmt"
)

type Fetcher interface {
	// Fetch returns a slice of Page URLs and
	// a slice of resource URLs, which is not secure, found on that page.
	Fetch(url string) (pageUrls []string, nonHTTPSResourceUrls []string, err error)
}

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

	pageUrls, resourceUrls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
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
	Crawl("http://golang.org/", 4, fetcher)
}

// UrlFetcher is Fetcher that returns canned results.
type UrlFetcher map[string]*FetchResult

type FetchResult struct {
	pageUrls     []string
	resourceUrls []string
}

func (f UrlFetcher) Fetch(url string) ([]string, []string, error) {
	if res, ok := f[url]; ok {
		return res.pageUrls, res.resourceUrls, nil
	}
	return nil, nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = UrlFetcher{
	"http://golang.org/": &FetchResult{
		[]string{
			"https://golang.org/image/secure.jpg",
			"http://golang.org/image/insecure.jpg",
		},
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &FetchResult{
		[]string{
			"https://golang.org/image/secure.jpg",
			"http://golang.org/image/insecure.jpg",
		},
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &FetchResult{
		[]string{
			"https://golang.org/image/secure.jpg",
			"http://golang.org/image/insecure.jpg",
		},
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &FetchResult{
		[]string{
			"https://golang.org/image/secure.jpg",
			"http://golang.org/image/insecure.jpg",
		},
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
