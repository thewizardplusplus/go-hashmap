package hashmap

// SynchronizedConfig ...
type SynchronizedConfig struct {
	innerMap Storage
}

// SynchronizedOption ...
type SynchronizedOption func(options *SynchronizedConfig)
