package hashmap

import (
	"sync"
)

// SynchronizedHashMap ...
type SynchronizedHashMap struct {
	innerMap *HashMap
	lock     sync.RWMutex
}
