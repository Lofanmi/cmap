package cmap

import (
	"bytes"
	"cmp"
	"encoding/gob"
	"encoding/json"
	"sync"

	"github.com/bytedance/sonic"
	jsoniter "github.com/json-iterator/go"
)

// ---------------------------------------------------------------------------------------------------------------------

// Tuple 用于序列化的键值对
type Tuple[K cmp.Ordered, V any] struct {
	Key   K `json:"key"`
	Value V `json:"value"`
}

// SerializableData 可序列化的数据结构
type SerializableData[K cmp.Ordered, V any] struct {
	Items []Tuple[K, V] `json:"items"`
}

// ---------------------------------------------------------------------------------------------------------------------

type SerializerFunc struct {
	NameFunc      func() string
	MarshalFunc   func(v interface{}) ([]byte, error)
	UnmarshalFunc func(data []byte, v interface{}) error
	isJSON        bool // 标记是否为JSON序列化器
}

func (s *SerializerFunc) Name() string {
	return s.NameFunc()
}

func (s *SerializerFunc) Marshal(v interface{}) ([]byte, error) {
	return s.MarshalFunc(v)
}

func (s *SerializerFunc) Unmarshal(data []byte, v interface{}) error {
	return s.UnmarshalFunc(data, v)
}

// IsJSON 实现JSONCompatible接口
func (s *SerializerFunc) IsJSON() bool {
	return s.isJSON
}

// ---------------------------------------------------------------------------------------------------------------------

func JsonSerializer() *SerializerFunc {
	return &SerializerFunc{
		NameFunc:      func() string { return "json" },
		MarshalFunc:   json.Marshal,
		UnmarshalFunc: json.Unmarshal,
		isJSON:        true,
	}
}

func JsoniterSerializer() *SerializerFunc {
	return &SerializerFunc{
		NameFunc:      func() string { return "jsoniter" },
		MarshalFunc:   jsoniter.Marshal,
		UnmarshalFunc: jsoniter.Unmarshal,
		isJSON:        true,
	}
}

func SonicSerializer() *SerializerFunc {
	return &SerializerFunc{
		NameFunc:      func() string { return "sonic" },
		MarshalFunc:   sonic.Marshal,
		UnmarshalFunc: sonic.Unmarshal,
		isJSON:        true,
	}
}

// gobBufferPool 用于复用bytes.Buffer
var gobBufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func GobSerializer() *SerializerFunc {
	return &SerializerFunc{
		NameFunc: func() string { return "gob" },
		MarshalFunc: func(v interface{}) ([]byte, error) {
			// 从池中获取buffer
			buf := gobBufferPool.Get().(*bytes.Buffer)
			defer func() {
				buf.Reset()
				gobBufferPool.Put(buf)
			}()

			// 创建encoder
			enc := gob.NewEncoder(buf)
			err := enc.Encode(v)
			if err != nil {
				return nil, err
			}

			// 返回buffer的副本，避免数据竞争
			return append([]byte(nil), buf.Bytes()...), nil
		},
		UnmarshalFunc: func(data []byte, v interface{}) error {
			// 从池中获取buffer
			buf := gobBufferPool.Get().(*bytes.Buffer)
			defer func() {
				buf.Reset()
				gobBufferPool.Put(buf)
			}()

			// 将数据写入buffer
			buf.Write(data)

			// 创建decoder
			dec := gob.NewDecoder(buf)
			return dec.Decode(v)
		},
	}
}

// ---------------------------------------------------------------------------------------------------------------------
