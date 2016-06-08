package main

import (
	"fmt"
	"github.com/kkomelin/insecres/interfaces"
	"math/rand"
	"strings"
	"time"
)

// Goroutine callback, which fetches and parses the passed url
// in order to find insecure resources and next urls to fetch from.
func processPage(url string, queue chan string, registrar interfaces.Registrar, fetcher interfaces.Fetcher, parser interfaces.Parser, reporter interfaces.Reporter) {

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

	reportPageResources(url, insecureResourceUrls, reporter)

	for _, url := range pageUrls {
		// Random pause before sending to the main thread.
		delayBetweenRequests()
		queue <- url
	}
}

// Reports page resources.
func reportPageResources(url string, resources []string, reporter interfaces.Reporter) error {
	if len(resources) == 0 {
		return nil
	}

	if !reporter.IsEmpty() {
		for i, insecureResourceUrl := range resources {
			resources[i] = url + ", " + insecureResourceUrl
		}

		return reporter.WriteLines(resources)
	}

	fmt.Printf("\n%s:\n", url)
	for _, insecureResourceUrl := range resources {
		fmt.Printf("- %s\n", insecureResourceUrl)
	}
	return nil
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
func Crawl(url, reportFile string) {

	url = strings.TrimSuffix(url, "/")

	report := &Report{}

	// Print results to the file.
	if reportFile != "" {
		err := report.Open(reportFile)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else { // Print results to console.
		fmt.Println("-----")
		fmt.Println("Insecure resources (grouped by page):")
		fmt.Println("-----")
	}

	registry := &Processed{processed: make(map[string]int)}
	finder := &ResourceAndLinkFinder{}

	queue := make(chan string)

	go processPage(url, queue, registry, finder, finder, report)

	tick := time.Tick(time.Duration(BeforeEngDelay) * time.Millisecond)
	flag := false
	for {
		select {
		case url := <-queue:
			flag = false

			go processPage(url, queue, registry, finder, finder, report)
		case <-tick:
			if flag {
				// TODO: Implement a verbose mode when all crawled pages are also displayed.
				//if false {
				//	fmt.Println("-----")
				//	fmt.Println("Analized pages:")
				//	fmt.Println("-----")
				//	fmt.Println(registry)
				//}
				fmt.Println("")

				// Close report.
				report.Close()

				return
			}
			flag = true
		}
	}
}

func displayHelp() {
	fmt.Printf(`usage: insecres [-h|-f="path/to/report.csv"] <url>
ARGUMENTS
  url
    A url to start from, e.g. https://example.com"
OPTIONS
  -h
    Show this help message.
  -f
    Define the location of the CSV file with the results.
    If it is not set, results are printed to the console.
`)
}
