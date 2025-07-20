package cmap

import (
	"testing"
)

func TestConstructors(t *testing.T) {
	// Test string constructors
	t.Run("String constructors", func(t *testing.T) {
		// Test HashMap
		hashMap := NewStringHashMap[int]()
		if hashMap == nil {
			t.Error("NewStringHashMap returned nil")
		}
		hashMap.Put("key", 42)
		if val, ok := hashMap.Get("key"); !ok || val != 42 {
			t.Errorf("HashMap Get failed, got %v, want 42", val)
		}

		// Test TreeMap
		treeMap := NewStringTreeMap[int]()
		if treeMap == nil {
			t.Error("NewStringTreeMap returned nil")
		}
		treeMap.Put("key", 42)
		if val, ok := treeMap.Get("key"); !ok || val != 42 {
			t.Errorf("TreeMap Get failed, got %v, want 42", val)
		}

		// Test LinkedHashMap
		linkedMap := NewStringLinkedHashMap[int]()
		if linkedMap == nil {
			t.Error("NewStringLinkedHashMap returned nil")
		}
		linkedMap.Put("key", 42)
		if val, ok := linkedMap.Get("key"); !ok || val != 42 {
			t.Errorf("LinkedHashMap Get failed, got %v, want 42", val)
		}
	})

	t.Run("Int constructors", func(t *testing.T) {
		// Test HashMap
		hashMap := NewIntHashMap[string]()
		if hashMap == nil {
			t.Error("NewIntHashMap returned nil")
		}
		hashMap.Put(1, "value")
		if val, ok := hashMap.Get(1); !ok || val != "value" {
			t.Errorf("IntHashMap Get failed, got %v, want 'value'", val)
		}

		// Test TreeMap
		treeMap := NewIntTreeMap[string]()
		if treeMap == nil {
			t.Error("NewIntTreeMap returned nil")
		}
		treeMap.Put(1, "value")
		if val, ok := treeMap.Get(1); !ok || val != "value" {
			t.Errorf("IntTreeMap Get failed, got %v, want 'value'", val)
		}

		// Test LinkedHashMap
		linkedMap := NewIntLinkedHashMap[string]()
		if linkedMap == nil {
			t.Error("NewIntLinkedHashMap returned nil")
		}
		linkedMap.Put(1, "value")
		if val, ok := linkedMap.Get(1); !ok || val != "value" {
			t.Errorf("IntLinkedHashMap Get failed, got %v, want 'value'", val)
		}
	})

	t.Run("Generic constructors", func(t *testing.T) {
		// Test generic HashMap
		genHashMap := NewHashMap[string, int]()
		if genHashMap == nil {
			t.Error("NewHashMap returned nil")
		}
		genHashMap.Put("key", 42)
		if val, ok := genHashMap.Get("key"); !ok || val != 42 {
			t.Errorf("Generic HashMap Get failed, got %v, want 42", val)
		}

		// Test generic TreeMap
		genTreeMap := NewTreeMap[string, int]()
		if genTreeMap == nil {
			t.Error("NewTreeMap returned nil")
		}
		genTreeMap.Put("key", 42)
		if val, ok := genTreeMap.Get("key"); !ok || val != 42 {
			t.Errorf("Generic TreeMap Get failed, got %v, want 42", val)
		}

		// Test generic LinkedHashMap
		genLinkedMap := NewLinkedHashMap[string, int]()
		if genLinkedMap == nil {
			t.Error("NewLinkedHashMap returned nil")
		}
		genLinkedMap.Put("key", 42)
		if val, ok := genLinkedMap.Get("key"); !ok || val != 42 {
			t.Errorf("Generic LinkedHashMap Get failed, got %v, want 42", val)
		}
	})

	t.Run("With options", func(t *testing.T) {
		// Test with custom shard count
		customMap := NewStringHashMap[int](WithShardCount(128))
		if customMap == nil {
			t.Error("NewStringHashMap with options returned nil")
		}

		// Test with serialization format option
		customMap2 := NewStringHashMap[int](WithSerializer(GobSerializer()))
		if customMap2 == nil {
			t.Error("NewStringHashMap with serializer option returned nil")
		}
	})
}

func TestInt64Constructors(t *testing.T) {
	// Test int64 constructors
	if m := NewInt64HashMap[string](); m == nil {
		t.Error("NewInt64HashMap returned nil")
	}

	if m := NewInt64TreeMap[string](); m == nil {
		t.Error("NewInt64TreeMap returned nil")
	}

	if m := NewInt64LinkedHashMap[string](); m == nil {
		t.Error("NewInt64LinkedHashMap returned nil")
	}
}

