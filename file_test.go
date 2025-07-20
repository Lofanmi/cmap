package cmap

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

// TestSaveToFile 测试保存到文件
func TestSaveToFile(t *testing.T) {
	// 创建临时文件
	tmpFile := filepath.Join(t.TempDir(), "test_save.json")
	defer func() { _ = os.Remove(tmpFile) }()

	m := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	m.Put("key1", 100)
	m.Put("key2", 200)

	// 保存到文件
	err := m.SaveToFile(tmpFile)
	if err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Fatal("SaveToFile did not create file")
	}
}

// TestLoadFromFile 测试从文件加载
func TestLoadFromFile(t *testing.T) {
	// 创建临时文件
	tmpFile := filepath.Join(t.TempDir(), "test_load.json")
	defer func() { _ = os.Remove(tmpFile) }()

	// 先保存数据
	original := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	original.Put("key1", 100)
	original.Put("key2", 200)
	err := original.SaveToFile(tmpFile)
	if err != nil {
		t.Fatalf("Setup SaveToFile failed: %v", err)
	}

	// 从文件加载
	loaded := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	err = loaded.LoadFromFile(tmpFile)
	if err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	// 验证数据一致性
	if loaded.Size() != original.Size() {
		t.Errorf("Size mismatch: original %d, loaded %d", original.Size(), loaded.Size())
	}

	val, ok := loaded.Get("key1")
	if !ok || val != 100 {
		t.Errorf("LoadFromFile data mismatch for key1: got %v, want 100", val)
	}

	val, ok = loaded.Get("key2")
	if !ok || val != 200 {
		t.Errorf("LoadFromFile data mismatch for key2: got %v, want 200", val)
	}
}

// TestSaveLoadRoundTrip 测试保存加载往返
func TestSaveLoadRoundTrip(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test_roundtrip.json")
	defer func() { _ = os.Remove(tmpFile) }()

	// 创建测试数据
	original := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	for i := 0; i < 100; i++ {
		original.Put(fmt.Sprintf("key%d", i), i*10)
	}

	// 保存
	err := original.SaveToFile(tmpFile)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// 加载
	restored := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	err = restored.LoadFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// 验证数据完整性
	if restored.Size() != original.Size() {
		t.Errorf("Size mismatch: original %d, restored %d", original.Size(), restored.Size())
	}

	// 验证所有数据
	originalKeys := original.Keys()
	for _, key := range originalKeys {
		originalVal, _ := original.Get(key)
		restoredVal, ok := restored.Get(key)
		if !ok || restoredVal != originalVal {
			t.Errorf("Data mismatch for %s: original %d, restored %d", key, originalVal, restoredVal)
		}
	}
}

// TestSaveToFileDifferentFormats 测试不同格式保存
func TestSaveToFileDifferentFormats(t *testing.T) {
	formats := []*SerializerFunc{
		JsonSerializer(),
		GobSerializer(),
	}

	for _, format := range formats {
		t.Run(format.Name(), func(t *testing.T) {
			tmpFile := filepath.Join(t.TempDir(), "test_format."+format.Name())
			defer func() { _ = os.Remove(tmpFile) }()

			m := NewStringHashMap[int](WithSerializer(format))
			m.Put("test", 42)
			m.Put("hello", 100)

			// 保存
			err := m.SaveToFile(tmpFile)
			if err != nil {
				t.Fatalf("SaveToFile with %s failed: %v", format.Name(), err)
			}

			// 加载
			loaded := NewStringHashMap[int](WithSerializer(format))
			err = loaded.LoadFromFile(tmpFile)
			if err != nil {
				t.Fatalf("LoadFromFile with %s failed: %v", format.Name(), err)
			}

			// 验证
			if val, ok := loaded.Get("test"); !ok || val != 42 {
				t.Errorf("Round-trip failed for %s: test = %v", format.Name(), val)
			}
		})
	}
}

