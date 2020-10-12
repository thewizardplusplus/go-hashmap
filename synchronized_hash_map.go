package hashmap

import (
	"sync"
)

// SynchronizedHashMap ...
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
	config := SynchronizedConfig{innerMap: NewHashMap()}
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

// Iterate ...
func (hashMap *SynchronizedHashMap) Iterate(handler Handler) bool {
	hashMap.lock.RLock()
	defer hashMap.lock.RUnlock()

	return hashMap.innerMap.Iterate(func(key Key, value interface{}) bool {
		hashMap.lock.RUnlock()
		defer hashMap.lock.RLock()

		return handler(key, value)
	})
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
