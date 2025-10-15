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
type Storage[KeyType Key, ValueType any] interface {
	Get(key KeyType) (value ValueType, ok bool)
	All() iter.Seq2[KeyType, ValueType]
	Set(key KeyType, value ValueType)
	Delete(key KeyType)
}
