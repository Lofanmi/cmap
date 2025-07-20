package cmap

import (
	"cmp"
	"strconv"
	"strings"
	"sync"

	"github.com/emirpasic/gods/v2/maps"
)

// Map 并发安全的Map实现
type Map[K cmp.Ordered, V any] struct {
	shards []shard[K, V]
	mask   uint32
	dirty  bool
	mu     *sync.RWMutex // 用于保护dirty字段
	hasher hasher[K]     // 哈希器
	opts   *Options
}

// shard 分片结构
type shard[K cmp.Ordered, V any] struct {
	m    maps.Map[K, V]
	mu   *sync.RWMutex
	opts *Options
}

// Put 插入键值对
func (m *Map[K, V]) Put(key K, value V) {
	sh := m.getShard(key)
	sh.mu.Lock()
	sh.m.Put(key, value)
	sh.mu.Unlock()

	m.mu.Lock()
	m.dirty = true
	m.mu.Unlock()
}

// Get 获取值
func (m *Map[K, V]) Get(key K) (value V, found bool) {
	sh := m.getShard(key)
	sh.mu.RLock()
	value, found = sh.m.Get(key)
	sh.mu.RUnlock()
	return
}

// Remove 删除键
func (m *Map[K, V]) Remove(key K) {
	sh := m.getShard(key)
	sh.mu.Lock()
	_, ok := sh.m.Get(key)
	if ok {
		sh.m.Remove(key)
	}
	sh.mu.Unlock()

	if ok {
		m.mu.Lock()
		m.dirty = true
		m.mu.Unlock()
	}
}

// Empty 检查是否为空
func (m *Map[K, V]) Empty() bool {
	for i := range m.shards {
		m.shards[i].mu.RLock()
		empty := m.shards[i].m.Empty()
		m.shards[i].mu.RUnlock()
		if !empty {
			return false
		}
	}
	return true
}

// Size 获取大小
func (m *Map[K, V]) Size() int {
	size := 0
	for i := range m.shards {
		m.shards[i].mu.RLock()
		size += m.shards[i].m.Size()
		m.shards[i].mu.RUnlock()
	}
	return size
}

// Clear 清空映射
func (m *Map[K, V]) Clear() {
	for i := range m.shards {
		m.shards[i].mu.Lock()
		m.shards[i].m.Clear()
		m.shards[i].mu.Unlock()
	}

	m.mu.Lock()
	m.dirty = true
	m.mu.Unlock()
}

// Keys 获取所有键
func (m *Map[K, V]) Keys() []K {
	var keys []K
	for i := range m.shards {
		m.shards[i].mu.RLock()
		shardKeys := m.shards[i].m.Keys()
		keys = append(keys, shardKeys...)
		m.shards[i].mu.RUnlock()
	}
	return keys
}

// Values 获取所有值
func (m *Map[K, V]) Values() []V {
	var values []V
	for i := range m.shards {
		m.shards[i].mu.RLock()
		shardValues := m.shards[i].m.Values()
		values = append(values, shardValues...)
		m.shards[i].mu.RUnlock()
	}
	return values
}

// String 字符串表示
func (m *Map[K, V]) String() string {
	var b strings.Builder
	b.WriteString("CMap:\n")

	for i, sh := range m.shards {
		sh.mu.RLock()
		if !sh.m.Empty() {
			b.WriteString("- Shard ")
			b.WriteString(strconv.Itoa(i))
			b.WriteString(": ")
			b.WriteString(sh.m.String())
			b.WriteString("\n")
		}
		sh.mu.RUnlock()
	}

	return b.String()
}

// IsDirty 检查映射是否被修改过
func (m *Map[K, V]) IsDirty() (dirty bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	dirty = m.dirty
	return
}

// getShard 获取键对应的分片
func (m *Map[K, V]) getShard(key K) *shard[K, V] {
	hash := m.hasher.Hash(key)
	return &m.shards[hash&m.mask]
}

// roundUpToPowerOf2 向上取整到2的幂
func roundUpToPowerOf2(n uint32) uint32 {
	if n == 0 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}
