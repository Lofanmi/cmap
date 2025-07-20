package cmap

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestConcurrentStress 高强度并发压力测试
func TestConcurrentStress(t *testing.T) {
	cm := New[string, int]()
	const numGoroutines = 100
	const operationsPerGoroutine = 1000

	var wg sync.WaitGroup
	var totalOperations int64

	// 启动多个goroutine进行并发操作
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := "key" + string(rune(id)) + "_" + string(rune(j))
				value := id*operationsPerGoroutine + j

				// 随机操作类型
				switch j % 4 {
				case 0:
					cm.Put(key, value)
					atomic.AddInt64(&totalOperations, 1)
				case 1:
					cm.Get(key)
					atomic.AddInt64(&totalOperations, 1)
				case 2:
					cm.Remove(key)
					atomic.AddInt64(&totalOperations, 1)
				case 3:
					cm.Size()
					atomic.AddInt64(&totalOperations, 1)
				}
			}
		}(i)
	}

	wg.Wait()

	// 验证最终状态
	finalSize := cm.Size()
	t.Logf("Total operations: %d, Final size: %d", totalOperations, finalSize)
}

// TestConcurrentReadWrite 读写并发测试
func TestConcurrentReadWrite(t *testing.T) {
	cm := New[string, int]()
	const numReaders = 50
	const numWriters = 10
	const duration = 100 * time.Millisecond

	var wg sync.WaitGroup
	var stop int32

	// 启动写者
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			j := 0
			for atomic.LoadInt32(&stop) == 0 {
				key := "writer" + string(rune(id)) + "_" + string(rune(j))
				cm.Put(key, j)
				j++
			}
		}(i)
	}

	// 启动读者
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for atomic.LoadInt32(&stop) == 0 {
				key := "writer" + string(rune(id%numWriters)) + "_" + string(rune(id))
				cm.Get(key)
			}
		}(i)
	}

	// 运行一段时间
	time.Sleep(duration)
	atomic.StoreInt32(&stop, 1)
	wg.Wait()

	t.Logf("Final map size: %d", cm.Size())
}

// TestConcurrentClear 并发清空测试
func TestConcurrentClear(t *testing.T) {
	cm := New[string, int]()
	const numGoroutines = 20

	// 预填充数据
	for i := 0; i < 1000; i++ {
		cm.Put("key"+string(rune(i)), i)
	}

	var wg sync.WaitGroup
	var clearCount int32

	// 多个goroutine同时调用Clear
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cm.Clear()
			atomic.AddInt32(&clearCount, 1)
		}()
	}

	wg.Wait()

	// 验证最终状态
	finalSize := cm.Size()
	if finalSize != 0 {
		t.Errorf("Expected empty map after concurrent clear, got size: %d", finalSize)
	}
	t.Logf("Clear operations: %d", clearCount)
}

// TestConcurrentIteration 并发迭代测试
func TestConcurrentIteration(t *testing.T) {
	cm := New[string, int]()
	const numItems = 1000
	const numIterators = 10

	// 预填充数据
	for i := 0; i < numItems; i++ {
		cm.Put("key"+string(rune(i)), i)
	}

	var wg sync.WaitGroup
	var totalIterated int64

	// 多个goroutine同时获取Keys和Values
	for i := 0; i < numIterators; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			keys := cm.Keys()
			values := cm.Values()
			atomic.AddInt64(&totalIterated, int64(len(keys)+len(values)))
		}()
	}

	wg.Wait()

	// 验证结果
	expectedTotal := int64(numItems * numIterators * 2) // Keys + Values
	if totalIterated != expectedTotal {
		t.Errorf("Expected %d total iterations, got %d", expectedTotal, totalIterated)
	}
}

// TestConcurrentBatchOperations 并发批量操作测试
func TestConcurrentBatchOperations(t *testing.T) {
	cm := New[string, int]()
	const numGoroutines = 20
	const batchSize = 100

	var wg sync.WaitGroup

	// 并发批量添加
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			data := make(map[string]int, batchSize)
			for j := 0; j < batchSize; j++ {
				key := "batch" + string(rune(id)) + "_" + string(rune(j))
				data[key] = j
			}
			cm.PutAll(data)
		}(i)
	}

	wg.Wait()

	// 验证结果
	expectedSize := numGoroutines * batchSize
	actualSize := cm.Size()
	if actualSize != expectedSize {
		t.Errorf("Expected size %d after batch operations, got %d", expectedSize, actualSize)
	}

	// 并发批量获取
	keys := make([]string, batchSize)
	for i := 0; i < batchSize; i++ {
		keys[i] = "batch0_" + string(rune(i))
	}

	var results []map[string]int
	var mu sync.Mutex

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := cm.GetMultiple(keys)
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}()
	}

	wg.Wait()

	// 验证批量获取结果
	if len(results) != numGoroutines {
		t.Errorf("Expected %d batch get results, got %d", numGoroutines, len(results))
	}
}

