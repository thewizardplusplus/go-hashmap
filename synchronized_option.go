package hashmap

// SynchronizedConfig ...
type SynchronizedConfig struct {
	innerMap Storage
}

// SynchronizedOption ...
type SynchronizedOption func(options *SynchronizedConfig)

// WithInnerMap ...
//
// Default: an instance of the HashMap structure with default options.
//
func WithInnerMap(innerMap Storage) SynchronizedOption {
	return func(options *SynchronizedConfig) {
		options.innerMap = innerMap
	}
}
