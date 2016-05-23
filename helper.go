package main

import (
	"flag"
	"fmt"
	"github.com/kkomelin/insecres/interfaces"
	"math/rand"
	"strings"
	"time"
)

// Goroutine callback, which fetches and parses the passed url
// in order to find insecure resources and next urls to fetch from.
func fetchPage(url string, queue chan string, registrar interfaces.Registrar, fetcher interfaces.Fetcher, parser interfaces.Parser) {

	// Ignore processed urls.
	if !registrar.IsNew(url) {
		return
	}
	// Lock url so that no one other goroutine can process it.
	registrar.Register(url)

	fmt.Print(".")

	responseBody, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Printf("Error occured: %s\n", err)
		return
	}

	defer responseBody.Close()

	insecureResourceUrls, pageUrls, err := parser.Parse(url, responseBody)
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
	if len(resources) > 0 {
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
	randNum := randomInRange(MinDelayBetweenRequests, MaxDelayBetweenRequests)
	time.Sleep(time.Duration(randNum) * time.Millisecond)
}

// Returns a random number in a given range.
// The idea has been borrowed from http://golangcookbook.blogspot.ru/2012/11/generate-random-number-in-given-range.html
// and improved.
func randomInRange(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())

	if (min == 0 || max == 0) || (min > max) {
		return 0
	}

	if min == max {
		return min
	}

	return rand.Intn(max-min) + min
}

// Crawl pages starting from the passed url and find insecure resources.
func Crawl(url string) {

	url = strings.TrimSuffix(url, "/")

	registry := &ProcessedUrls{processed: make(map[string]int)}
	finder := &ResourceAndLinkFinder{}

	queue := make(chan string)

	go fetchPage(url, queue, registry, finder, finder)

	fmt.Println("-----")
	fmt.Println("Insecure resources (grouped by page):")
	fmt.Println("-----")

	tick := time.Tick(time.Duration(BeforeEngDelay) * time.Millisecond)
	flag := false
	for {
		select {
		case url := <-queue:
			flag = false

			go fetchPage(url, queue, registry, finder, finder)
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
