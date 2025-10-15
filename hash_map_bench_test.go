package hashmap

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math/rand/v2"
	"testing"
)

type IntKey int

func (key IntKey) Hash() int {
	hash := fnv.New32()
	binary.Write(hash, binary.LittleEndian, int32(key)) // nolint: errcheck

	return int(hash.Sum32())
}

func (key IntKey) Equals(other Key) bool {
	return key == other.(IntKey)
}

func BenchmarkBuiltinMap(benchmark *testing.B) {
	for _, data := range []struct {
		name      string
		prepare   func(size int) map[int]int
		benchmark func(size int, builtinMap map[int]int)
	}{
		{
			name: "Get",
			prepare: func(size int) map[int]int {
				builtinMap := make(map[int]int)
				for i := 0; i < size; i++ {
					builtinMap[i] = i
				}

				return builtinMap
			},
			benchmark: func(size int, builtinMap map[int]int) {
				_ = builtinMap[rand.IntN(size)]
			},
		},
		{
			name: "All",
			prepare: func(size int) map[int]int {
				builtinMap := make(map[int]int)
				for i := 0; i < size; i++ {
					builtinMap[i] = i
				}

				return builtinMap
			},
			benchmark: func(size int, builtinMap map[int]int) {
				for range builtinMap {
				}
			},
		},
		{
			name:    "Set",
			prepare: func(size int) map[int]int { return make(map[int]int) },
			benchmark: func(size int, builtinMap map[int]int) {
				for i := 0; i < size; i++ {
					builtinMap[i] = i
				}
			},
		},
		{
			name: "Delete",
			prepare: func(size int) map[int]int {
				builtinMap := make(map[int]int)
				for i := 0; i < size; i++ {
					builtinMap[i] = i
				}

				return builtinMap
			},
			benchmark: func(size int, builtinMap map[int]int) {
				delete(builtinMap, rand.IntN(size))
			},
		},
	} {
		for size := 10; size <= 1e6; size *= 10 {
			name := fmt.Sprintf("%s/%d", data.name, size)
			benchmark.Run(name, func(benchmark *testing.B) {
				builtinMap := data.prepare(size)
				benchmark.ResetTimer()

				for i := 0; i < benchmark.N; i++ {
					data.benchmark(size, builtinMap)
				}
			})
		}
	}
}

func BenchmarkHashMap(benchmark *testing.B) {
	for _, data := range []struct {
		name      string
		prepare   func(size int) *HashMap[Key, any]
		benchmark func(size int, hashMap *HashMap[Key, any])
	}{
		{
			name: "Get",
			prepare: func(size int) *HashMap[Key, any] {
				hashMap := NewHashMap[Key, any]()
				for i := 0; i < size; i++ {
					hashMap.Set(IntKey(i), i)
				}

				return hashMap
			},
			benchmark: func(size int, hashMap *HashMap[Key, any]) {
				hashMap.Get(IntKey(rand.IntN(size)))
			},
		},
		{
			name: "All",
			prepare: func(size int) *HashMap[Key, any] {
				hashMap := NewHashMap[Key, any]()
				for i := 0; i < size; i++ {
					hashMap.Set(IntKey(i), i)
				}

				return hashMap
			},
			benchmark: func(size int, hashMap *HashMap[Key, any]) {
				for range hashMap.All() {
				}
			},
		},
		{
			name:    "Set",
			prepare: func(size int) *HashMap[Key, any] { return NewHashMap[Key, any]() },
			benchmark: func(size int, hashMap *HashMap[Key, any]) {
				for i := 0; i < size; i++ {
					hashMap.Set(IntKey(i), i)
				}
			},
		},
		{
			name: "Delete",
			prepare: func(size int) *HashMap[Key, any] {
				hashMap := NewHashMap[Key, any]()
				for i := 0; i < size; i++ {
					hashMap.Set(IntKey(i), i)
				}

				return hashMap
			},
			benchmark: func(size int, hashMap *HashMap[Key, any]) {
				hashMap.Delete(IntKey(rand.IntN(size)))
			},
		},
	} {
		for size := 10; size <= 1e6; size *= 10 {
			name := fmt.Sprintf("%s/%d", data.name, size)
			benchmark.Run(name, func(benchmark *testing.B) {
				hashMap := data.prepare(size)
				benchmark.ResetTimer()

				for i := 0; i < benchmark.N; i++ {
					data.benchmark(size, hashMap)
				}
			})
		}
	}
}
