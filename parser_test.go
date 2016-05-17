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
</body>`)

	var page Page

	links, err := page.Parse(reader)
	if err != nil {
		t.Error("Error: %v", err)
	}

	fmt.Println(links)

	if len(links) != 1 {
		t.Errorf("Wrong number of links. Found %d of %d", len(links), 1)
	}

	if links[0] != "http://example.com/images/test.png" {
		t.Errorf("The first link is incorrect. Expected: %s, Given: %s", "http://example.com/images/test.png", links[0])
	}
}
