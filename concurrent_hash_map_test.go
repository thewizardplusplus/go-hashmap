package hashmap

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewConcurrentHashMap(test *testing.T) {
	type args struct {
		options func(test *testing.T) []ConcurrentOption
	}

	for _, data := range []struct {
		name string
		args args
		want func(test *testing.T) ConcurrentHashMap
	}{
		{
			name: "with the default config",
			args: args{
				options: func(test *testing.T) []ConcurrentOption { return nil },
			},
			want: func(test *testing.T) ConcurrentHashMap {
				return ConcurrentHashMap{
					config: func() ConcurrentConfig {
						config := defaultConcurrentConfig
						config.segmentFactory = nil
						config.permuter = nil

						return config
					}(),
					segments: func() []Storage {
						var segments []Storage
						for i := 0; i < defaultConcurrentConfig.concurrencyLevel; i++ {
							segments = append(segments, &SynchronizedHashMap{
								innerMap: &HashMap[Key, any]{
									config: func() Config {
										config := defaultConfig
										config.permuter = nil

										return config
									}(),
									buckets: make([]*bucket[Key, any], defaultConfig.initialCapacity),
									size:    0,
								},
							})
						}

						return segments
					}(),
				}
			},
		},
		{
			name: "with the set concurrency level",
			args: args{
				options: func(test *testing.T) []ConcurrentOption {
					return []ConcurrentOption{WithConcurrencyLevel(23)}
				},
			},
			want: func(test *testing.T) ConcurrentHashMap {
				return ConcurrentHashMap{
					config: func() ConcurrentConfig {
						config := defaultConcurrentConfig
						config.concurrencyLevel = 23
						config.segmentFactory = nil
						config.permuter = nil

						return config
					}(),
					segments: func() []Storage {
						var segments []Storage
						for i := 0; i < 23; i++ {
							segments = append(segments, &SynchronizedHashMap{
								innerMap: &HashMap[Key, any]{
									config: func() Config {
										config := defaultConfig
										config.permuter = nil

										return config
									}(),
									buckets: make([]*bucket[Key, any], defaultConfig.initialCapacity),
									size:    0,
								},
							})
						}

						return segments
					}(),
				}
			},
		},
		{
			name: "with the set segment factory",
			args: args{
				options: func(test *testing.T) []ConcurrentOption {
					return []ConcurrentOption{
						WithSegmentFactory(func() Storage { return NewMockStorage(test) }),
					}
				},
			},
			want: func(test *testing.T) ConcurrentHashMap {
				return ConcurrentHashMap{
					config: func() ConcurrentConfig {
						config := defaultConcurrentConfig
						config.segmentFactory = nil
						config.permuter = nil

						return config
					}(),
					segments: func() []Storage {
						var segments []Storage
						for i := 0; i < defaultConcurrentConfig.concurrencyLevel; i++ {
							segments = append(segments, NewMockStorage(test))
						}

						return segments
					}(),
				}
			},
		},
		{
			name: "with the set config",
			args: args{
				options: func(test *testing.T) []ConcurrentOption {
					return []ConcurrentOption{
						WithConcurrencyLevel(23),
						WithSegmentFactory(func() Storage { return NewMockStorage(test) }),
						WithConcurrentPermuter(rand.New(rand.NewPCG(0, 0)).Perm),
					}
				},
			},
			want: func(test *testing.T) ConcurrentHashMap {
				return ConcurrentHashMap{
					config: func() ConcurrentConfig {
						config := defaultConcurrentConfig
						config.concurrencyLevel = 23
						config.segmentFactory = nil
						config.permuter = nil

						return config
					}(),
					segments: func() []Storage {
						var segments []Storage
						for i := 0; i < 23; i++ {
							segments = append(segments, NewMockStorage(test))
						}

						return segments
					}(),
				}
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewConcurrentHashMap(data.args.options(test)...)

			for _, segment := range got.segments {
				if typedSegment, ok := segment.(*SynchronizedHashMap); ok {
					if innerMap, ok := typedSegment.innerMap.(*HashMap[Key, any]); ok {
						innerMap.config.permuter = nil
					}
				}
			}

			assert.NotNil(test, got.config.segmentFactory)
			got.config.segmentFactory = nil

			assert.NotNil(test, got.config.permuter)
			got.config.permuter = nil

			assert.Equal(test, data.want(test), got)
		})
	}
}

