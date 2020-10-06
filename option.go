package hashmap

// Config ...
type Config struct {
	initialCapacity int
	maxLoadFactor   float64
	growFactor      int
}

// nolint: gochecknoglobals
var (
	defaultConfig = Config{
		initialCapacity: 16,
		maxLoadFactor:   0.75,
		growFactor:      2,
	}
)

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

// WithGrowFactor ...
//
// Default: 2.
//
func WithGrowFactor(growFactor int) Option {
	return func(options *Config) {
		options.growFactor = growFactor
	}
}
