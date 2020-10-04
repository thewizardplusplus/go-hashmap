package hashmap

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

func BenchmarkConcurrentHashMap(benchmark *testing.B) {
	for _, data := range []struct {
		name      string
		prepare   func() ConcurrentHashMap
		benchmark func(hashMap ConcurrentHashMap)
	}{
		{
			name: "Get",
			prepare: func() ConcurrentHashMap {
				hashMap := NewConcurrentHashMap(16, func() Storage {
					return NewSynchronizedHashMap(NewHashMap())
				})
				for i := 0; i < sizeForSyncBench; i++ {
					hashMap.Set(IntKey(i), i)
				}

				return hashMap
			},
			benchmark: func(hashMap ConcurrentHashMap) {
				hashMap.Get(IntKey(rand.Intn(sizeForSyncBench)))
			},
		},
		{
			name: "Iterate",
			prepare: func() ConcurrentHashMap {
				hashMap := NewConcurrentHashMap(16, func() Storage {
					return NewSynchronizedHashMap(NewHashMap())
				})
				for i := 0; i < sizeForSyncBench; i++ {
					hashMap.Set(IntKey(i), i)
				}

				return hashMap
			},
			benchmark: func(hashMap ConcurrentHashMap) {
				hashMap.Iterate(func(key Key, value interface{}) bool { return true })
			},
		},
		{
			name: "Set",
			prepare: func() ConcurrentHashMap {
				return NewConcurrentHashMap(16, func() Storage {
					return NewSynchronizedHashMap(NewHashMap())
				})
			},
			benchmark: func(hashMap ConcurrentHashMap) {
				for i := 0; i < sizeForSyncBench; i++ {
					hashMap.Set(IntKey(i), i)
				}
			},
		},
		{
			name: "Delete",
			prepare: func() ConcurrentHashMap {
				hashMap := NewConcurrentHashMap(16, func() Storage {
					return NewSynchronizedHashMap(NewHashMap())
				})
				for i := 0; i < sizeForSyncBench; i++ {
					hashMap.Set(IntKey(i), i)
				}

				return hashMap
			},
			benchmark: func(hashMap ConcurrentHashMap) {
				hashMap.Delete(IntKey(rand.Intn(sizeForSyncBench)))
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
