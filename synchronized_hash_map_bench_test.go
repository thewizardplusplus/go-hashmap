package hashmap

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"sync"
	"testing"
)

type SynchronizedBuiltinMap struct {
	lock     sync.RWMutex
	innerMap map[int]int
}

func NewSynchronizedBuiltinMap() *SynchronizedBuiltinMap {
	innerMap := make(map[int]int)
	return &SynchronizedBuiltinMap{innerMap: innerMap}
}

func (builtinMap *SynchronizedBuiltinMap) Get(key int) int {
	builtinMap.lock.RLock()
	defer builtinMap.lock.RUnlock()

	return builtinMap.innerMap[key]
}

func (builtinMap *SynchronizedBuiltinMap) All() iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		builtinMap.lock.RLock()
		defer builtinMap.lock.RUnlock()

		for key, value := range builtinMap.innerMap {
			if ok := func() bool {
				builtinMap.lock.RUnlock()
				defer builtinMap.lock.RLock()

				return yield(key, value)
			}(); !ok {
				break
			}
		}
	}
}

func (builtinMap *SynchronizedBuiltinMap) Set(key int, value int) {
	builtinMap.lock.Lock()
	defer builtinMap.lock.Unlock()

	builtinMap.innerMap[key] = value
}

func (builtinMap *SynchronizedBuiltinMap) Delete(key int) {
	builtinMap.lock.Lock()
	defer builtinMap.lock.Unlock()

	delete(builtinMap.innerMap, key)
}

const (
	sizeForSyncBench = 1000
)

func BenchmarkSynchronizedBuiltinMap(benchmark *testing.B) {
	for _, data := range []struct {
		name      string
		prepare   func() *SynchronizedBuiltinMap
		benchmark func(builtinMap *SynchronizedBuiltinMap)
	}{
		{
			name: "Get",
			prepare: func() *SynchronizedBuiltinMap {
				builtinMap := NewSynchronizedBuiltinMap()
				for i := 0; i < sizeForSyncBench; i++ {
					builtinMap.Set(i, i)
				}

				return builtinMap
			},
			benchmark: func(builtinMap *SynchronizedBuiltinMap) {
				builtinMap.Get(rand.IntN(sizeForSyncBench))
			},
		},
		{
			name: "All",
			prepare: func() *SynchronizedBuiltinMap {
				builtinMap := NewSynchronizedBuiltinMap()
				for i := 0; i < sizeForSyncBench; i++ {
					builtinMap.Set(i, i)
				}

				return builtinMap
			},
			benchmark: func(builtinMap *SynchronizedBuiltinMap) {
				for range builtinMap.All() {
				}
			},
		},
		{
			name: "Set",
			prepare: func() *SynchronizedBuiltinMap {
				return NewSynchronizedBuiltinMap()
			},
			benchmark: func(builtinMap *SynchronizedBuiltinMap) {
				for i := 0; i < sizeForSyncBench; i++ {
					builtinMap.Set(i, i)
				}
			},
		},
		{
			name: "Delete",
			prepare: func() *SynchronizedBuiltinMap {
				builtinMap := NewSynchronizedBuiltinMap()
				for i := 0; i < sizeForSyncBench; i++ {
					builtinMap.Set(i, i)
				}

				return builtinMap
			},
			benchmark: func(builtinMap *SynchronizedBuiltinMap) {
				builtinMap.Delete(rand.IntN(sizeForSyncBench))
			},
		},
	} {
		for threads := 1; threads <= 1e3; threads *= 10 {
			name := fmt.Sprintf("%s/%d/%d", data.name, sizeForSyncBench, threads)
			benchmark.Run(name, func(benchmark *testing.B) {
				builtinMap := data.prepare()
				benchmark.ResetTimer()

				for i := 0; i < benchmark.N; i++ {
					var waiter sync.WaitGroup
					waiter.Add(threads)

					for j := 0; j < threads; j++ {
						go func() {
							defer waiter.Done()
							data.benchmark(builtinMap)
						}()
					}

					waiter.Wait()
				}
			})
		}
	}
}

