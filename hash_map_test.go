package hashmap

import (
	"math/rand"
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
		{
			name: "without buckets",
			fields: fields{
				makeBuckets: func() []*bucket { return make([]*bucket, initialCapacity) },
			},
			args: args{
				makeKey: func() Key {
					key := new(MockKey)
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
					fiveKey := new(MockKey)
					fiveKey.On("Equals", mock.Anything).Return(true)

					buckets := make([]*bucket, initialCapacity)
					buckets[5] = &bucket{key: fiveKey, value: "five"}
					buckets[6] = &bucket{key: new(MockKey), value: "six"}
					buckets[7] = &bucket{key: new(MockKey), value: "seven"}

					return buckets
				},
			},
			args: args{
				makeKey: func() Key {
					key := new(MockKey)
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
					fiveKey := new(MockKey)
					fiveKey.On("Equals", mock.Anything).Return(false)

					sixKey := new(MockKey)
					sixKey.On("Equals", mock.Anything).Return(false)

					sevenKey := new(MockKey)
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
					key := new(MockKey)
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
					fiveKey := new(MockKey)
					fiveKey.On("Equals", mock.Anything).Return(false)

					sixKey := new(MockKey)
					sixKey.On("Equals", mock.Anything).Return(false)

					sevenKey := new(MockKey)
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
					key := new(MockKey)
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
			assert.Equal(test, data.wantValue, gotValue)
			data.wantOk(test, gotOk)
		})
	}
}

func TestHashMap_Iterate(test *testing.T) {
	type fields struct {
		buckets []*bucket
	}

	for _, data := range []struct {
		name             string
		fields           fields
		interruptOnCount int
		wantBuckets      []bucket
		wantOk           assert.BoolAssertionFunc
	}{
		{
			name:             "without buckets",
			fields:           fields{buckets: make([]*bucket, initialCapacity)},
			interruptOnCount: 10,
			wantBuckets:      nil,
			wantOk:           assert.True,
		},
		{
			name: "with few buckets and without an interrupt",
			fields: fields{
				buckets: []*bucket{
					5: {key: new(MockKey), value: "five"},
					6: {key: new(MockKey), value: "six"},
					7: {key: new(MockKey), value: "seven"},
				},
			},
			interruptOnCount: 10,
			wantBuckets: []bucket{
				{key: new(MockKey), value: "five"},
				{key: new(MockKey), value: "six"},
				{key: new(MockKey), value: "seven"},
			},
			wantOk: assert.True,
		},
		{
			name: "with few buckets and with an interrupt",
			fields: fields{
				buckets: []*bucket{
					5: {key: new(MockKey), value: "five"},
					6: {key: new(MockKey), value: "six"},
					7: {key: new(MockKey), value: "seven"},
				},
			},
			interruptOnCount: 2,
			wantBuckets: []bucket{
				{key: new(MockKey), value: "five"},
				{key: new(MockKey), value: "six"},
			},
			wantOk: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			// reset the random generator to make tests deterministic
			rand.Seed(1)

			var gotBuckets []bucket
			hashMap := HashMap{buckets: data.fields.buckets}
			gotOk := hashMap.Iterate(func(key Key, value interface{}) bool {
				gotBuckets = append(gotBuckets, bucket{key, value})
				// interrupt after a specified count of got buckets
				return len(gotBuckets) < data.interruptOnCount
			})

			for _, bucket := range data.fields.buckets {
				if bucket != nil {
					mock.AssertExpectationsForObjects(test, bucket.key)
				}
			}
			assert.Equal(test, data.wantBuckets, gotBuckets)
			data.wantOk(test, gotOk)
		})
	}
}

func TestHashMap_Iterate_order(test *testing.T) {
	hashMap := HashMap{
		buckets: []*bucket{
			5: {key: new(MockKey), value: "five"},
			6: {key: new(MockKey), value: "six"},
			7: {key: new(MockKey), value: "seven"},
		},
	}

	var gotBucketsOne []bucket
	rand.Seed(1)
	gotOkOne := hashMap.Iterate(func(key Key, value interface{}) bool {
		gotBucketsOne = append(gotBucketsOne, bucket{key, value})
		return true
	})

	var gotBucketsTwo []bucket
	rand.Seed(2)
	gotOkTwo := hashMap.Iterate(func(key Key, value interface{}) bool {
		gotBucketsTwo = append(gotBucketsTwo, bucket{key, value})
		return true
	})

	for _, bucket := range hashMap.buckets {
		if bucket != nil {
			mock.AssertExpectationsForObjects(test, bucket.key)
		}
	}

	wantBucketsOne := []bucket{
		{key: new(MockKey), value: "five"},
		{key: new(MockKey), value: "six"},
		{key: new(MockKey), value: "seven"},
	}
	assert.Equal(test, wantBucketsOne, gotBucketsOne)
	assert.True(test, gotOkOne)

	wantBucketsTwo := []bucket{
		{key: new(MockKey), value: "six"},
		{key: new(MockKey), value: "seven"},
		{key: new(MockKey), value: "five"},
	}
	assert.Equal(test, wantBucketsTwo, gotBucketsTwo)
	assert.True(test, gotOkTwo)
}

