package cmap

// JSONCompatible 标记支持JSON序列化的接口
type JSONCompatible interface {
	IsJSON() bool
}

// MarshalJSON JSON序列化
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	serializer := m.getJSONSerializer()
	return m.MarshalWith(serializer)
}

// UnmarshalJSON JSON反序列化
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	serializer := m.getJSONSerializer()
	return m.UnmarshalWith(data, serializer)
}

// getJSONSerializer 获取JSON序列化器，使用接口判断
func (m *Map[K, V]) getJSONSerializer() *SerializerFunc {
	serializer := m.opts.Serializer
	if serializer == nil {
		return JsonSerializer()
	}

	// 如果实现了JSONCompatible接口，就直接使用
	if _, ok := any(serializer).(JSONCompatible); ok {
		return serializer
	}

	// 默认回退到标准JSON序列化器
	return JsonSerializer()
}

// MarshalWith 使用指定序列化器进行序列化
func (m *Map[K, V]) MarshalWith(serializer *SerializerFunc) ([]byte, error) {
	items := make([]Tuple[K, V], 0, m.Size())
	for i := range m.shards {
		m.shards[i].mu.RLock()
		keys := m.shards[i].m.Keys()
		for _, key := range keys {
			value, _ := m.shards[i].m.Get(key)
			items = append(items, Tuple[K, V]{Key: key, Value: value})
		}
		m.shards[i].mu.RUnlock()
	}

	data := SerializableData[K, V]{Items: items}

	return serializer.Marshal(data)
}

// UnmarshalWith 使用指定序列化器进行反序列化
func (m *Map[K, V]) UnmarshalWith(data []byte, serializer *SerializerFunc) error {
	var serializableData SerializableData[K, V]
	if err := serializer.Unmarshal(data, &serializableData); err != nil {
		return err
	}

	// 清空现有数据
	m.Clear()

	// 加载数据
	for _, tuple := range serializableData.Items {
		m.Put(tuple.Key, tuple.Value)
	}

	// 加载完成后标记为未修改
	m.mu.Lock()
	m.dirty = false
	m.mu.Unlock()

	return nil
}
