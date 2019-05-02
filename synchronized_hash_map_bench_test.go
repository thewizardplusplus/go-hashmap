package hashmap

import (
	"fmt"
	"sync"
	"testing"
)

const (
	sizeForSyncBench = 1000
)

func BenchmarkSynchronizedHashMap(benchmark *testing.B) {
	for _, data := range []struct {
		name      string
		prepare   func() *SynchronizedHashMap
		benchmark func(hashMap *SynchronizedHashMap)
	}{
		// TODO: add benchmark cases
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