func TestRefactoredConstructorsDemonstration(t *testing.T) {
	// Demonstrate the refactored constructors maintain backward compatibility
	// while being much more concise in implementation

	// Original type-specific constructors still work
	t.Run("TypeSpecificConstructors", func(t *testing.T) {
		stringHashMap := NewStringHashMap[int]()
		if stringHashMap == nil {
			t.Fatal("NewStringHashMap returned nil")
		}
		stringHashMap.Put("key1", 100)
		if val, ok := stringHashMap.Get("key1"); !ok || val != 100 {
			t.Errorf("Expected 100, got %v", val)
		}

		intTreeMap := NewIntTreeMap[string]()
		if intTreeMap == nil {
			t.Fatal("NewIntTreeMap returned nil")
		}
		intTreeMap.Put(42, "value")
		if val, ok := intTreeMap.Get(42); !ok || val != "value" {
			t.Errorf("Expected 'value', got %v", val)
		}

		int64LinkedMap := NewInt64LinkedHashMap[float64]()
		if int64LinkedMap == nil {
			t.Fatal("NewInt64LinkedHashMap returned nil")
		}
		int64LinkedMap.Put(int64(999), 3.14)
		if val, ok := int64LinkedMap.Get(int64(999)); !ok || val != 3.14 {
			t.Errorf("Expected 3.14, got %v", val)
		}
	})

	// Generic constructors work with any comparable type
	t.Run("GenericConstructors", func(t *testing.T) {
		// String keys
		strMap := NewHashMap[string, int]()
		if strMap == nil {
			t.Fatal("NewHashMap returned nil")
		}
		strMap.Put("generic", 200)
		if val, ok := strMap.Get("generic"); !ok || val != 200 {
			t.Errorf("Expected 200, got %v", val)
		}

		// Int keys
		intMap := NewTreeMap[int, string]()
		if intMap == nil {
			t.Fatal("NewTreeMap returned nil")
		}
		intMap.Put(1, "one")
		intMap.Put(3, "three")
		intMap.Put(2, "two")
		if val, ok := intMap.Get(2); !ok || val != "two" {
			t.Errorf("Expected 'two', got %v", val)
		}

		// Custom type keys
		type customKey string
		customMap := NewLinkedHashMap[customKey, bool]()
		if customMap == nil {
			t.Fatal("NewLinkedHashMap returned nil")
		}
		customMap.Put("test", true)
		if val, ok := customMap.Get("test"); !ok || val != true {
			t.Errorf("Expected true, got %v", val)
		}
	})

	// Demonstrate reduced code duplication
	t.Run("CodeReduction", func(t *testing.T) {
		maps := []interface{}{
			NewStringHashMap[int](),
			NewStringTreeMap[int](),
			NewStringLinkedHashMap[int](),
			NewIntHashMap[string](),
			NewIntTreeMap[string](),
			NewIntLinkedHashMap[string](),
			NewInt64HashMap[string](),
			NewInt64TreeMap[string](),
			NewInt64LinkedHashMap[string](),
		}
		for i, m := range maps {
			if m == nil {
				t.Errorf("Constructor %d returned nil", i)
			}
		}
	})
}

func TestNewFunction(t *testing.T) {
	// 测试 New 函数
	t.Run("New_function", func(t *testing.T) {
		cm := New[string, int]()
		if cm == nil {
			t.Fatal("New function returned nil")
		}

		// 测试基本操作
		cm.Put("test", 123)
		if value, found := cm.Get("test"); !found || value != 123 {
			t.Errorf("Expected value 123, got %d, found: %v", value, found)
		}
	})

	// 测试带选项的构造函数
	t.Run("New_with_options", func(t *testing.T) {
		cm := New[string, int](WithShardCount(8))
		if cm == nil {
			t.Fatal("New function with options returned nil")
		}

		// 验证分片数量
		if len(cm.shards) != 8 {
			t.Errorf("Expected 8 shards, got %d", len(cm.shards))
		}
	})

	// 测试不同键值类型
	t.Run("different_types", func(t *testing.T) {
		// 测试 int 键
		cm1 := New[int, string]()
		cm1.Put(1, "one")
		if value, found := cm1.Get(1); !found || value != "one" {
			t.Errorf("Expected 'one', got %s, found: %v", value, found)
		}

		// 测试 float64 键
		cm2 := New[float64, bool]()
		cm2.Put(3.14, true)
		if value, found := cm2.Get(3.14); !found || value != true {
			t.Errorf("Expected true, got %v, found: %v", value, found)
		}
	})
}
