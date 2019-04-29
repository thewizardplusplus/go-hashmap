package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHashMap(test *testing.T) {
	hashMap := NewHashMap()
	require.NotNil(test, hashMap)
	assert.Len(test, hashMap.buckets, initialCapacity)
	assert.Zero(test, hashMap.size)
}
