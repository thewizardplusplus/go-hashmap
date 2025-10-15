package hashmap

import (
	"iter"
)

type bucket[KeyType Key, ValueType any] struct {
	key   KeyType
	value ValueType
}

// HashMap ...
//
// It's not safe for concurrent access.
type HashMap[KeyType Key, ValueType any] struct {
	config  Config
	buckets []*bucket[KeyType, ValueType]
	size    int
}

// NewHashMap ...
func NewHashMap[KeyType Key, ValueType any](
	options ...Option,
) *HashMap[KeyType, ValueType] {
	config := defaultConfig
	for _, option := range options {
		option(&config)
	}

	return newHashMapWithCapacity[KeyType, ValueType](
		config,
		config.initialCapacity,
	)
}

// Get ...
func (hashMap HashMap[KeyType, ValueType]) Get(
	key KeyType,
) (value ValueType, ok bool) {
	index, ok := hashMap.find(key)
	if !ok {
		var zeroValue ValueType
		return zeroValue, false
	}

	return hashMap.buckets[index].value, true
}

// All ...
//
// It randomizes of iteration order.
func (hashMap HashMap[KeyType, ValueType]) All() iter.Seq2[KeyType, ValueType] {
	return func(yield func(KeyType, ValueType) bool) {
		for _, index := range hashMap.config.permuter(len(hashMap.buckets)) {
			bucket := hashMap.buckets[index]
			if bucket == nil {
				continue
			}

			if ok := yield(bucket.key, bucket.value); !ok {
				break
			}
		}
	}
}

// Set ...
func (hashMap *HashMap[KeyType, ValueType]) Set(key KeyType, value ValueType) {
	index, ok := hashMap.find(key)
	if ok {
		hashMap.buckets[index].value = value
		return
	}

	hashMap.buckets[index] = &bucket[KeyType, ValueType]{key, value}
	hashMap.size++

	loadFactor := float64(hashMap.size) / float64(len(hashMap.buckets))
	if loadFactor > hashMap.config.maxLoadFactor {
		hashMap.rehash()
	}
}

// Delete ...
func (hashMap *HashMap[KeyType, ValueType]) Delete(key KeyType) {
	index, ok := hashMap.find(key)
	if !ok {
		return
	}

	hashMap.buckets[index] = nil
	hashMap.size--
}

func (hashMap HashMap[KeyType, ValueType]) find(
	key KeyType,
) (index int, ok bool) {
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

func (hashMap *HashMap[KeyType, ValueType]) rehash() {
	newCapacity := int(float64(len(hashMap.buckets)) * hashMap.config.growFactor)
	newHashMap := newHashMapWithCapacity[KeyType, ValueType](
		hashMap.config,
		newCapacity,
	)
	for key, value := range hashMap.All() {
		newHashMap.Set(key, value)
	}

	*hashMap = *newHashMap
}

func newHashMapWithCapacity[KeyType Key, ValueType any](
	config Config,
	capacity int,
) *HashMap[KeyType, ValueType] {
	buckets := make([]*bucket[KeyType, ValueType], capacity)
	return &HashMap[KeyType, ValueType]{config: config, buckets: buckets, size: 0}
}