func BenchmarkSyncMap(benchmark *testing.B) {
	for _, data := range []struct {
		name      string
		prepare   func() *sync.Map
		benchmark func(syncMap *sync.Map)
	}{
		{
			name: "Get",
			prepare: func() *sync.Map {
				syncMap := new(sync.Map)
				for i := 0; i < sizeForSyncBench; i++ {
					syncMap.Store(i, i)
				}

				return syncMap
			},
			benchmark: func(syncMap *sync.Map) {
				syncMap.Load(rand.IntN(sizeForSyncBench))
			},
		},
		{
			name: "All",
			prepare: func() *sync.Map {
				syncMap := new(sync.Map)
				for i := 0; i < sizeForSyncBench; i++ {
					syncMap.Store(i, i)
				}

				return syncMap
			},
			benchmark: func(syncMap *sync.Map) {
				for range syncMap.Range {
				}
			},
		},
		{
			name:    "Set",
			prepare: func() *sync.Map { return new(sync.Map) },
			benchmark: func(syncMap *sync.Map) {
				for i := 0; i < sizeForSyncBench; i++ {
					syncMap.Store(i, i)
				}
			},
		},
		{
			name: "Delete",
			prepare: func() *sync.Map {
				syncMap := new(sync.Map)
				for i := 0; i < sizeForSyncBench; i++ {
					syncMap.Store(i, i)
				}

				return syncMap
			},
			benchmark: func(syncMap *sync.Map) {
				syncMap.Delete(rand.IntN(sizeForSyncBench))
			},
		},
	} {
		for threads := 1; threads <= 1e3; threads *= 10 {
			name := fmt.Sprintf("%s/%d/%d", data.name, sizeForSyncBench, threads)
			benchmark.Run(name, func(benchmark *testing.B) {
				syncMap := data.prepare()
				benchmark.ResetTimer()

				for i := 0; i < benchmark.N; i++ {
					var waiter sync.WaitGroup
					waiter.Add(threads)

					for j := 0; j < threads; j++ {
						go func() {
							defer waiter.Done()
							data.benchmark(syncMap)
						}()
					}

					waiter.Wait()
				}
			})
		}
	}
}

func BenchmarkSynchronizedHashMap(benchmark *testing.B) {
	for _, data := range []struct {
		name      string
		prepare   func() *SynchronizedHashMap
		benchmark func(hashMap *SynchronizedHashMap)
	}{
		{
			name: "Get",
			prepare: func() *SynchronizedHashMap {
				hashMap := NewSynchronizedHashMap()
				for i := 0; i < sizeForSyncBench; i++ {
					hashMap.Set(IntKey(i), i)
				}

				return hashMap
			},
			benchmark: func(hashMap *SynchronizedHashMap) {
				hashMap.Get(IntKey(rand.IntN(sizeForSyncBench)))
			},
		},
		{
			name: "All",
			prepare: func() *SynchronizedHashMap {
				hashMap := NewSynchronizedHashMap()
				for i := 0; i < sizeForSyncBench; i++ {
					hashMap.Set(IntKey(i), i)
				}

				return hashMap
			},
			benchmark: func(hashMap *SynchronizedHashMap) {
				for range hashMap.All() {
				}
			},
		},
		{
			name:    "Set",
			prepare: func() *SynchronizedHashMap { return NewSynchronizedHashMap() },
			benchmark: func(hashMap *SynchronizedHashMap) {
				for i := 0; i < sizeForSyncBench; i++ {
					hashMap.Set(IntKey(i), i)
				}
			},
		},
		{
			name: "Delete",
			prepare: func() *SynchronizedHashMap {
				hashMap := NewSynchronizedHashMap()
				for i := 0; i < sizeForSyncBench; i++ {
					hashMap.Set(IntKey(i), i)
				}

				return hashMap
			},
			benchmark: func(hashMap *SynchronizedHashMap) {
				hashMap.Delete(IntKey(rand.IntN(sizeForSyncBench)))
			},
		},
	} {
		for threads := 1; threads <= 1e3; threads *= 10 {
			name := fmt.Sprintf("%s/%d/%d", data.name, sizeForSyncBench, threads)
			benchmark.Run(name, func(benchmark *testing.B) {
				hashMap := data.prepare()
				benchmark.ResetTimer()

				for i := 0; i < benchmark.N; i++ {
					var waiter sync.WaitGroup
					waiter.Add(threads)

					for j := 0; j < threads; j++ {
						go func() {
							defer waiter.Done()
							data.benchmark(hashMap)
						}()
					}

					waiter.Wait()
				}
			})
		}
	}
}
