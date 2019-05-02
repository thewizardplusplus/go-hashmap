package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewConcurrentHashMap(test *testing.T) {
	hashMap := NewConcurrentHashMap()
	assert.Len(test, hashMap.segments, concurrencyLevel)
	for _, segment := range hashMap.segments {
		assert.NotNil(test, segment)
	}
}

func TestConcurrentHashMap(test *testing.T) {
	for _, data := range []struct {
		name        string
		makeHashMap func() ConcurrentHashMap
		makeKey     func() Key
		wantValue   interface{}
		wantOk      assert.BoolAssertionFunc
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			hashMap := data.makeHashMap()
			key := data.makeKey()
			gotValue, gotOk := hashMap.Get(key)

			for _, segment := range hashMap.segments {
				for _, bucket := range segment.innerMap.buckets {
					if bucket != nil {
						mock.AssertExpectationsForObjects(test, bucket.key)
					}
				}
			}
			mock.AssertExpectationsForObjects(test, key)
			assert.Equal(test, gotValue, data.wantValue)
			data.wantOk(test, gotOk)
		})
	}
}
