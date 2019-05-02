package hashmap

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

const (
	sizeForSyncBench = 1000
)

func BenchmarkSyncMap(benchmark *testing.B) {
	for _, data := range []struct {
		name      string
		prepare   func() *sync.Map
		benchmark func(syncMap *sync.Map)
	}{
		// TODO: add benchmark cases
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
				hashMap.Get(IntKey(rand.Intn(sizeForSyncBench)))
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
				for i := 0; i < sizeForSyncBench; i++ {
					hashMap.Delete(IntKey(i))
				}
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
