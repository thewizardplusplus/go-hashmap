package hashmap

import (
	"math/rand"
)

type bucket struct {
	key   Key
	value interface{}
}

// HashMap ...
type HashMap struct {
	config  Config
	buckets []*bucket
	size    int
}

// NewHashMap ...
func NewHashMap(options ...Option) *HashMap {
	config := defaultConfig
	for _, option := range options {
		option(&config)
	}

	return newHashMapWithCapacity(config, config.initialCapacity)
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
	for _, index := range rand.Perm(len(hashMap.buckets)) {
		bucket := hashMap.buckets[index]
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
	if loadFactor > hashMap.config.maxLoadFactor {
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
	newCapacity := int(float64(len(hashMap.buckets)) * hashMap.config.growFactor)
	newHashMap := newHashMapWithCapacity(hashMap.config, newCapacity)
	hashMap.Iterate(func(key Key, value interface{}) bool {
		newHashMap.Set(key, value)
		return true
	})

	*hashMap = *newHashMap
}

func newHashMapWithCapacity(config Config, capacity int) *HashMap {
	buckets := make([]*bucket, capacity)
	return &HashMap{config: config, buckets: buckets, size: 0}
}
