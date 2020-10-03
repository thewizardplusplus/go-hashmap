package hashmap

// Storage ...
type Storage interface {
	Get(key Key) (value interface{}, ok bool)
	Iterate(handler Handler) bool
	Set(key Key, value interface{})
	Delete(key Key)
}
