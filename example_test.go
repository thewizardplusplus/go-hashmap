package hashmap

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

func Example() {
	timeZones := NewSynchronizedHashMap()
	timeZones.Set(StringKey("EST"), -5*60*60)
	timeZones.Set(StringKey("CST"), -6*60*60)
	timeZones.Set(StringKey("MST"), -7*60*60)

	estOffset, ok := timeZones.Get(StringKey("EST"))
	fmt.Println(estOffset, ok)
	// Output:
	// -18000 true
}
