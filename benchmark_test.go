package cmap

import (
	"fmt"
	"sync"
	"testing"
)

// BenchmarkMapOperations 基准测试基本操作
func BenchmarkMapOperations(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		cm := New[string, int]()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cm.Put("key"+string(rune(i)), i)
		}
	})

	b.Run("Get", func(b *testing.B) {
		cm := New[string, int]()
		// 预填充数据
		for i := 0; i < 1000; i++ {
			cm.Put("key"+string(rune(i)), i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cm.Get("key" + string(rune(i%1000)))
		}
	})

	b.Run("Remove", func(b *testing.B) {
		cm := New[string, int]()
		// 预填充数据
		for i := 0; i < b.N; i++ {
			cm.Put("key"+string(rune(i)), i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cm.Remove("key" + string(rune(i)))
		}
	})
}

// BenchmarkConcurrentOperations 并发操作基准测试
func BenchmarkConcurrentOperations(b *testing.B) {
	b.Run("Concurrent_Put", func(b *testing.B) {
		cm := New[string, int]()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				cm.Put("key"+string(rune(i)), i)
				i++
			}
		})
	})

	b.Run("Concurrent_Get", func(b *testing.B) {
		cm := New[string, int]()
		// 预填充数据
		for i := 0; i < 1000; i++ {
			cm.Put("key"+string(rune(i)), i)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				cm.Get("key" + string(rune(i%1000)))
				i++
			}
		})
	})

	b.Run("Concurrent_Mixed", func(b *testing.B) {
		cm := New[string, int]()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				if i%3 == 0 {
					cm.Put("key"+string(rune(i)), i)
				} else if i%3 == 1 {
					cm.Get("key" + string(rune(i%1000)))
				} else {
					cm.Remove("key" + string(rune(i%1000)))
				}
				i++
			}
		})
	})
}

// BenchmarkDifferentShardCounts 不同分片数量的性能对比
func BenchmarkDifferentShardCounts(b *testing.B) {
	shardCounts := []uint32{1, 4, 16, 64, 256}

	for _, shardCount := range shardCounts {
		b.Run(fmt.Sprintf("Shards_%d", shardCount), func(b *testing.B) {
			cm := New[string, int](WithShardCount(shardCount))
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					cm.Put("key"+string(rune(i)), i)
					cm.Get("key" + string(rune(i%1000)))
					i++
				}
			})
		})
	}
}

// BenchmarkMapTypes 不同Map类型的性能对比
func BenchmarkMapTypes(b *testing.B) {
	b.Run("HashMap", func(b *testing.B) {
		cm := NewHashMap[string, int]()
		benchmarkMapType(b, cm)
	})

	b.Run("TreeMap", func(b *testing.B) {
		cm := NewTreeMap[string, int]()
		benchmarkMapType(b, cm)
	})

	b.Run("LinkedHashMap", func(b *testing.B) {
		cm := NewLinkedHashMap[string, int]()
		benchmarkMapType(b, cm)
	})
}

func benchmarkMapType(b *testing.B, cm *Map[string, int]) {
	// 预填充数据
	for i := 0; i < 1000; i++ {
		cm.Put("key"+string(rune(i)), i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cm.Put("key"+string(rune(i)), i)
			cm.Get("key" + string(rune(i%1000)))
			i++
		}
	})
}

// BenchmarkBatchOperations 批量操作基准测试
func BenchmarkBatchOperations(b *testing.B) {
	b.Run("PutAll", func(b *testing.B) {
		cm := New[string, int]()
		data := make(map[string]int, 1000)
		for i := 0; i < 1000; i++ {
			data["key"+string(rune(i))] = i
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cm.PutAll(data)
			cm.Clear()
		}
	})

	b.Run("GetMultiple", func(b *testing.B) {
		cm := New[string, int]()
		// 预填充数据
		for i := 0; i < 1000; i++ {
			cm.Put("key"+string(rune(i)), i)
		}
		keys := make([]string, 100)
		for i := 0; i < 100; i++ {
			keys[i] = "key" + string(rune(i))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cm.GetMultiple(keys)
		}
	})
}

// BenchmarkVsStandardMap 与标准map的性能对比
func BenchmarkVsStandardMap(b *testing.B) {
	b.Run("CMap_Put", func(b *testing.B) {
		cm := New[string, int]()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cm.Put("key"+string(rune(i)), i)
		}
	})

	b.Run("StandardMap_Put", func(b *testing.B) {
		m := make(map[string]int)
		var mu sync.RWMutex
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mu.Lock()
			m["key"+string(rune(i))] = i
			mu.Unlock()
		}
	})

	b.Run("CMap_Get", func(b *testing.B) {
		cm := New[string, int]()
		// 预填充数据
		for i := 0; i < 1000; i++ {
			cm.Put("key"+string(rune(i)), i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cm.Get("key" + string(rune(i%1000)))
		}
	})

	b.Run("StandardMap_Get", func(b *testing.B) {
		m := make(map[string]int)
		var mu sync.RWMutex
		// 预填充数据
		for i := 0; i < 1000; i++ {
			m["key"+string(rune(i))] = i
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mu.RLock()
			_ = m["key"+string(rune(i%1000))]
			mu.RUnlock()
		}
	})
}