func TestHashMap_Set(test *testing.T) {
	type fields struct {
		makeBuckets func() []*bucket
		size        int
	}
	type args struct {
		makeKey func() Key
	}

	for _, data := range []struct {
		name         string
		fields       fields
		args         args
		wantSize     int
		wantCapacity int
	}{
		{
			name: "without buckets",
			fields: fields{
				makeBuckets: func() []*bucket { return make([]*bucket, initialCapacity) },
				size:        0,
			},
			args: args{
				makeKey: func() Key {
					key := new(MockKey)
					key.On("Hash").Return(5)
					// it's called inside the HashMap.Get() method below
					key.On("Equals", mock.Anything).Return(true)

					return key
				},
			},
			wantSize:     1,
			wantCapacity: initialCapacity,
		},
		{
			name: "with few buckets and a match at the start",
			fields: fields{
				makeBuckets: func() []*bucket {
					fiveKey := new(MockKey)
					fiveKey.On("Equals", mock.Anything).Return(true)

					buckets := make([]*bucket, initialCapacity)
					buckets[5] = &bucket{key: fiveKey, value: "five"}
					buckets[6] = &bucket{key: new(MockKey), value: "six"}
					buckets[7] = &bucket{key: new(MockKey), value: "seven"}

					return buckets
				},
				size: 3,
			},
			args: args{
				makeKey: func() Key {
					key := new(MockKey)
					key.On("Hash").Return(5)

					return key
				},
			},
			wantSize:     3,
			wantCapacity: initialCapacity,
		},
		{
			name: "with few buckets and a match at the end",
			fields: fields{
				makeBuckets: func() []*bucket {
					fiveKey := new(MockKey)
					fiveKey.On("Equals", mock.Anything).Return(false)

					sixKey := new(MockKey)
					sixKey.On("Equals", mock.Anything).Return(false)

					sevenKey := new(MockKey)
					sevenKey.On("Equals", mock.Anything).Return(true)

					buckets := make([]*bucket, initialCapacity)
					buckets[5] = &bucket{key: fiveKey, value: "five"}
					buckets[6] = &bucket{key: sixKey, value: "six"}
					buckets[7] = &bucket{key: sevenKey, value: "seven"}

					return buckets
				},
				size: 3,
			},
			args: args{
				makeKey: func() Key {
					key := new(MockKey)
					key.On("Hash").Return(5)

					return key
				},
			},
			wantSize:     3,
			wantCapacity: initialCapacity,
		},
		{
			name: "with few buckets and no match",
			fields: fields{
				makeBuckets: func() []*bucket {
					fiveKey := new(MockKey)
					fiveKey.On("Equals", mock.Anything).Return(false)

					sixKey := new(MockKey)
					sixKey.On("Equals", mock.Anything).Return(false)

					sevenKey := new(MockKey)
					sevenKey.On("Equals", mock.Anything).Return(false)

					buckets := make([]*bucket, initialCapacity)
					buckets[5] = &bucket{key: fiveKey, value: "five"}
					buckets[6] = &bucket{key: sixKey, value: "six"}
					buckets[7] = &bucket{key: sevenKey, value: "seven"}

					return buckets
				},
				size: 3,
			},
			args: args{
				makeKey: func() Key {
					key := new(MockKey)
					key.On("Hash").Return(5)
					// it's called inside the HashMap.Get() method below
					key.On("Equals", mock.Anything).Return(true)

					return key
				},
			},
			wantSize:     4,
			wantCapacity: initialCapacity,
		},
		{
			name: "with a load factor over the maximum and a match",
			fields: fields{
				makeBuckets: func() []*bucket {
					threeKey := new(MockKey)
					threeKey.On("Equals", mock.Anything).Return(true)

					buckets := make([]*bucket, 5)
					buckets[0] = &bucket{key: new(MockKey), value: "zero"}
					buckets[1] = &bucket{key: new(MockKey), value: "one"}
					buckets[2] = &bucket{key: new(MockKey), value: "two"}
					buckets[3] = &bucket{key: threeKey, value: "three"}

					return buckets
				},
				size: 4,
			},
			args: args{
				makeKey: func() Key {
					key := new(MockKey)
					key.On("Hash").Return(3)

					return key
				},
			},
			wantSize:     4,
			wantCapacity: 5,
		},
		{
			name: "with a load factor over the maximum and no match",
			fields: fields{
				makeBuckets: func() []*bucket {
					zeroKey := new(MockKey)
					zeroKey.On("Hash").Return(0)

					oneKey := new(MockKey)
					oneKey.On("Hash").Return(1)

					twoKey := new(MockKey)
					twoKey.On("Hash").Return(2)

					threeKey := new(MockKey)
					threeKey.On("Hash").Return(3)
					threeKey.On("Equals", mock.Anything).Return(false)

					buckets := make([]*bucket, 5)
					buckets[0] = &bucket{key: zeroKey, value: "zero"}
					buckets[1] = &bucket{key: oneKey, value: "one"}
					buckets[2] = &bucket{key: twoKey, value: "two"}
					buckets[3] = &bucket{key: threeKey, value: "three"}

					return buckets
				},
				size: 4,
			},
			args: args{
				makeKey: func() Key {
					key := new(MockKey)
					key.On("Hash").Return(3)
					// it's called inside the HashMap.Get() method below
					key.On("Equals", mock.Anything).Return(true)

					return key
				},
			},
			wantSize:     5,
			wantCapacity: 10,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			buckets := data.fields.makeBuckets()
			key := data.args.makeKey()
			value := rand.Int()

			hashMap := HashMap{buckets: buckets, size: data.fields.size}
			hashMap.Set(key, value)

			gotValue, gotOk := hashMap.Get(key)

			for _, bucket := range buckets {
				if bucket != nil {
					mock.AssertExpectationsForObjects(test, bucket.key)
				}
			}
			mock.AssertExpectationsForObjects(test, key)
			assert.Equal(test, data.wantSize, hashMap.size)
			assert.Equal(test, data.wantCapacity, len(hashMap.buckets))
			assert.Equal(test, value, gotValue)
			assert.True(test, gotOk)
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
	}{
		{
			name: "without buckets",
			fields: fields{
				makeBuckets: func() []*bucket { return make([]*bucket, initialCapacity) },
				size:        0,
			},
			args: args{
				makeKey: func() Key {
					key := new(MockKey)
					key.On("Hash").Return(5)

					return key
				},
			},
			wantSize: 0,
		},
		{
			name: "with few buckets and a match at the start",
			fields: fields{
				makeBuckets: func() []*bucket {
					fiveKey := new(MockKey)
					fiveKey.On("Equals", mock.Anything).Return(true)

					buckets := make([]*bucket, initialCapacity)
					buckets[5] = &bucket{key: fiveKey, value: "five"}
					buckets[6] = &bucket{key: new(MockKey), value: "six"}
					buckets[7] = &bucket{key: new(MockKey), value: "seven"}

					return buckets
				},
				size: 3,
			},
			args: args{
				makeKey: func() Key {
					key := new(MockKey)
					key.On("Hash").Return(5)

					return key
				},
			},
			wantSize: 2,
		},
		{
			name: "with few buckets and a match at the end",
			fields: fields{
				makeBuckets: func() []*bucket {
					fiveKey := new(MockKey)
					fiveKey.On("Equals", mock.Anything).Return(false)

					sixKey := new(MockKey)
					sixKey.On("Equals", mock.Anything).Return(false)

					sevenKey := new(MockKey)
					sevenKey.On("Equals", mock.Anything).Return(true)

					buckets := make([]*bucket, initialCapacity)
					buckets[5] = &bucket{key: fiveKey, value: "five"}
					buckets[6] = &bucket{key: sixKey, value: "six"}
					buckets[7] = &bucket{key: sevenKey, value: "seven"}

					return buckets
				},
				size: 3,
			},
			args: args{
				makeKey: func() Key {
					key := new(MockKey)
					key.On("Hash").Return(5)

					return key
				},
			},
			wantSize: 2,
		},
		{
			name: "with few buckets and no match",
			fields: fields{
				makeBuckets: func() []*bucket {
					fiveKey := new(MockKey)
					fiveKey.On("Equals", mock.Anything).Return(false)

					sixKey := new(MockKey)
					sixKey.On("Equals", mock.Anything).Return(false)

					sevenKey := new(MockKey)
					sevenKey.On("Equals", mock.Anything).Return(false)

					buckets := make([]*bucket, initialCapacity)
					buckets[5] = &bucket{key: fiveKey, value: "five"}
					buckets[6] = &bucket{key: sixKey, value: "six"}
					buckets[7] = &bucket{key: sevenKey, value: "seven"}

					return buckets
				},
				size: 3,
			},
			args: args{
				makeKey: func() Key {
					key := new(MockKey)
					key.On("Hash").Return(5)

					return key
				},
			},
			wantSize: 3,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			buckets := data.fields.makeBuckets()
			key := data.args.makeKey()

			hashMap := HashMap{buckets: buckets, size: data.fields.size}
			hashMap.Delete(key)

			_, gotGetOk := hashMap.Get(key)

			for _, bucket := range buckets {
				if bucket != nil {
					mock.AssertExpectationsForObjects(test, bucket.key)
				}
			}
			mock.AssertExpectationsForObjects(test, key)
			assert.Equal(test, data.wantSize, hashMap.size)
			assert.False(test, gotGetOk)
		})
	}
}
