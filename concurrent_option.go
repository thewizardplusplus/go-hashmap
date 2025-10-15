package hashmap

import (
	"math/rand/v2"
)

// StorageFactory ...
type StorageFactory func() Storage[Key, any]

// ConcurrentConfig ...
type ConcurrentConfig struct {
	concurrencyLevel int
	segmentFactory   StorageFactory
	permuter         Permuter
}

// nolint: gochecknoglobals
var (
	defaultConcurrentConfig = ConcurrentConfig{
		concurrencyLevel: 16,
		segmentFactory: func() Storage[Key, any] {
			return NewSynchronizedHashMap[Key, any]()
		},
		permuter: rand.Perm,
	}
)

// ConcurrentOption ...
type ConcurrentOption func(options *ConcurrentConfig)

// WithConcurrencyLevel ...
//
// Default: 16.
func WithConcurrencyLevel(concurrencyLevel int) ConcurrentOption {
	return func(options *ConcurrentConfig) {
		options.concurrencyLevel = concurrencyLevel
	}
}

// WithSegmentFactory ...
//
// Default: a factory that produces an instance
// of the SynchronizedHashMap structure with default options.
func WithSegmentFactory(segmentFactory StorageFactory) ConcurrentOption {
	return func(options *ConcurrentConfig) {
		options.segmentFactory = segmentFactory
	}
}

// WithConcurrentPermuter ...
//
// Default: math/rand/v2.Perm.
func WithConcurrentPermuter(permuter Permuter) ConcurrentOption {
	return func(options *ConcurrentConfig) {
		options.permuter = permuter
	}
}
