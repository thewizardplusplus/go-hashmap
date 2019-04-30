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

// Set ...
func (hashMap *HashMap) Set(key Key, value interface{}) {
	for index := key.Hash(); ; index++ {
		modIndex := index % len(hashMap.buckets)
		b := hashMap.buckets[modIndex]
		if b == nil {
			hashMap.buckets[modIndex] = &bucket{key, value}
			hashMap.size++

			return
		}
		if b.key.Equals(key) {
			hashMap.buckets[modIndex].value = value
			return
		}
	}
}

// Delete ...
func (hashMap *HashMap) Delete(key Key) (ok bool) {
	for index := key.Hash(); ; index++ {
		modIndex := index % len(hashMap.buckets)
		bucket := hashMap.buckets[modIndex]
		if bucket == nil {
			return false
		}
		if bucket.key.Equals(key) {
			hashMap.buckets[modIndex] = nil
			hashMap.size--

			return true
		}
	}
}
