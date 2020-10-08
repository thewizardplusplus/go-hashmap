package hashmap

// ConcurrentConfig ...
type ConcurrentConfig struct {
	concurrencyLevel int
	segmentFactory   StorageFactory
}

// ConcurrentOption ...
type ConcurrentOption func(options *ConcurrentConfig)

// WithConcurrencyLevel ...
//
// Default: 16.
//
func WithConcurrencyLevel(concurrencyLevel int) ConcurrentOption {
	return func(options *ConcurrentConfig) {
		options.concurrencyLevel = concurrencyLevel
	}
}
