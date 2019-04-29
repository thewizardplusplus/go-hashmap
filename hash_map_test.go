package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewHashMap(test *testing.T) {
	hashMap := NewHashMap()
	require.NotNil(test, hashMap)
	assert.Len(test, hashMap.buckets, initialCapacity)
	assert.Zero(test, hashMap.size)
}

func TestHashMap_Get(test *testing.T) {
	type fields struct {
		makeBuckets func() []*bucket
	}
	type args struct {
		makeKey func() Key
	}

	for _, data := range []struct {
		name      string
		fields    fields
		args      args
		wantValue interface{}
		wantOk    assert.BoolAssertionFunc
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			buckets := data.fields.makeBuckets()
			key := data.args.makeKey()
			gotValue, gotOk := HashMap{buckets: buckets}.Get(key)

			for _, bucket := range buckets {
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
