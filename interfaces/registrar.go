package interfaces

// Registrar is the interface that wraps methods necessary for storing processed urls.
// All methods of this interface should be thread-safe.
type Registrar interface {
	// Register adds processed url to the registry.
	Register(url string)
	// IsNew checks whether the passed url is new or not.
	IsNew(url string) bool
}
