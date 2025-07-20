package cmap

import (
	"sync"
	"testing"
)

// TestMapBasicOperations 测试基本操作
func TestMapBasicOperations(t *testing.T) {
	m := NewStringHashMap[int]()

	// Test Put and Get
	m.Put("key1", 100)
	val, ok := m.Get("key1")
	if !ok || val != 100 {
		t.Errorf("Put/Get failed, got %v, want 100", val)
	}

	// Test Get non-existent
	_, ok = m.Get("nonexistent")
	if ok {
		t.Error("Get should return false for non-existent key")
	}

	// Test Remove
	m.Put("key2", 200)
	m.Remove("key2")
	_, ok = m.Get("key2")
	if ok {
		t.Error("Remove failed, key still exists")
	}

	// Test Size
	m.Put("key3", 300)
	if m.Size() != 2 {
		t.Errorf("Size failed, got %d, want 2", m.Size())
	}

	// Test Empty
	if m.Empty() {
		t.Error("Empty should return false for non-empty map")
	}

	// Test Clear
	m.Clear()
	if !m.Empty() || m.Size() != 0 {
		t.Error("Clear failed, map should be empty")
	}
}

// TestMapKeysValues 测试Keys和Values方法
func TestMapKeysValues(t *testing.T) {
	m := NewStringHashMap[int]()

	// 空map
	keys := m.Keys()
	if len(keys) != 0 {
		t.Errorf("Keys on empty map should return empty slice, got %v", keys)
	}

	values := m.Values()
	if len(values) != 0 {
		t.Errorf("Values on empty map should return empty slice, got %v", values)
	}

	// 添加数据
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)

	keys = m.Keys()
	if len(keys) != 3 {
		t.Errorf("Keys should return 3 keys, got %d", len(keys))
	}

	values = m.Values()
	if len(values) != 3 {
		t.Errorf("Values should return 3 values, got %d", len(values))
	}
}

// TestMapConcurrent 测试并发安全
func TestMapConcurrent(t *testing.T) {
	m := NewStringHashMap[int]()
	const goroutines = 100
	const operations = 1000

	var wg sync.WaitGroup
	wg.Add(goroutines)

	// 并发写入
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				key := string(rune('a' + (id+j)%26))
				m.Put(key, id+j)
			}
		}(i)
	}

	wg.Wait()

	// 验证数据一致性
	if m.Size() == 0 {
		t.Error("Concurrent operations should produce non-empty map")
	}
}

// TestMapString 测试String方法
func TestMapString(t *testing.T) {
	m := NewStringHashMap[int]()

	// 空map
	str := m.String()
	if str == "" {
		t.Error("String should not be empty")
	}

	// 有数据的map
	m.Put("test", 123)
	str = m.String()
	if str == "" {
		t.Error("String should not be empty for non-empty map")
	}
}

// TestMapWithCustomTypes 测试自定义类型
func TestMapWithCustomTypes(t *testing.T) {
	type CustomKey string
	type CustomValue struct {
		Name  string
		Value int
	}

	m := NewHashMap[CustomKey, CustomValue]()
	key := CustomKey("test")
	value := CustomValue{Name: "test", Value: 42}

	m.Put(key, value)
	retrieved, ok := m.Get(key)
	if !ok || retrieved.Name != "test" || retrieved.Value != 42 {
		t.Errorf("Custom type test failed, got %v, want %v", retrieved, value)
	}
}

// TestMapWithFloatKeys 测试浮点类型键
func TestMapWithFloatKeys(t *testing.T) {
	m := NewHashMap[float64, string]()
	m.Put(3.14, "pi")
	m.Put(2.71, "e")

	val, ok := m.Get(3.14)
	if !ok || val != "pi" {
		t.Errorf("Float key test failed, got %v", val)
	}
}

// TestMapEdgeCases 测试边界情况
func TestMapEdgeCases(t *testing.T) {
	m := NewStringHashMap[int]()

	// 测试空字符串键
	m.Put("", 0)
	val, ok := m.Get("")
	if !ok || val != 0 {
		t.Errorf("Empty string key test failed")
	}

	// 测试零值
	m.Put("zero", 0)
	val, ok = m.Get("zero")
	if !ok || val != 0 {
		t.Errorf("Zero value test failed")
	}

	// 测试覆盖现有键
	m.Put("key", 100)
	m.Put("key", 200)
	val, ok = m.Get("key")
	if !ok || val != 200 {
		t.Errorf("Overwrite test failed, got %v, want 200", val)
	}
}
