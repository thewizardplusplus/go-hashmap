package hashmap

// ConcurrentHashMap ...
type ConcurrentHashMap struct {
	segments []*SynchronizedHashMap
}

const (
	concurrencyLevel = 16
)

// NewConcurrentHashMap ...
func NewConcurrentHashMap() ConcurrentHashMap {
	var segments []*SynchronizedHashMap
	for i := 0; i < concurrencyLevel; i++ {
		segments = append(segments, NewSynchronizedHashMap())
	}

	return ConcurrentHashMap{segments: segments}
}

// Get ...
func (hashMap ConcurrentHashMap) Get(key Key) (value interface{}, ok bool) {
	index := key.Hash() % len(hashMap.segments)
	return hashMap.segments[index].Get(key)
}

// Set ...
func (hashMap ConcurrentHashMap) Set(key Key, value interface{}) {
	index := key.Hash() % len(hashMap.segments)
	hashMap.segments[index].Set(key, value)
}
