package hashmap

import (
	"iter"
	"sync"
)

// SynchronizedHashMap ...
//
// It's safe for concurrent access because it uses a mutex lock to access
// the inner map.
type SynchronizedHashMap struct {
	lock     sync.RWMutex
	innerMap Storage
}

// NewSynchronizedHashMap ...
func NewSynchronizedHashMap(
	options ...SynchronizedOption,
) *SynchronizedHashMap {
	// you can't move the default synchronized config into a global variable
	// because the default inner map should be new every time
	config := SynchronizedConfig{innerMap: NewHashMap[Key, any]()}
	for _, option := range options {
		option(&config)
	}

	return &SynchronizedHashMap{innerMap: config.innerMap}
}

// Get ...
func (hashMap *SynchronizedHashMap) Get(key Key) (value interface{}, ok bool) {
	hashMap.lock.RLock()
	defer hashMap.lock.RUnlock()

	return hashMap.innerMap.Get(key)
}

// All ...
//
// It randomizes of iteration order.
//
// A mutex lock is using only for iteration, not for handling (the handler
// is called out of lock).
func (hashMap *SynchronizedHashMap) All() iter.Seq2[Key, interface{}] {
	return func(yield func(Key, interface{}) bool) {
		hashMap.lock.RLock()
		defer hashMap.lock.RUnlock()

		for key, value := range hashMap.innerMap.All() {
			if ok := func() bool {
				hashMap.lock.RUnlock()
				defer hashMap.lock.RLock()

				return yield(key, value)
			}(); !ok {
				break
			}
		}
	}
}

// Set ...
func (hashMap *SynchronizedHashMap) Set(key Key, value interface{}) {
	hashMap.lock.Lock()
	defer hashMap.lock.Unlock()

	hashMap.innerMap.Set(key, value)
}

// Delete ...
func (hashMap *SynchronizedHashMap) Delete(key Key) {
	hashMap.lock.Lock()
	defer hashMap.lock.Unlock()

	hashMap.innerMap.Delete(key)
}
