package hashmap

const (
	defaultConcurrencyLevel = 16
)

// StorageFactory ...
type StorageFactory[KeyType Key, ValueType any] func() Storage[KeyType, ValueType] //nolint:lll

// ConcurrentConfig ...
type ConcurrentConfig[KeyType Key, ValueType any] struct {
	concurrencyLevel int
	segmentFactory   StorageFactory[KeyType, ValueType]
	permuter         Permuter
}

// ConcurrentOption ...
type ConcurrentOption[KeyType Key, ValueType any] func(
	options *ConcurrentConfig[KeyType, ValueType],
)

// WithConcurrencyLevel ...
//
// Default: 16.
func WithConcurrencyLevel[KeyType Key, ValueType any](
	concurrencyLevel int,
) ConcurrentOption[KeyType, ValueType] {
	return func(options *ConcurrentConfig[KeyType, ValueType]) {
		options.concurrencyLevel = concurrencyLevel
	}
}

// WithSegmentFactory ...
//
// Default: a factory that produces an instance
// of the SynchronizedHashMap structure with default options.
func WithSegmentFactory[KeyType Key, ValueType any](
	segmentFactory StorageFactory[KeyType, ValueType],
) ConcurrentOption[KeyType, ValueType] {
	return func(options *ConcurrentConfig[KeyType, ValueType]) {
		options.segmentFactory = segmentFactory
	}
}

// WithConcurrentPermuter ...
//
// Default: math/rand/v2.Perm.
func WithConcurrentPermuter[KeyType Key, ValueType any](
	permuter Permuter,
) ConcurrentOption[KeyType, ValueType] {
	return func(options *ConcurrentConfig[KeyType, ValueType]) {
		options.permuter = permuter
	}
}
