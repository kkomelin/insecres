package main

import (
	"sync"
)

// Thread-safe registry for processed urls.
type Registry struct {
	processed map[string]int
	mux       sync.Mutex
}

// Add a processed url to the registry.
func (r *Registry) MarkAsProcessed(url string) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.processed[url] = 1
}

// Check whether a url is new.
func (r *Registry) IsNew(url string) bool {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, ok := r.processed[url]; ok {
		return false
	}

	return true
}

// Define our own way to output the processed urls.
// [url1]\n
// [url2]\n
func (r *Registry) String() string {
	r.mux.Lock()
	defer r.mux.Unlock()

	output := ""
	for url, _ := range r.processed {
		output += url + "\n"
	}
	return output
}
