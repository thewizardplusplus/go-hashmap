package hashmap

// Key ...
type Key interface {
	Hash() int
	Equals(key interface{}) bool
}

type bucket struct {
	key   Key
	value interface{}
}

// HashMap ...
type HashMap struct {
	buckets []*bucket
	size    int
}
