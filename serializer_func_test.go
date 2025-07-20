package cmap

import (
	"fmt"
	"reflect"
	"testing"
)

// TestSerializerFunc 测试序列化函数包装器
func TestSerializerFunc(t *testing.T) {
	serializer := JsonSerializer()

	// 测试基本功能
	if serializer.Name() != "json" {
		t.Errorf("Expected name 'json', got '%s'", serializer.Name())
	}

	// 测试序列化
	testData := map[string]interface{}{
		"test":  123,
		"hello": "world",
	}

	data, err := serializer.Marshal(testData)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("Marshal returned empty data")
	}

	// 测试反序列化
	var restored map[string]interface{}
	err = serializer.Unmarshal(data, &restored)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if restored["test"] != 123.0 { // JSON解码float64
		t.Errorf("Unmarshal data mismatch")
	}
}

// TestJsonSerializer 测试标准JSON序列化器
func TestJsonSerializer(t *testing.T) {
	serializer := JsonSerializer()

	if serializer.Name() != "json" {
		t.Errorf("Expected name 'json', got '%s'", serializer.Name())
	}

	// 测试JSON兼容性接口
	if jsonSerializer, ok := any(serializer).(JSONCompatible); ok {
		if !jsonSerializer.IsJSON() {
			t.Error("JsonSerializer should return true for IsJSON()")
		}
	} else {
		t.Error("JsonSerializer should implement JSONCompatible interface")
	}
}

// TestJsoniterSerializer 测试Jsoniter序列化器
func TestJsoniterSerializer(t *testing.T) {
	serializer := JsoniterSerializer()

	if serializer.Name() != "jsoniter" {
		t.Errorf("Expected name 'jsoniter', got '%s'", serializer.Name())
	}

	// 测试JSON兼容性接口
	if jsonSerializer, ok := any(serializer).(JSONCompatible); ok {
		if !jsonSerializer.IsJSON() {
			t.Error("JsoniterSerializer should return true for IsJSON()")
		}
	} else {
		t.Error("JsoniterSerializer should implement JSONCompatible interface")
	}
}

// TestSonicSerializer 测试Sonic序列化器
func TestSonicSerializer(t *testing.T) {
	serializer := SonicSerializer()

	if serializer.Name() != "sonic" {
		t.Errorf("Expected name 'sonic', got '%s'", serializer.Name())
	}

	// 测试JSON兼容性接口
	if jsonSerializer, ok := any(serializer).(JSONCompatible); ok {
		if !jsonSerializer.IsJSON() {
			t.Error("SonicSerializer should return true for IsJSON()")
		}
	} else {
		t.Error("SonicSerializer should implement JSONCompatible interface")
	}
}

// TestGobSerializer 测试Gob序列化器
func TestGobSerializer(t *testing.T) {
	serializer := GobSerializer()

	if serializer.Name() != "gob" {
		t.Errorf("Expected name 'gob', got '%s'", serializer.Name())
	}

	// 测试Gob序列化器不实现JSON兼容性接口
	if jsonSerializer, ok := any(serializer).(JSONCompatible); ok {
		if jsonSerializer.IsJSON() {
			t.Error("GobSerializer should return false for IsJSON()")
		}
	} else {
		t.Error("GobSerializer should implement JSONCompatible interface but return false")
	}
}

// TestSerializerRoundTrip 测试各序列化器往返
func TestSerializerRoundTrip(t *testing.T) {
	serializers := []*SerializerFunc{
		JsonSerializer(),
		JsoniterSerializer(),
		SonicSerializer(),
		GobSerializer(),
	}

	testData := SerializableData[string, int]{
		Items: []Tuple[string, int]{
			{Key: "test1", Value: 100},
			{Key: "test2", Value: 200},
			{Key: "test3", Value: 300},
		},
	}

	for _, serializer := range serializers {
		t.Run(serializer.Name(), func(t *testing.T) {
			// 序列化
			data, err := serializer.Marshal(testData)
			if err != nil {
				t.Fatalf("%s marshal failed: %v", serializer.Name(), err)
			}

			if len(data) == 0 {
				t.Error("Marshal returned empty data")
			}

			// 反序列化
			var restored SerializableData[string, int]
			err = serializer.Unmarshal(data, &restored)
			if err != nil {
				t.Fatalf("%s unmarshal failed: %v", serializer.Name(), err)
			}

			// 验证数据一致性
			if len(restored.Items) != len(testData.Items) {
				t.Errorf("%s: item count mismatch", serializer.Name())
			}

			for i, item := range testData.Items {
				if i < len(restored.Items) {
					if restored.Items[i].Key != item.Key || restored.Items[i].Value != item.Value {
						t.Errorf("%s: data mismatch at index %d", serializer.Name(), i)
					}
				}
			}
		})
	}
}

