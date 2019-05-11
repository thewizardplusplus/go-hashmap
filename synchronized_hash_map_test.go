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
		name        string
		fields      fields
		wantBuckets []bucket
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			var gotBuckets []bucket
			innerMap := HashMap{buckets: data.fields.buckets}
			hashMap := SynchronizedHashMap{innerMap: &innerMap}
			hashMap.Iterate(func(key Key, value interface{}) {
				gotBuckets = append(gotBuckets, bucket{key, value})
			})

			for _, bucket := range hashMap.innerMap.buckets {
				if bucket != nil {
					mock.AssertExpectationsForObjects(test, bucket.key)
				}
			}
			assert.ElementsMatch(test, data.wantBuckets, gotBuckets)
		})
	}
}
