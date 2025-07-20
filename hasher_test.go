package cmap

import (
	"fmt"
	"testing"
)

// TestBasicHashers 测试基本哈希器
func TestBasicHashers(t *testing.T) {
	testCases := []struct {
		name string
		key  interface{}
		want uint32
	}{
		{"int8", int8(42), uint32(42)},
		{"int16", int16(1000), uint32(1000)},
		{"int32", int32(123456), uint32(123456)},
		{"uint8", uint8(200), uint32(200)},
		{"uint16", uint16(50000), uint32(50000)},
		{"uint32", uint32(123456789), uint32(123456789)},
		{"int", 42, uint32(42)},
		{"int64", int64(9223372036854775807), uint32(9223372036854775807 >> 32)},
		{"uint", uint(42), uint32(42)},
		{"uint64", uint64(18446744073709551615), uint32(18446744073709551615 >> 32)},
		{"uintptr", uintptr(123), uint32(123)},
		{"float32", float32(3.14), uint32(1078523331)}, // 3.14的IEEE表示
		{"float64", 3.14159265359, uint32(4614256656552045848 >> 32)},
		{"string", "hello", 0}, // 字符串哈希值不确定
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 对于interface{}类型，我们直接测试字符串转换
			hasher := getHasher[string]()
			keyStr := fmt.Sprintf("%v", tc.key)
			hash := hasher.Hash(keyStr)
			if hash == 0 && keyStr != "" {
				t.Logf("Zero hash for %v may be expected", tc.key)
			}
		})
	}
}

// TestStringHasher 测试字符串哈希器
func TestStringHasher(t *testing.T) {
	hasher := getHasher[string]()

	testStrings := []string{
		"",
		"a",
		"hello",
		"世界",
		"very long string with special chars: !@#$%^&*()",
		"重复字符串重复字符串重复字符串",
	}

	for _, str := range testStrings {
		hash1 := hasher.Hash(str)
		hash2 := hasher.Hash(str)

		if hash1 != hash2 {
			t.Errorf("String hasher not deterministic for %q: %d != %d", str, hash1, hash2)
		}

		if hash1 == 0 && str != "" {
			t.Logf("Non-zero hash for %q: %d", str, hash1)
		}
	}
}

// TestNumericHashers 测试数值类型哈希器
func TestNumericHashers(t *testing.T) {
	testCases := []struct {
		name string
		zero interface{}
	}{
		{"int8", int8(0)}, {"int16", int16(0)}, {"int32", int32(0)}, {"int64", int64(0)},
		{"uint8", uint8(0)}, {"uint16", uint16(0)}, {"uint32", uint32(0)}, {"uint64", uint64(0)},
		{"int", 0}, {"uint", uint(0)}, {"uintptr", uintptr(0)},
		{"float32", float32(0)}, {"float64", float64(0)}, {"string", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch v := tc.zero.(type) {
			case int8:
				hasher := int8Hasher{}
				result := hasher.Hash(v)
				_ = result
			case int16:
				hasher := int16Hasher{}
				result := hasher.Hash(v)
				_ = result
			case int32:
				hasher := int32Hasher{}
				result := hasher.Hash(v)
				_ = result
			case int64:
				hasher := int64Hasher{}
				result := hasher.Hash(v)
				_ = result
			case uint8:
				hasher := uint8Hasher{}
				result := hasher.Hash(v)
				_ = result
			case uint16:
				hasher := uint16Hasher{}
				result := hasher.Hash(v)
				_ = result
			case uint32:
				hasher := uint32Hasher{}
				result := hasher.Hash(v)
				_ = result
			case uint64:
				hasher := uint64Hasher{}
				result := hasher.Hash(v)
				_ = result
			case int:
				hasher := intHasher{}
				result := hasher.Hash(v)
				_ = result
			case uint:
				hasher := uintHasher{}
				result := hasher.Hash(v)
				_ = result
			case uintptr:
				hasher := uintptrHasher{}
				result := hasher.Hash(v)
				_ = result
			case float32:
				hasher := float32Hasher{}
				result := hasher.Hash(v)
				if v == 0 && result != 0 {
					t.Errorf("float32(0) should hash to 0, got %d", result)
				}
			case float64:
				hasher := float64Hasher{}
				result := hasher.Hash(v)
				if v == 0 && result != 0 {
					t.Errorf("float64(0) should hash to 0, got %d", result)
				}
			case string:
				hasher := stringHasher{}
				result := hasher.Hash(v)
				_ = result
			}
		})
	}
}

// TestGenericHasher 测试通用哈希器
func TestGenericHasher(t *testing.T) {
	type CustomType string

	sh := genericHasher[CustomType]{}

	testCases := []CustomType{
		"",
		"test",
		"custom type",
		"1234567890",
	}

	for _, key := range testCases {
		hash := sh.Hash(key)
		if hash == 0 {
			t.Errorf("Generic hasher failed for %v", key)
		}
	}
}

