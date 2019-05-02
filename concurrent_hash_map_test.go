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
		{
			name:        "getting by a nonexistent key",
			makeHashMap: func() ConcurrentHashMap { return NewConcurrentHashMap() },
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
			makeHashMap: func() ConcurrentHashMap {
				key := new(MockKey)
				key.On("Hash").Return(5)
				// it's called inside the HashMap.Get() method below
				key.On("Equals", mock.Anything).Return(true)

				hashMap := NewConcurrentHashMap()
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
			makeHashMap: func() ConcurrentHashMap {
				key := new(MockKey)
				key.On("Hash").Return(5)
				key.On("Equals", mock.Anything).Return(true)

				hashMap := NewConcurrentHashMap()
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
			makeHashMap: func() ConcurrentHashMap {
				key := new(MockKey)
				key.On("Hash").Return(5)

				hashMap := NewConcurrentHashMap()
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
			makeHashMap: func() ConcurrentHashMap {
				key := new(MockKey)
				key.On("Hash").Return(5)
				key.On("Equals", mock.Anything).Return(true)

				hashMap := NewConcurrentHashMap()
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
