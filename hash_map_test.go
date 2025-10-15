package hashmap

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewHashMap(test *testing.T) {
	type args struct {
		options []Option
	}

	for _, data := range []struct {
		name string
		args args
		want *HashMap[Key, any]
	}{
		{
			name: "with the default config",
			args: args{
				options: nil,
			},
			want: &HashMap[Key, any]{
				config: func() Config {
					config := defaultConfig
					config.permuter = nil

					return config
				}(),
				buckets: make([]*bucket[Key, any], defaultConfig.initialCapacity),
				size:    0,
			},
		},
		{
			name: "with the set initial capacity",
			args: args{
				options: []Option{WithInitialCapacity(23)},
			},
			want: &HashMap[Key, any]{
				config: func() Config {
					config := defaultConfig
					config.initialCapacity = 23
					config.permuter = nil

					return config
				}(),
				buckets: make([]*bucket[Key, any], 23),
				size:    0,
			},
		},
		{
			name: "with the set maximal load factor",
			args: args{
				options: []Option{WithMaxLoadFactor(23)},
			},
			want: &HashMap[Key, any]{
				config: func() Config {
					config := defaultConfig
					config.maxLoadFactor = 23
					config.permuter = nil

					return config
				}(),
				buckets: make([]*bucket[Key, any], defaultConfig.initialCapacity),
				size:    0,
			},
		},
		{
			name: "with the set grow factor",
			args: args{
				options: []Option{WithGrowFactor(23)},
			},
			want: &HashMap[Key, any]{
				config: func() Config {
					config := defaultConfig
					config.growFactor = 23
					config.permuter = nil

					return config
				}(),
				buckets: make([]*bucket[Key, any], defaultConfig.initialCapacity),
				size:    0,
			},
		},
		{
			name: "with the set config",
			args: args{
				options: []Option{
					WithInitialCapacity(12),
					WithMaxLoadFactor(23),
					WithGrowFactor(42),
					WithPermuter(rand.New(rand.NewPCG(0, 0)).Perm),
				},
			},
			want: &HashMap[Key, any]{
				config: Config{
					initialCapacity: 12,
					maxLoadFactor:   23,
					growFactor:      42,
					permuter:        nil,
				},
				buckets: make([]*bucket[Key, any], 12),
				size:    0,
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewHashMap[Key, any](data.args.options...)

			assert.NotNil(test, got.config.permuter)
			got.config.permuter = nil

			assert.Equal(test, data.want, got)
		})
	}
}

func TestHashMap_Get(test *testing.T) {
	type fields struct {
		buckets func(test *testing.T) []*bucket[Key, any]
	}
	type args struct {
		key func(test *testing.T) Key
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
				buckets: func(test *testing.T) []*bucket[Key, any] {
					return make([]*bucket[Key, any], defaultConfig.initialCapacity)
				},
			},
			args: args{
				key: func(test *testing.T) Key {
					key := NewMockKey(test)
					key.EXPECT().Hash().Return(5)

					return key
				},
			},
			wantValue: nil,
			wantOk:    assert.False,
		},
		{
			name: "with few buckets and a match at the start",
			fields: fields{
				buckets: func(test *testing.T) []*bucket[Key, any] {
					fiveKey := NewMockKey(test)
					fiveKey.EXPECT().Equals(mock.Anything).Return(true)

					buckets := make([]*bucket[Key, any], defaultConfig.initialCapacity)
					buckets[5] = &bucket[Key, any]{key: fiveKey, value: "five"}
					buckets[6] = &bucket[Key, any]{key: NewMockKey(test), value: "six"}
					buckets[7] = &bucket[Key, any]{key: NewMockKey(test), value: "seven"}

					return buckets
				},
			},
			args: args{
				key: func(test *testing.T) Key {
					key := NewMockKey(test)
					key.EXPECT().Hash().Return(5)

					return key
				},
			},
			wantValue: "five",
			wantOk:    assert.True,
		},
		{
			name: "with few buckets and a match at the end",
			fields: fields{
				buckets: func(test *testing.T) []*bucket[Key, any] {
					fiveKey := NewMockKey(test)
					fiveKey.EXPECT().Equals(mock.Anything).Return(false)

					sixKey := NewMockKey(test)
					sixKey.EXPECT().Equals(mock.Anything).Return(false)

					sevenKey := NewMockKey(test)
					sevenKey.EXPECT().Equals(mock.Anything).Return(true)

					buckets := make([]*bucket[Key, any], defaultConfig.initialCapacity)
					buckets[5] = &bucket[Key, any]{key: fiveKey, value: "five"}
					buckets[6] = &bucket[Key, any]{key: sixKey, value: "six"}
					buckets[7] = &bucket[Key, any]{key: sevenKey, value: "seven"}

					return buckets
				},
			},
			args: args{
				key: func(test *testing.T) Key {
					key := NewMockKey(test)
					key.EXPECT().Hash().Return(5)

					return key
				},
			},
			wantValue: "seven",
			wantOk:    assert.True,
		},
		{
			name: "with few buckets and no match",
			fields: fields{
				buckets: func(test *testing.T) []*bucket[Key, any] {
					fiveKey := NewMockKey(test)
					fiveKey.EXPECT().Equals(mock.Anything).Return(false)

					sixKey := NewMockKey(test)
					sixKey.EXPECT().Equals(mock.Anything).Return(false)

					sevenKey := NewMockKey(test)
					sevenKey.EXPECT().Equals(mock.Anything).Return(false)

					buckets := make([]*bucket[Key, any], defaultConfig.initialCapacity)
					buckets[5] = &bucket[Key, any]{key: fiveKey, value: "five"}
					buckets[6] = &bucket[Key, any]{key: sixKey, value: "six"}
					buckets[7] = &bucket[Key, any]{key: sevenKey, value: "seven"}

					return buckets
				},
			},
			args: args{
				key: func(test *testing.T) Key {
					key := NewMockKey(test)
					key.EXPECT().Hash().Return(5)

					return key
				},
			},
			wantValue: nil,
			wantOk:    assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			hashMap := HashMap[Key, any]{buckets: data.fields.buckets(test)}
			gotValue, gotOk := hashMap.Get(data.args.key(test))

			assert.Equal(test, data.wantValue, gotValue)
			data.wantOk(test, gotOk)
		})
	}
}

