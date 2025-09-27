package hashmap

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewSynchronizedHashMap(test *testing.T) {
	type args struct {
		options []SynchronizedOption
	}

	for _, data := range []struct {
		name string
		args args
		want *SynchronizedHashMap
	}{
		{
			name: "with the default config",
			args: args{
				options: nil,
			},
			want: &SynchronizedHashMap{
				innerMap: &HashMap{
					config: func() Config {
						config := defaultConfig
						config.permuter = nil

						return config
					}(),
					buckets: make([]*bucket, defaultConfig.initialCapacity),
					size:    0,
				},
			},
		},
		{
			name: "with the set inner map",
			args: args{
				options: []SynchronizedOption{WithInnerMap(new(MockStorage))},
			},
			want: &SynchronizedHashMap{
				innerMap: new(MockStorage),
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewSynchronizedHashMap(data.args.options...)
			if innerMap, ok := got.innerMap.(*HashMap); ok {
				innerMap.config.permuter = nil
			}

			_, ok := got.innerMap.(interface {
				AssertExpectations(assert.TestingT) bool // nolint: staticcheck
			})
			if ok {
				mock.AssertExpectationsForObjects(test, got.innerMap)
			}
			assert.Equal(test, data.want, got)
		})
	}
}

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

			for _, bucket := range hashMap.innerMap.(*HashMap).buckets {
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

func TestSynchronizedHashMap_All(test *testing.T) {
	type fields struct {
		buckets []*bucket
	}

	for _, data := range []struct {
		name             string
		fields           fields
		interruptOnCount int
		wantBuckets      []bucket
	}{
		{
			name: "without buckets",
			fields: fields{
				buckets: make([]*bucket, defaultConfig.initialCapacity),
			},
			interruptOnCount: 10,
			wantBuckets:      nil,
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
				{key: new(MockKey), value: "six"},
				{key: new(MockKey), value: "five"},
				{key: new(MockKey), value: "seven"},
			},
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
				{key: new(MockKey), value: "six"},
				{key: new(MockKey), value: "five"},
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			var gotBuckets []bucket
			innerMap := HashMap{
				config:  Config{permuter: rand.New(rand.NewPCG(0, 0)).Perm},
				buckets: data.fields.buckets,
			}
			hashMap := SynchronizedHashMap{innerMap: &innerMap}
			for key, value := range hashMap.All() {
				gotBuckets = append(gotBuckets, bucket{key, value})
				// interrupt after a specified count of got buckets
				if len(gotBuckets) >= data.interruptOnCount {
					break
				}
			}

			for _, bucket := range data.fields.buckets {
				if bucket != nil {
					mock.AssertExpectationsForObjects(test, bucket.key)
				}
			}
			assert.Equal(test, data.wantBuckets, gotBuckets)
		})
	}
}

func TestSynchronizedHashMap_All_order(test *testing.T) {
	innerMap := HashMap{
		config: Config{permuter: rand.New(rand.NewPCG(0, 0)).Perm},
		buckets: []*bucket{
			5: {key: new(MockKey), value: "five"},
			6: {key: new(MockKey), value: "six"},
			7: {key: new(MockKey), value: "seven"},
		},
	}
	hashMap := SynchronizedHashMap{innerMap: &innerMap}

	var gotBucketsOne []bucket
	for key, value := range hashMap.All() {
		gotBucketsOne = append(gotBucketsOne, bucket{key, value})
	}

	var gotBucketsTwo []bucket
	for key, value := range hashMap.All() {
		gotBucketsTwo = append(gotBucketsTwo, bucket{key, value})
	}

	for _, bucket := range hashMap.innerMap.(*HashMap).buckets {
		if bucket != nil {
			mock.AssertExpectationsForObjects(test, bucket.key)
		}
	}

	wantBucketsOne := []bucket{
		{key: new(MockKey), value: "six"},
		{key: new(MockKey), value: "five"},
		{key: new(MockKey), value: "seven"},
	}
	assert.Equal(test, wantBucketsOne, gotBucketsOne)

	wantBucketsTwo := []bucket{
		{key: new(MockKey), value: "seven"},
		{key: new(MockKey), value: "six"},
		{key: new(MockKey), value: "five"},
	}
	assert.Equal(test, wantBucketsTwo, gotBucketsTwo)
}
