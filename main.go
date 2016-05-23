package main

import (
	"fmt"
	"os"
)

const (
	// BeforeEngDelay defines time in milliseconds, which the program waits before exit
	// so that all goroutines can finish and return results.
	BeforeEngDelay int = 2000
	// MinDelayBetweenRequests is minimum time in milliseconds,
	// which the program waits before processing any new url.
	// We wait for some random time (between MinDelayBetweenRequests and MaxDelayBetweenRequests)
	// to prevent blacklisting by the server.
	MinDelayBetweenRequests int = 500
	// MaxDelayBetweenRequests is maximum time in milliseconds,
	// which the program waits before processing any new url.
	MaxDelayBetweenRequests int = 1000
)

func main() {

	startUrl, err := startUrl()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("-----")
	fmt.Println("Insecure resources (grouped by page):")
	fmt.Println("-----")

	Crawl(startUrl, ResourceAndLinkFetcher{})
}
