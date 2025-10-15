package hashmap

import (
	"iter"
	"sync"
)

// SynchronizedHashMap ...
//
// It's safe for concurrent access because it uses a mutex lock to access
// the inner map.
type SynchronizedHashMap[KeyType Key, ValueType any] struct {
	lock     sync.RWMutex
	innerMap Storage[KeyType, ValueType]
}

// NewSynchronizedHashMap ...
func NewSynchronizedHashMap[KeyType Key, ValueType any](
	options ...SynchronizedOption[KeyType, ValueType],
) *SynchronizedHashMap[KeyType, ValueType] {
	// you can't move the default synchronized config into a global variable
	// because the default inner map should be new every time
	config := SynchronizedConfig[KeyType, ValueType]{
		innerMap: NewHashMap[KeyType, ValueType](),
	}
	for _, option := range options {
		option(&config)
	}

	return &SynchronizedHashMap[KeyType, ValueType]{innerMap: config.innerMap}
}

// Get ...
func (hashMap *SynchronizedHashMap[KeyType, ValueType]) Get(
	key KeyType,
) (value ValueType, ok bool) {
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
func (hashMap *SynchronizedHashMap[KeyType, ValueType]) All() iter.Seq2[KeyType, ValueType] { //nolint:lll
	return func(yield func(KeyType, ValueType) bool) {
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
func (hashMap *SynchronizedHashMap[KeyType, ValueType]) Set(
	key KeyType,
	value ValueType,
) {
	hashMap.lock.Lock()
	defer hashMap.lock.Unlock()

	hashMap.innerMap.Set(key, value)
}

// Delete ...
func (hashMap *SynchronizedHashMap[KeyType, ValueType]) Delete(key KeyType) {
	hashMap.lock.Lock()
	defer hashMap.lock.Unlock()

	hashMap.innerMap.Delete(key)
}
