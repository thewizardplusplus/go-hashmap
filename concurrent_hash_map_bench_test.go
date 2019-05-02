package hashmap

import (
	"fmt"
	"sync"
	"testing"
)

func BenchmarkConcurrentHashMap(benchmark *testing.B) {
	for _, data := range []struct {
		name      string
		prepare   func() ConcurrentHashMap
		benchmark func(hashMap ConcurrentHashMap)
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
