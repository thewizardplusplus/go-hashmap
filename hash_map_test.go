package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-hashmap/mocks"
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
		{
			name: "without buckets",
			fields: fields{
				makeBuckets: func() []*bucket { return make([]*bucket, initialCapacity) },
			},
			args: args{
				makeKey: func() Key {
					key := new(mocks.Key)
					key.On("Hash").Return(5)

					return key
				},
			},
			wantValue: nil,
			wantOk:    assert.False,
		},
		{
			name: "with few buckets and a match at the start",
			fields: fields{
				makeBuckets: func() []*bucket {
					fiveKey := new(mocks.Key)
					fiveKey.On("Equals", mock.Anything).Return(true)

					buckets := make([]*bucket, initialCapacity)
					buckets[5] = &bucket{key: fiveKey, value: "five"}
					buckets[6] = &bucket{key: new(mocks.Key), value: "six"}
					buckets[7] = &bucket{key: new(mocks.Key), value: "seven"}

					return buckets
				},
			},
			args: args{
				makeKey: func() Key {
					key := new(mocks.Key)
					key.On("Hash").Return(5)

					return key
				},
			},
			wantValue: "five",
			wantOk:    assert.True,
		},
		{
			name: "with few buckets and a match at the end",
			fields: fields{
				makeBuckets: func() []*bucket {
					fiveKey := new(mocks.Key)
					fiveKey.On("Equals", mock.Anything).Return(false)

					sixKey := new(mocks.Key)
					sixKey.On("Equals", mock.Anything).Return(false)

					sevenKey := new(mocks.Key)
					sevenKey.On("Equals", mock.Anything).Return(true)

					buckets := make([]*bucket, initialCapacity)
					buckets[5] = &bucket{key: fiveKey, value: "five"}
					buckets[6] = &bucket{key: sixKey, value: "six"}
					buckets[7] = &bucket{key: sevenKey, value: "seven"}

					return buckets
				},
			},
			args: args{
				makeKey: func() Key {
					key := new(mocks.Key)
					key.On("Hash").Return(5)

					return key
				},
			},
			wantValue: "seven",
			wantOk:    assert.True,
		},
		{
			name: "with few buckets and no match",
			fields: fields{
				makeBuckets: func() []*bucket {
					fiveKey := new(mocks.Key)
					fiveKey.On("Equals", mock.Anything).Return(false)

					sixKey := new(mocks.Key)
					sixKey.On("Equals", mock.Anything).Return(false)

					sevenKey := new(mocks.Key)
					sevenKey.On("Equals", mock.Anything).Return(false)

					buckets := make([]*bucket, initialCapacity)
					buckets[5] = &bucket{key: fiveKey, value: "five"}
					buckets[6] = &bucket{key: sixKey, value: "six"}
					buckets[7] = &bucket{key: sevenKey, value: "seven"}

					return buckets
				},
			},
			args: args{
				makeKey: func() Key {
					key := new(mocks.Key)
					key.On("Hash").Return(5)

					return key
				},
			},
			wantValue: nil,
			wantOk:    assert.False,
		},
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

func TestHashMap_Delete(test *testing.T) {
	type fields struct {
		makeBuckets func() []*bucket
		size        int
	}
	type args struct {
		makeKey func() Key
	}

	for _, data := range []struct {
		name     string
		fields   fields
		args     args
		wantSize int
		wantOk   assert.BoolAssertionFunc
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			buckets := data.fields.makeBuckets()
			key := data.args.makeKey()

			hashMap := HashMap{buckets: buckets, size: data.fields.size}
			gotDeleteOk := hashMap.Delete(key)
			_, gotGetOk := hashMap.Get(key)

			for _, bucket := range buckets {
				if bucket != nil {
					mock.AssertExpectationsForObjects(test, bucket.key)
				}
			}
			mock.AssertExpectationsForObjects(test, key)
			assert.Equal(test, data.wantSize, hashMap.size)
			data.wantOk(test, gotDeleteOk)
			assert.False(test, gotGetOk)
		})
	}
}
