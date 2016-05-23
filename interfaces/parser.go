package interfaces

import "io"

// Parser is the interface that wraps the Parse method.
type Parser interface {
	// Parse parses passed response body and finds urls of resources and pages.
	Parse(baseUrl string, httpBody io.Reader) (resourceUrls []string, linkUrls []string, err error)
}
