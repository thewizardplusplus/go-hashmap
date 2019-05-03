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
	type result struct {
		value interface{}
		ok    bool
	}

	for _, data := range []struct {
		name                string
		makeHashMap         func() ConcurrentHashMap
		makeKeys            func() []Key
		wantTouchedSegments map[int]struct{}
		wantResults         []result
	}{
		{
			name:        "getting by a nonexistent key",
			makeHashMap: func() ConcurrentHashMap { return NewConcurrentHashMap() },
			makeKeys: func() []Key {
				key := new(MockKey)
				key.On("Hash").Return(5)

				return []Key{key}
			},
			wantTouchedSegments: nil,
			wantResults:         []result{{nil, false}},
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
			makeKeys: func() []Key {
				key := new(MockKey)
				key.On("Hash").Return(5)

				return []Key{key}
			},
			wantTouchedSegments: map[int]struct{}{5: struct{}{}},
			wantResults:         []result{{"five", true}},
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
			makeKeys: func() []Key {
				key := new(MockKey)
				key.On("Hash").Return(5)

				return []Key{key}
			},
			wantTouchedSegments: map[int]struct{}{5: struct{}{}},
			wantResults:         []result{{"five #2", true}},
		},
		{
			name: "setting by keys touched different segments",
			makeHashMap: func() ConcurrentHashMap {
				fiveKey := new(MockKey)
				fiveKey.On("Hash").Return(5)
				fiveKey.On("Equals", mock.Anything).Return(true)

				sixKey := new(MockKey)
				sixKey.On("Hash").Return(6)
				sixKey.On("Equals", mock.Anything).Return(true)

				hashMap := NewConcurrentHashMap()
				hashMap.Set(fiveKey, "five")
				hashMap.Set(sixKey, "six")

				return hashMap
			},
			makeKeys: func() []Key {
				fiveKey := new(MockKey)
				fiveKey.On("Hash").Return(5)

				sixKey := new(MockKey)
				sixKey.On("Hash").Return(6)

				return []Key{fiveKey, sixKey}
			},
			wantTouchedSegments: map[int]struct{}{5: struct{}{}, 6: struct{}{}},
			wantResults:         []result{{"five", true}, {"six", true}},
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
			makeKeys: func() []Key {
				key := new(MockKey)
				key.On("Hash").Return(5)

				return []Key{key}
			},
			wantTouchedSegments: nil,
			wantResults:         []result{{nil, false}},
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
			makeKeys: func() []Key {
				key := new(MockKey)
				key.On("Hash").Return(5)

				return []Key{key}
			},
			wantTouchedSegments: nil,
			wantResults:         []result{{nil, false}},
		},
		{
			name: "deleting by keys touched different segments",
			makeHashMap: func() ConcurrentHashMap {
				fiveKey := new(MockKey)
				fiveKey.On("Hash").Return(5)
				fiveKey.On("Equals", mock.Anything).Return(true)

				sixKey := new(MockKey)
				sixKey.On("Hash").Return(6)
				sixKey.On("Equals", mock.Anything).Return(true)

				hashMap := NewConcurrentHashMap()
				hashMap.Set(fiveKey, "five")
				hashMap.Set(sixKey, "six")
				hashMap.Delete(fiveKey)
				hashMap.Delete(sixKey)

				return hashMap
			},
			makeKeys: func() []Key {
				fiveKey := new(MockKey)
				fiveKey.On("Hash").Return(5)

				sixKey := new(MockKey)
				sixKey.On("Hash").Return(6)

				return []Key{fiveKey, sixKey}
			},
			wantTouchedSegments: nil,
			wantResults:         []result{{nil, false}, {nil, false}},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			hashMap := data.makeHashMap()
			keys := data.makeKeys()

			var results []result
			for _, key := range keys {
				gotValue, gotOk := hashMap.Get(key)
				results = append(results, result{gotValue, gotOk})
			}

			for index, segment := range hashMap.segments {
				for _, bucket := range segment.innerMap.buckets {
					if bucket != nil {
						mock.AssertExpectationsForObjects(test, bucket.key)
					}
				}
				if _, ok := data.wantTouchedSegments[index]; ok {
					assert.NotZero(test, segment.innerMap.size)
				} else {
					assert.Zero(test, segment.innerMap.size)
				}
			}
			for _, key := range keys {
				mock.AssertExpectationsForObjects(test, key)
			}
			assert.Equal(test, data.wantResults, results)
		})
	}
}
