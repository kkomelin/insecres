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
<img src="//example.com/images/test.png" />
<a href="#">Anchor (ignored)</a>
<a href="/article/test1">Relative link</a>
<a href="http://example.com/test2">Absolute HTTP link</a>
<a href="https://example.com/test3">Absolute HTTPS link</a>
<a href="http://www.example.com/test3">Absolute HTTPS link</a>
<a href="https://www.youtube.com/watch?v=yIhJEO6QvFA">External link</a>
<a href="//www.youtube.com/watch?v=o4cM2KUdfTg">Reproduces bug in Go url.isAbs()</a>
<iframe width="560" height="315" src="https://www.youtube.com/embed/0sRPY3WWSNc" frameborder="0" allowfullscreen></iframe>
<iframe width="560" height="315" src="http://www.youtube.com/embed/0sRPY3WWSNc" frameborder="0" allowfullscreen></iframe>
<iframe width="560" height="315" src="//www.youtube.com/embed/0sRPY3WWSNc" frameborder="0" allowfullscreen></iframe>
</body>`)

	expected_resources := []string{
		0: "http://example.com/images/test.png",
		1: "http://www.youtube.com/embed/0sRPY3WWSNc",
	}

	expected_links := []string{
		0: "https://example.com/article/test1",
		1: "http://example.com/test2",
		2: "https://example.com/test3",
		3: "http://www.example.com/test3",
	}

	fetcher := InsecureResourceFetcher{}

	resources, links, err := fetcher.Parse("https://example.com/", reader)
	if err != nil {
		t.Error("Error: %v", err)
	}

	// Check resources.
	fmt.Printf("Resources: %q\n", resources)

	if len(resources) != len(expected_resources) {
		t.Errorf("Wrong number of resources. Found %d of %d", len(resources), len(expected_resources))
	} else {
		for i := 0; i < len(expected_resources); i++ {
			if resources[i] != expected_resources[i] {
				t.Errorf("Resource url %d is incorrect. Expected: %s, Given: %s", i, expected_resources[i], resources[i])
			}
		}
	}

	// Check links.
	fmt.Printf("Links: %q\n", links)

	if len(links) != len(expected_links) {
		t.Errorf("Wrong number of links. Found %d of %d", len(links), len(expected_links))

	} else {
		for i := 0; i < len(expected_links); i++ {
			if links[i] != expected_links[i] {
				t.Errorf("Link url %d is incorrect. Expected: %s, Given: %s", i, expected_links[i], links[i])
			}
		}
	}
}
