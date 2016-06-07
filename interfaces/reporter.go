package interfaces

// Reporter is the interface that wraps methods for reporting results.
type Reporter interface {
	// Init prepares file to report to.
	Open(filePath string) error
	// WriteLines dumps slice of strings to the report.
	WriteLines(lines []string) error
	// Close releases file.
	Close() error
	// Close releases file.
	IsEmpty() bool
}
