package hashmap

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math/rand"
	"testing"
)

type IntKey int

func (key IntKey) Hash() int {
	hash := fnv.New32()
	binary.Write(hash, binary.LittleEndian, int32(key)) // nolint: errcheck

	return int(hash.Sum32())
}

func (key IntKey) Equals(other interface{}) bool {
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
				_ = builtinMap[rand.Intn(size)]
			},
		},
		{
			name: "Iterate",
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
				delete(builtinMap, rand.Intn(size))
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
		prepare   func(size int) *HashMap
		benchmark func(size int, hashMap *HashMap)
	}{
		{
			name: "Get",
			prepare: func(size int) *HashMap {
				hashMap := NewHashMap()
				for i := 0; i < size; i++ {
					hashMap.Set(IntKey(i), i)
				}

				return hashMap
			},
			benchmark: func(size int, hashMap *HashMap) {
				hashMap.Get(IntKey(rand.Intn(size)))
			},
		},
		{
			name: "Iterate",
			prepare: func(size int) *HashMap {
				hashMap := NewHashMap()
				for i := 0; i < size; i++ {
					hashMap.Set(IntKey(i), i)
				}

				return hashMap
			},
			benchmark: func(size int, hashMap *HashMap) {
				hashMap.Iterate(func(key Key, value interface{}) bool { return true })
			},
		},
		{
			name:    "Set",
			prepare: func(size int) *HashMap { return NewHashMap() },
			benchmark: func(size int, hashMap *HashMap) {
				for i := 0; i < size; i++ {
					hashMap.Set(IntKey(i), i)
				}
			},
		},
		{
			name: "Delete",
			prepare: func(size int) *HashMap {
				hashMap := NewHashMap()
				for i := 0; i < size; i++ {
					hashMap.Set(IntKey(i), i)
				}

				return hashMap
			},
			benchmark: func(size int, hashMap *HashMap) {
				hashMap.Delete(IntKey(rand.Intn(size)))
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
