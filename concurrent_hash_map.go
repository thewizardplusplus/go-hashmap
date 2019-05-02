package hashmap

// ConcurrentHashMap ...
type ConcurrentHashMap struct {
	segments []*SynchronizedHashMap
}

const (
	concurrencyLevel = 16
)

// NewConcurrentHashMap ...
func NewConcurrentHashMap() ConcurrentHashMap {
	var segments []*SynchronizedHashMap
	for i := 0; i < concurrencyLevel; i++ {
		segments = append(segments, NewSynchronizedHashMap())
	}

	return ConcurrentHashMap{segments: segments}
}
