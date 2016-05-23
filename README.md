# insecRes - Insecure Resource Finder
[![Build Status](https://travis-ci.org/kkomelin/insecres.svg)](https://travis-ci.org/kkomelin/insecres)
[![Go Report Card](https://goreportcard.com/badge/github.com/kkomelin/insecres)](https://goreportcard.com/report/github.com/kkomelin/insecres)
[![GoDoc](https://godoc.org/github.com/kkomelin/insecres?status.png)](http://godoc.org/github.com/kkomelin/insecres)

A console tool that finds insecure resources on HTTPS sites.
It is written in Go language and uses the power of "multi-threading" (goroutines) to crawl and parse site pages.

## The motivation

Some time ago, I switched my site to HTTPS. _And you should too!_
All went well except the fact that my pages contained images, embedded videos and other resources,
which pointed to HTTP content and made browsers display warnings about the insecure content on the pages.
After some research of existing tools, which did not fit my needs, I decided to create my own one.

## Features

- Crawls all site pages in parallel
- Finds the following resources with absolute HTTP (insecure) urls:
  - IMG
  - IFRAME
  - OBJECT
  - AUDIO, VIDEO, SOURCE, TRACK
- Uses a random delay between requests to prevent blacklisting

## Installation

First of all, [install Go](https://golang.org/doc/install).

After that, run the following command:

```
go get github.com/kkomelin/insecres
```

## Usage

```
$GOPATH/bin/insecres https://example.com
```

## Roadmap

- [ ] Implement verbose mode and hide redundant information from display by default
- [ ] Print results to a file (CSV?)
- [x] Add random delay between requests to prevent blacklisting
- [x] Ignore trailing slashes (https://example.com and https://example.com/ are considered equivalent)
- [x] Handle domains w/ and w/o WWW
- [x] Support IFRAME tags
- [x] Support OBJECT tags
- [x] Support AUDIO, VIDEO, SOURCE and TRACK tags
