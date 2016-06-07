package main

import (
	"bufio"
	"sync"
	"os"
)

type Report struct {
	file *os.File
	writer *bufio.Writer
	mux sync.Mutex
}

func (r *Report) Open(filePath string) error {
	var err error

	r.file, err = os.Create(filePath)
	if err != nil {
		return err
	}

	r.writer = bufio.NewWriter(r.file)

	return nil
}

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

func (r *Report) Close() error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.IsEmpty() {
		return nil
	}

	return r.file.Close()
}

func (r *Report) IsEmpty() bool {
	return (r.file == nil)
}