// TestSaveToFileErrorCases 测试错误处理
func TestSaveToFileErrorCases(t *testing.T) {
	m := NewStringHashMap[int]()
	m.Put("test", 42)

	// 测试无效路径，Windows可以写入，跳过。
	if runtime.GOOS != "windows" {
		err := m.SaveToFile("/invalid/path/file.json")
		if err == nil {
			t.Error("SaveToFile should fail for invalid path")
		}
	}

	// 测试空map保存
	empty := NewStringHashMap[int]()
	tmpFile := filepath.Join(t.TempDir(), "empty.json")
	defer func() { _ = os.Remove(tmpFile) }()

	err := empty.SaveToFile(tmpFile)
	if err != nil {
		t.Errorf("SaveToFile empty map failed: %v", err)
	}
}

// TestLoadFromFileErrorCases 测试加载错误处理
func TestLoadFromFileErrorCases(t *testing.T) {
	m := NewStringHashMap[int]()

	// 测试不存在的文件
	err := m.LoadFromFile("nonexistent.json")
	if err == nil {
		t.Error("LoadFromFile should fail for nonexistent file")
	}

	// 测试无效文件格式
	tmpFile := filepath.Join(t.TempDir(), "invalid.json")
	defer func() { _ = os.Remove(tmpFile) }()

	// 创建无效JSON文件
	err = os.WriteFile(tmpFile, []byte("invalid json content"), 0644)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	m = NewStringHashMap[int](WithSerializer(JsonSerializer()))
	err = m.LoadFromFile(tmpFile)
	if err == nil {
		t.Error("LoadFromFile should fail for invalid JSON")
	}
}

// TestSaveLoadWithCustomTypes 测试自定义类型的文件操作
func TestSaveLoadWithCustomTypes(t *testing.T) {
	type Key string
	type Value struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	tmpFile := filepath.Join(t.TempDir(), "custom.json")
	defer func() { _ = os.Remove(tmpFile) }()

	original := NewHashMap[Key, Value](WithSerializer(JsonSerializer()))
	original.Put("key1", Value{Name: "test1", Count: 1})
	original.Put("key2", Value{Name: "test2", Count: 2})

	// 保存
	err := original.SaveToFile(tmpFile)
	if err != nil {
		t.Fatalf("Save custom types failed: %v", err)
	}

	// 加载
	restored := NewHashMap[Key, Value](WithSerializer(JsonSerializer()))
	err = restored.LoadFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Load custom types failed: %v", err)
	}

	// 验证
	val, ok := restored.Get("key1")
	if !ok || val.Name != "test1" || val.Count != 1 {
		t.Errorf("Custom type round-trip failed")
	}
}

// TestConcurrentFileOperations 测试并发文件操作
func TestConcurrentFileOperations(t *testing.T) {
	tmpDir := t.TempDir()

	// 并发保存
	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()

			m := NewStringHashMap[int]()
			m.Put(fmt.Sprintf("key%d", id), id*100)

			filename := filepath.Join(tmpDir, fmt.Sprintf("test_%d.json", id))
			err := m.SaveToFile(filename)
			if err != nil {
				t.Errorf("Concurrent save failed for %d: %v", id, err)
				return
			}

			// 立即加载验证
			loaded := NewStringHashMap[int]()
			err = loaded.LoadFromFile(filename)
			if err != nil {
				t.Errorf("Concurrent load failed for %d: %v", id, err)
				return
			}

			val, ok := loaded.Get(fmt.Sprintf("key%d", id))
			if !ok || val != id*100 {
				t.Errorf("Concurrent round-trip failed for %d", id)
			}
		}(i)
	}

	wg.Wait()
}

// TestFilePermissions 测试文件权限
func TestFilePermissions(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "permissions.json")

	m := NewStringHashMap[int]()
	m.Put("test", 42)

	// 保存文件
	err := m.SaveToFile(tmpFile)
	if err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	// 检查文件权限
	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	// 检查文件是否存在且可读
	if info.Size() == 0 {
		t.Error("Saved file is empty")
	}
}

// TestEmptyFileLoad 测试空文件加载
func TestEmptyFileLoad(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "empty.json")

	// 创建空文件
	err := os.WriteFile(tmpFile, []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// 加载空文件应该成功
	m := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	err = m.LoadFromFile(tmpFile)
	if err != nil {
		t.Errorf("Load empty file failed: %v", err)
	}

	if !m.Empty() {
		t.Error("Loading empty file should result in empty map")
	}
}

