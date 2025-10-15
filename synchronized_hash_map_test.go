package hashmap

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewSynchronizedHashMap(test *testing.T) {
	type args struct {
		options func(test *testing.T) []SynchronizedOption
	}

	for _, data := range []struct {
		name string
		args args
		want func(test *testing.T) *SynchronizedHashMap
	}{
		{
			name: "with the default config",
			args: args{
				options: func(test *testing.T) []SynchronizedOption { return nil },
			},
			want: func(test *testing.T) *SynchronizedHashMap {
				return &SynchronizedHashMap{
					innerMap: &HashMap[Key, any]{
						config: func() Config {
							config := defaultConfig
							config.permuter = nil

							return config
						}(),
						buckets: make([]*bucket[Key, any], defaultConfig.initialCapacity),
						size:    0,
					},
				}
			},
		},
		{
			name: "with the set inner map",
			args: args{
				options: func(test *testing.T) []SynchronizedOption {
					return []SynchronizedOption{WithInnerMap(NewMockStorage(test))}
				},
			},
			want: func(test *testing.T) *SynchronizedHashMap {
				return &SynchronizedHashMap{
					innerMap: NewMockStorage(test),
				}
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewSynchronizedHashMap(data.args.options(test)...)
			if innerMap, ok := got.innerMap.(*HashMap[Key, any]); ok {
				innerMap.config.permuter = nil
			}

			assert.Equal(test, data.want(test), got)
		})
	}
}

func TestSynchronizedHashMap(test *testing.T) {
	for _, data := range []struct {
		name      string
		hashMap   func(test *testing.T) *SynchronizedHashMap
		key       func(test *testing.T) Key
		wantValue interface{}
		wantOk    assert.BoolAssertionFunc
	}{
		{
			name: "getting by a nonexistent key",
			hashMap: func(test *testing.T) *SynchronizedHashMap {
				return NewSynchronizedHashMap()
			},
			key: func(test *testing.T) Key {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)

				return key
			},
			wantValue: nil,
			wantOk:    assert.False,
		},
		{
			name: "setting by a nonexistent key",
			hashMap: func(test *testing.T) *SynchronizedHashMap {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)
				// it's called inside the HashMap.Get() method below
				key.EXPECT().Equals(mock.Anything).Return(true)

				hashMap := NewSynchronizedHashMap()
				hashMap.Set(key, "five")

				return hashMap
			},
			key: func(test *testing.T) Key {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)

				return key
			},
			wantValue: "five",
			wantOk:    assert.True,
		},
		{
			name: "setting by an existing key",
			hashMap: func(test *testing.T) *SynchronizedHashMap {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)
				key.EXPECT().Equals(mock.Anything).Return(true)

				hashMap := NewSynchronizedHashMap()
				hashMap.Set(key, "five #1")
				hashMap.Set(key, "five #2")

				return hashMap
			},
			key: func(test *testing.T) Key {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)

				return key
			},
			wantValue: "five #2",
			wantOk:    assert.True,
		},
		{
			name: "deleting by a nonexistent key",
			hashMap: func(test *testing.T) *SynchronizedHashMap {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)

				hashMap := NewSynchronizedHashMap()
				hashMap.Delete(key)

				return hashMap
			},
			key: func(test *testing.T) Key {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)

				return key
			},
			wantValue: nil,
			wantOk:    assert.False,
		},
		{
			name: "deleting by an existing key",
			hashMap: func(test *testing.T) *SynchronizedHashMap {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)
				key.EXPECT().Equals(mock.Anything).Return(true)

				hashMap := NewSynchronizedHashMap()
				hashMap.Set(key, "five")
				hashMap.Delete(key)

				return hashMap
			},
			key: func(test *testing.T) Key {
				key := NewMockKey(test)
				key.EXPECT().Hash().Return(5)

				return key
			},
			wantValue: nil,
			wantOk:    assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotValue, gotOk := data.hashMap(test).Get(data.key(test))

			assert.Equal(test, data.wantValue, gotValue)
			data.wantOk(test, gotOk)
		})
	}
}

func TestSynchronizedHashMap_All(test *testing.T) {
	type fields struct {
		buckets func(test *testing.T) []*bucket[Key, any]
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
				buckets: func(test *testing.T) []*bucket[Key, any] {
					return make([]*bucket[Key, any], defaultConfig.initialCapacity)
				},
			},
			interruptOnCount: 10,
			wantBuckets:      func(test *testing.T) []bucket[Key, any] { return nil },
		},
		{
			name: "with few buckets and without an interrupt",
			fields: fields{
				buckets: func(test *testing.T) []*bucket[Key, any] {
					return []*bucket[Key, any]{
						5: {key: NewMockKey(test), value: "five"},
						6: {key: NewMockKey(test), value: "six"},
						7: {key: NewMockKey(test), value: "seven"},
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
			name: "with few buckets and with an interrupt",
			fields: fields{
				buckets: func(test *testing.T) []*bucket[Key, any] {
					return []*bucket[Key, any]{
						5: {key: NewMockKey(test), value: "five"},
						6: {key: NewMockKey(test), value: "six"},
						7: {key: NewMockKey(test), value: "seven"},
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
			var gotBuckets []bucket[Key, any]
			innerMap := HashMap[Key, any]{
				config:  Config{permuter: rand.New(rand.NewPCG(0, 0)).Perm},
				buckets: data.fields.buckets(test),
			}
			hashMap := SynchronizedHashMap{innerMap: &innerMap}
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

func TestSynchronizedHashMap_All_order(test *testing.T) {
	innerMap := HashMap[Key, any]{
		config: Config{permuter: rand.New(rand.NewPCG(0, 0)).Perm},
		buckets: []*bucket[Key, any]{
			5: {key: NewMockKey(test), value: "five"},
			6: {key: NewMockKey(test), value: "six"},
			7: {key: NewMockKey(test), value: "seven"},
		},
	}
	hashMap := SynchronizedHashMap{innerMap: &innerMap}

	var gotBucketsOne []bucket[Key, any]
	for key, value := range hashMap.All() {
		gotBucketsOne = append(gotBucketsOne, bucket[Key, any]{key, value})
	}

	var gotBucketsTwo []bucket[Key, any]
	for key, value := range hashMap.All() {
		gotBucketsTwo = append(gotBucketsTwo, bucket[Key, any]{key, value})
	}

	wantBucketsOne := []bucket[Key, any]{
		{key: NewMockKey(test), value: "six"},
		{key: NewMockKey(test), value: "five"},
		{key: NewMockKey(test), value: "seven"},
	}
	assert.Equal(test, wantBucketsOne, gotBucketsOne)

	wantBucketsTwo := []bucket[Key, any]{
		{key: NewMockKey(test), value: "seven"},
		{key: NewMockKey(test), value: "six"},
		{key: NewMockKey(test), value: "five"},
	}
	assert.Equal(test, wantBucketsTwo, gotBucketsTwo)
}
