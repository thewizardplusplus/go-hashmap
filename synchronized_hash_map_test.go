package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSynchronizedHashMap(test *testing.T) {
	for _, data := range []struct {
		name        string
		makeHashMap func() *SynchronizedHashMap
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

			for _, bucket := range hashMap.innerMap.buckets {
				if bucket != nil {
					mock.AssertExpectationsForObjects(test, bucket.key)
				}
			}
			mock.AssertExpectationsForObjects(test, key)
			assert.Equal(test, gotValue, data.wantValue)
			data.wantOk(test, gotOk)
		})
	}
}
