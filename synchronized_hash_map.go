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
