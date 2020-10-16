package hashmap

import (
	"math/rand"
)

// ConcurrentHashMap ...
//
// It's partially safe for concurrent access because it uses data sharding.
// Each segment should take care of concurrent access safety itself.
//
type ConcurrentHashMap struct {
	segments []Storage
}

// NewConcurrentHashMap ...
func NewConcurrentHashMap(options ...ConcurrentOption) ConcurrentHashMap {
	config := defaultConcurrentConfig
	for _, option := range options {
		option(&config)
	}

	var segments []Storage
	for i := 0; i < config.concurrencyLevel; i++ {
		segment := config.segmentFactory()
		segments = append(segments, segment)
	}

	return ConcurrentHashMap{segments: segments}
}

// Get ...
func (hashMap ConcurrentHashMap) Get(key Key) (value interface{}, ok bool) {
	return hashMap.selectSegment(key).Get(key)
}

// Iterate ...
//
// If the handler returns false, iteration is broken.
//
// It randomizes of iteration order over items and their keys and over segments.
//
func (hashMap ConcurrentHashMap) Iterate(handler Handler) bool {
	for _, index := range rand.Perm(len(hashMap.segments)) {
		segment := hashMap.segments[index]
		if ok := segment.Iterate(handler); !ok {
			return false
		}
	}

	return true
}

// Set ...
func (hashMap ConcurrentHashMap) Set(key Key, value interface{}) {
	hashMap.selectSegment(key).Set(key, value)
}

// Delete ...
func (hashMap ConcurrentHashMap) Delete(key Key) {
	hashMap.selectSegment(key).Delete(key)
}

func (hashMap ConcurrentHashMap) selectSegment(key Key) Storage {
	index := key.Hash() % len(hashMap.segments)
	return hashMap.segments[index]
}
