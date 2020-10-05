package hashmap

// Config ...
type Config struct {
	initialCapacity int
	maxLoadFactor   float64
	growFactor      int
}

// Option ...
type Option func(options *Config)

// WithInitialCapacity ...
//
// Default: 16.
//
func WithInitialCapacity(initialCapacity int) Option {
	return func(options *Config) {
		options.initialCapacity = initialCapacity
	}
}

// WithMaxLoadFactor ...
//
// Default: 0.75.
//
func WithMaxLoadFactor(maxLoadFactor float64) Option {
	return func(options *Config) {
		options.maxLoadFactor = maxLoadFactor
	}
}
