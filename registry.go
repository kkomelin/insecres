package main

import (
	"sync"
)

// Thread-safe processed urls registry.
type Registry struct {
	processed map[string]int
	queue []string
	mux  sync.Mutex
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
// [number of times found]: [url]
func (r *Registry) String() string {
	r.mux.Lock()
	defer r.mux.Unlock()

	output := ""
	for url, _ := range r.processed {
		output += url + "\n"
	}
	return output
}

func (r *Registry) AddToQueue(url string) {
	r.mux.Lock()
	defer r.mux.Unlock()

	// TODO: Check for existence. It is not enough to check just r.processed and we have to also check queue.

	r.queue = append(r.queue, url)
}

func (r *Registry) GetFromQueue() (string, bool) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if len(r.queue) == 0 {
		return "", false
	}

	url := r.queue[0]
	// Remove the top element from the slice.
	r.queue = r.queue[1:]

	return url, true
}
