package hashmap

import (
	"fmt"
	"testing"
)

func BenchmarkHashMap(benchmark *testing.B) {
	for _, data := range []struct {
		name      string
		prepare   func(size int) *HashMap
		benchmark func(size int, hashMap *HashMap)
	}{
		// TODO: add benchmark cases
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