// TestConcurrentSerialization 并发序列化测试
func TestConcurrentSerialization(t *testing.T) {
	cm := New[string, int]()
	const numGoroutines = 10

	// 预填充数据
	for i := 0; i < 100; i++ {
		cm.Put("key"+string(rune(i)), i)
	}

	var wg sync.WaitGroup
	var serializationCount int32

	// 并发序列化
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data, err := cm.MarshalJSON()
			if err != nil {
				t.Errorf("Serialization failed: %v", err)
				return
			}
			if len(data) == 0 {
				t.Error("Serialized data is empty")
				return
			}
			atomic.AddInt32(&serializationCount, 1)
		}()
	}

	wg.Wait()

	if serializationCount != int32(numGoroutines) {
		t.Errorf("Expected %d successful serializations, got %d", numGoroutines, serializationCount)
	}
}

// TestConcurrentFileOperationsStress 并发文件操作压力测试
func TestConcurrentFileOperationsStress(t *testing.T) {
	cm := New[string, int](WithSerializer(JsonSerializer()))
	const numGoroutines = 5

	// 预填充数据
	for i := 0; i < 100; i++ {
		cm.Put("key"+string(rune(i)), i)
	}

	var wg sync.WaitGroup
	var saveCount int32

	// 并发保存到不同文件
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			filename := fmt.Sprintf("test_concurrent_%d.json", id)
			err := cm.SaveToFile(filename)
			if err != nil {
				t.Errorf("Failed to save file %s: %v", filename, err)
				return
			}
			atomic.AddInt32(&saveCount, 1)
		}(i)
	}

	wg.Wait()

	if saveCount != int32(numGoroutines) {
		t.Errorf("Expected %d successful saves, got %d", numGoroutines, saveCount)
	}
}

// TestConcurrentEdgeCases 并发边界条件测试
func TestConcurrentEdgeCases(t *testing.T) {
	cm := New[string, int](WithShardCount(1)) // 使用单个分片增加竞争

	const numGoroutines = 100
	const operations = 1000

	var wg sync.WaitGroup
	var putCount, getCount, removeCount int64

	// 并发操作相同的键
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				key := "contested_key"
				value := j

				switch j % 3 {
				case 0:
					cm.Put(key, value)
					atomic.AddInt64(&putCount, 1)
				case 1:
					cm.Get(key)
					atomic.AddInt64(&getCount, 1)
				case 2:
					cm.Remove(key)
					atomic.AddInt64(&removeCount, 1)
				}
			}
		}()
	}

	wg.Wait()

	t.Logf("Operations completed - Put: %d, Get: %d, Remove: %d", putCount, getCount, removeCount)
	t.Logf("Final map size: %d", cm.Size())
}

// TestConcurrentMemoryLeak 并发内存泄漏测试
func TestConcurrentMemoryLeak(t *testing.T) {
	const iterations = 10
	const numGoroutines = 50
	const operationsPerIteration = 1000

	for iter := 0; iter < iterations; iter++ {
		cm := New[string, int]()

		var wg sync.WaitGroup

		// 创建大量goroutine进行并发操作
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < operationsPerIteration; j++ {
					key := "key" + string(rune(id)) + "_" + string(rune(j))
					value := j

					cm.Put(key, value)
					cm.Get(key)
					if j%10 == 0 {
						cm.Remove(key)
					}
				}
			}(i)
		}

		wg.Wait()

		// 验证最终状态
		finalSize := cm.Size()
		expectedSize := numGoroutines * operationsPerIteration * 9 / 10 // 90% of operations (10% removed)
		if finalSize != expectedSize {
			t.Errorf("Iteration %d: Expected size %d, got %d", iter, expectedSize, finalSize)
		}
	}
}
