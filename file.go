package cmap

import (
	"fmt"
	"os"
	"path/filepath"
)

// SaveToFile 保存到文件
func (m *Map[K, V]) SaveToFile(filename string) (err error) {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	// 检查文件路径是否有效
	if !filepath.IsAbs(filename) {
		// 对于相对路径，确保目录存在
		dir := filepath.Dir(filename)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		}
	}

	if m.opts.Serializer == nil || m.opts.Serializer.MarshalFunc == nil {
		return fmt.Errorf("no serializer configured for marshaling")
	}

	size := m.Size()
	if size > 0 {
		m.mu.RLock()
		isDirty := m.dirty
		m.mu.RUnlock()
		if !isDirty {
			// 如果未修改，则不保存，避免频繁的文件写入
			return nil
		}
	}

	data, err := m.MarshalWith(m.opts.Serializer)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// 先写入临时文件
	tempFile := filename + ".tmp"
	if err = os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file %s: %w", tempFile, err)
	}

	// 原子性重命名
	if err = os.Rename(tempFile, filename); err != nil {
		_ = os.Remove(tempFile)
		return fmt.Errorf("failed to rename temporary file to %s: %w", filename, err)
	}

	// 标记为未修改
	m.mu.Lock()
	m.dirty = false
	m.mu.Unlock()

	return nil
}

// LoadFromFile 从文件加载
func (m *Map[K, V]) LoadFromFile(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	if m.opts.Serializer == nil || m.opts.Serializer.UnmarshalFunc == nil {
		return fmt.Errorf("no serializer configured for unmarshaling")
	}

	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", filename)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	// 检查文件是否为空
	if len(data) == 0 {
		// 空文件是合法的，清空当前映射
		m.Clear()
		m.mu.Lock()
		m.dirty = false
		m.mu.Unlock()
		return nil
	}

	err = m.UnmarshalWith(data, m.opts.Serializer)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data from %s: %w", filename, err)
	}

	// 加载完成后标记为未修改，因为数据与文件同步
	m.mu.Lock()
	m.dirty = false
	m.mu.Unlock()

	return nil
}
