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
	timeZones := NewHashMap()
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