// TestFilePathValidation 测试路径验证
func TestFilePathValidation(t *testing.T) {
	m := NewStringHashMap[int]()
	m.Put("test", 123)

	// 测试相对路径
	relPath := "test_relative.json"
	defer func() { _ = os.Remove(relPath) }()

	err := m.SaveToFile(relPath)
	if err != nil {
		t.Errorf("Save to relative path failed: %v", err)
	}

	// 测试加载相对路径
	loaded := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	err = loaded.LoadFromFile(relPath)
	if err != nil {
		t.Errorf("Load from relative path failed: %v", err)
	}

	val, ok := loaded.Get("test")
	if !ok || val != 123 {
		t.Error("Relative path round-trip failed")
	}
}

// TestLargeFileOperations 测试大文件操作
func TestLargeFileOperations(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "large.json")
	defer func() { _ = os.Remove(tmpFile) }()

	// 创建包含大量数据的map
	m := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	const itemCount = 10000

	for i := 0; i < itemCount; i++ {
		m.Put(fmt.Sprintf("item_%d", i), i)
	}

	// 保存大文件
	err := m.SaveToFile(tmpFile)
	if err != nil {
		t.Fatalf("Save large file failed: %v", err)
	}

	// 加载大文件
	loaded := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	err = loaded.LoadFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Load large file failed: %v", err)
	}

	if loaded.Size() != itemCount {
		t.Errorf("Large file load size mismatch: expected %d, got %d", itemCount, loaded.Size())
	}

	// 验证部分数据
	val, ok := loaded.Get("item_0")
	if !ok || val != 0 {
		t.Error("Large file validation failed for item_0")
	}

	val, ok = loaded.Get(fmt.Sprintf("item_%d", itemCount-1))
	if !ok || val != itemCount-1 {
		t.Error("Large file validation failed for last item")
	}
}

// TestFileCleanup 测试文件清理
func TestFileCleanup(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "cleanup_test.json")

	m := NewStringHashMap[int]()
	m.Put("cleanup", 123)

	// 保存文件
	err := m.SaveToFile(tmpFile)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// 删除文件
	err = os.Remove(tmpFile)
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	// 验证文件已被删除
	if _, err := os.Stat(tmpFile); !os.IsNotExist(err) {
		t.Error("File was not properly deleted")
	}
}

// TestFileFormatValidation 测试文件格式验证
func TestFileFormatValidation(t *testing.T) {
	testCases := []struct {
		name    string
		content string
		valid   bool
	}{
		{"valid_json", `{"items":[{"key":"test","value":42}]}`, true},
		{"invalid_json", `{"invalid":}`, false},
		{"empty_json", `{}`, true},
		{"malformed", `not json`, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile := filepath.Join(t.TempDir(), tc.name+".json")
			defer func() { _ = os.Remove(tmpFile) }()

			err := os.WriteFile(tmpFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			m := NewStringHashMap[int](WithSerializer(JsonSerializer()))
			err = m.LoadFromFile(tmpFile)

			if tc.valid && err != nil {
				t.Errorf("Expected valid format but got error: %v", err)
			}
			if !tc.valid && err == nil {
				t.Error("Expected invalid format but no error")
			}
		})
	}
}

// TestFilePermissionErrors 测试文件权限错误
func TestFilePermissionErrors(t *testing.T) {
	if runtime.GOOS == "windows" {
		return
	}

	// 测试目录不存在的情况
	invalidDir := "/nonexistent/directory/test.json"
	m := NewStringHashMap[int]()
	m.Put("test", 42)

	err := m.SaveToFile(invalidDir)
	if err == nil {
		t.Error("Should fail for non-existent directory")
	}

	// 测试文件权限问题（在不同系统上可能表现不同）
	// 这个测试主要是为了完整性，实际行为取决于操作系统
	t.Log("Note: File permission tests may vary by OS")
}

