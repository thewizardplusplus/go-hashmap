package hashmap

// Key ...
//go:generate mockery -name=Key -case=underscore
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

// Get ...
func (hashMap HashMap) Get(key Key) (value interface{}, ok bool) {
	for index := key.Hash(); ; index++ {
		modIndex := index % len(hashMap.buckets)
		bucket := hashMap.buckets[modIndex]
		if bucket == nil {
			return nil, false
		}
		if bucket.key.Equals(key) {
			return bucket.value, true
		}
	}
}
