package main

import (
	"strings"
	"testing"
	"fmt"
)

func TestParser(t *testing.T) {

	reader := strings.NewReader(`<!DOCTYPE html>
<html lang="en"><head><title></title></head>
<body>
<img src="/images/test.png">
<img src="http://example.com/images/test.png" />
<img src="https://example.com/images/test.png" />
<a href="#">Anchor (ignored)</a>
<a href="/article/test1">Relative link</a>
<a href="http://example.com/test2">Absolute HTTP link</a>
<a href="https://example.com/test3">Absolute HTTPS link</a>
</body>`)

	expected_resources := []string{
		"http://example.com/images/test.png",
	}

	expected_links := []string{
		"https://example.com/article/test1",
		"http://example.com/test2",
		"https://example.com/test3",
	}

	fetcher := UrlFetcher{}

	resources, links, err := fetcher.Parse("https://example.com/", reader)
	if err != nil {
		t.Error("Error: %v", err)
	}

	fmt.Printf("Resources: %q\n", resources)

	if len(resources) != len(expected_resources) {
		t.Errorf("Wrong number of links. Found %d of %d", len(resources), len(expected_resources))
	} else {
		for i := 0; i < len(expected_resources); i++ {
			if resources[i] != expected_resources[i] {
				t.Errorf("Resource url %d is incorrect. Expected: %s, Given: %s", i, expected_resources[i], resources[i])
			}
		}
	}

	fmt.Printf("Links: %q\n", links)

	if len(links) != len(expected_links) {
		t.Errorf("Wrong number of links. Found %d of %d", len(links), len(expected_links))

	} else {
		for i := 0; i < len(expected_links); i++ {
			if links[i] != expected_links[i] {
				t.Errorf("Link url %d is incorrect. Expected: %s, Given: %s", i, expected_links[0], links[0])
			}
		}
	}
}
