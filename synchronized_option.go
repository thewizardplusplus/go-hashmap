package hashmap

// SynchronizedConfig ...
type SynchronizedConfig[KeyType Key, ValueType any] struct {
	innerMap Storage[KeyType, ValueType]
}

// SynchronizedOption ...
type SynchronizedOption[KeyType Key, ValueType any] func(
	options *SynchronizedConfig[KeyType, ValueType],
)

// WithInnerMap ...
//
// Default: an instance of the HashMap structure with default options.
func WithInnerMap[KeyType Key, ValueType any](
	innerMap Storage[KeyType, ValueType],
) SynchronizedOption[KeyType, ValueType] {
	return func(options *SynchronizedConfig[KeyType, ValueType]) {
		options.innerMap = innerMap
	}
}
