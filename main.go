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

	// Lock url so that no one other goroutine can process it.
	registry.MarkAsProcessed(url)

	fetcher := ResourceAndLinkFetcher{}

	insecureResourceUrls, pageUrls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Printf("Error occured: %s\n", err)
		return
	}

	for _, insecureResourceUrl := range insecureResourceUrls {
		fmt.Printf("%s: %s\n", url, insecureResourceUrl)
	}

	for _, url := range pageUrls {
		// Random pause before sending to the main thread.
		delayBetweenRequests()
		queue <- url
	}
}

// Implement random pause before sending the next request to
// (no more than half of beforeEngTimeout constant).
// It is one of the measures to prevent banning by the server.
func delayBetweenRequests() {
	randNum := rand.Intn(beforeEngTimeout/2)
	time.Sleep(time.Duration(randNum) * time.Millisecond)
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
	fmt.Println("Insecure resources (page: resource):")
	fmt.Println("-----")

	crawl(startUrl, ResourceAndLinkFetcher{})
}
