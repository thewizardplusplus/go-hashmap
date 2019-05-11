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
	return hashMap.selectSegment(key).Get(key)
}

// Iterate ...
func (hashMap ConcurrentHashMap) Iterate(
	handler func(key Key, value interface{}),
) {
	for _, segment := range hashMap.segments {
		segment.Iterate(handler)
	}
}

// Set ...
func (hashMap ConcurrentHashMap) Set(key Key, value interface{}) {
	hashMap.selectSegment(key).Set(key, value)
}

// Delete ...
func (hashMap ConcurrentHashMap) Delete(key Key) {
	hashMap.selectSegment(key).Delete(key)
}

func (hashMap ConcurrentHashMap) selectSegment(key Key) *SynchronizedHashMap {
	index := key.Hash() % len(hashMap.segments)
	return hashMap.segments[index]
}
