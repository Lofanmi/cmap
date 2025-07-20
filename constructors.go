package cmap

import (
	"cmp"
	"runtime"
	"sync"

	"github.com/emirpasic/gods/v2/maps"
	"github.com/emirpasic/gods/v2/maps/hashmap"
	"github.com/emirpasic/gods/v2/maps/linkedhashmap"
	"github.com/emirpasic/gods/v2/maps/treemap"
)

// ---------------------------------------------------------------------------------------------------------------------

// New 创建默认的并发映射（使用HashMap作为底层实现）
func New[K cmp.Ordered, V any](options ...Option) *Map[K, V] {
	return NewHashMap[K, V](options...)
}

// ---------------------------------------------------------------------------------------------------------------------

// NewHashMap 创建HashMap类型的并发映射
func NewHashMap[K cmp.Ordered, V any](options ...Option) *Map[K, V] {
	return createMap(func() maps.Map[K, V] {
		return hashmap.New[K, V]()
	}, options...)
}

// NewTreeMap 创建TreeMap类型的并发映射
func NewTreeMap[K cmp.Ordered, V any](options ...Option) *Map[K, V] {
	return createMap(func() maps.Map[K, V] {
		return treemap.New[K, V]()
	}, options...)
}

// NewLinkedHashMap 创建LinkedHashMap类型的并发映射
func NewLinkedHashMap[K cmp.Ordered, V any](options ...Option) *Map[K, V] {
	return createMap(func() maps.Map[K, V] {
		return linkedhashmap.New[K, V]()
	}, options...)
}

// ---------------------------------------------------------------------------------------------------------------------

// NewStringHashMap 创建使用string键的HashMap
func NewStringHashMap[V any](options ...Option) *Map[string, V] {
	return createMap(func() maps.Map[string, V] {
		return hashmap.New[string, V]()
	}, options...)
}

// NewStringTreeMap 创建使用string键的TreeMap
func NewStringTreeMap[V any](options ...Option) *Map[string, V] {
	return createMap(func() maps.Map[string, V] {
		return treemap.New[string, V]()
	}, options...)
}

// NewStringLinkedHashMap 创建使用string键的LinkedHashMap
func NewStringLinkedHashMap[V any](options ...Option) *Map[string, V] {
	return createMap(func() maps.Map[string, V] {
		return linkedhashmap.New[string, V]()
	}, options...)
}

// ---------------------------------------------------------------------------------------------------------------------

// NewIntHashMap 创建使用int键的HashMap
func NewIntHashMap[V any](options ...Option) *Map[int, V] {
	return createMap(func() maps.Map[int, V] {
		return hashmap.New[int, V]()
	}, options...)
}

// NewIntTreeMap 创建使用int键的TreeMap
func NewIntTreeMap[V any](options ...Option) *Map[int, V] {
	return createMap(func() maps.Map[int, V] {
		return treemap.New[int, V]()
	}, options...)
}

// NewIntLinkedHashMap 创建使用int键的LinkedHashMap
func NewIntLinkedHashMap[V any](options ...Option) *Map[int, V] {
	return createMap(func() maps.Map[int, V] {
		return linkedhashmap.New[int, V]()
	}, options...)
}

// ---------------------------------------------------------------------------------------------------------------------

// NewInt64HashMap 创建使用int64键的HashMap
func NewInt64HashMap[V any](options ...Option) *Map[int64, V] {
	return createMap(func() maps.Map[int64, V] {
		return hashmap.New[int64, V]()
	}, options...)
}

// NewInt64TreeMap 创建使用int64键的TreeMap
func NewInt64TreeMap[V any](options ...Option) *Map[int64, V] {
	return createMap(func() maps.Map[int64, V] {
		return treemap.New[int64, V]()
	}, options...)
}

// NewInt64LinkedHashMap 创建使用int64键的LinkedHashMap
func NewInt64LinkedHashMap[V any](options ...Option) *Map[int64, V] {
	return createMap(func() maps.Map[int64, V] {
		return linkedhashmap.New[int64, V]()
	}, options...)
}

// ---------------------------------------------------------------------------------------------------------------------

// createMap is a generic helper function that creates a new Map instance
func createMap[K cmp.Ordered, V any](
	createUnderlyingMap func() maps.Map[K, V],
	options ...Option,
) *Map[K, V] {
	opts := &Options{
		ShardCount: uint32(runtime.NumCPU() * 16),
		Serializer: JsonSerializer(),
	}
	for _, option := range options {
		option(opts)
	}

	shardCount := roundUpToPowerOf2(opts.ShardCount)
	shards := make([]shard[K, V], shardCount)
	for i := range shards {
		shards[i] = shard[K, V]{
			m:    createUnderlyingMap(),
			mu:   &sync.RWMutex{},
			opts: opts,
		}
	}

	return &Map[K, V]{
		shards: shards,
		mask:   shardCount - 1,
		dirty:  false,
		mu:     &sync.RWMutex{},
		hasher: getHasher[K](),
		opts:   opts,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
