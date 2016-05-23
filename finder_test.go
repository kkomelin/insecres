package main

import (
	"fmt"
	"strings"
	"testing"
)

// TestParse tests fetcher.Parse method and its results.
func TestParse(t *testing.T) {

	fmt.Println("TestParse")

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
<audio src="http://www.example.com/audio.ogg" autoplay>
  Your browser does not support HTML5 audio tag.
  <track kind="captions" src="http://www.example.com/audio_track.vtt" srclang="en" label="English">
  <source src="http://www.example.com/audio_in_source.ogg" type="audio/ogg">
</audio>
<video src="http://www.example.com/video.mp4" poster="http://www.example.com/poster.jpg" autoplay>
   Your browser doesn't support HTML5 video tag.
   <track kind="subtitles" src="http://www.example.com/video_track.vtt" srclang="en" label="English">
   <source src="http://www.example.com/video_in_source.mp4" type="video/mp4">
</video>
</body>`)

	expectedResources := map[string]int{
		// img[src]
		"http://example.com/images/test.png": 0,
		// iframe[src]
		"http://www.youtube.com/embed/0sRPY3WWSNc": 0,
		// object[data]
		"http://www.example.com/flash/insecure.swf": 0,
		// audio[src]
		"http://www.example.com/audio.ogg": 0,
		// audio track[src]
		"http://www.example.com/audio_track.vtt": 0,
		// audio source[src]
		"http://www.example.com/audio_in_source.ogg": 0,
		// video[src]
		"http://www.example.com/video.mp4": 0,
		// FIXME: Currently golang.org/x/net/html library ignores video[poster] attribute for some reasons. We need to do something with that.
		// video[poster]
		// "http://www.example.com/poster.jpg": 0,
		// video track[src]
		"http://www.example.com/video_track.vtt": 0,
		// video source[src]
		"http://www.example.com/video_in_source.mp4": 0,
	}

	expectedLinks := map[string]int{
		"https://example.com/article/test1": 0,
		"http://example.com/test2":          0,
		"https://example.com/test3":         0,
		"http://www.example.com/test3":      0,
		"http://www.example.com/test4":      0,
	}

	resources, links, err := (ResourceAndLinkFinder{}).Parse("https://example.com/", reader)
	if err != nil {
		t.Fatalf("fetcher.Parse has returned error: %s\n", err)
	}

	// Check resources.
	fmt.Printf("Resources: %q\n", resources)

	if len(resources) != len(expectedResources) {
		t.Errorf("Wrong number of resources. Found %d of %d", len(resources), len(expectedResources))
	} else {
		for i := 0; i < len(resources); i++ {
			if _, ok := expectedResources[resources[i]]; !ok {
				t.Errorf("Resource url is not found in the expected values: %s", resources[i])
			}
		}
	}

	// Check links.
	fmt.Printf("Links: %q\n", links)

	if len(links) != len(expectedLinks) {
		t.Errorf("Wrong number of links. Found %d of %d", len(links), len(expectedLinks))

	} else {
		for i := 0; i < len(links); i++ {
			if _, ok := expectedLinks[links[i]]; !ok {
				t.Errorf("Link url is not found in the expected values: %s", links[i])
			}
		}
	}
}
