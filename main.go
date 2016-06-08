package main

import (
	"flag"
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
	var (
		helpFlag   bool
		reportFlag string
	)

	// Find options.
	flag.BoolVar(&helpFlag, "h", false, "")
	flag.StringVar(&reportFlag, "f", "", "")
	flag.Parse()

	// Find argument.
	args := flag.Args()
	if len(args) < 1 {
		displayHelp()
		return
	}

	// Display help.
	if helpFlag {
		displayHelp()
		return
	}

	// Run the crawler.
	Crawl(args[0], reportFlag)
}