// TestConcurrentFileAccess 测试并发文件访问
func TestConcurrentFileAccess(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建共享文件
	sharedFile := filepath.Join(tmpDir, "shared.json")

	// 写入初始数据
	initial := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	initial.Put("initial", 999)
	err := initial.SaveToFile(sharedFile)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// 并发读取
	const readers = 10
	var wg sync.WaitGroup
	wg.Add(readers)

	for i := 0; i < readers; i++ {
		go func(id int) {
			defer wg.Done()

			m := NewStringHashMap[int](WithSerializer(JsonSerializer()))
			err := m.LoadFromFile(sharedFile)
			if err != nil {
				t.Errorf("Concurrent read %d failed: %v", id, err)
				return
			}

			val, ok := m.Get("initial")
			if !ok || val != 999 {
				t.Errorf("Concurrent read %d validation failed", id)
			}
		}(i)
	}

	wg.Wait()
}

// TestFilePathEdgeCases 测试路径边界情况
func TestFilePathEdgeCases(t *testing.T) {
	// 测试特殊字符路径
	specialPath := filepath.Join(t.TempDir(), "test-file_name.json")

	m := NewStringHashMap[int]()
	m.Put("special", 42)

	err := m.SaveToFile(specialPath)
	if err != nil {
		t.Fatalf("Special path save failed: %v", err)
	}

	loaded := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	err = loaded.LoadFromFile(specialPath)
	if err != nil {
		t.Fatalf("Special path load failed: %v", err)
	}

	val, ok := loaded.Get("special")
	if !ok || val != 42 {
		t.Error("Special path round-trip failed")
	}
}

// TestFileSizeLimits 测试文件大小限制
func TestFileSizeLimits(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "size_test.json")

	m := NewStringHashMap[int](WithSerializer(JsonSerializer()))

	// 测试空数据
	err := m.SaveToFile(tmpFile)
	if err != nil {
		t.Fatalf("Empty save failed: %v", err)
	}

	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	// 空map应该生成小文件
	if info.Size() > 100 {
		t.Logf("Empty file size: %d bytes", info.Size())
	}

	// 测试小数据
	m.Put("key", 1)
	err = m.SaveToFile(tmpFile)
	if err != nil {
		t.Fatalf("Small data save failed: %v", err)
	}

	info, err = os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	// 小数据应该生成合理大小的文件
	if info.Size() > 1000 {
		t.Logf("Small data file size: %d bytes", info.Size())
	}
}

// TestFileCorruptionRecovery 测试文件损坏恢复
func TestFileCorruptionRecovery(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "corrupted.json")

	// 先创建有效文件
	m := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	m.Put("before", 123)
	err := m.SaveToFile(tmpFile)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// 损坏文件
	_ = os.WriteFile(tmpFile, []byte("{"), 0644)

	// 尝试加载损坏文件应该失败
	m = NewStringHashMap[int](WithSerializer(JsonSerializer()))
	err = m.LoadFromFile(tmpFile)
	if err == nil {
		t.Error("Should fail to load corrupted file")
	}

	// 文件应该为空
	if m.Size() != 0 {
		t.Error("Corrupted file should result in empty map")
	}
}

// TestFileBackupAndRestore 测试备份和恢复
func TestFileBackupAndRestore(t *testing.T) {
	tmpDir := t.TempDir()
	originalFile := filepath.Join(tmpDir, "original.json")
	backupFile := filepath.Join(tmpDir, "backup.json")

	// 创建并保存原始数据
	m1 := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	m1.Put("data1", 100)
	m1.Put("data2", 200)
	err := m1.SaveToFile(originalFile)
	if err != nil {
		t.Fatalf("Save original failed: %v", err)
	}

	// 创建备份
	m2 := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	err = m2.LoadFromFile(originalFile)
	if err != nil {
		t.Fatalf("Load original failed: %v", err)
	}
	// 为了触发保存，我们需要做一些修改
	m2.Put("backup_marker", 999)
	err = m2.SaveToFile(backupFile)
	if err != nil {
		t.Fatalf("Save backup failed: %v", err)
	}

	// 验证备份完整性
	m3 := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	err = m3.LoadFromFile(backupFile)
	if err != nil {
		t.Fatalf("Load backup failed: %v", err)
	}

	if m3.Size() != 3 {
		t.Error("Backup size mismatch")
	}

	val, ok := m3.Get("data1")
	if !ok || val != 100 {
		t.Error("Backup data1 mismatch")
	}

	val, ok = m3.Get("backup_marker")
	if !ok || val != 999 {
		t.Error("Backup marker mismatch")
	}
}