func TestConcurrentHashMap(test *testing.T) {
	type result struct {
		value interface{}
		ok    bool
	}

	for _, data := range []struct {
		name                string
		hashMap             func(test *testing.T) ConcurrentHashMap
		keys                func(test *testing.T) []Key
		wantTouchedSegments map[int]struct{}
		wantResults         []result
	}{
		{
			name: "getting by a nonexistent key",
			hashMap: func(test *testing.T) ConcurrentHashMap {
				return NewConcurrentHashMap()
			},
			keys: func(test *testing.T) []Key {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)

				return []Key{key}
			},
			wantTouchedSegments: nil,
			wantResults:         []result{{nil, false}},
		},
		{
			name: "setting by a nonexistent key",
			hashMap: func(test *testing.T) ConcurrentHashMap {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)
				// it's called inside the HashMap.Get() method below
				key.EXPECT().Equals(mock.Anything).Return(true)

				hashMap := NewConcurrentHashMap()
				hashMap.Set(key, "five")

				return hashMap
			},
			keys: func(test *testing.T) []Key {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)

				return []Key{key}
			},
			wantTouchedSegments: map[int]struct{}{5: {}},
			wantResults:         []result{{"five", true}},
		},
		{
			name: "setting by an existing key",
			hashMap: func(test *testing.T) ConcurrentHashMap {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)
				key.EXPECT().Equals(mock.Anything).Return(true)

				hashMap := NewConcurrentHashMap()
				hashMap.Set(key, "five #1")
				hashMap.Set(key, "five #2")

				return hashMap
			},
			keys: func(test *testing.T) []Key {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)

				return []Key{key}
			},
			wantTouchedSegments: map[int]struct{}{5: {}},
			wantResults:         []result{{"five #2", true}},
		},
		{
			name: "setting by keys touched different segments",
			hashMap: func(test *testing.T) ConcurrentHashMap {
				fiveKey := NewMockKey(test)
				fiveKey.EXPECT().Hash().Return(5)
				fiveKey.EXPECT().Equals(mock.Anything).Return(true)

				sixKey := NewMockKey(test)
				sixKey.EXPECT().Hash().Return(6)
				sixKey.EXPECT().Equals(mock.Anything).Return(true)

				hashMap := NewConcurrentHashMap()
				hashMap.Set(fiveKey, "five")
				hashMap.Set(sixKey, "six")

				return hashMap
			},
			keys: func(test *testing.T) []Key {
				fiveKey := NewMockKey(test)
				fiveKey.EXPECT().Hash().Return(5)

				sixKey := NewMockKey(test)
				sixKey.EXPECT().Hash().Return(6)

				return []Key{fiveKey, sixKey}
			},
			wantTouchedSegments: map[int]struct{}{5: {}, 6: {}},
			wantResults:         []result{{"five", true}, {"six", true}},
		},
		{
			name: "deleting by a nonexistent key",
			hashMap: func(test *testing.T) ConcurrentHashMap {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)

				hashMap := NewConcurrentHashMap()
				hashMap.Delete(key)

				return hashMap
			},
			keys: func(test *testing.T) []Key {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)

				return []Key{key}
			},
			wantTouchedSegments: nil,
			wantResults:         []result{{nil, false}},
		},
		{
			name: "deleting by an existing key",
			hashMap: func(test *testing.T) ConcurrentHashMap {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)
				key.EXPECT().Equals(mock.Anything).Return(true)

				hashMap := NewConcurrentHashMap()
				hashMap.Set(key, "five")
				hashMap.Delete(key)

				return hashMap
			},
			keys: func(test *testing.T) []Key {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)

				return []Key{key}
			},
			wantTouchedSegments: nil,
			wantResults:         []result{{nil, false}},
		},
		{
			name: "deleting by keys touched different segments",
			hashMap: func(test *testing.T) ConcurrentHashMap {
				fiveKey := NewMockKey(test)
				fiveKey.EXPECT().Hash().Return(5)
				fiveKey.EXPECT().Equals(mock.Anything).Return(true)

				sixKey := NewMockKey(test)
				sixKey.EXPECT().Hash().Return(6)
				sixKey.EXPECT().Equals(mock.Anything).Return(true)

				hashMap := NewConcurrentHashMap()
				hashMap.Set(fiveKey, "five")
				hashMap.Set(sixKey, "six")
				hashMap.Delete(fiveKey)
				hashMap.Delete(sixKey)

				return hashMap
			},
			keys: func(test *testing.T) []Key {
				fiveKey := NewMockKey(test)
				fiveKey.EXPECT().Hash().Return(5)

				sixKey := NewMockKey(test)
				sixKey.EXPECT().Hash().Return(6)

				return []Key{fiveKey, sixKey}
			},
			wantTouchedSegments: nil,
			wantResults:         []result{{nil, false}, {nil, false}},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			hashMap := data.hashMap(test)

			var results []result
			for _, key := range data.keys(test) {
				gotValue, gotOk := hashMap.Get(key)
				results = append(results, result{gotValue, gotOk})
			}

			for index, segment := range hashMap.segments {
				innerMap := segment.(*SynchronizedHashMap).innerMap.(*HashMap[Key, any])
				if _, ok := data.wantTouchedSegments[index]; ok {
					assert.NotZero(test, innerMap.size)
				} else {
					assert.Zero(test, innerMap.size)
				}
			}
			assert.Equal(test, data.wantResults, results)
		})
	}
}