func TestHashMap_All(test *testing.T) {
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
			hashMap := HashMap[Key, any]{
				config:  Config{permuter: rand.New(rand.NewPCG(0, 0)).Perm},
				buckets: data.fields.buckets(test),
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

func TestHashMap_All_order(test *testing.T) {
	hashMap := HashMap[Key, any]{
		config: Config{permuter: rand.New(rand.NewPCG(0, 0)).Perm},
		buckets: []*bucket[Key, any]{
			5: {key: NewMockKey(test), value: "five"},
			6: {key: NewMockKey(test), value: "six"},
			7: {key: NewMockKey(test), value: "seven"},
		},
	}

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

func TestHashMap_Set(test *testing.T) {
	type fields struct {
		config  Config
		buckets func(test *testing.T) []*bucket[Key, any]
		size    int
	}
	type args struct {
		key func(test *testing.T) Key
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
				config: func() Config {
					config := defaultConfig
					config.permuter = nil

					return config
				}(),
				buckets: func(test *testing.T) []*bucket[Key, any] {
					return make([]*bucket[Key, any], defaultConfig.initialCapacity)
				},
				size: 0,
			},
			args: args{
				key: func(test *testing.T) Key {
					key := NewMockKey(test)
					key.EXPECT().Hash().Return(5)
					// it's called inside the HashMap.Get() method below
					key.EXPECT().Equals(mock.Anything).Return(true)

					return key
				},
			},
			wantSize:     1,
			wantCapacity: defaultConfig.initialCapacity,
		},
		{
			name: "with few buckets and a match at the start",
			fields: fields{
				config: func() Config {
					config := defaultConfig
					config.permuter = nil

					return config
				}(),
				buckets: func(test *testing.T) []*bucket[Key, any] {
					fiveKey := NewMockKey(test)
					fiveKey.EXPECT().Equals(mock.Anything).Return(true)

					buckets := make([]*bucket[Key, any], defaultConfig.initialCapacity)
					buckets[5] = &bucket[Key, any]{key: fiveKey, value: "five"}
					buckets[6] = &bucket[Key, any]{key: NewMockKey(test), value: "six"}
					buckets[7] = &bucket[Key, any]{key: NewMockKey(test), value: "seven"}

					return buckets
				},
				size: 3,
			},
			args: args{
				key: func(test *testing.T) Key {
					key := NewMockKey(test)
					key.EXPECT().Hash().Return(5)

					return key
				},
			},
			wantSize:     3,
			wantCapacity: defaultConfig.initialCapacity,
		},
		{
			name: "with few buckets and a match at the end",
			fields: fields{
				config: func() Config {
					config := defaultConfig
					config.permuter = nil

					return config
				}(),
				buckets: func(test *testing.T) []*bucket[Key, any] {
					fiveKey := NewMockKey(test)
					fiveKey.EXPECT().Equals(mock.Anything).Return(false)

					sixKey := NewMockKey(test)
					sixKey.EXPECT().Equals(mock.Anything).Return(false)

					sevenKey := NewMockKey(test)
					sevenKey.EXPECT().Equals(mock.Anything).Return(true)

					buckets := make([]*bucket[Key, any], defaultConfig.initialCapacity)
					buckets[5] = &bucket[Key, any]{key: fiveKey, value: "five"}
					buckets[6] = &bucket[Key, any]{key: sixKey, value: "six"}
					buckets[7] = &bucket[Key, any]{key: sevenKey, value: "seven"}

					return buckets
				},
				size: 3,
			},
			args: args{
				key: func(test *testing.T) Key {
					key := NewMockKey(test)
					key.EXPECT().Hash().Return(5)

					return key
				},
			},
			wantSize:     3,
			wantCapacity: defaultConfig.initialCapacity,
		},
		{
			name: "with few buckets and no match",
			fields: fields{
				config: func() Config {
					config := defaultConfig
					config.permuter = nil

					return config
				}(),
				buckets: func(test *testing.T) []*bucket[Key, any] {
					fiveKey := NewMockKey(test)
					fiveKey.EXPECT().Equals(mock.Anything).Return(false)

					sixKey := NewMockKey(test)
					sixKey.EXPECT().Equals(mock.Anything).Return(false)

					sevenKey := NewMockKey(test)
					sevenKey.EXPECT().Equals(mock.Anything).Return(false)

					buckets := make([]*bucket[Key, any], defaultConfig.initialCapacity)
					buckets[5] = &bucket[Key, any]{key: fiveKey, value: "five"}
					buckets[6] = &bucket[Key, any]{key: sixKey, value: "six"}
					buckets[7] = &bucket[Key, any]{key: sevenKey, value: "seven"}

					return buckets
				},
				size: 3,
			},
			args: args{
				key: func(test *testing.T) Key {
					key := NewMockKey(test)
					key.EXPECT().Hash().Return(5)
					// it's called inside the HashMap.Get() method below
					key.EXPECT().Equals(mock.Anything).Return(true)

					return key
				},
			},
			wantSize:     4,
			wantCapacity: defaultConfig.initialCapacity,
		},
		{
			name: "with a load factor over the maximum and a match",
			fields: fields{
				config: func() Config {
					config := defaultConfig
					config.permuter = nil

					return config
				}(),
				buckets: func(test *testing.T) []*bucket[Key, any] {
					threeKey := NewMockKey(test)
					threeKey.EXPECT().Equals(mock.Anything).Return(true)

					buckets := make([]*bucket[Key, any], 5)
					buckets[0] = &bucket[Key, any]{key: NewMockKey(test), value: "zero"}
					buckets[1] = &bucket[Key, any]{key: NewMockKey(test), value: "one"}
					buckets[2] = &bucket[Key, any]{key: NewMockKey(test), value: "two"}
					buckets[3] = &bucket[Key, any]{key: threeKey, value: "three"}

					return buckets
				},
				size: 4,
			},
			args: args{
				key: func(test *testing.T) Key {
					key := NewMockKey(test)
					key.EXPECT().Hash().Return(3)

					return key
				},
			},
			wantSize:     4,
			wantCapacity: 5,
		},
		{
			name: "with a load factor over the maximum and no match",
			fields: fields{
				config: func() Config {
					config := defaultConfig
					config.permuter = rand.New(rand.NewPCG(0, 0)).Perm

					return config
				}(),
				buckets: func(test *testing.T) []*bucket[Key, any] {
					zeroKey := NewMockKey(test)
					zeroKey.EXPECT().Hash().Return(0)

					oneKey := NewMockKey(test)
					oneKey.EXPECT().Hash().Return(1)

					twoKey := NewMockKey(test)
					twoKey.EXPECT().Hash().Return(2)

					threeKey := NewMockKey(test)
					threeKey.EXPECT().Hash().Return(3)
					threeKey.EXPECT().Equals(mock.Anything).Return(false)

					buckets := make([]*bucket[Key, any], 5)
					buckets[0] = &bucket[Key, any]{key: zeroKey, value: "zero"}
					buckets[1] = &bucket[Key, any]{key: oneKey, value: "one"}
					buckets[2] = &bucket[Key, any]{key: twoKey, value: "two"}
					buckets[3] = &bucket[Key, any]{key: threeKey, value: "three"}

					return buckets
				},
				size: 4,
			},
			args: args{
				key: func(test *testing.T) Key {
					key := NewMockKey(test)
					key.EXPECT().Hash().Return(3)
					// it's called inside the HashMap.Get() method below
					key.EXPECT().Equals(mock.Anything).Return(true)

					return key
				},
			},
			wantSize:     5,
			wantCapacity: 10,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			key := data.args.key(test)
			value := rand.Int()

			hashMap := HashMap[Key, any]{
				config:  data.fields.config,
				buckets: data.fields.buckets(test),
				size:    data.fields.size,
			}
			hashMap.Set(key, value)

			gotValue, gotOk := hashMap.Get(key)

			assert.Equal(test, data.wantSize, hashMap.size)
			assert.Equal(test, data.wantCapacity, len(hashMap.buckets))
			assert.Equal(test, value, gotValue)
			assert.True(test, gotOk)
		})
	}
}