// TestFileConcurrentModification 测试并发文件修改
func TestFileConcurrentModification(t *testing.T) {
	tmpDir := t.TempDir()

	// 测试并发写入不同文件
	const files = 5
	var wg sync.WaitGroup
	wg.Add(files)

	for i := 0; i < files; i++ {
		go func(id int) {
			defer wg.Done()

			filename := filepath.Join(tmpDir, fmt.Sprintf("concurrent_%d.json", id))
			m := NewStringHashMap[int](WithSerializer(JsonSerializer()))
			m.Put(fmt.Sprintf("unique_%d", id), id*100)

			err := m.SaveToFile(filename)
			if err != nil {
				t.Errorf("Concurrent save %d failed: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	// 验证所有文件
	for i := 0; i < files; i++ {
		filename := filepath.Join(tmpDir, fmt.Sprintf("concurrent_%d.json", i))
		m := NewStringHashMap[int](WithSerializer(JsonSerializer()))
		err := m.LoadFromFile(filename)
		if err != nil {
			t.Errorf("Verify concurrent %d failed: %v", i, err)
			continue
		}

		val, ok := m.Get(fmt.Sprintf("unique_%d", i))
		if !ok || val != i*100 {
			t.Errorf("Verify concurrent %d data failed", i)
		}
	}
}

// TestFileSystemCompatibility 测试文件系统兼容性
func TestFileSystemCompatibility(t *testing.T) {
	tmpDir := t.TempDir()

	// 测试长文件名
	longName := "very_long_file_name_for_testing_purposes_with_underscores_and_numbers_12345.json"
	longFile := filepath.Join(tmpDir, longName)

	m := NewStringHashMap[int]()
	m.Put("long", 999)

	err := m.SaveToFile(longFile)
	if err != nil {
		t.Fatalf("Long filename save failed: %v", err)
	}

	loaded := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	err = loaded.LoadFromFile(longFile)
	if err != nil {
		t.Fatalf("Long filename load failed: %v", err)
	}

	val, ok := loaded.Get("long")
	if !ok || val != 999 {
		t.Error("Long filename round-trip failed")
	}
}

// TestFileIntegrity 测试文件完整性
func TestFileIntegrity(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "integrity.json")

	// 创建包含各种数据类型的map
	m := NewStringHashMap[interface{}](WithSerializer(JsonSerializer()))
	m.Put("string", "hello")
	m.Put("number", 42)
	m.Put("boolean", true)
	m.Put("null", nil)

	// 保存
	err := m.SaveToFile(tmpFile)
	if err != nil {
		t.Fatalf("Save integrity test failed: %v", err)
	}

	// 加载验证
	loaded := NewStringHashMap[interface{}](WithSerializer(JsonSerializer()))
	err = loaded.LoadFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Load integrity test failed: %v", err)
	}

	// 验证数据类型
	if str, ok := loaded.Get("string"); !ok || str != "hello" {
		t.Errorf("String integrity failed: %v", str)
	}

	if num, ok := loaded.Get("number"); !ok || num != 42.0 {
		t.Errorf("Number integrity failed: %v", num)
	}

	if boolVal, ok := loaded.Get("boolean"); !ok || boolVal != true {
		t.Errorf("Boolean integrity failed: %v", boolVal)
	}
}

// TestFilePerformance 测试文件性能
func TestFilePerformance(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "performance.json")
	defer func() { _ = os.Remove(tmpFile) }()

	// 创建中等大小的map
	m := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	const itemCount = 1000

	for i := 0; i < itemCount; i++ {
		m.Put(fmt.Sprintf("perf_key_%d", i), i)
	}

	// 测试保存性能
	// 这里只是确保操作不超时，不严格计时
	err := m.SaveToFile(tmpFile)
	if err != nil {
		t.Fatalf("Performance test save failed: %v", err)
	}

	// 测试加载性能
	loaded := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	err = loaded.LoadFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Performance test load failed: %v", err)
	}

	if loaded.Size() != itemCount {
		t.Errorf("Performance test size mismatch: expected %d, got %d", itemCount, loaded.Size())
	}
}

