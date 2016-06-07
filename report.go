package main

import (
	"bufio"
	"os"
	"sync"
)

// Report is a thread-safe data reporting tool.
type Report struct {
	file   *os.File
	writer *bufio.Writer
	mux    sync.Mutex
}

// Open opens or creates a file and initializes buffered writer.
func (r *Report) Open(filePath string) error {
	var err error

	r.file, err = os.Create(filePath)
	if err != nil {
		return err
	}

	r.writer = bufio.NewWriter(r.file)

	return nil
}

// WriteLines dump slice of strings to the file. It also adds trailing endline marker to each string.
func (r *Report) WriteLines(lines []string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	var err error

	for _, line := range lines {
		_, err = r.writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return r.writer.Flush()
}

// Close closes file handler in case it is not empty.
func (r *Report) Close() error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.IsEmpty() {
		return nil
	}

	return r.file.Close()
}

// IsEmpty check whether the file handler is initialized or not.
func (r *Report) IsEmpty() bool {
	return (r.file == nil)
}