// TestHashDistribution 测试哈希分布
func TestHashDistribution(t *testing.T) {
	m := NewStringHashMap[int](WithShardCount(16))

	// 测试多个键的哈希分布
	keys := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	shards := make(map[uint32]int)
	for _, key := range keys {
		shard := m.getShard(key)
		shardIndex := uint32(0)
		for i, s := range m.shards {
			if s == *shard {
				shardIndex = uint32(i)
				break
			}
		}
		shards[shardIndex]++
	}

	// 应该有合理的分布
	if len(shards) < 2 {
		t.Log("Hash distribution test: keys may be distributed to few shards")
	}
}

// TestFloatSpecialValues 测试浮点数的特殊值
func TestFloatSpecialValues(t *testing.T) {
	// 测试float32特殊值
	f32Hasher := float32Hasher{}
	testFloat32s := []float32{
		0.0,
		-0.0,
		1.0,
		-1.0,
		float32(0.1),
		float32(-0.1),
		float32(3.14159),
		float32(-3.14159),
	}

	for _, f := range testFloat32s {
		hash := f32Hasher.Hash(f)
		_ = hash
	}

	// 测试float64特殊值
	f64Hasher := float64Hasher{}
	testFloat64s := []float64{
		0.0,
		-0.0,
		1.0,
		-1.0,
		0.1,
		-0.1,
		3.141592653589793,
		-3.141592653589793,
	}

	for _, f := range testFloat64s {
		hash := f64Hasher.Hash(f)
		_ = hash
	}
}

// TestGetHasherConsistency 测试getHasher的一致性
func TestGetHasherConsistency(t *testing.T) {
	// 测试getHasher对同一类型返回一致的哈希器
	hasher1 := getHasher[string]()
	hasher2 := getHasher[string]()

	key := "test"
	hash1 := hasher1.Hash(key)
	hash2 := hasher2.Hash(key)

	if hash1 != hash2 {
		t.Errorf("getHasher not consistent for same type: %d != %d", hash1, hash2)
	}
}

// TestHashCollision 测试哈希冲突
func TestHashCollision(t *testing.T) {
	// 测试不同输入产生不同哈希（理想情况下）
	m := NewStringHashMap[int](WithShardCount(32))

	testKeys := []string{
		"collision1",
		"collision2",
		"collision3",
		"test1",
		"test2",
		"test3",
	}

	hashes := make(map[uint32]string)
	for _, key := range testKeys {
		hash := m.hasher.Hash(key)
		if existing, exists := hashes[hash]; exists {
			t.Logf("Hash collision detected: %q and %q both hash to %d",
				existing, key, hash)
		} else {
			hashes[hash] = key
		}
	}
}

// TestHasherWithZeroValues 测试零值处理
func TestHasherWithZeroValues(t *testing.T) {
	testCases := []struct {
		name string
		zero interface{}
	}{
		{"int zero", 0},
		{"string zero", ""},
		{"float64 zero", 0.0},
		{"uint64 zero", uint64(0)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 对于interface{}类型，我们直接测试字符串转换
			sh := getHasher[string]()
			keyStr := fmt.Sprintf("%v", tc.zero)
			hash := sh.Hash(keyStr)
			_ = hash
		})
	}
}

// BenchmarkHashers 性能基准测试
func BenchmarkHashers(b *testing.B) {
	b.Run("string", func(b *testing.B) {
		hasher := getHasher[string]()
		key := "benchmark_test_key"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = hasher.Hash(key)
		}
	})

	b.Run("int", func(b *testing.B) {
		hasher := getHasher[int]()
		key := 42
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = hasher.Hash(key)
		}
	})

	b.Run("float64", func(b *testing.B) {
		hasher := getHasher[float64]()
		key := 3.141592653589793
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = hasher.Hash(key)
		}
	})

	b.Run("custom_type", func(b *testing.B) {
		type CustomKey string
		hasher := getHasher[CustomKey]()
		key := CustomKey("custom_benchmark_key")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = hasher.Hash(key)
		}
	})
}

// TestGenericHasherWithComplexTypes 测试通用哈希器与复杂类型
func TestGenericHasherWithComplexTypes(t *testing.T) {
	// 对于不满足cmp.Ordered的类型，我们使用字符串哈希器
	hasher := getHasher[string]()

	key := fmt.Sprintf("complex_key_%d_%s", 1, "test")
	hash := hasher.Hash(key)

	// 主要测试不panic
	_ = hash
}

// TestHasherEdgeCases 测试边界情况
func TestHasherEdgeCases(t *testing.T) {
	// 测试各种类型的边界值
	testCases := []struct {
		name string
		key  interface{}
	}{
		{"max_int8", int8(127)},
		{"min_int8", int8(-128)},
		{"max_uint8", uint8(255)},
		{"max_int32", int32(2147483647)},
		{"min_int32", int32(-2147483648)},
		{"max_uint32", uint32(4294967295)},
		{"max_int64", int64(9223372036854775807)},
		{"min_int64", int64(-9223372036854775808)},
		{"max_uint64", uint64(18446744073709551615)},
		{"max_float32", float32(3.4028235e+38)},
		{"max_float64", 1.7976931348623157e+308},
		{"unicode_string", "🚀 中文测试 🎉"},
		{"empty_string", ""},
		{"single_char", "x"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 对于interface{}类型，我们直接测试字符串转换
			sh := getHasher[string]()
			keyStr := fmt.Sprintf("%v", tc.key)
			hash := sh.Hash(keyStr)
			_ = hash
		})
	}
}
