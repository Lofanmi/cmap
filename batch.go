package cmap

// GetMultiple 批量获取
func (m *Map[K, V]) GetMultiple(keys []K) map[K]V {
	result := make(map[K]V)
	for _, key := range keys {
		if value, found := m.Get(key); found {
			result[key] = value
		}
	}
	return result
}

// RemoveMultiple 批量删除
func (m *Map[K, V]) RemoveMultiple(keys []K) {
	for _, key := range keys {
		m.Remove(key)
	}
}

// PutAll 批量插入
func (m *Map[K, V]) PutAll(data map[K]V) {
	for key, value := range data {
		m.Put(key, value)
	}
}
