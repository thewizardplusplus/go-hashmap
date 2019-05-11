package hashmap

import (
	"sync"
)

// SynchronizedHashMap ...
type SynchronizedHashMap struct {
	innerMap *HashMap
	lock     sync.RWMutex
}

// NewSynchronizedHashMap ...
func NewSynchronizedHashMap() *SynchronizedHashMap {
	return &SynchronizedHashMap{innerMap: NewHashMap()}
}

// Get ...
func (hashMap *SynchronizedHashMap) Get(key Key) (value interface{}, ok bool) {
	hashMap.lock.RLock()
	defer hashMap.lock.RUnlock()

	return hashMap.innerMap.Get(key)
}

// Iterate ...
func (hashMap *SynchronizedHashMap) Iterate(handler Handler) {
	hashMap.lock.RLock()
	defer hashMap.lock.RUnlock()

	hashMap.innerMap.Iterate(func(key Key, value interface{}) {
		hashMap.lock.RUnlock()
		defer hashMap.lock.RLock()

		handler(key, value)
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
