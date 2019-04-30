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

## License

The MIT License (MIT)

Copyright &copy; 2019 thewizardplusplus
