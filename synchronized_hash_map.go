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
