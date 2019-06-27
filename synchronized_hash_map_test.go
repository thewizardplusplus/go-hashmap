package hashmap

import (
	"math/rand"
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
		{
			name:        "getting by a nonexistent key",
			makeHashMap: func() *SynchronizedHashMap { return NewSynchronizedHashMap() },
			makeKey: func() Key {
				key := new(MockKey)
				key.On("Hash").Return(5)

				return key
			},
			wantValue: nil,
			wantOk:    assert.False,
		},
		{
			name: "setting by a nonexistent key",
			makeHashMap: func() *SynchronizedHashMap {
				key := new(MockKey)
				key.On("Hash").Return(5)
				// it's called inside the HashMap.Get() method below
				key.On("Equals", mock.Anything).Return(true)

				hashMap := NewSynchronizedHashMap()
				hashMap.Set(key, "five")

				return hashMap
			},
			makeKey: func() Key {
				key := new(MockKey)
				key.On("Hash").Return(5)

				return key
			},
			wantValue: "five",
			wantOk:    assert.True,
		},
		{
			name: "setting by an existing key",
			makeHashMap: func() *SynchronizedHashMap {
				key := new(MockKey)
				key.On("Hash").Return(5)
				key.On("Equals", mock.Anything).Return(true)

				hashMap := NewSynchronizedHashMap()
				hashMap.Set(key, "five #1")
				hashMap.Set(key, "five #2")

				return hashMap
			},
			makeKey: func() Key {
				key := new(MockKey)
				key.On("Hash").Return(5)

				return key
			},
			wantValue: "five #2",
			wantOk:    assert.True,
		},
		{
			name: "deleting by a nonexistent key",
			makeHashMap: func() *SynchronizedHashMap {
				key := new(MockKey)
				key.On("Hash").Return(5)

				hashMap := NewSynchronizedHashMap()
				hashMap.Delete(key)

				return hashMap
			},
			makeKey: func() Key {
				key := new(MockKey)
				key.On("Hash").Return(5)

				return key
			},
			wantValue: nil,
			wantOk:    assert.False,
		},
		{
			name: "deleting by an existing key",
			makeHashMap: func() *SynchronizedHashMap {
				key := new(MockKey)
				key.On("Hash").Return(5)
				key.On("Equals", mock.Anything).Return(true)

				hashMap := NewSynchronizedHashMap()
				hashMap.Set(key, "five")
				hashMap.Delete(key)

				return hashMap
			},
			makeKey: func() Key {
				key := new(MockKey)
				key.On("Hash").Return(5)

				return key
			},
			wantValue: nil,
			wantOk:    assert.False,
		},
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
			assert.Equal(test, data.wantValue, gotValue)
			data.wantOk(test, gotOk)
		})
	}
}

func TestSynchronizedHashMap_Iterate(test *testing.T) {
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
			innerMap := HashMap{buckets: data.fields.buckets}
			hashMap := SynchronizedHashMap{innerMap: &innerMap}
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
