package hashmap

import (
	"iter"
)

// Key ...
type Key interface {
	Hash() int
	Equals(key Key) bool
}

// Storage ...
type Storage interface {
	Get(key Key) (value interface{}, ok bool)
	All() iter.Seq2[Key, interface{}]
	Set(key Key, value interface{})
	Delete(key Key)
}
