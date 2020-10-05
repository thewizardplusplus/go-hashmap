package hashmap

//go:generate mockery -name=Key -inpkg -case=underscore -testonly

// Key ...
type Key interface {
	Hash() int
	Equals(key Key) bool
}

// Handler ...
type Handler func(key Key, value interface{}) bool

// Storage ...
type Storage interface {
	Get(key Key) (value interface{}, ok bool)
	Iterate(handler Handler) bool
	Set(key Key, value interface{})
	Delete(key Key)
}
