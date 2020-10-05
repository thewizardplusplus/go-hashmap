package hashmap

// Config ...
type Config struct {
	initialCapacity int
	maxLoadFactor   float64
	growFactor      int
}

// Option ...
type Option func(options *Config)
