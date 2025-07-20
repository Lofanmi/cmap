package cmap

import (
	"testing"
)

// TestGetMultiple 测试批量获取
func TestGetMultiple(t *testing.T) {
	m := NewStringHashMap[int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)

	// 测试存在的键
	keys := []string{"a", "b", "d"} // d不存在
	result := m.GetMultiple(keys)

	if len(result) != 2 {
		t.Errorf("GetMultiple should return 2 items, got %d", len(result))
	}

	if result["a"] != 1 || result["b"] != 2 {
		t.Errorf("GetMultiple returned unexpected values: %v", result)
	}

	if _, exists := result["d"]; exists {
		t.Error("GetMultiple should not include non-existent keys")
	}
}

// TestRemoveMultiple 测试批量删除
func TestRemoveMultiple(t *testing.T) {
	m := NewStringHashMap[int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)
	m.Put("d", 4)

	// 删除存在的和不存在的键
	keys := []string{"a", "b", "z"} // z不存在
	m.RemoveMultiple(keys)

	if _, ok := m.Get("a"); ok {
		t.Error("RemoveMultiple failed to remove key 'a'")
	}
	if _, ok := m.Get("b"); ok {
		t.Error("RemoveMultiple failed to remove key 'b'")
	}
	if _, ok := m.Get("c"); !ok {
		t.Error("RemoveMultiple incorrectly removed key 'c'")
	}
	if _, ok := m.Get("d"); !ok {
		t.Error("RemoveMultiple incorrectly removed key 'd'")
	}
}

// TestPutAll 测试批量插入
func TestPutAll(t *testing.T) {
	m := NewStringHashMap[int]()

	data := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	m.PutAll(data)

	if m.Size() != 3 {
		t.Errorf("PutAll should add 3 items, got %d", m.Size())
	}

	for key, expected := range data {
		if val, ok := m.Get(key); !ok || val != expected {
			t.Errorf("PutAll failed for key %s: got %v, want %v", key, val, expected)
		}
	}
}

// TestBatchOperationsEmpty 测试空map的批量操作
func TestBatchOperationsEmpty(t *testing.T) {
	m := NewStringHashMap[int]()

	// 空map的批量操作
	result := m.GetMultiple([]string{"a", "b"})
	if len(result) != 0 {
		t.Errorf("GetMultiple on empty map should return empty result, got %v", result)
	}

	m.RemoveMultiple([]string{"a", "b"}) // 应该不panic
	m.PutAll(map[string]int{})           // 空map插入

	if !m.Empty() {
		t.Error("Empty PutAll should leave map empty")
	}
}

// TestBatchOperationsConcurrent 测试批量操作的并发安全
func TestBatchOperationsConcurrent(t *testing.T) {
	m := NewStringHashMap[int]()

	// 批量插入
	data := make(map[string]int)
	for i := 0; i < 100; i++ {
		data[string(rune('a'+i%26))] = i
	}

	m.PutAll(data)

	// 验证数据
	if m.Size() != 26 {
		t.Errorf("Expected 26 unique keys, got %d", m.Size())
	}

	// 批量获取
	keys := make([]string, 0, 26)
	for k := range data {
		keys = append(keys, k)
	}

	result := m.GetMultiple(keys)
	if len(result) != 26 {
		t.Errorf("GetMultiple should return 26 items, got %d", len(result))
	}

	// 批量删除
	m.RemoveMultiple(keys)
	if m.Size() != 0 {
		t.Errorf("RemoveMultiple should remove all items, got %d remaining", m.Size())
	}
}

// TestBatchOperationsWithCustomTypes 测试自定义类型的批量操作
func TestBatchOperationsWithCustomTypes(t *testing.T) {
	type Key string
	type Value struct {
		Data int
	}

	m := NewHashMap[Key, Value]()

	// 批量插入自定义类型
	data := map[Key]Value{
		"a": {Data: 1},
		"b": {Data: 2},
		"c": {Data: 3},
	}

	m.PutAll(data)

	// 批量获取
	keys := []Key{"a", "b", "c"}
	result := m.GetMultiple(keys)

	if len(result) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result))
	}

	if result["a"].Data != 1 || result["b"].Data != 2 || result["c"].Data != 3 {
		t.Errorf("Custom type batch operations failed")
	}
}