// TestSerializerErrorHandling 测试错误处理
func TestSerializerErrorHandling(t *testing.T) {
	serializers := []*SerializerFunc{
		JsonSerializer(),
		JsoniterSerializer(),
		SonicSerializer(),
		GobSerializer(),
	}

	for _, serializer := range serializers {
		t.Run(serializer.Name(), func(t *testing.T) {
			// 测试无效数据反序列化
			var result interface{}
			err := serializer.Unmarshal([]byte("invalid data"), &result)
			if err == nil {
				t.Errorf("%s should fail on invalid data", serializer.Name())
			}

			// 测试nil目标
			err = serializer.Unmarshal([]byte("{}"), nil)
			if err == nil {
				t.Errorf("%s should fail on nil target", serializer.Name())
			}
		})
	}
}

// TestSerializerEmptyData 测试空数据处理
func TestSerializerEmptyData(t *testing.T) {
	serializers := []*SerializerFunc{
		JsonSerializer(),
		JsoniterSerializer(),
		SonicSerializer(),
		GobSerializer(),
	}

	for _, serializer := range serializers {
		t.Run(serializer.Name(), func(t *testing.T) {
			// 测试空结构序列化
			emptyData := SerializableData[string, int]{}
			data, err := serializer.Marshal(emptyData)
			if err != nil {
				t.Fatalf("%s marshal empty failed: %v", serializer.Name(), err)
			}

			// 测试空结构反序列化
			var restored SerializableData[string, int]
			err = serializer.Unmarshal(data, &restored)
			if err != nil {
				t.Fatalf("%s unmarshal empty failed: %v", serializer.Name(), err)
			}

			if len(restored.Items) != 0 {
				t.Errorf("%s: expected empty result, got %d items",
					serializer.Name(), len(restored.Items))
			}
		})
	}
}

// TestSerializerCustomTypes 测试自定义类型序列化
func TestSerializerCustomTypes(t *testing.T) {
	type CustomKey string
	type CustomValue struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	serializers := []*SerializerFunc{
		JsonSerializer(),
		JsoniterSerializer(),
		SonicSerializer(),
	}

	for _, serializer := range serializers {
		t.Run(serializer.Name(), func(t *testing.T) {
			data := SerializableData[CustomKey, CustomValue]{
				Items: []Tuple[CustomKey, CustomValue]{
					{Key: "key1", Value: CustomValue{Name: "test1", Count: 1}},
					{Key: "key2", Value: CustomValue{Name: "test2", Count: 2}},
				},
			}

			// 序列化
			serialized, err := serializer.Marshal(data)
			if err != nil {
				t.Fatalf("%s marshal custom types failed: %v", serializer.Name(), err)
			}

			// 反序列化
			var restored SerializableData[CustomKey, CustomValue]
			err = serializer.Unmarshal(serialized, &restored)
			if err != nil {
				t.Fatalf("%s unmarshal custom types failed: %v", serializer.Name(), err)
			}

			// 验证数据
			if len(restored.Items) != 2 {
				t.Errorf("%s: expected 2 items, got %d", serializer.Name(), len(restored.Items))
			}

			if restored.Items[0].Value.Name != "test1" || restored.Items[0].Value.Count != 1 {
				t.Errorf("%s: custom type data mismatch", serializer.Name())
			}
		})
	}
}

// TestGobSerializerComplexTypes 测试Gob序列化复杂类型
func TestGobSerializerComplexTypes(t *testing.T) {
	serializer := GobSerializer()

	// 测试复杂结构
	type ComplexValue struct {
		Data    map[string]int
		Numbers []int
		Flag    bool
	}

	data := SerializableData[string, ComplexValue]{
		Items: []Tuple[string, ComplexValue]{
			{
				Key: "complex",
				Value: ComplexValue{
					Data:    map[string]int{"a": 1, "b": 2},
					Numbers: []int{1, 2, 3, 4, 5},
					Flag:    true,
				},
			},
		},
	}

	// 序列化
	serialized, err := serializer.Marshal(data)
	if err != nil {
		t.Fatalf("Gob marshal complex failed: %v", err)
	}

	// 反序列化
	var restored SerializableData[string, ComplexValue]
	err = serializer.Unmarshal(serialized, &restored)
	if err != nil {
		t.Fatalf("Gob unmarshal complex failed: %v", err)
	}

	// 验证复杂数据
	if len(restored.Items) != 1 {
		t.Fatal("Expected 1 item")
	}

	restoredValue := restored.Items[0].Value
	if !reflect.DeepEqual(restoredValue.Data, data.Items[0].Value.Data) {
		t.Error("Complex map data mismatch")
	}
	if !reflect.DeepEqual(restoredValue.Numbers, data.Items[0].Value.Numbers) {
		t.Error("Complex slice data mismatch")
	}
	if restoredValue.Flag != data.Items[0].Value.Flag {
		t.Error("Complex bool data mismatch")
	}
}

