package interfaces

import "io"

// Fetcher is the interface that wraps the Fetch method.
type Fetcher interface {
	// Fetch fetches page by url and returns the response body.
	Fetch(url string) (responseBody io.ReadCloser, err error)
}
