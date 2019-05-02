package hashmap

// ConcurrentHashMap ...
type ConcurrentHashMap struct {
	segments []*SynchronizedHashMap
}
