package hashmap

//go:generate mockery -name=Key -inpkg -case=underscore -testonly

// Key ...
type Key interface {
	Hash() int
	Equals(key interface{}) bool
}

// Handler ...
type Handler func(key Key, value interface{}) bool

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
	maxLoadFactor   = 0.75
	growFactor      = 2
)

// NewHashMap ...
func NewHashMap() *HashMap {
	return newHashMapWithCapacity(initialCapacity)
}

// Get ...
func (hashMap HashMap) Get(key Key) (value interface{}, ok bool) {
	index, ok := hashMap.find(key)
	if !ok {
		return nil, false
	}

	return hashMap.buckets[index].value, true
}

// Iterate ...
func (hashMap HashMap) Iterate(handler Handler) bool {
	for _, bucket := range hashMap.buckets {
		if bucket == nil {
			continue
		}

		if ok := handler(bucket.key, bucket.value); !ok {
			return false
		}
	}

	return true
}

// Set ...
func (hashMap *HashMap) Set(key Key, value interface{}) {
	index, ok := hashMap.find(key)
	if ok {
		hashMap.buckets[index].value = value
		return
	}

	hashMap.buckets[index] = &bucket{key, value}
	hashMap.size++

	loadFactor := float64(hashMap.size) / float64(len(hashMap.buckets))
	if loadFactor > maxLoadFactor {
		hashMap.rehash()
	}
}

// Delete ...
func (hashMap *HashMap) Delete(key Key) {
	index, ok := hashMap.find(key)
	if !ok {
		return
	}

	hashMap.buckets[index] = nil
	hashMap.size--
}

func (hashMap HashMap) find(key Key) (index int, ok bool) {
	for index := key.Hash(); ; index++ {
		modIndex := index % len(hashMap.buckets)
		bucket := hashMap.buckets[modIndex]
		if bucket == nil {
			return modIndex, false
		}
		if bucket.key.Equals(key) {
			return modIndex, true
		}
	}
}

func (hashMap *HashMap) rehash() {
	newHashMap := newHashMapWithCapacity(len(hashMap.buckets) * growFactor)
	hashMap.Iterate(func(key Key, value interface{}) bool {
		newHashMap.Set(key, value)
		return true
	})

	*hashMap = *newHashMap
}

func newHashMapWithCapacity(capacity int) *HashMap {
	buckets := make([]*bucket, capacity)
	return &HashMap{buckets: buckets, size: 0}
}