func TestHashMap_Delete(test *testing.T) {
	type fields struct {
		buckets func(test *testing.T) []*bucket[Key, any]
		size    int
	}
	type args struct {
		key func(test *testing.T) Key
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
				buckets: func(test *testing.T) []*bucket[Key, any] {
					return make([]*bucket[Key, any], defaultConfig.initialCapacity)
				},
				size: 0,
			},
			args: args{
				key: func(test *testing.T) Key {
					key := NewMockKey(test)
					key.EXPECT().Hash().Return(5)

					return key
				},
			},
			wantSize: 0,
		},
		{
			name: "with few buckets and a match at the start",
			fields: fields{
				buckets: func(test *testing.T) []*bucket[Key, any] {
					fiveKey := NewMockKey(test)
					fiveKey.EXPECT().Equals(mock.Anything).Return(true)

					buckets := make([]*bucket[Key, any], defaultConfig.initialCapacity)
					buckets[5] = &bucket[Key, any]{key: fiveKey, value: "five"}
					buckets[6] = &bucket[Key, any]{key: NewMockKey(test), value: "six"}
					buckets[7] = &bucket[Key, any]{key: NewMockKey(test), value: "seven"}

					return buckets
				},
				size: 3,
			},
			args: args{
				key: func(test *testing.T) Key {
					key := NewMockKey(test)
					key.EXPECT().Hash().Return(5)

					return key
				},
			},
			wantSize: 2,
		},
		{
			name: "with few buckets and a match at the end",
			fields: fields{
				buckets: func(test *testing.T) []*bucket[Key, any] {
					fiveKey := NewMockKey(test)
					fiveKey.EXPECT().Equals(mock.Anything).Return(false)

					sixKey := NewMockKey(test)
					sixKey.EXPECT().Equals(mock.Anything).Return(false)

					sevenKey := NewMockKey(test)
					sevenKey.EXPECT().Equals(mock.Anything).Return(true)

					buckets := make([]*bucket[Key, any], defaultConfig.initialCapacity)
					buckets[5] = &bucket[Key, any]{key: fiveKey, value: "five"}
					buckets[6] = &bucket[Key, any]{key: sixKey, value: "six"}
					buckets[7] = &bucket[Key, any]{key: sevenKey, value: "seven"}

					return buckets
				},
				size: 3,
			},
			args: args{
				key: func(test *testing.T) Key {
					key := NewMockKey(test)
					key.EXPECT().Hash().Return(5)

					return key
				},
			},
			wantSize: 2,
		},
		{
			name: "with few buckets and no match",
			fields: fields{
				buckets: func(test *testing.T) []*bucket[Key, any] {
					fiveKey := NewMockKey(test)
					fiveKey.EXPECT().Equals(mock.Anything).Return(false)

					sixKey := NewMockKey(test)
					sixKey.EXPECT().Equals(mock.Anything).Return(false)

					sevenKey := NewMockKey(test)
					sevenKey.EXPECT().Equals(mock.Anything).Return(false)

					buckets := make([]*bucket[Key, any], defaultConfig.initialCapacity)
					buckets[5] = &bucket[Key, any]{key: fiveKey, value: "five"}
					buckets[6] = &bucket[Key, any]{key: sixKey, value: "six"}
					buckets[7] = &bucket[Key, any]{key: sevenKey, value: "seven"}

					return buckets
				},
				size: 3,
			},
			args: args{
				key: func(test *testing.T) Key {
					key := NewMockKey(test)
					key.EXPECT().Hash().Return(5)

					return key
				},
			},
			wantSize: 3,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			key := data.args.key(test)

			hashMap := HashMap[Key, any]{
				buckets: data.fields.buckets(test),
				size:    data.fields.size,
			}
			hashMap.Delete(key)

			_, gotGetOk := hashMap.Get(key)

			assert.Equal(test, data.wantSize, hashMap.size)
			assert.False(test, gotGetOk)
		})
	}
}