func TestConcurrentHashMap_All(test *testing.T) {
	type fields struct {
		buckets func(test *testing.T) [][]*bucket[Key, any]
	}

	for _, data := range []struct {
		name             string
		fields           fields
		interruptOnCount int
		wantBuckets      func(test *testing.T) []bucket[Key, any]
	}{
		{
			name: "without buckets",
			fields: fields{
				buckets: func(test *testing.T) [][]*bucket[Key, any] {
					return [][]*bucket[Key, any]{
						make([]*bucket[Key, any], defaultConfig.initialCapacity),
						make([]*bucket[Key, any], defaultConfig.initialCapacity),
					}
				},
			},
			interruptOnCount: 10,
			wantBuckets:      func(test *testing.T) []bucket[Key, any] { return nil },
		},
		{
			name: "with few buckets in the same segment and without an interrupt",
			fields: fields{
				buckets: func(test *testing.T) [][]*bucket[Key, any] {
					return [][]*bucket[Key, any]{
						5: {
							5: {key: NewMockKey(test), value: "five #1"},
							6: {key: NewMockKey(test), value: "five #2"},
							7: {key: NewMockKey(test), value: "five #3"},
						},
					}
				},
			},
			interruptOnCount: 10,
			wantBuckets: func(test *testing.T) []bucket[Key, any] {
				return []bucket[Key, any]{
					{key: NewMockKey(test), value: "five #2"},
					{key: NewMockKey(test), value: "five #1"},
					{key: NewMockKey(test), value: "five #3"},
				}
			},
		},
		{
			name: "with few buckets in the same segment and with an interrupt",
			fields: fields{
				buckets: func(test *testing.T) [][]*bucket[Key, any] {
					return [][]*bucket[Key, any]{
						5: {
							5: {key: NewMockKey(test), value: "five #1"},
							6: {key: NewMockKey(test), value: "five #2"},
							7: {key: NewMockKey(test), value: "five #3"},
						},
					}
				},
			},
			interruptOnCount: 2,
			wantBuckets: func(test *testing.T) []bucket[Key, any] {
				return []bucket[Key, any]{
					{key: NewMockKey(test), value: "five #2"},
					{key: NewMockKey(test), value: "five #1"},
				}
			},
		},
		{
			name: "with few buckets in different segments and without an interrupt",
			fields: fields{
				buckets: func(test *testing.T) [][]*bucket[Key, any] {
					return [][]*bucket[Key, any]{
						5: {5: {key: NewMockKey(test), value: "five"}},
						6: {6: {key: NewMockKey(test), value: "six"}},
						7: {7: {key: NewMockKey(test), value: "seven"}},
					}
				},
			},
			interruptOnCount: 10,
			wantBuckets: func(test *testing.T) []bucket[Key, any] {
				return []bucket[Key, any]{
					{key: NewMockKey(test), value: "six"},
					{key: NewMockKey(test), value: "five"},
					{key: NewMockKey(test), value: "seven"},
				}
			},
		},
		{
			name: "with few buckets in different segments and with an interrupt",
			fields: fields{
				buckets: func(test *testing.T) [][]*bucket[Key, any] {
					return [][]*bucket[Key, any]{
						5: {5: {key: NewMockKey(test), value: "five"}},
						6: {6: {key: NewMockKey(test), value: "six"}},
						7: {7: {key: NewMockKey(test), value: "seven"}},
					}
				},
			},
			interruptOnCount: 2,
			wantBuckets: func(test *testing.T) []bucket[Key, any] {
				return []bucket[Key, any]{
					{key: NewMockKey(test), value: "six"},
					{key: NewMockKey(test), value: "five"},
				}
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			var segments []Storage
			for _, buckets := range data.fields.buckets(test) {
				innerMap := &HashMap[Key, any]{
					config:  Config{permuter: rand.New(rand.NewPCG(0, 0)).Perm},
					buckets: buckets,
				}
				segment := &SynchronizedHashMap{innerMap: innerMap}
				segments = append(segments, segment)
			}

			var gotBuckets []bucket[Key, any]
			hashMap := ConcurrentHashMap{
				config:   ConcurrentConfig{permuter: rand.New(rand.NewPCG(0, 0)).Perm},
				segments: segments,
			}
			for key, value := range hashMap.All() {
				gotBuckets = append(gotBuckets, bucket[Key, any]{key, value})
				// interrupt after a specified count of got buckets
				if len(gotBuckets) >= data.interruptOnCount {
					break
				}
			}

			assert.Equal(test, data.wantBuckets(test), gotBuckets)
		})
	}
}

