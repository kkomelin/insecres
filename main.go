package main

import (
	"fmt"
	"os"
)

const (
	// The goal of beforeEngTimeout is to wait for some milliseconds before exit
	// so that all goroutines can finish.
	beforeEngTimeout int = 2000
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

	crawl(startUrl, ResourceAndLinkFetcher{})
}