// TestSerializerConsistency 测试序列化器一致性
func TestSerializerConsistency(t *testing.T) {
	serializers := []*SerializerFunc{
		JsonSerializer(),
		JsoniterSerializer(),
		SonicSerializer(),
		GobSerializer(),
	}

	for _, serializer := range serializers {
		t.Run(serializer.Name(), func(t *testing.T) {
			data := SerializableData[string, int]{
				Items: []Tuple[string, int]{
					{Key: "consistency", Value: 123},
				},
			}

			// 多次序列化应该得到相同结果（对于确定性序列化器）
			data1, err1 := serializer.Marshal(data)
			data2, err2 := serializer.Marshal(data)

			if err1 != nil || err2 != nil {
				t.Fatalf("Consistency test failed: %v, %v", err1, err2)
			}

			// JSON序列化器应该是确定性的
			if serializer.Name() == "json" || serializer.Name() == "jsoniter" {
				if string(data1) != string(data2) {
					t.Errorf("%s: inconsistent serialization", serializer.Name())
				}
			}
		})
	}
}

// TestSerializerNameConsistency 测试名称一致性
func TestSerializerNameConsistency(t *testing.T) {
	serializers := map[string]*SerializerFunc{
		"json":     JsonSerializer(),
		"jsoniter": JsoniterSerializer(),
		"sonic":    SonicSerializer(),
		"gob":      GobSerializer(),
	}

	for expectedName, serializer := range serializers {
		if serializer.Name() != expectedName {
			t.Errorf("Expected name %s, got %s", expectedName, serializer.Name())
		}
	}
}

// TestSerializerNilSafety 测试nil安全
func TestSerializerNilSafety(t *testing.T) {
	serializers := []*SerializerFunc{
		JsonSerializer(),
		JsoniterSerializer(),
		SonicSerializer(),
		GobSerializer(),
	}

	for _, serializer := range serializers {
		t.Run(serializer.Name(), func(t *testing.T) {
			// 测试nil数据序列化
			_, err := serializer.Marshal(nil)
			if err == nil {
				t.Logf("%s: nil serialization succeeded (may be expected)", serializer.Name())
			}

			// 测试空数据反序列化
			var result interface{}
			err = serializer.Unmarshal([]byte("{}"), &result)
			if err != nil {
				t.Logf("%s: empty object deserialization failed: %v", serializer.Name(), err)
			}
		})
	}
}

// BenchmarkSerializers 性能基准测试
func BenchmarkSerializers(b *testing.B) {
	const size = 10000
	data := SerializableData[string, int]{
		Items: make([]Tuple[string, int], size),
	}

	for i := 0; i < size; i++ {
		data.Items[i] = Tuple[string, int]{Key: fmt.Sprintf("key%d", i), Value: i}
	}

	serializers := []*SerializerFunc{
		JsonSerializer(),
		JsoniterSerializer(),
		SonicSerializer(),
		GobSerializer(),
	}

	for _, serializer := range serializers {
		b.Run(serializer.Name(), func(b *testing.B) {
			b.Run("Marshal", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = serializer.Marshal(data)
				}
			})

			b.Run("Unmarshal", func(b *testing.B) {
				serialized, _ := serializer.Marshal(data)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					var restored SerializableData[string, int]
					_ = serializer.Unmarshal(serialized, &restored)
				}
			})
		})
	}
}

// TestSerializerLargeData 测试大数据处理
func TestSerializerLargeData(t *testing.T) {
	serializers := []*SerializerFunc{
		JsonSerializer(),
		JsoniterSerializer(),
		SonicSerializer(),
		GobSerializer(),
	}

	data := SerializableData[string, int]{
		Items: make([]Tuple[string, int], 1000),
	}

	for i := 0; i < 1000; i++ {
		data.Items[i] = Tuple[string, int]{Key: fmt.Sprintf("large_key_%d", i), Value: i}
	}

	for _, serializer := range serializers {
		t.Run(serializer.Name(), func(t *testing.T) {
			// 序列化大数据
			serialized, err := serializer.Marshal(data)
			if err != nil {
				t.Fatalf("%s marshal large data failed: %v", serializer.Name(), err)
			}

			// 反序列化大数据
			var restored SerializableData[string, int]
			err = serializer.Unmarshal(serialized, &restored)
			if err != nil {
				t.Fatalf("%s unmarshal large data failed: %v", serializer.Name(), err)
			}

			if len(restored.Items) != len(data.Items) {
				t.Errorf("%s: large data item count mismatch", serializer.Name())
			}
		})
	}
}
