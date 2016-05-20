# httpsrescheck
A console tool that finds insecure resources on HTTPS sites.
It is written in Go language and uses the power of "multi-threading" (goroutines) to crawl and parse site pages.

## The motivation

Some time ago, I switched my site to HTTPS. And you should too!
All went well except the fact that my pages contained images, embedded videos and other resources,
which pointed to HTTP content and made browsers display warnings about the insecure content on the pages.
After some research of existing tools, which did not fit my needs, I decided to create my own one.

## Features

- Crawl all site pages in parallel
- Find IMG and IFRAME resources with absolute HTTP urls

## Installation

```
go get github.com/kkomelin/httpsrescheck
```

## Usage

```
./httpsrescheck https://example.com
```

## Roadmap

- Check OBJECT tag
- Implement debug option and hide url log by default
- Improve output format (CSV?)
- Handle trailing slash (https://example.com and https://example.com/)
- Handle domains w/ and w/o WWW
- Add some Sleep between requests to prevent blacklisting
