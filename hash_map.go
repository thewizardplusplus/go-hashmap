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

const (
	initialCapacity = 16
)

// NewHashMap ...
func NewHashMap() *HashMap {
	buckets := make([]*bucket, initialCapacity)
	return &HashMap{buckets: buckets, size: 0}
}
