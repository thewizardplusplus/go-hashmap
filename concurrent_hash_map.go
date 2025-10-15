package hashmap

import (
	"iter"
	"math/rand/v2"
)

// ConcurrentHashMap ...
//
// It's partially safe for concurrent access because it uses data sharding.
// Each segment should take care of concurrent access safety itself.
type ConcurrentHashMap[KeyType Key, ValueType any] struct {
	config   ConcurrentConfig[KeyType, ValueType]
	segments []Storage[KeyType, ValueType]
}

// NewConcurrentHashMap ...
func NewConcurrentHashMap[KeyType Key, ValueType any](
	options ...ConcurrentOption[KeyType, ValueType],
) ConcurrentHashMap[KeyType, ValueType] {
	// you can't move the default concurrent config into a global variable
	// because the generic type parameters must be resolved inside the constructor
	config := ConcurrentConfig[KeyType, ValueType]{
		concurrencyLevel: defaultConcurrencyLevel,
		segmentFactory: func() Storage[KeyType, ValueType] {
			return NewSynchronizedHashMap[KeyType, ValueType]()
		},
		permuter: rand.Perm,
	}
	for _, option := range options {
		option(&config)
	}

	var segments []Storage[KeyType, ValueType]
	for i := 0; i < config.concurrencyLevel; i++ {
		segment := config.segmentFactory()
		segments = append(segments, segment)
	}

	return ConcurrentHashMap[KeyType, ValueType]{
		config:   config,
		segments: segments,
	}
}

// Get ...
func (hashMap ConcurrentHashMap[KeyType, ValueType]) Get(
	key KeyType,
) (value ValueType, ok bool) {
	return hashMap.selectSegment(key).Get(key)
}

// All ...
//
// It randomizes of iteration order over items and their keys and over segments.
func (hashMap ConcurrentHashMap[KeyType, ValueType]) All() iter.Seq2[KeyType, ValueType] { //nolint:lll
	return func(yield func(KeyType, ValueType) bool) {
	outerLoop:
		for _, index := range hashMap.config.permuter(len(hashMap.segments)) {
			segment := hashMap.segments[index]
			for key, value := range segment.All() {
				if ok := yield(key, value); !ok {
					break outerLoop
				}
			}
		}
	}
}

// Set ...
func (hashMap ConcurrentHashMap[KeyType, ValueType]) Set(
	key KeyType,
	value ValueType,
) {
	hashMap.selectSegment(key).Set(key, value)
}

// Delete ...
func (hashMap ConcurrentHashMap[KeyType, ValueType]) Delete(key KeyType) {
	hashMap.selectSegment(key).Delete(key)
}

func (hashMap ConcurrentHashMap[KeyType, ValueType]) selectSegment(
	key KeyType,
) Storage[KeyType, ValueType] {
	index := key.Hash() % len(hashMap.segments)
	return hashMap.segments[index]
}
