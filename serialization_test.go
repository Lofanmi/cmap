package cmap

import (
	"encoding/json"
	"fmt"
	"testing"
)

// TestMarshalJSON 测试JSON序列化
func TestMarshalJSON(t *testing.T) {
	m := NewStringHashMap[int]()
	m.Put("key1", 100)
	m.Put("key2", 200)

	// 测试标准JSON序列化器
	data, err := m.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("MarshalJSON returned empty data")
	}

	// 验证JSON格式
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Errorf("Invalid JSON format: %v", err)
	}
}

// TestUnmarshalJSON 测试JSON反序列化
func TestUnmarshalJSON(t *testing.T) {
	m := NewStringHashMap[int]()

	// 准备测试数据
	jsonData := `{"items":[{"key":"key1","value":100},{"key":"key2","value":200}]}`

	err := m.UnmarshalJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if m.Size() != 2 {
		t.Errorf("Expected 2 items after unmarshal, got %d", m.Size())
	}

	val, ok := m.Get("key1")
	if !ok || val != 100 {
		t.Errorf("Unmarshal failed for key1: got %v, want 100", val)
	}

	val, ok = m.Get("key2")
	if !ok || val != 200 {
		t.Errorf("Unmarshal failed for key2: got %v, want 200", val)
	}
}

// TestMarshalUnmarshalRoundTrip 测试序列化反序列化往返
func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	// 原始map
	original := NewStringHashMap[int]()
	original.Put("a", 1)
	original.Put("b", 2)
	original.Put("c", 3)

	// 序列化
	data, err := original.MarshalJSON()
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化到新map
	restored := NewStringHashMap[int]()
	err = restored.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 验证数据一致性
	if restored.Size() != original.Size() {
		t.Errorf("Size mismatch: original %d, restored %d", original.Size(), restored.Size())
	}

	// 验证所有键值对
	originalKeys := original.Keys()
	for _, key := range originalKeys {
		originalVal, _ := original.Get(key)
		restoredVal, ok := restored.Get(key)
		if !ok || restoredVal != originalVal {
			t.Errorf("Data mismatch for key %s: original %v, restored %v", key, originalVal, restoredVal)
		}
	}
}

// TestMarshalWithDifferentSerializers 测试不同序列化器
func TestMarshalWithDifferentSerializers(t *testing.T) {
	m := NewStringHashMap[int]()
	m.Put("test", 42)

	serializers := []*SerializerFunc{
		JsonSerializer(),
		JsoniterSerializer(),
		SonicSerializer(),
	}

	for _, serializer := range serializers {
		data, err := m.MarshalWith(serializer)
		if err != nil {
			t.Errorf("MarshalWith %s failed: %v", serializer.Name(), err)
			continue
		}

		// 测试反序列化
		restored := NewStringHashMap[int]()
		err = restored.UnmarshalWith(data, serializer)
		if err != nil {
			t.Errorf("UnmarshalWith %s failed: %v", serializer.Name(), err)
			continue
		}

		val, ok := restored.Get("test")
		if !ok || val != 42 {
			t.Errorf("Round-trip failed for %s", serializer.Name())
		}
	}
}

// TestMarshalEmptyMap 测试空map的序列化
func TestMarshalEmptyMap(t *testing.T) {
	m := NewStringHashMap[int]()

	data, err := m.MarshalJSON()
	if err != nil {
		t.Fatalf("Marshal empty map failed: %v", err)
	}

	// 反序列化空数据
	restored := NewStringHashMap[int]()
	err = restored.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("Unmarshal empty map failed: %v", err)
	}

	if !restored.Empty() {
		t.Error("Restored empty map should be empty")
	}
}

// TestMarshalWithGob 测试Gob序列化
func TestMarshalWithGob(t *testing.T) {
	m := NewStringHashMap[int]()
	m.Put("gob", 999)

	serializer := GobSerializer()

	data, err := m.MarshalWith(serializer)
	if err != nil {
		t.Fatalf("Gob marshal failed: %v", err)
	}

	restored := NewStringHashMap[int]()
	err = restored.UnmarshalWith(data, serializer)
	if err != nil {
		t.Fatalf("Gob unmarshal failed: %v", err)
	}

	val, ok := restored.Get("gob")
	if !ok || val != 999 {
		t.Errorf("Gob round-trip failed, got %v, want 999", val)
	}
}

// TestJSONCompatibleInterface 测试JSON兼容接口
func TestJSONCompatibleInterface(t *testing.T) {
	serializers := []*SerializerFunc{
		JsonSerializer(),
		JsoniterSerializer(),
		SonicSerializer(),
		GobSerializer(),
	}

	expectedJSON := []bool{true, true, true, false}

	for i, serializer := range serializers {
		if jsonSerializer, ok := any(serializer).(JSONCompatible); ok {
			isJSON := jsonSerializer.IsJSON()
			if isJSON != expectedJSON[i] {
				t.Errorf("%s IsJSON() returned %v, expected %v",
					serializer.Name(), isJSON, expectedJSON[i])
			}
		} else if expectedJSON[i] {
			t.Errorf("%s should implement JSONCompatible interface", serializer.Name())
		}
	}
}

// TestMarshalInvalidData 测试错误处理
func TestMarshalInvalidData(t *testing.T) {
	// 测试反序列化无效JSON
	m := NewStringHashMap[int]()
	invalidJSON := `{"invalid": json}`

	err := m.UnmarshalJSON([]byte(invalidJSON))
	if err == nil {
		t.Error("Should fail on invalid JSON")
	}

	// 测试反序列化空数据
	err = m.UnmarshalJSON([]byte(""))
	if err == nil {
		t.Error("Should fail on empty data")
	}
}

// TestMarshalCustomTypes 测试自定义类型的序列化
func TestMarshalCustomTypes(t *testing.T) {
	type Key string
	type Value struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	m := NewHashMap[Key, Value]()
	m.Put("key1", Value{Name: "test1", Count: 1})
	m.Put("key2", Value{Name: "test2", Count: 2})

	data, err := m.MarshalJSON()
	if err != nil {
		t.Fatalf("Custom type marshal failed: %v", err)
	}

	restored := NewHashMap[Key, Value]()
	err = restored.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("Custom type unmarshal failed: %v", err)
	}

	val, ok := restored.Get("key1")
	if !ok || val.Name != "test1" || val.Count != 1 {
		t.Errorf("Custom type round-trip failed")
	}
}

// TestMarshalLargeData 测试大数据量的序列化
func TestMarshalLargeData(t *testing.T) {
	m := NewStringHashMap[int]()

	// 插入大量数据，确保键唯一
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		m.Put(key, i)
	}

	data, err := m.MarshalJSON()
	if err != nil {
		t.Fatalf("Large data marshal failed: %v", err)
	}

	restored := NewStringHashMap[int]()
	err = restored.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("Large data unmarshal failed: %v", err)
	}

	if restored.Size() != 1000 {
		t.Errorf("Large data round-trip failed, got %d items, want 1000", restored.Size())
	}
}