func TestConcurrentHashMap_All_order(test *testing.T) {
	type fields struct {
		buckets func(test *testing.T) [][]*bucket[Key, any]
	}

	for _, data := range []struct {
		name           string
		fields         fields
		randomSeedOne  int64
		randomSeedTwo  int64
		wantBucketsOne func(test *testing.T) []bucket[Key, any]
		wantBucketsTwo func(test *testing.T) []bucket[Key, any]
	}{
		{
			name: "with few buckets in the same segment",
			fields: fields{
				buckets: func(test *testing.T) [][]*bucket[Key, any] {
					return [][]*bucket[Key, any]{
						5: {
							5: {key: NewMockKey(test), value: "five #1"},
							6: {key: NewMockKey(test), value: "five #2"},
							7: {key: NewMockKey(test), value: "five #3"},
						},
					}
				},
			},
			randomSeedOne: 1,
			randomSeedTwo: 2,
			wantBucketsOne: func(test *testing.T) []bucket[Key, any] {
				return []bucket[Key, any]{
					{key: NewMockKey(test), value: "five #2"},
					{key: NewMockKey(test), value: "five #1"},
					{key: NewMockKey(test), value: "five #3"},
				}
			},
			wantBucketsTwo: func(test *testing.T) []bucket[Key, any] {
				return []bucket[Key, any]{
					{key: NewMockKey(test), value: "five #3"},
					{key: NewMockKey(test), value: "five #2"},
					{key: NewMockKey(test), value: "five #1"},
				}
			},
		},
		{
			name: "with few buckets in different segments",
			fields: fields{
				buckets: func(test *testing.T) [][]*bucket[Key, any] {
					return [][]*bucket[Key, any]{
						5: {5: {key: NewMockKey(test), value: "five"}},
						6: {6: {key: NewMockKey(test), value: "six"}},
						7: {7: {key: NewMockKey(test), value: "seven"}},
					}
				},
			},
			randomSeedOne: 1,
			randomSeedTwo: 2,
			wantBucketsOne: func(test *testing.T) []bucket[Key, any] {
				return []bucket[Key, any]{
					{key: NewMockKey(test), value: "six"},
					{key: NewMockKey(test), value: "five"},
					{key: NewMockKey(test), value: "seven"},
				}
			},
			wantBucketsTwo: func(test *testing.T) []bucket[Key, any] {
				return []bucket[Key, any]{
					{key: NewMockKey(test), value: "seven"},
					{key: NewMockKey(test), value: "six"},
					{key: NewMockKey(test), value: "five"},
				}
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			var segments []Storage
			for _, buckets := range data.fields.buckets(test) {
				innerMap := &HashMap[Key, any]{
					config:  Config{permuter: rand.New(rand.NewPCG(0, 0)).Perm},
					buckets: buckets,
				}
				segment := &SynchronizedHashMap{innerMap: innerMap}
				segments = append(segments, segment)
			}

			hashMap := ConcurrentHashMap{
				config:   ConcurrentConfig{permuter: rand.New(rand.NewPCG(0, 0)).Perm},
				segments: segments,
			}

			var gotBucketsOne []bucket[Key, any]
			for key, value := range hashMap.All() {
				gotBucketsOne = append(gotBucketsOne, bucket[Key, any]{key, value})
			}

			var gotBucketsTwo []bucket[Key, any]
			for key, value := range hashMap.All() {
				gotBucketsTwo = append(gotBucketsTwo, bucket[Key, any]{key, value})
			}

			assert.Equal(test, data.wantBucketsOne(test), gotBucketsOne)
			assert.Equal(test, data.wantBucketsTwo(test), gotBucketsTwo)
		})
	}
}
