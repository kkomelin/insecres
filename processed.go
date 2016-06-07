package main

import (
	"sync"
)

// Processed is a thread-safe storage for processed urls.
type Processed struct {
	processed map[string]int
	mux       sync.Mutex
}

// Register adds a processed url to the registry.
func (r *Processed) Register(url string) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.processed[url] = 1
}

// IsNew checks whether the url is new.
func (r *Processed) IsNew(url string) bool {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, ok := r.processed[url]; ok {
		return false
	}

	return true
}

// String defines our own way to output the processed urls.
// [url1]\n
// [url2]\n
func (r *Processed) String() string {
	r.mux.Lock()
	defer r.mux.Unlock()

	output := ""
	for url := range r.processed {
		output += url + "\n"
	}
	return output
}