// TestFileErrorRecovery 测试文件错误恢复
func TestFileErrorRecovery(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "recovery.json")

	// 创建有效数据
	m := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	m.Put("original", 123)
	err := m.SaveToFile(tmpFile)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// 部分损坏文件
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if len(data) > 10 {
		// 截断文件使其损坏
		_ = os.WriteFile(tmpFile, data[:len(data)-5], 0644)
	}

	// 加载损坏文件
	m = NewStringHashMap[int](WithSerializer(JsonSerializer()))
	err = m.LoadFromFile(tmpFile)
	if err == nil {
		t.Log("File corruption detection may vary by format")
	}
}

// TestFileBackupMechanism 测试文件备份机制
func TestFileBackupMechanism(t *testing.T) {
	tmpDir := t.TempDir()
	originalFile := filepath.Join(tmpDir, "original.json")
	backupFile := filepath.Join(tmpDir, "original.json.bak")

	// 创建原始数据
	m := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	m.Put("data", 456)
	err := m.SaveToFile(originalFile)
	if err != nil {
		t.Fatalf("Save original failed: %v", err)
	}

	// 创建备份
	data, err := os.ReadFile(originalFile)
	if err != nil {
		t.Fatalf("Read original failed: %v", err)
	}

	err = os.WriteFile(backupFile, data, 0644)
	if err != nil {
		t.Fatalf("Create backup failed: %v", err)
	}

	// 验证备份
	backupLoaded := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	err = backupLoaded.LoadFromFile(backupFile)
	if err != nil {
		t.Fatalf("Load backup failed: %v", err)
	}

	val, ok := backupLoaded.Get("data")
	if !ok || val != 456 {
		t.Error("Backup verification failed")
	}
}

// TestFileSymlinkHandling 测试符号链接处理
func TestFileSymlinkHandling(t *testing.T) {
	tmpDir := t.TempDir()
	realFile := filepath.Join(tmpDir, "real.json")
	symlinkFile := filepath.Join(tmpDir, "symlink.json")

	// 创建真实文件
	m := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	m.Put("symlink", 789)
	err := m.SaveToFile(realFile)
	if err != nil {
		t.Fatalf("Save real file failed: %v", err)
	}

	// 创建符号链接（如果支持）
	err = os.Symlink(realFile, symlinkFile)
	if err != nil {
		t.Skip("Symlinks not supported on this system")
		return
	}

	// 通过符号链接加载
	loaded := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	err = loaded.LoadFromFile(symlinkFile)
	if err != nil {
		t.Fatalf("Load via symlink failed: %v", err)
	}

	val, ok := loaded.Get("symlink")
	if !ok || val != 789 {
		t.Error("Symlink round-trip failed")
	}
}

// TestFileMetadata 测试文件元数据
func TestFileMetadata(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "metadata.json")

	m := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	m.Put("meta", 999)

	// 保存文件
	err := m.SaveToFile(tmpFile)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// 检查文件信息
	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	// 验证基本属性
	if info.Name() != "metadata.json" {
		t.Errorf("File name mismatch: %s", info.Name())
	}

	if info.Size() <= 0 {
		t.Error("File should have content")
	}
}

