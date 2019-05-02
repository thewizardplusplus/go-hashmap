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
	timeZones := NewConcurrentHashMap()
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

### SynchronizedHashMap & ConcurrentHashMap

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
BenchmarkConcurrentHashMap/Get/1000/1-4         	 3000000	      4316 ns/op	      55 B/op	       7 allocs/op
BenchmarkConcurrentHashMap/Get/1000/10-4        	 1000000	     20194 ns/op	     415 B/op	      70 allocs/op
BenchmarkConcurrentHashMap/Get/1000/100-4       	  200000	    101014 ns/op	    4015 B/op	     700 allocs/op
BenchmarkConcurrentHashMap/Get/1000/1000-4      	   20000	    643696 ns/op	   40085 B/op	    6999 allocs/op
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
BenchmarkConcurrentHashMap/Set/1000/1-4         	   20000	    632797 ns/op	   48006 B/op	    7997 allocs/op
BenchmarkConcurrentHashMap/Set/1000/10-4        	    5000	   3333305 ns/op	  479906 B/op	   79962 allocs/op
BenchmarkConcurrentHashMap/Set/1000/100-4       	     500	  33006847 ns/op	 4799931 B/op	  799631 allocs/op
BenchmarkConcurrentHashMap/Set/1000/1000-4      	      20	 667340402 ns/op	48085449 B/op	 7997552 allocs/op
```

```
BenchmarkSynchronizedBuiltinMap/Delete/1000/1-4         	10000000	      1401 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Delete/1000/10-4        	 2000000	      9781 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Delete/1000/100-4       	  200000	     67847 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Delete/1000/1000-4      	   50000	    393990 ns/op	      18 B/op	       1 allocs/op
BenchmarkSyncMap/Delete/1000/1-4                        	10000000	      1329 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Delete/1000/10-4                       	 2000000	      8929 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Delete/1000/100-4                      	  200000	     61734 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Delete/1000/1000-4                     	   50000	    374220 ns/op	      18 B/op	       1 allocs/op
BenchmarkSynchronizedHashMap/Delete/1000/1-4            	10000000	      1845 ns/op	      39 B/op	       4 allocs/op
BenchmarkSynchronizedHashMap/Delete/1000/10-4           	 1000000	     13879 ns/op	     256 B/op	      40 allocs/op
BenchmarkSynchronizedHashMap/Delete/1000/100-4          	  100000	    128004 ns/op	    2420 B/op	     400 allocs/op
BenchmarkSynchronizedHashMap/Delete/1000/1000-4         	   10000	   1034500 ns/op	   24281 B/op	    4002 allocs/op
BenchmarkConcurrentHashMap/Delete/1000/1-4      	 5000000	      3330 ns/op	      55 B/op	       7 allocs/op
BenchmarkConcurrentHashMap/Delete/1000/10-4     	 1000000	     17399 ns/op	     415 B/op	      70 allocs/op
BenchmarkConcurrentHashMap/Delete/1000/100-4    	  200000	     96695 ns/op	    4016 B/op	     700 allocs/op
BenchmarkConcurrentHashMap/Delete/1000/1000-4   	   20000	    595056 ns/op	   40126 B/op	    6999 allocs/op
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
BenchmarkBuiltinMap/Delete/10-4     	300000000	        49.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Delete/100-4    	300000000	        49.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Delete/1000-4   	300000000	        48.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Delete/10000-4  	300000000	        47.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Delete/100000-4 	300000000	        47.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Delete/1000000-4         	200000000	        64.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashMap/Delete/10-4                 	50000000	       249 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Delete/100-4                	50000000	       247 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Delete/1000-4               	50000000	       250 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Delete/10000-4              	50000000	       258 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Delete/100000-4             	50000000	       303 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Delete/1000000-4            	30000000	       351 ns/op	      24 B/op	       3 allocs/op
```

## License

The MIT License (MIT)

Copyright &copy; 2019 thewizardplusplus
