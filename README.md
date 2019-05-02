# go-hashmap

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-hashmap?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-hashmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-hashmap)](https://goreportcard.com/report/github.com/thewizardplusplus/go-hashmap)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-hashmap.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-hashmap)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-hashmap/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-hashmap)

## Installation

```
$ go get github.com/thewizardplusplus/go-hashmap
```

## Example

```go
package main

import (
	"fmt"
	"hash/fnv"
	"io"
)

type StringKey string

func (key StringKey) Hash() int {
	hash := fnv.New32()
	io.WriteString(hash, string(key))

	return int(hash.Sum32())
}

func (key StringKey) Equals(other interface{}) bool {
	return key == other.(StringKey)
}

func main() {
	timeZones := NewSynchronizedHashMap()
	timeZones.Set(StringKey("EST"), -5*60*60)
	timeZones.Set(StringKey("CST"), -6*60*60)
	timeZones.Set(StringKey("MST"), -7*60*60)

	estOffset, ok := timeZones.Get(StringKey("EST"))
	fmt.Println(estOffset, ok)
	// Output:
	// -18000 true
}
```

## Benchmarks

### SynchronizedHashMap

```
BenchmarkSynchronizedBuiltinMap/Get/1000/1-4         	10000000	      1459 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Get/1000/10-4        	 1000000	     10359 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Get/1000/100-4       	  200000	     67092 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Get/1000/1000-4      	   50000	    381497 ns/op	      18 B/op	       1 allocs/op
BenchmarkSyncMap/Get/1000/1-4                        	10000000	      1409 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Get/1000/10-4                       	 2000000	      9741 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Get/1000/100-4                      	  200000	     66314 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Get/1000/1000-4                     	   50000	    386608 ns/op	      18 B/op	       1 allocs/op
BenchmarkSynchronizedHashMap/Get/1000/1-4            	10000000	      1841 ns/op	      39 B/op	       4 allocs/op
BenchmarkSynchronizedHashMap/Get/1000/10-4           	 1000000	     14205 ns/op	     255 B/op	      40 allocs/op
BenchmarkSynchronizedHashMap/Get/1000/100-4          	  200000	     81408 ns/op	    2415 B/op	     400 allocs/op
BenchmarkSynchronizedHashMap/Get/1000/1000-4         	   30000	    477555 ns/op	   24063 B/op	    3999 allocs/op
```

```
BenchmarkSynchronizedBuiltinMap/Set/1000/1-4         	  200000	    121075 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Set/1000/10-4        	    5000	   3924351 ns/op	      42 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Set/1000/100-4       	     300	  40956121 ns/op	     356 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Set/1000/1000-4      	      30	 432632968 ns/op	    8193 B/op	      62 allocs/op
BenchmarkSyncMap/Set/1000/1-4                        	   50000	    255414 ns/op	   32002 B/op	    2999 allocs/op
BenchmarkSyncMap/Set/1000/10-4                       	    5000	   3604140 ns/op	  319904 B/op	   29981 allocs/op
BenchmarkSyncMap/Set/1000/100-4                      	     500	  37029045 ns/op	 3199333 B/op	  299812 allocs/op
BenchmarkSyncMap/Set/1000/1000-4                     	      20	 547466007 ns/op	32080577 B/op	 2999260 allocs/op
BenchmarkSynchronizedHashMap/Set/1000/1-4            	   50000	    348997 ns/op	   32002 B/op	    4998 allocs/op
BenchmarkSynchronizedHashMap/Set/1000/10-4           	    3000	   6011785 ns/op	  319918 B/op	   49973 allocs/op
BenchmarkSynchronizedHashMap/Set/1000/100-4          	     200	  60215359 ns/op	 3199397 B/op	  499740 allocs/op
BenchmarkSynchronizedHashMap/Set/1000/1000-4         	      10	1055409411 ns/op	32054326 B/op	 4998449 allocs/op
```

```
BenchmarkSynchronizedBuiltinMap/Delete/1000/1-4         	  200000	     84456 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Delete/1000/10-4        	    5000	   2994515 ns/op	      25 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Delete/1000/100-4       	     500	  32285601 ns/op	      45 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Delete/1000/1000-4      	      50	 241843645 ns/op	    2929 B/op	      35 allocs/op
BenchmarkSyncMap/Delete/1000/1-4                        	  500000	     40940 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Delete/1000/10-4                       	   10000	   1638023 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Delete/1000/100-4                      	    1000	  17531780 ns/op	      24 B/op	       1 allocs/op
BenchmarkSyncMap/Delete/1000/1000-4                     	     200	  96635575 ns/op	     702 B/op	       9 allocs/op
BenchmarkSynchronizedHashMap/Delete/1000/1-4            	   50000	    309958 ns/op	   24000 B/op	    3999 allocs/op
BenchmarkSynchronizedHashMap/Delete/1000/10-4           	    3000	   5388247 ns/op	  239949 B/op	   39981 allocs/op
BenchmarkSynchronizedHashMap/Delete/1000/100-4          	     300	  54078650 ns/op	 2399484 B/op	  399804 allocs/op
BenchmarkSynchronizedHashMap/Delete/1000/1000-4         	      10	1075765782 ns/op	24031883 B/op	 3998525 allocs/op
```

### HashMap

```
BenchmarkBuiltinMap/Get/10-4         	200000000	        83.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Get/100-4        	200000000	        85.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Get/1000-4       	200000000	        82.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Get/10000-4      	200000000	        88.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Get/100000-4     	100000000	       110 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Get/1000000-4    	100000000	       232 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashMap/Get/10-4            	50000000	       274 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Get/100-4           	50000000	       274 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Get/1000-4          	50000000	       281 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Get/10000-4         	50000000	       342 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Get/100000-4        	30000000	       560 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Get/1000000-4       	20000000	       654 ns/op	      23 B/op	       3 allocs/op
```

```
BenchmarkBuiltinMap/Set/10-4         	50000000	       319 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Set/100-4        	 5000000	      3568 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Set/1000-4       	  500000	     32145 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Set/10000-4      	   50000	    382942 ns/op	      13 B/op	       0 allocs/op
BenchmarkBuiltinMap/Set/100000-4     	    2000	   5859511 ns/op	    2869 B/op	       1 allocs/op
BenchmarkBuiltinMap/Set/1000000-4    	     100	 136460549 ns/op	  879443 B/op	     384 allocs/op
BenchmarkHashMap/Set/10-4            	 5000000	      2457 ns/op	     304 B/op	      47 allocs/op
BenchmarkHashMap/Set/100-4           	  500000	     25061 ns/op	    3184 B/op	     497 allocs/op
BenchmarkHashMap/Set/1000-4          	   50000	    256028 ns/op	   31986 B/op	    4997 allocs/op
BenchmarkHashMap/Set/10000-4         	    5000	   3133410 ns/op	  320219 B/op	   50008 allocs/op
BenchmarkHashMap/Set/100000-4        	     300	  50005020 ns/op	 3256090 B/op	  502951 allocs/op
BenchmarkHashMap/Set/1000000-4       	      20	 551907702 ns/op	39052584 B/op	 5364570 allocs/op
```

```
BenchmarkBuiltinMap/Delete/10-4         	300000000	        48.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Delete/100-4        	30000000	       445 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Delete/1000-4       	 3000000	      4338 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Delete/10000-4      	  300000	     43228 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Delete/100000-4     	   30000	    432701 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Delete/1000000-4    	    3000	   4378122 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashMap/Delete/10-4            	10000000	      2067 ns/op	     232 B/op	      38 allocs/op
BenchmarkHashMap/Delete/100-4           	 1000000	     20780 ns/op	    2392 B/op	     398 allocs/op
BenchmarkHashMap/Delete/1000-4          	  100000	    207031 ns/op	   23992 B/op	    3998 allocs/op
BenchmarkHashMap/Delete/10000-4         	   10000	   2170782 ns/op	  239992 B/op	   39998 allocs/op
BenchmarkHashMap/Delete/100000-4        	     500	  25951595 ns/op	 2400000 B/op	  399998 allocs/op
BenchmarkHashMap/Delete/1000000-4       	      50	 297262710 ns/op	24000005 B/op	 3999998 allocs/op
```

## License

The MIT License (MIT)

Copyright &copy; 2019 thewizardplusplus
