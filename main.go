package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"math/rand"
)

const (
	// The goal is to wait for some milliseconds before exit so that all goroutines can finish.
	beforeEngTimeout int = 2000
)

// Goroutine function fetches and parses the passed url in order to find insecure resources and next urls to fetch from.
func fetchUrl(url string, queue chan string, registry *Registry) {

	fmt.Print(".")

	// Lock url so that no one other goroutine can process it.
	registry.MarkAsProcessed(url)

	insecureResourceUrls, pageUrls, err := (ResourceAndLinkFetcher{}).Fetch(url)
	if err != nil {
		fmt.Printf("Error occured: %s\n", err)
		return
	}

	displayPageResources(url, insecureResourceUrls)

	for _, url := range pageUrls {
		// Random pause before sending to the main thread.
		delayBetweenRequests()
		queue <- url
	}
}

// Displays page resources.
func displayPageResources(url string, resources []string) {
	if (len(resources) > 0) {
		fmt.Printf("\n%s:\n", url)
		for _, insecureResourceUrl := range resources {
			fmt.Printf("- %s\n", insecureResourceUrl)
		}
	}
}

// Implement random pause before sending the next request to
// (no more than beforeEngTimeout/ and no less than beforeEngTimeout/4 constant).
// It is one of the measures to prevent banning by the server.
func delayBetweenRequests() {
	randNum := randomInRange(beforeEngTimeout/4, beforeEngTimeout/2)
	time.Sleep(time.Duration(randNum) * time.Millisecond)
}

// Returns a random number in a given range.
// The idea has been borrowed from http://golangcookbook.blogspot.ru/2012/11/generate-random-number-in-given-range.html
// and improved.
func randomInRange(min, max int) int {
	rand.Seed(time.Now().Unix() + rand.Int63n(200))
	return rand.Intn(max - min) + min
}

// Crawl pages starting with url and find insecure resources.
func crawl(url string, fetcher Fetcher) {

	url = strings.TrimSuffix(url, "/")

	registry := &Registry{processed: make(map[string]int)}

	queue := make(chan string)

	go fetchUrl(url, queue, registry)

	tick := time.Tick(time.Duration(beforeEngTimeout) * time.Millisecond)

	flag := false
	for {
		select {
		case url := <-queue:
			flag = false

			// Ignore processed urls.
			if !registry.IsNew(url) {
				continue
			}

			go fetchUrl(url, queue, registry)
		case <-tick:
			if flag {
				fmt.Println("-----")
				fmt.Println("Analized pages:")
				fmt.Println("-----")
				fmt.Println(registry)
				return
			}
			flag = true
		}
	}
}

// Get start url from the command line arguments.
func startUrl() (string, error) {
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		return "", fmt.Errorf("Please specify a starting point, e.g. https://example.com")
	}

	return args[0], nil
}

func main() {

	startUrl, err := startUrl()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("-----")
	fmt.Println("Insecure resources (grouped by page):")
	fmt.Println("-----")

	crawl(startUrl, ResourceAndLinkFetcher{})
}
