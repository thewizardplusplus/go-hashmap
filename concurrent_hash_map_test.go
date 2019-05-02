package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConcurrentHashMap(test *testing.T) {
	hashMap := NewConcurrentHashMap()
	assert.Len(test, hashMap.segments, concurrencyLevel)
	for _, segment := range hashMap.segments {
		assert.NotNil(test, segment)
	}
}
