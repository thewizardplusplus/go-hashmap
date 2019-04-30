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
	index, ok := hashMap.find(key)
	if !ok {
		return nil, false
	}

	return hashMap.buckets[index].value, true
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
}

// Delete ...
func (hashMap *HashMap) Delete(key Key) (ok bool) {
	index, ok := hashMap.find(key)
	if !ok {
		return false
	}

	hashMap.buckets[index] = nil
	hashMap.size--

	return true
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
