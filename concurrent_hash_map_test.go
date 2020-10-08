package hashmap

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewConcurrentHashMap(test *testing.T) {
	hashMap := NewConcurrentHashMap()
	assert.Len(test, hashMap.segments, defaultConcurrentConfig.concurrencyLevel)
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
			wantTouchedSegments: map[int]struct{}{5: {}},
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
			wantTouchedSegments: map[int]struct{}{5: {}},
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
			wantTouchedSegments: map[int]struct{}{5: {}, 6: {}},
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
				innerMap := segment.(*SynchronizedHashMap).innerMap.(*HashMap)
				for _, bucket := range innerMap.buckets {
					if bucket != nil {
						mock.AssertExpectationsForObjects(test, bucket.key)
					}
				}
				if _, ok := data.wantTouchedSegments[index]; ok {
					assert.NotZero(test, innerMap.size)
				} else {
					assert.Zero(test, innerMap.size)
				}
			}
			for _, key := range keys {
				mock.AssertExpectationsForObjects(test, key)
			}
			assert.Equal(test, data.wantResults, results)
		})
	}
}

func TestConcurrentHashMap_Iterate(test *testing.T) {
	type fields struct {
		buckets [][]*bucket
	}

	for _, data := range []struct {
		name             string
		fields           fields
		interruptOnCount int
		wantBuckets      []bucket
		wantOk           assert.BoolAssertionFunc
	}{
		{
			name: "without buckets",
			fields: fields{
				buckets: [][]*bucket{
					make([]*bucket, defaultConfig.initialCapacity),
					make([]*bucket, defaultConfig.initialCapacity),
				},
			},
			interruptOnCount: 10,
			wantBuckets:      nil,
			wantOk:           assert.True,
		},
		{
			name: "with few buckets in the same segment and without an interrupt",
			fields: fields{
				buckets: [][]*bucket{
					5: {
						5: {key: new(MockKey), value: "five #1"},
						6: {key: new(MockKey), value: "five #2"},
						7: {key: new(MockKey), value: "five #3"},
					},
				},
			},
			interruptOnCount: 10,
			wantBuckets: []bucket{
				{key: new(MockKey), value: "five #3"},
				{key: new(MockKey), value: "five #2"},
				{key: new(MockKey), value: "five #1"},
			},
			wantOk: assert.True,
		},
		{
			name: "with few buckets in the same segment and with an interrupt",
			fields: fields{
				buckets: [][]*bucket{
					5: {
						5: {key: new(MockKey), value: "five #1"},
						6: {key: new(MockKey), value: "five #2"},
						7: {key: new(MockKey), value: "five #3"},
					},
				},
			},
			interruptOnCount: 2,
			wantBuckets: []bucket{
				{key: new(MockKey), value: "five #3"},
				{key: new(MockKey), value: "five #2"},
			},
			wantOk: assert.False,
		},
		{
			name: "with few buckets in different segments and without an interrupt",
			fields: fields{
				buckets: [][]*bucket{
					5: {5: {key: new(MockKey), value: "five"}},
					6: {6: {key: new(MockKey), value: "six"}},
					7: {7: {key: new(MockKey), value: "seven"}},
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
			name: "with few buckets in different segments and with an interrupt",
			fields: fields{
				buckets: [][]*bucket{
					5: {5: {key: new(MockKey), value: "five"}},
					6: {6: {key: new(MockKey), value: "six"}},
					7: {7: {key: new(MockKey), value: "seven"}},
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

			var segments []Storage
			for _, buckets := range data.fields.buckets {
				innerMap := &HashMap{buckets: buckets}
				segment := &SynchronizedHashMap{innerMap: innerMap}
				segments = append(segments, segment)
			}

			var gotBuckets []bucket
			hashMap := ConcurrentHashMap{segments: segments}
			gotOk := hashMap.Iterate(func(key Key, value interface{}) bool {
				gotBuckets = append(gotBuckets, bucket{key, value})
				// interrupt after a specified count of got buckets
				return len(gotBuckets) < data.interruptOnCount
			})

			for _, buckets := range data.fields.buckets {
				for _, bucket := range buckets {
					if bucket != nil {
						mock.AssertExpectationsForObjects(test, bucket.key)
					}
				}
			}
			assert.Equal(test, data.wantBuckets, gotBuckets)
			data.wantOk(test, gotOk)
		})
	}
}

func TestConcurrentHashMap_Iterate_order(test *testing.T) {
	type fields struct {
		buckets [][]*bucket
	}

	for _, data := range []struct {
		name           string
		fields         fields
		randomSeedOne  int64
		randomSeedTwo  int64
		wantBucketsOne []bucket
		wantBucketsTwo []bucket
	}{
		{
			name: "with few buckets in the same segment",
			fields: fields{
				buckets: [][]*bucket{
					5: {
						5: {key: new(MockKey), value: "five #1"},
						6: {key: new(MockKey), value: "five #2"},
						7: {key: new(MockKey), value: "five #3"},
					},
				},
			},
			randomSeedOne: 1,
			randomSeedTwo: 2,
			wantBucketsOne: []bucket{
				{key: new(MockKey), value: "five #3"},
				{key: new(MockKey), value: "five #2"},
				{key: new(MockKey), value: "five #1"},
			},
			wantBucketsTwo: []bucket{
				{key: new(MockKey), value: "five #3"},
				{key: new(MockKey), value: "five #2"},
				{key: new(MockKey), value: "five #1"},
			},
		},
		{
			name: "with few buckets in different segments",
			fields: fields{
				buckets: [][]*bucket{
					5: {5: {key: new(MockKey), value: "five"}},
					6: {6: {key: new(MockKey), value: "six"}},
					7: {7: {key: new(MockKey), value: "seven"}},
				},
			},
			randomSeedOne: 1,
			randomSeedTwo: 2,
			wantBucketsOne: []bucket{
				{key: new(MockKey), value: "five"},
				{key: new(MockKey), value: "six"},
				{key: new(MockKey), value: "seven"},
			},
			wantBucketsTwo: []bucket{
				{key: new(MockKey), value: "six"},
				{key: new(MockKey), value: "seven"},
				{key: new(MockKey), value: "five"},
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			var segments []Storage
			for _, buckets := range data.fields.buckets {
				innerMap := &HashMap{buckets: buckets}
				segment := &SynchronizedHashMap{innerMap: innerMap}
				segments = append(segments, segment)
			}

			hashMap := ConcurrentHashMap{segments: segments}

			var gotBucketsOne []bucket
			rand.Seed(data.randomSeedOne)
			gotOkOne := hashMap.Iterate(func(key Key, value interface{}) bool {
				gotBucketsOne = append(gotBucketsOne, bucket{key, value})
				return true
			})

			var gotBucketsTwo []bucket
			rand.Seed(data.randomSeedTwo)
			gotOkTwo := hashMap.Iterate(func(key Key, value interface{}) bool {
				gotBucketsTwo = append(gotBucketsTwo, bucket{key, value})
				return true
			})

			for _, buckets := range data.fields.buckets {
				for _, bucket := range buckets {
					if bucket != nil {
						mock.AssertExpectationsForObjects(test, bucket.key)
					}
				}
			}

			assert.Equal(test, data.wantBucketsOne, gotBucketsOne)
			assert.True(test, gotOkOne)

			assert.Equal(test, data.wantBucketsTwo, gotBucketsTwo)
			assert.True(test, gotOkTwo)
		})
	}
}
