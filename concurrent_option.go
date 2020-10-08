package hashmap

// ConcurrentConfig ...
type ConcurrentConfig struct {
	concurrencyLevel int
	segmentFactory   StorageFactory
}

// ConcurrentOption ...
type ConcurrentOption func(options *ConcurrentConfig)
