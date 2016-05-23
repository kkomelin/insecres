package main

import (
	"fmt"
	"os"
)

const (
	// The goal of BeforeEngTimeout is to wait for some time before exit
	// so that all goroutines can finish.
	BeforeEngDelay int = 2000
	// Before processing any new url we wait for some random time
	// (between MinDelayBetweenRequests and MaxDelayBetweenRequests)
	// to prevent blacklisting by the server.
	MinDelayBetweenRequests int = 500
	// See MinDelayBetweenRequests.
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
