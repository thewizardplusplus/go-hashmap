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

func (key StringKey) Equals(other Key) bool {
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
BenchmarkSynchronizedBuiltinMap/Get/1000/1-4    	10000000	      1403 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Get/1000/10-4   	 2000000	     10086 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Get/1000/100-4  	  200000	     67435 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Get/1000/1000-4 	   50000	    385837 ns/op	      18 B/op	       1 allocs/op
BenchmarkSyncMap/Get/1000/1-4                   	10000000	      1373 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Get/1000/10-4                  	 2000000	      9594 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Get/1000/100-4                 	  200000	     65672 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Get/1000/1000-4                	   50000	    381484 ns/op	      18 B/op	       1 allocs/op
BenchmarkSynchronizedHashMap/Get/1000/1-4       	10000000	      1905 ns/op	      39 B/op	       4 allocs/op
BenchmarkSynchronizedHashMap/Get/1000/10-4      	 1000000	     14374 ns/op	     255 B/op	      40 allocs/op
BenchmarkSynchronizedHashMap/Get/1000/100-4     	  200000	     83661 ns/op	    2415 B/op	     400 allocs/op
BenchmarkSynchronizedHashMap/Get/1000/1000-4    	   30000	    494123 ns/op	   24065 B/op	    3999 allocs/op
BenchmarkConcurrentHashMap/Get/1000/1-4         	 3000000	      4145 ns/op	      55 B/op	       7 allocs/op
BenchmarkConcurrentHashMap/Get/1000/10-4        	 1000000	     19554 ns/op	     415 B/op	      70 allocs/op
BenchmarkConcurrentHashMap/Get/1000/100-4       	  200000	    100572 ns/op	    4015 B/op	     700 allocs/op
BenchmarkConcurrentHashMap/Get/1000/1000-4      	   20000	    638640 ns/op	   40076 B/op	    6998 allocs/op
```

```
BenchmarkSynchronizedBuiltinMap/Iterate/1000/1-4    	  300000	     42922 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Iterate/1000/10-4   	   30000	    445263 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Iterate/1000/100-4  	    3000	   4047902 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Iterate/1000/1000-4 	     300	  40699462 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Iterate/1000/1-4                   	  500000	     32798 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Iterate/1000/10-4                  	  100000	    192020 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Iterate/1000/100-4                 	   10000	   1486500 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Iterate/1000/1000-4                	    1000	  14725244 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedHashMap/Iterate/1000/1-4       	  100000	    207407 ns/op	   16400 B/op	       2 allocs/op
BenchmarkSynchronizedHashMap/Iterate/1000/10-4      	    3000	   4468235 ns/op	  163871 B/op	      11 allocs/op
BenchmarkSynchronizedHashMap/Iterate/1000/100-4     	     300	  45543225 ns/op	 1643317 B/op	     174 allocs/op
BenchmarkSynchronizedHashMap/Iterate/1000/1000-4    	      50	 351090497 ns/op	16484443 B/op	    2394 allocs/op
BenchmarkConcurrentHashMap/Iterate/1000/1-4         	  100000	    205957 ns/op	   16528 B/op	      18 allocs/op
BenchmarkConcurrentHashMap/Iterate/1000/10-4        	    5000	   3937290 ns/op	  165158 B/op	     171 allocs/op
BenchmarkConcurrentHashMap/Iterate/1000/100-4       	     300	  41970197 ns/op	 1651729 B/op	    1707 allocs/op
BenchmarkConcurrentHashMap/Iterate/1000/1000-4      	      30	 412794971 ns/op	16609136 B/op	   18373 allocs/op
```

```
BenchmarkSynchronizedBuiltinMap/Set/1000/1-4    	  100000	    118490 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Set/1000/10-4   	    5000	   4004046 ns/op	      34 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Set/1000/100-4  	     300	  41541398 ns/op	     359 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Set/1000/1000-4 	      30	 460350938 ns/op	    7329 B/op	      58 allocs/op
BenchmarkSyncMap/Set/1000/1-4                   	   50000	    253130 ns/op	   32002 B/op	    2999 allocs/op
BenchmarkSyncMap/Set/1000/10-4                  	    5000	   3618081 ns/op	  319901 B/op	   29981 allocs/op
BenchmarkSyncMap/Set/1000/100-4                 	     500	  37144669 ns/op	 3199366 B/op	  299813 allocs/op
BenchmarkSyncMap/Set/1000/1000-4                	      30	 537376345 ns/op	32073810 B/op	 2999177 allocs/op
BenchmarkSynchronizedHashMap/Set/1000/1-4       	   50000	    350498 ns/op	   32002 B/op	    4998 allocs/op
BenchmarkSynchronizedHashMap/Set/1000/10-4      	    3000	   5968663 ns/op	  319917 B/op	   49973 allocs/op
BenchmarkSynchronizedHashMap/Set/1000/100-4     	     300	  59637726 ns/op	 3199157 B/op	  499728 allocs/op
BenchmarkSynchronizedHashMap/Set/1000/1000-4    	      10	1097406544 ns/op	32042888 B/op	 4998294 allocs/op
BenchmarkConcurrentHashMap/Set/1000/1-4         	   20000	    598574 ns/op	   48006 B/op	    7997 allocs/op
BenchmarkConcurrentHashMap/Set/1000/10-4        	    5000	   3239388 ns/op	  479908 B/op	   79962 allocs/op
BenchmarkConcurrentHashMap/Set/1000/100-4       	     500	  32427188 ns/op	 4799851 B/op	  799630 allocs/op
BenchmarkConcurrentHashMap/Set/1000/1000-4      	      20	 678751115 ns/op	48080325 B/op	 7997485 allocs/op
```

```
BenchmarkSynchronizedBuiltinMap/Delete/1000/1-4    	10000000	      1392 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Delete/1000/10-4   	 2000000	      9892 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Delete/1000/100-4  	  200000	     68712 ns/op	      16 B/op	       1 allocs/op
BenchmarkSynchronizedBuiltinMap/Delete/1000/1000-4 	   50000	    399065 ns/op	      18 B/op	       1 allocs/op
BenchmarkSyncMap/Delete/1000/1-4                   	10000000	      1325 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Delete/1000/10-4                  	 2000000	      8876 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Delete/1000/100-4                 	  200000	     63251 ns/op	      16 B/op	       1 allocs/op
BenchmarkSyncMap/Delete/1000/1000-4                	   50000	    385645 ns/op	      18 B/op	       1 allocs/op
BenchmarkSynchronizedHashMap/Delete/1000/1-4       	10000000	      1862 ns/op	      39 B/op	       4 allocs/op
BenchmarkSynchronizedHashMap/Delete/1000/10-4      	 1000000	     13655 ns/op	     256 B/op	      40 allocs/op
BenchmarkSynchronizedHashMap/Delete/1000/100-4     	  100000	    128966 ns/op	    2420 B/op	     400 allocs/op
BenchmarkSynchronizedHashMap/Delete/1000/1000-4    	   10000	   1037180 ns/op	   24337 B/op	    4003 allocs/op
BenchmarkConcurrentHashMap/Delete/1000/1-4         	 5000000	      3149 ns/op	      55 B/op	       7 allocs/op
BenchmarkConcurrentHashMap/Delete/1000/10-4        	 1000000	     16714 ns/op	     415 B/op	      70 allocs/op
BenchmarkConcurrentHashMap/Delete/1000/100-4       	  200000	     94336 ns/op	    4016 B/op	     700 allocs/op
BenchmarkConcurrentHashMap/Delete/1000/1000-4      	   30000	    575415 ns/op	   40134 B/op	    6999 allocs/op
```

### HashMap

```
BenchmarkBuiltinMap/Get/10-4         	200000000	        84 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Get/100-4        	200000000	        86 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Get/1000-4       	200000000	        85 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Get/10000-4      	200000000	        91 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Get/100000-4     	100000000	       113 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Get/1000000-4    	100000000	       236 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashMap/Get/10-4            	 50000000	       264 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Get/100-4           	 50000000	       270 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Get/1000-4          	 50000000	       275 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Get/10000-4         	 50000000	       340 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Get/100000-4        	 30000000	       564 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Get/1000000-4       	 20000000	       662 ns/op	      23 B/op	       3 allocs/op
```

```
BenchmarkBuiltinMap/Iterate/10-4         	50000000	       294 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Iterate/100-4        	10000000	      2314 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Iterate/1000-4       	  500000	     25152 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Iterate/10000-4      	  100000	    234488 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Iterate/100000-4     	   10000	   2308275 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Iterate/1000000-4    	     500	  27197563 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashMap/Iterate/10-4            	20000000	      1068 ns/op	     128 B/op	       1 allocs/op
BenchmarkHashMap/Iterate/100-4           	 1000000	     15677 ns/op	    2048 B/op	       1 allocs/op
BenchmarkHashMap/Iterate/1000-4          	  100000	    128244 ns/op	   16384 B/op	       1 allocs/op
BenchmarkHashMap/Iterate/10000-4         	   10000	   1153831 ns/op	  131072 B/op	       1 allocs/op
BenchmarkHashMap/Iterate/100000-4        	    1000	  23338095 ns/op	 2097154 B/op	       1 allocs/op
BenchmarkHashMap/Iterate/1000000-4       	      50	 369991586 ns/op	16777216 B/op	       1 allocs/op
```

```
BenchmarkBuiltinMap/Set/10-4         	 50000000	       331 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Set/100-4        	  5000000	      3697 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Set/1000-4       	   500000	     33502 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Set/10000-4      	    30000	    399097 ns/op	      22 B/op	       0 allocs/op
BenchmarkBuiltinMap/Set/100000-4     	     2000	   6047494 ns/op	    2879 B/op	       2 allocs/op
BenchmarkBuiltinMap/Set/1000000-4    	      100	 144642324 ns/op	  879348 B/op	     383 allocs/op
BenchmarkHashMap/Set/10-4            	 10000000	      2343 ns/op	     304 B/op	      47 allocs/op
BenchmarkHashMap/Set/100-4           	  1000000	     24307 ns/op	    3184 B/op	     497 allocs/op
BenchmarkHashMap/Set/1000-4          	    50000	    245038 ns/op	   31986 B/op	    4997 allocs/op
BenchmarkHashMap/Set/10000-4         	     5000	   3013202 ns/op	  320219 B/op	   50008 allocs/op
BenchmarkHashMap/Set/100000-4        	      300	  47913658 ns/op	 3256090 B/op	  502951 allocs/op
BenchmarkHashMap/Set/1000000-4       	       20	 526220479 ns/op	39052582 B/op	 5364570 allocs/op
```

```
BenchmarkBuiltinMap/Delete/10-4         	300000000	        48 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Delete/100-4        	300000000	        48 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Delete/1000-4       	300000000	        48 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Delete/10000-4      	300000000	        48 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Delete/100000-4     	300000000	        48 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuiltinMap/Delete/1000000-4    	200000000	        60 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashMap/Delete/10-4            	 50000000	       245 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Delete/100-4           	 50000000	       246 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Delete/1000-4          	 50000000	       248 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Delete/10000-4         	 50000000	       254 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Delete/100000-4        	 50000000	       296 ns/op	      23 B/op	       3 allocs/op
BenchmarkHashMap/Delete/1000000-4       	 30000000	       353 ns/op	      24 B/op	       3 allocs/op
```

## License

The MIT License (MIT)

Copyright &copy; 2019 thewizardplusplus