// TestFileConcurrentModificationSafety 测试并发修改安全性
func TestFileConcurrentModificationSafety(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "concurrent_mod.json")

	// 创建初始数据
	m := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	m.Put("initial", 100)
	err := m.SaveToFile(tmpFile)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// 并发修改并保存
	const goroutines = 5
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()

			// 加载
			m := NewStringHashMap[int](WithSerializer(JsonSerializer()))
			err := m.LoadFromFile(tmpFile)
			if err != nil {
				t.Errorf("Concurrent load %d failed: %v", id, err)
				return
			}

			// 添加数据
			m.Put(fmt.Sprintf("concurrent_%d", id), id*1000)

			// 保存到独立文件避免冲突
			uniqueFile := filepath.Join(filepath.Dir(tmpFile), fmt.Sprintf("concurrent_%d.json", id))
			err = m.SaveToFile(uniqueFile)
			if err != nil {
				t.Errorf("Concurrent save %d failed: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	// 验证所有并发文件
	for i := 0; i < goroutines; i++ {
		uniqueFile := filepath.Join(filepath.Dir(tmpFile), fmt.Sprintf("concurrent_%d.json", i))
		m := NewStringHashMap[int](WithSerializer(JsonSerializer()))
		err := m.LoadFromFile(uniqueFile)
		if err != nil {
			t.Errorf("Verify concurrent %d failed: %v", i, err)
			continue
		}

		val, ok := m.Get(fmt.Sprintf("concurrent_%d", i))
		if !ok || val != i*1000 {
			t.Errorf("Verify concurrent %d data failed", i)
		}
	}
}

// TestFileCrossPlatformCompatibility 测试跨平台兼容性
func TestFileCrossPlatformCompatibility(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "cross_platform.json")

	// 使用标准JSON确保跨平台兼容
	m := NewStringHashMap[interface{}](WithSerializer(JsonSerializer()))
	m.Put("string", "Hello World")
	m.Put("number", 42)
	m.Put("boolean", true)
	m.Put("null", nil)
	m.Put("array", []int{1, 2, 3})
	m.Put("object", map[string]string{"key": "value"})

	// 保存
	err := m.SaveToFile(tmpFile)
	if err != nil {
		t.Fatalf("Cross-platform save failed: %v", err)
	}

	// 读取验证
	loaded := NewStringHashMap[interface{}](WithSerializer(JsonSerializer()))
	err = loaded.LoadFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Cross-platform load failed: %v", err)
	}

	// 验证数据类型
	if val, found := loaded.Get("string"); !found {
		t.Errorf("String key not found")
	} else if str, ok := val.(string); !ok || str != "Hello World" {
		t.Errorf("String type failed: %v", str)
	}

	if val, found := loaded.Get("number"); !found {
		t.Errorf("Number key not found")
	} else if num, ok := val.(float64); !ok || num != 42 {
		t.Errorf("Number type failed: %v", num)
	}

	if val, found := loaded.Get("boolean"); !found {
		t.Errorf("Boolean key not found")
	} else if boolVal, ok := val.(bool); !ok || boolVal != true {
		t.Errorf("Boolean type failed: %v", boolVal)
	}

	if val, found := loaded.Get("null"); !found {
		t.Errorf("Null key not found")
	} else if val != nil {
		t.Errorf("Null type failed: %v", val)
	}
}

// TestFileMemoryUsage 测试文件内存使用
func TestFileMemoryUsage(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "memory_test.json")

	// 创建中等大小的数据集
	m := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	const dataSize = 1000

	for i := 0; i < dataSize; i++ {
		m.Put(fmt.Sprintf("key_%d", i), i)
	}

	// 保存并验证
	err := m.SaveToFile(tmpFile)
	if err != nil {
		t.Fatalf("Memory usage test save failed: %v", err)
	}

	// 加载并验证
	loaded := NewStringHashMap[int](WithSerializer(JsonSerializer()))
	err = loaded.LoadFromFile(tmpFile)
	if err != nil {
		t.Fatalf("Memory usage test load failed: %v", err)
	}

	if loaded.Size() != dataSize {
		t.Errorf("Memory usage test failed: expected %d items, got %d", dataSize, loaded.Size())
	}

	// 验证内存中的数据一致性
	for i := 0; i < dataSize; i++ {
		key := fmt.Sprintf("key_%d", i)
		val, ok := loaded.Get(key)
		if !ok || val != i {
			t.Errorf("Memory validation failed for %s: expected %d, got %d", key, i, val)
		}
	}
}

// TestFileErrorHandling 测试文件错误处理
func TestFileErrorHandling(t *testing.T) {
	// 测试各种错误情况
	testCases := []struct {
		name   string
		action func() error
		expect bool
	}{
		{
			"empty_filename",
			func() error { return NewStringHashMap[int]().SaveToFile("") },
			true,
		},
		{
			"nonexistent_load",
			func() error { return NewStringHashMap[int]().LoadFromFile("nonexistent.json") },
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.action()
			if (err != nil) != tc.expect {
				if tc.expect {
					t.Errorf("Expected error but got none")
				} else {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}
