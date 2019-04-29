package hashmap

// Key ...
type Key interface {
	Hash() int
	Equals(key interface{}) bool
}
