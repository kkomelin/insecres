package main

import (
	"fmt"
	"strings"
	"testing"
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
<a href="http://www.example.com/test4/">Ignoring trailing slash</a>
<iframe width="560" height="315" src="https://www.youtube.com/embed/0sRPY3WWSNc" frameborder="0" allowfullscreen></iframe>
<iframe width="560" height="315" src="http://www.youtube.com/embed/0sRPY3WWSNc" frameborder="0" allowfullscreen></iframe>
<iframe width="560" height="315" src="//www.youtube.com/embed/0sRPY3WWSNc" frameborder="0" allowfullscreen></iframe>
<object type="application/x-shockwave-flash" data="http://www.example.com/flash/insecure.swf" width="400" height="300">
    <param name="quality" value="high">
    <param name="wmode" value="opaque">
</object>
<object type="application/x-shockwave-flash" data="https://www.example.com/flash/secure.swf" width="400" height="300">
    <param name="quality" value="high">
    <param name="wmode" value="opaque">
</object>
</body>`)

	expected_resources := map[string]int{
		"http://example.com/images/test.png":        0,
		"http://www.youtube.com/embed/0sRPY3WWSNc":  1,
		"http://www.example.com/flash/insecure.swf": 2,
	}

	expected_links := map[string]int{
		"https://example.com/article/test1": 0,
		"http://example.com/test2":          1,
		"https://example.com/test3":         2,
		"http://www.example.com/test3":      3,
		"http://www.example.com/test4":      4,
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
		for i := 0; i < len(resources); i++ {
			if _, ok := expected_resources[resources[i]]; !ok {
				t.Errorf("Resource url is not found in the expected values: %s", resources[i])
			}
		}
	}

	// Check links.
	fmt.Printf("Links: %q\n", links)

	if len(links) != len(expected_links) {
		t.Errorf("Wrong number of links. Found %d of %d", len(links), len(expected_links))

	} else {
		for i := 0; i < len(links); i++ {
			if _, ok := expected_links[links[i]]; !ok {
				t.Errorf("Link url is not found in the expected values: %s", links[i])
			}
		}
	}
}
